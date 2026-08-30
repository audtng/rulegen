#!/usr/bin/env bash
#
# Script to parse ghsa_golang.json, extract vulnerable and fixed version pairs,
# fix commits, and generate direct git diff endpoints with comprehensive validation.
#
set -euo pipefail

# Help message
if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    echo "Usage: $0 [INPUT_GOLANG_JSON] [OUTPUT_GIT_DIFFS_JSON]"
    echo ""
    echo "Extracts vulnerable and fixed version diffs, commit SHAs, and direct diff URLs"
    echo "with multi-layer validation from INPUT_GOLANG_JSON into OUTPUT_GIT_DIFFS_JSON."
    echo ""
    echo "Arguments:"
    echo "  INPUT_GOLANG_JSON       Path to Golang JSON advisory file (default: /src/ghsa_golang.json)"
    echo "  OUTPUT_GIT_DIFFS_JSON   Path for the extracted git diffs JSON file (default: /src/ghsa_golang_git_diffs.json)"
    echo ""
    exit 0
fi

INPUT_FILE="${1:-/src/ghsa_golang.json}"
OUTPUT_FILE="${2:-/src/ghsa_golang_git_diffs.json}"

# Resolve relative path if default does not exist in absolute path
if [ ! -f "$INPUT_FILE" ]; then
    if [ -f "./ghsa_golang.json" ]; then
        INPUT_FILE="./ghsa_golang.json"
    else
        echo "Error: Input file '$INPUT_FILE' not found." >&2
        exit 1
    fi
fi

# Ensure destination directory exists
mkdir -p "$(dirname "$OUTPUT_FILE")"

echo "=================================================="
echo "GHSA Golang Git Diff & Version Extractor"
echo "=================================================="
echo "Input File  : $INPUT_FILE"
echo "Output File : $OUTPUT_FILE"
echo "Validation  : Multi-Layer Syntax & Order Validation"
echo "--------------------------------------------------"

# Execute extractor in Perl
perl -MTime::HiRes=time -MJSON::PP -e '
    use strict;
    use warnings;

    my ($input_path, $output_path) = @ARGV;

    print "Reading and parsing $input_path...\n";
    my $t0 = time();

    open(my $in, "<", $input_path) or die "Cannot open input file $input_path: $!\n";
    local $/;
    my $raw_json = <$in>;
    close($in);

    my $advisories = decode_json($raw_json);
    die "Expected JSON array in $input_path\n" unless ref($advisories) eq "ARRAY";

    my $total_advisories = scalar @$advisories;
    printf "Loaded %d advisories in %.2fs. Extracting validated Git diff data...\n", $total_advisories, time() - $t0;

    my %diff_index;
    my $stats = {
        total_advisories        => $total_advisories,
        with_any_diff           => 0,
        with_commit_diffs       => 0,
        with_version_diffs      => 0,
        with_pr_diffs           => 0,
        total_commit_diffs      => 0,
        total_version_diffs     => 0,
        total_pr_diffs          => 0,
        without_diff_source     => 0,
    };

    for my $adv (@$advisories) {
        my $adv_id = $adv->{id} // "UNKNOWN";
        my $aliases = $adv->{aliases} // [];
        my $summary = $adv->{summary} // "";
        my $details = $adv->{details} // "";
        my $references = $adv->{references} // [];

        # Determine primary package & repository URL
        my $primary_pkg = "";
        my $repo_url = "";
        my @warnings;

        my $affected = $adv->{affected} // [];
        if (ref($affected) eq "ARRAY") {
            for my $aff (@$affected) {
                my $pkg = $aff->{package} // {};
                my $eco = $pkg->{ecosystem} // "";
                my $name = $pkg->{name} // "";
                if ($eco =~ /^go(lang)?$/i || $name =~ /^(stdlib|toolchain|cmd\/go)$/i || $eco eq "") {
                    $primary_pkg = $name if !$primary_pkg;
                    if ($name =~ m{^github\.com/([^/]+)/([^/]+)}) {
                        $repo_url = "https://github.com/$1/$2" if !$repo_url;
                    } elsif ($name =~ m{^gitlab\.com/([^/]+)/([^/]+)}) {
                        $repo_url = "https://gitlab.com/$1/$2" if !$repo_url;
                    } elsif ($name =~ m{^golang\.org/x/([^/]+)}) {
                        $repo_url = "https://github.com/golang/$1" if !$repo_url;
                    } elsif ($name =~ m{^gopkg\.in/(?:([^/]+)/)?([^/.]+)\.v\d+}) {
                        my $owner = $1 // "go-$2";
                        my $repo = $2;
                        $repo_url = "https://github.com/$owner/$repo" if !$repo_url;
                    } elsif ($name eq "stdlib" || $name eq "cmd/go" || $name eq "toolchain") {
                        $repo_url = "https://github.com/golang/go" if !$repo_url;
                    }
                }
            }
        }

        # If repo URL not found from package name, search references
        if (!$repo_url) {
            for my $ref (@$references) {
                my $url = $ref->{url} // "";
                if ($url =~ m{^(https?://github\.com/[^/]+/[^/]+)}) {
                    $repo_url = $1;
                    last;
                } elsif ($url =~ m{^(https?://gitlab\.com/[^/]+/[^/]+)}) {
                    $repo_url = $1;
                    last;
                }
            }
        }

        # Normalize repo URL (strip trailing .git, slashes)
        if ($repo_url) {
            $repo_url =~ s/\.git$//;
            $repo_url =~ s/\/+$//;
        }

        # 1. Extract & Validate Fix Commits
        my @fix_commits;
        my %seen_commits;

        for my $ref (@$references) {
            my $url = $ref->{url} // "";
            if ($url =~ m{github\.com/([^/]+)/([^/]+)/commit/([0-9a-fA-F]{7,40})}i) {
                my ($owner, $repo, $sha) = ($1, $2, $3);
                $sha = lc($sha);
                next if $seen_commits{$sha}++;

                # Validate SHA hexadecimal format
                if ($sha !~ /^[0-9a-f]{7,40}$/) {
                    push @warnings, "Invalid commit SHA format: $sha";
                    next;
                }

                my $commit_base = "https://github.com/$owner/$repo/commit/$sha";
                push @fix_commits, {
                    commit_sha    => $sha,
                    commit_url    => $commit_base,
                    diff_url      => "$commit_base.diff",
                    patch_url     => "$commit_base.patch",
                    git_show_cmd  => "git show $sha",
                };
            } elsif ($url =~ m{gitlab\.com/([^/]+)/([^/]+)/-/commit/([0-9a-fA-F]{7,40})}i) {
                my ($owner, $repo, $sha) = ($1, $2, $3);
                $sha = lc($sha);
                next if $seen_commits{$sha}++;

                my $commit_base = "https://gitlab.com/$owner/$repo/-/commit/$sha";
                push @fix_commits, {
                    commit_sha    => $sha,
                    commit_url    => $commit_base,
                    diff_url      => "$commit_base.diff",
                    patch_url     => "$commit_base.patch",
                    git_show_cmd  => "git show $sha",
                };
            } elsif ($url =~ m{go\.dev/cl/(\d+)}i || $url =~ m{golang\.org/cl/(\d+)}i) {
                my $cl_num = $1;
                push @fix_commits, {
                    gerrit_cl     => $cl_num,
                    commit_url    => "https://go-review.googlesource.com/c/go/+/$cl_num",
                    diff_url      => "https://go-review.googlesource.com/changes/go~$cl_num/revisions/current/patch",
                    patch_url     => "https://go-review.googlesource.com/changes/go~$cl_num/revisions/current/patch",
                    git_show_cmd  => "# Gerrit Change-Id CL $cl_num",
                };
            }
        }

        # 2. Extract & Validate Pull Requests
        my @pull_requests;
        my %seen_prs;

        for my $ref (@$references) {
            my $url = $ref->{url} // "";
            if ($url =~ m{github\.com/([^/]+)/([^/]+)/pull/(\d+)}i) {
                my ($owner, $repo, $pr_num) = ($1, $2, $3);
                next if $seen_prs{"$owner/$repo/$pr_num"}++;

                my $pr_base = "https://github.com/$owner/$repo/pull/$pr_num";
                push @pull_requests, {
                    pr_number   => int($pr_num),
                    pr_url      => $pr_base,
                    diff_url    => "$pr_base.diff",
                    patch_url   => "$pr_base.patch",
                };
            } elsif ($url =~ m{gitlab\.com/([^/]+)/([^/]+)/-/merge_requests/(\d+)}i) {
                my ($owner, $repo, $mr_num) = ($1, $2, $3);
                next if $seen_prs{"$owner/$repo/$mr_num"}++;

                my $mr_base = "https://gitlab.com/$owner/$repo/-/merge_requests/$mr_num";
                push @pull_requests, {
                    pr_number   => int($mr_num),
                    pr_url      => $mr_base,
                    diff_url    => "$mr_base.diff",
                    patch_url   => "$mr_base.patch",
                };
            }
        }

        # 3. Extract & Validate Version Range Diffs
        my @version_diffs;
        my %seen_version_pairs;

        if (ref($affected) eq "ARRAY") {
            for my $aff (@$affected) {
                my $pkg = $aff->{package} // {};
                my $pkg_eco = $pkg->{ecosystem} // "";
                my $pkg_name = $pkg->{name} // "";
                next unless $pkg_eco =~ /^go(lang)?$/i || $pkg_name =~ /^(stdlib|toolchain|cmd\/go)$/i || $pkg_eco eq "";

                my $target_repo = $repo_url;
                if ($pkg_name =~ m{^github\.com/([^/]+)/([^/]+)}) {
                    $target_repo = "https://github.com/$1/$2";
                }

                if (ref($aff->{ranges}) eq "ARRAY") {
                    for my $r (@{$aff->{ranges}}) {
                        my $events = $r->{events} // [];
                        my $introduced = undef;
                        my $fixed = undef;
                        my $last_affected = undef;

                        for my $e (@$events) {
                            $introduced = $e->{introduced} if exists $e->{introduced};
                            $fixed = $e->{fixed} if exists $e->{fixed};
                            $last_affected = $e->{last_affected} if exists $e->{last_affected};
                        }

                        if ($fixed) {
                            # Version normalization for tag comparisons
                            my $vulnerable_tag = $last_affected // $introduced;
                            
                            # Validation: check for zero/unspecified intro
                            if (!defined $vulnerable_tag || $vulnerable_tag eq "0") {
                                $vulnerable_tag = "previous-release";
                            }

                            my $pair_key = "$pkg_name:$vulnerable_tag:$fixed";
                            next if $seen_version_pairs{$pair_key}++;

                            my $fixed_tag = $fixed;
                            $fixed_tag = "v$fixed_tag" if $fixed_tag =~ /^\d+\./ && $target_repo =~ /github\.com/;

                            my $vuln_tag_formatted = $vulnerable_tag;
                            $vuln_tag_formatted = "v$vuln_tag_formatted" if $vuln_tag_formatted =~ /^\d+\./ && $target_repo =~ /github\.com/;

                            my $compare_url = "";
                            my $diff_url = "";
                            my $git_cmd = "";

                            if ($target_repo && $vulnerable_tag ne "previous-release") {
                                $compare_url = "$target_repo/compare/$vuln_tag_formatted...$fixed_tag";
                                $diff_url = "$target_repo/compare/$vuln_tag_formatted...$fixed_tag.diff";
                                $git_cmd = "git diff $vuln_tag_formatted..$fixed_tag";
                            } elsif ($target_repo) {
                                $compare_url = "$target_repo/releases/tag/$fixed_tag";
                                $diff_url = "$target_repo/commit/$fixed_tag.diff";
                                $git_cmd = "git show $fixed_tag";
                            }

                            # Semantic order validation check
                            if (defined $introduced && $introduced ne "0" && defined $fixed) {
                                if ($introduced eq $fixed) {
                                    push @warnings, "Identical introduced and fixed version: $introduced == $fixed";
                                }
                            }

                            push @version_diffs, {
                                package                  => $pkg_name,
                                introduced_version       => $introduced // "0",
                                fixed_version            => $fixed,
                                last_vulnerable_version  => $last_affected // ($introduced ne "0" ? $introduced : undef),
                                fixed_tag                => $fixed_tag,
                                compare_url              => $compare_url,
                                diff_url                 => $diff_url,
                                git_diff_cmd             => $git_cmd,
                            };
                        }
                    }
                }
            }
        }

        # 4. Compile Record Validation Metadata
        my $has_commit = scalar @fix_commits > 0 ? 1 : 0;
        my $has_version = scalar @version_diffs > 0 ? 1 : 0;
        my $has_pr = scalar @pull_requests > 0 ? 1 : 0;
        my $has_any = ($has_commit || $has_version || $has_pr) ? 1 : 0;

        if ($has_any) {
            $stats->{with_any_diff}++;
        } else {
            $stats->{without_diff_source}++;
            push @warnings, "No direct commit, PR, or fixed version comparison diff found in advisory";
        }

        $stats->{with_commit_diffs}++ if $has_commit;
        $stats->{with_version_diffs}++ if $has_version;
        $stats->{with_pr_diffs}++ if $has_pr;
        $stats->{total_commit_diffs} += scalar @fix_commits;
        $stats->{total_version_diffs} += scalar @version_diffs;
        $stats->{total_pr_diffs} += scalar @pull_requests;

        $diff_index{$adv_id} = {
            aliases         => $aliases,
            package         => $primary_pkg,
            repository_url  => $repo_url,
            summary         => $summary,
            validation      => {
                has_diff_source   => $has_any ? JSON::PP::true : JSON::PP::false,
                has_commit_diff   => $has_commit ? JSON::PP::true : JSON::PP::false,
                has_version_diff  => $has_version ? JSON::PP::true : JSON::PP::false,
                has_pr_diff       => $has_pr ? JSON::PP::true : JSON::PP::false,
                warnings          => \@warnings,
            },
            version_diffs   => \@version_diffs,
            fix_commits     => \@fix_commits,
            pull_requests   => \@pull_requests,
        };
    }

    print "Encoding and writing output to $output_path...\n";
    my $json_encoder = JSON::PP->new->utf8->pretty->canonical;
    my $encoded_json = $json_encoder->encode(\%diff_index);

    open(my $out, ">", $output_path) or die "Cannot create output file $output_path: $!\n";
    print $out $encoded_json;
    close($out);

    my $total_elapsed = time() - $t0;
    my $size_bytes = -s $output_path;
    my $size_mb = $size_bytes / (1024 * 1024);

    print "--------------------------------------------------\n";
    print "VALIDATION & EXTRACTION AUDIT REPORT:\n";
    printf "  Total Go Advisories Processed : %d\n", $stats->{total_advisories};
    printf "  Advisories With Diff Sources  : %d (%.1f%%)\n",
        $stats->{with_any_diff}, ($stats->{with_any_diff} / $stats->{total_advisories}) * 100;
    printf "  Advisories With Fix Commits   : %d (Total commits: %d)\n",
        $stats->{with_commit_diffs}, $stats->{total_commit_diffs};
    printf "  Advisories With Version Diffs : %d (Total version pairs: %d)\n",
        $stats->{with_version_diffs}, $stats->{total_version_diffs};
    printf "  Advisories With PR Diffs      : %d (Total PRs: %d)\n",
        $stats->{with_pr_diffs}, $stats->{total_pr_diffs};
    printf "  Advisories Without Diff Source: %d\n", $stats->{without_diff_source};
    printf "Output file size: %.2f MB (%d bytes) in %.2fs\n", $size_mb, $size_bytes, $total_elapsed;
' "$INPUT_FILE" "$OUTPUT_FILE"

echo "--------------------------------------------------"
echo "Extraction & Validation complete!"
echo "=================================================="
