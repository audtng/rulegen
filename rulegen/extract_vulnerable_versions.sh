#!/usr/bin/env bash
#
# Script to parse ghsa_golang.json and generate a package-centric index
# of all vulnerable Go dependencies and their affected versions (Schema A).
#
set -euo pipefail

# Help message
if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    echo "Usage: $0 [INPUT_GOLANG_JSON] [OUTPUT_VERSIONS_JSON]"
    echo ""
    echo "Parses INPUT_GOLANG_JSON and generates a package-centric mapping of vulnerable Go versions in OUTPUT_VERSIONS_JSON."
    echo ""
    echo "Arguments:"
    echo "  INPUT_GOLANG_JSON     Path to Golang JSON advisory file (default: /src/ghsa_golang.json)"
    echo "  OUTPUT_VERSIONS_JSON  Path for the extracted versions JSON file (default: /src/ghsa_golang_vulnerable_versions.json)"
    echo ""
    exit 0
fi

INPUT_FILE="${1:-/src/ghsa_golang.json}"
OUTPUT_FILE="${2:-/src/ghsa_golang_vulnerable_versions.json}"

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
echo "Golang Vulnerable Dependency Versions Extractor"
echo "=================================================="
echo "Input File  : $INPUT_FILE"
echo "Output File : $OUTPUT_FILE"
echo "Schema      : Package-Centric Index (Schema A)"
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
    printf "Loaded %d advisories in %.2fs. Processing vulnerable Go packages...\n", $total_advisories, time() - $t0;

    my %package_index;
    my $total_vulnerability_entries = 0;

    for my $adv (@$advisories) {
        my $adv_id = $adv->{id} // "UNKNOWN";
        my $aliases = $adv->{aliases} // [];
        my $summary = $adv->{summary} // "";
        my $details = $adv->{details} // "";

        my $affected = $adv->{affected} // [];
        if (ref($affected) eq "ARRAY") {
            for my $aff (@{$affected}) {
                my $pkg = $aff->{package} // {};
                my $pkg_name = $pkg->{name} // "";
                my $eco = $pkg->{ecosystem} // "";

                # Check if this specific package entry belongs to the Go ecosystem
                my $is_go_pkg = 0;
                if ($eco =~ /^go(lang)?$/i) {
                    $is_go_pkg = 1;
                } elsif ($pkg_name =~ /^(stdlib|toolchain|cmd\/go)$/i) {
                    $is_go_pkg = 1;
                } elsif ($eco eq "") {
                    # For unreviewed advisories without explicit ecosystem, check if package name is a Go module path
                    if ($pkg_name =~ m{^(golang\.org|github\.com/|gopkg\.in/|google\.golang\.org/|cloud\.google\.com/go)}i) {
                        $is_go_pkg = 1;
                    } elsif (!$pkg_name && $details =~ /package\s+([a-zA-Z0-9.\/_-]+)/i) {
                        $pkg_name = $1;
                        $is_go_pkg = 1 if $pkg_name =~ m{^(golang\.org|github\.com/|gopkg\.in/|google\.golang\.org/)}i;
                    }
                }

                next unless $is_go_pkg && $pkg_name;

                # Normalize package name (strip trailing slashes, whitespace)
                $pkg_name =~ s/^\s+//;
                $pkg_name =~ s/\s+$//;
                $pkg_name =~ s/\/+$//;

                my @fixed_versions;
                my @range_strings;
                my @normalized_ranges;

                if (ref($aff->{ranges}) eq "ARRAY") {
                    for my $r (@{$aff->{ranges}}) {
                        my $range_type = $r->{type} // "ECOSYSTEM";
                        my $events = $r->{events} // [];
                        my %range_obj = (
                            type => $range_type,
                            events => $events,
                        );
                        push @normalized_ranges, \%range_obj;

                        my $introduced = undef;
                        my $fixed = undef;
                        my $last_affected = undef;

                        for my $e (@$events) {
                            if (exists $e->{introduced}) {
                                $introduced = $e->{introduced};
                            }
                            if (exists $e->{fixed}) {
                                $fixed = $e->{fixed};
                                push @fixed_versions, $fixed;
                            }
                            if (exists $e->{last_affected}) {
                                $last_affected = $e->{last_affected};
                            }
                        }

                        # Construct human-readable range expression
                        if (defined $introduced && defined $fixed) {
                            if ($introduced eq "0") {
                                push @range_strings, "< $fixed";
                            } else {
                                push @range_strings, ">= $introduced, < $fixed";
                            }
                        } elsif (defined $introduced && defined $last_affected) {
                            if ($introduced eq "0") {
                                push @range_strings, "<= $last_affected";
                            } else {
                                push @range_strings, ">= $introduced, <= $last_affected";
                            }
                        } elsif (defined $fixed) {
                            push @range_strings, "< $fixed";
                        } elsif (defined $introduced && $introduced ne "0") {
                            push @range_strings, ">= $introduced";
                        } elsif (defined $last_affected) {
                            push @range_strings, "<= $last_affected";
                        }
                    }
                }

                my $explicit_versions = $aff->{versions} // [];
                my $last_known = $aff->{database_specific}{last_known_affected_version_range} // undef;

                my $computed_affected_range = "";
                if (@range_strings) {
                    $computed_affected_range = join("; ", @range_strings);
                } elsif ($last_known) {
                    $computed_affected_range = $last_known;
                } elsif (@$explicit_versions) {
                    $computed_affected_range = join(", ", @$explicit_versions);
                }

                # Deduplicate fixed versions
                my %seen_fixed;
                my @unique_fixed = grep { !$seen_fixed{$_}++ } @fixed_versions;

                my %entry = (
                    advisory_id => $adv_id,
                    aliases     => $aliases,
                    summary     => $summary,
                    affected_range => $computed_affected_range,
                    fixed_versions => \@unique_fixed,
                    explicit_versions => $explicit_versions,
                    ranges      => \@normalized_ranges,
                );

                if ($last_known) {
                    $entry{last_known_affected_version_range} = $last_known;
                }

                $package_index{$pkg_name} //= [];
                push @{$package_index{$pkg_name}}, \%entry;
                $total_vulnerability_entries++;
            }
        }
    }

    my $unique_packages = scalar keys %package_index;
    print "Extracted $total_vulnerability_entries Go vulnerability entries across $unique_packages unique Go packages.\n";
    print "Writing output to $output_path...\n";

    my $json_encoder = JSON::PP->new->utf8->pretty->canonical;
    my $encoded_json = $json_encoder->encode(\%package_index);

    open(my $out, ">", $output_path) or die "Cannot create output file $output_path: $!\n";
    print $out $encoded_json;
    close($out);

    my $total_elapsed = time() - $t0;
    my $size_bytes = -s $output_path;
    my $size_mb = $size_bytes / (1024 * 1024);

    print "--------------------------------------------------\n";
    printf "Extracted %d unique Go packages into %s in %.2fs\n", $unique_packages, $output_path, $total_elapsed;
    printf "Output file size: %.2f MB (%d bytes)\n", $size_mb, $size_bytes;
' "$INPUT_FILE" "$OUTPUT_FILE"

echo "--------------------------------------------------"
echo "Extraction complete!"
echo "=================================================="
