#!/usr/bin/env bash
#
# Unified Golang Semgrep Rule Engineering Pipeline Engine
# Single-entrypoint engine for target preparation, AST validation, rule synthesis, and ledger auditing.
#
# Usage:
#   bash /src/rulegen/pipeline.sh next [COUNT]              # Find and prepare next actionable advisory (or N advisories)
#   bash /src/rulegen/pipeline.sh prepare <GHSA_ID>         # Prepare workspace and testbed for a specific advisory
#   bash /src/rulegen/pipeline.sh validate <GHSA_ID> [RULE] # Deterministically validate rule against testbed
#   bash /src/rulegen/pipeline.sh report                    # Display metrics and progress
#   bash /src/rulegen/pipeline.sh select [OPTIONS]          # List unhandled candidate IDs
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEFAULT_DATASET="${SCRIPT_DIR}/ghsa_golang_git_diffs_med_high_crit.json"
if [ ! -f "$DEFAULT_DATASET" ]; then
    DEFAULT_DATASET="${SCRIPT_DIR}/ghsa_golang_git_diffs.json"
fi

DATASET_JSON="$DEFAULT_DATASET"
RULES_STORE_DIR="${SCRIPT_DIR}/rules/go"
LEDGER_FILE="${SCRIPT_DIR}/validation_ledger.json"
WORKSPACES_DIR="${SCRIPT_DIR}/workspaces"

mkdir -p "$RULES_STORE_DIR" "$WORKSPACES_DIR"

show_help() {
    cat << 'EOF'
Antigravity Golang Semgrep Pipeline Engine

Commands:
  next [COUNT]                  Single-pass scan of dataset to prepare the next COUNT (default: 1) actionable targets.
                                Emits structured JSON ready for subagent ingestion.
  prepare <GHSA_ID>             Fetch commit diff, verify actionability, run heuristic check,
                                and prepare workspace/prompt for the given advisory ID.
  validate <GHSA_ID> [RULE_FILE] Validate a Semgrep rule (vuln.go >= 1, fixed.go == 0) and commit to store.
  report                        Display current generation metrics, ledger stats, and progress.
  select [-n N] [-s SEV]        Query unhandled advisory candidates from the dataset.
EOF
}

cmd_select() {
    local limit=10
    local severity="ALL"
    local skip_approved=1
    local failed_only=0

    while [[ $# -gt 0 ]]; do
        case "$1" in
            -n|--limit) limit="$2"; shift 2 ;;
            -s|--severity) severity="$(echo "$2" | tr '[:lower:]' '[:upper:]')"; shift 2 ;;
            -a|--all) skip_approved=0; shift ;;
            -f|--failed-only) failed_only=1; shift ;;
            -d|--dataset) DATASET_JSON="$2"; shift 2 ;;
            *) shift ;;
        esac
    done

    perl -MJSON::PP -e '
        use strict;
        use warnings;

        my ($dataset_path, $rules_dir, $ledger_path, $limit, $severity_filter, $skip_approved, $failed_only) = @ARGV;
        my $dataset;
        {
            open(my $dfh, "<", $dataset_path) or die "Cannot open dataset $dataset_path: $!\n";
            local $/;
            $dataset = decode_json(<$dfh>);
            close($dfh);
        }

        my $ledger = {};
        if (-f $ledger_path && open(my $lfh, "<", $ledger_path)) {
            local $/;
            eval { $ledger = decode_json(<$lfh>); };
            close($lfh);
        }

        my @selected;
        my $count = 0;
        for my $id (sort keys %$dataset) {
            my $entry = $dataset->{$id};
            next unless ($entry->{validation}{has_commit_diff});
            next unless (ref($entry->{fix_commits}) eq "ARRAY" && @{$entry->{fix_commits}});
            my $patch_url = $entry->{fix_commits}[0]{patch_url} // $entry->{fix_commits}[0]{diff_url} // "";
            next if ($patch_url eq "");
            next unless ($patch_url =~ m{^https?://(?:github\.com|gitlab\.com|codeberg\.org)}i || $patch_url =~ m{\.(?:patch|diff)$}i);

            if ($severity_filter ne "ALL") {
                my $entry_sev = uc($entry->{severity} // $entry->{database_specific}{severity} // "");
                next unless ($entry_sev eq $severity_filter);
            }

            my $rule_file = "$rules_dir/$id.yaml";
            my $is_approved = (-f $rule_file || ($ledger->{$id} && $ledger->{$id}{passed})) ? 1 : 0;

            if ($failed_only) {
                my $is_failed = ($ledger->{$id} && !$ledger->{$id}{passed}) ? 1 : 0;
                next unless $is_failed;
            } elsif ($skip_approved && $is_approved) {
                next;
            }

            push @selected, $id;
            $count++;
            last if ($limit > 0 && $count >= $limit);
        }
        for my $id (@selected) {
            print "$id\n";
        }
    ' "$DATASET_JSON" "$RULES_STORE_DIR" "$LEDGER_FILE" "$limit" "$severity" "$skip_approved" "$failed_only"
}

cmd_validate() {
    local adv_id="$1"
    local ws_dir="${WORKSPACES_DIR}/${adv_id}"
    local rfile="${2:-${ws_dir}/rule.yaml}"
    local vfile="${ws_dir}/vuln.go"
    local ffile="${ws_dir}/fixed.go"
    local resfile="${ws_dir}/validation_result.json"

    if [ ! -f "$rfile" ]; then
        echo "Error: Rule file '$rfile' not found." >&2
        return 1
    fi
    if [ ! -f "$vfile" ] || [ ! -f "$ffile" ]; then
        echo "Error: Testbed files not found in '$ws_dir'." >&2
        return 1
    fi

    perl -MJSON::PP -e '
        use strict;
        use warnings;
        my ($ghsa_id, $rule_path, $vuln_path, $fixed_path, $result_path, $rules_store, $ledger_path) = @ARGV;

        my $rule_content;
        {
            open(my $rf, "<", $rule_path) or die $!;
            local $/;
            $rule_content = <$rf>;
            close($rf);
        }

        my ($vuln_code, $fixed_code);
        {
            open(my $vf, "<", $vuln_path) or die $!;
            local $/;
            $vuln_code = <$vf>;
            close($vf);

            open(my $ff, "<", $fixed_path) or die $!;
            local $/;
            $fixed_code = <$ff>;
            close($ff);
        }

        my @syntax_errors;
        push @syntax_errors, "Missing rules: declaration" if ($rule_content !~ /rules:\s*\n/);
        push @syntax_errors, "Missing id: field" if ($rule_content !~ /id:\s*[a-zA-Z0-9_-]+/);
        push @syntax_errors, "Missing languages: [go]" if ($rule_content !~ /languages:\s*\[.*go.*\]/i && $rule_content !~ /languages:\s*\n\s*-\s*go/i);

        my (@patterns, @pattern_nots, @pattern_insides, @pattern_not_insides);
        my @lines = split("\n", $rule_content);
        for my $l (@lines) {
            if ($l =~ /^\s*-\s*pattern:\s*(.*)/) {
                push @patterns, $1 if $1 && $1 !~ /^\|/;
            } elsif ($l =~ /^\s*-\s*pattern-not:\s*(.*)/) {
                push @pattern_nots, $1 if $1 && $1 !~ /^\|/;
            } elsif ($l =~ /^\s*-\s*pattern-inside:\s*(.*)/) {
                push @pattern_insides, $1 if $1 && $1 !~ /^\|/;
            } elsif ($l =~ /^\s*-\s*pattern-not-inside:\s*(.*)/) {
                push @pattern_not_insides, $1 if $1 && $1 !~ /^\|/;
            } elsif ($l =~ /^\s*pattern:\s*(.*)/) {
                push @patterns, $1 if $1 && $1 !~ /^\|/;
            }
        }

        sub pat2re {
            my ($pat) = @_;
            $pat =~ s/^\s+|\s+$//g;
            $pat =~ s/^["\x27]|["\x27]$//g;
            my $esc = quotemeta($pat);
            $esc =~ s/\\\.\\\.\\\./\.\*\?/g;
            $esc =~ s/\\\$[A-Z0-9_]+/[a-zA-Z0-9_.()"]+/g;
            return $esc;
        }

        my (@vuln_matches, @fixed_matches);
        my $semgrep = `which semgrep 2>/dev/null`;
        chomp($semgrep);

        if ($semgrep) {
            my $vout = `$semgrep --config \x27$rule_path\x27 \x27$vuln_path\x27 --json 2>/dev/null`;
            my $vj = eval { decode_json($vout) };
            if ($vj && $vj->{results}) {
                for my $r (@{$vj->{results}}) {
                    push @vuln_matches, { line => $r->{start}{line}, match => $r->{extra}{lines} };
                }
            }
            my $fout = `$semgrep --config \x27$rule_path\x27 \x27$fixed_path\x27 --json 2>/dev/null`;
            my $fj = eval { decode_json($fout) };
            if ($fj && $fj->{results}) {
                for my $r (@{$fj->{results}}) {
                    push @fixed_matches, { line => $r->{start}{line}, match => $r->{extra}{lines} };
                }
            }
        } else {
            for my $pat (@patterns) {
                my $re = pat2re($pat);
                my @v_lines = split("\n", $vuln_code);
                for my $idx (0 .. $#v_lines) {
                    if ($v_lines[$idx] =~ /$re/i) {
                        push @vuln_matches, { line => $idx + 1, match => $v_lines[$idx] };
                    }
                }
                my @f_lines = split("\n", $fixed_code);
                for my $idx (0 .. $#f_lines) {
                    if ($f_lines[$idx] =~ /$re/i) {
                        my $neg = 0;
                        for my $np (@pattern_nots, @pattern_not_insides) {
                            my $nre = pat2re($np);
                            if ($f_lines[$idx] =~ /$nre/i || $fixed_code =~ /$nre/i) {
                                $neg = 1;
                                last;
                            }
                        }
                        push @fixed_matches, { line => $idx + 1, match => $f_lines[$idx] } unless $neg;
                    }
                }
            }
        }

        my $status = "REJECTED";
        my @reasons;
        if (@syntax_errors) {
            $status = "REJECTED";
            push @reasons, "Syntax errors: " . join("; ", @syntax_errors);
        } elsif (scalar @vuln_matches == 0) {
            $status = "REJECTED";
            push @reasons, "FAILED_TRUE_POSITIVE: 0 matches in vuln.go";
        } elsif (scalar @fixed_matches > 0) {
            $status = "REJECTED";
            push @reasons, "FAILED_FALSE_POSITIVE: " . scalar(@fixed_matches) . " matches in fixed.go";
        } else {
            $status = "APPROVED";
            push @reasons, "PASS: Matched " . scalar(@vuln_matches) . " line(s) in vuln.go, 0 in fixed.go";
        }

        my %result = (
            advisory_id            => $ghsa_id,
            status                 => $status,
            passed                 => ($status eq "APPROVED") ? JSON::PP::true : JSON::PP::false,
            vulnerable_matches_cnt => scalar @vuln_matches,
            fixed_matches_cnt      => scalar @fixed_matches,
            reasons                => \@reasons,
            syntax_errors          => \@syntax_errors,
            timestamp              => scalar gmtime() . " UTC",
        );

        open(my $res_fh, ">", $result_path) or die $!;
        print $res_fh JSON::PP->new->utf8->pretty->encode(\%result);
        close($res_fh);

        if ($status eq "APPROVED") {
            my $target_rule = "$rules_store/$ghsa_id.yaml";
            open(my $tf, ">", $target_rule) or die $!;
            print $tf $rule_content;
            close($tf);
            print "SUCCESS: Rule APPROVED -> $target_rule\n";
        } else {
            print "REJECTED: " . join("; ", @reasons) . "\n";
        }

        my $ledger_data = {};
        if (-f $ledger_path && open(my $lf, "<", $ledger_path)) {
            local $/;
            $ledger_data = eval { decode_json(<$lf>) } // {};
            close($lf);
        }
        $ledger_data->{$ghsa_id} = \%result;
        if (open(my $lf, ">", $ledger_path)) {
            print $lf JSON::PP->new->utf8->pretty->encode($ledger_data);
            close($lf);
        }
        exit($status eq "APPROVED" ? 0 : 1);
    ' "$adv_id" "$rfile" "$vfile" "$ffile" "$resfile" "$RULES_STORE_DIR" "$LEDGER_FILE"
}

cmd_next() {
    local target_count="${1:-1}"
    local target_adv="${2:-}"

    perl -MJSON::PP -e '
        use strict;
        use warnings;

        $| = 1;

        my ($dataset_path, $rules_dir, $ledger_path, $workspaces_dir, $target_count, $specific_id) = @ARGV;

        my $dataset;
        {
            open(my $dfh, "<", $dataset_path) or die "Cannot open dataset $dataset_path: $!\n";
            local $/;
            $dataset = decode_json(<$dfh>);
            close($dfh);
        }

        my $ledger = {};
        if (-f $ledger_path && open(my $lfh, "<", $ledger_path)) {
            local $/;
            eval { $ledger = decode_json(<$lfh>); };
            close($lfh);
        }

        my @candidates;
        if ($specific_id) {
            push @candidates, $specific_id if $dataset->{$specific_id};
        } else {
            for my $id (sort keys %$dataset) {
                my $entry = $dataset->{$id};
                next unless ($entry->{validation}{has_commit_diff});
                next unless (ref($entry->{fix_commits}) eq "ARRAY" && @{$entry->{fix_commits}});
                my $patch_url = $entry->{fix_commits}[0]{patch_url} // $entry->{fix_commits}[0]{diff_url} // "";
                next if ($patch_url eq "");
                next unless ($patch_url =~ m{^https?://(?:github\.com|gitlab\.com|codeberg\.org)}i || $patch_url =~ m{\.(?:patch|diff)$}i);

                my $rule_file = "$rules_dir/$id.yaml";
                my $is_approved = (-f $rule_file || ($ledger->{$id} && $ledger->{$id}{passed})) ? 1 : 0;
                next if $is_approved;
                push @candidates, $id;
            }
        }

        my $prepared = 0;
        for my $id (@candidates) {
            my $entry = $dataset->{$id};
            my $patch_url = $entry->{fix_commits}[0]{patch_url} // $entry->{fix_commits}[0]{diff_url} // "";
            my $commit_sha = $entry->{fix_commits}[0]{commit_sha} // "";
            my $ws = "$workspaces_dir/$id";
            system("mkdir", "-p", $ws);

            my $patch_file = "$ws/commit.patch";
            if (! -s $patch_file) {
                my $ret = system("curl", "-sL", "--connect-timeout", "4", "--max-time", "8", "-o", $patch_file, $patch_url);
                next if ($ret != 0 || ! -s $patch_file);
            }

            open(my $pf, "<", $patch_file) or next;
            my @lines = <$pf>;
            close($pf);

            my @vuln;
            my @fixed;
            my @hunk_diff;
            my $in_go = 0;
            my $go_file_name = "";

            for my $l (@lines) {
                if ($l =~ /^diff --git/i) {
                    $in_go = ($l =~ /\.go/i) ? 1 : 0;
                    if ($l =~ /b\/(.+\.go)/i && !$go_file_name) {
                        $go_file_name = $1;
                        $go_file_name =~ s/\r$//;
                    }
                    push @hunk_diff, $l if $in_go;
                    next;
                }
                if ($in_go) {
                    push @hunk_diff, $l;
                    next if $l =~ /^(---|Index:|\+\+\+|index\s+|new file|deleted file|@@)/;
                    if ($l =~ /^\+(.*)/) {
                        push @fixed, "$1\n";
                    } elsif ($l =~ /^-(.*)/) {
                        push @vuln, "$1\n";
                    } elsif ($l =~ /^[ \t](.*)/) {
                        push @vuln, "$1\n";
                        push @fixed, "$1\n";
                    }
                }
            }

            my $vuln_code = join("", @vuln);
            my $fixed_code = join("", @fixed);
            my $diff_hunk = join("", @hunk_diff);

            open(my $vf, ">", "$ws/vuln.go") or next;
            print $vf "package main\n\n" . $vuln_code;
            close($vf);

            open(my $ff, ">", "$ws/fixed.go") or next;
            print $ff "package main\n\n" . $fixed_code;
            close($ff);

            my %meta = (
                advisory_id   => $id,
                package       => $entry->{package} // "",
                summary       => $entry->{summary} // "",
                severity      => $entry->{severity} // "MODERATE",
                cvss_score    => $entry->{cvss_score} // "",
                aliases       => $entry->{aliases} // [],
                commit_sha    => $commit_sha,
                patch_url     => $patch_url,
                file_modified => $go_file_name || "target.go",
                diff_hunk     => $diff_hunk,
            );

            open(my $mf, ">", "$ws/metadata.json") or next;
            print $mf JSON::PP->new->utf8->pretty->encode(\%meta);
            close($mf);

            # Actionability verification
            my $vsize = -s "$ws/vuln.go" // 0;
            my $fsize = -s "$ws/fixed.go" // 0;
            next if ($vsize <= 20 || $fsize <= 20);
            next if ($go_file_name =~ /go\.(mod|sum)$/i && $vsize < 100);

            my %fixed_set;
            for my $fl (split("\n", $fixed_code)) {
                $fl =~ s/^\s+|\s+$//g;
                $fixed_set{$fl} = 1 if length($fl) > 0;
            }

            my $has_diff = 0;
            for my $vl (split("\n", $vuln_code)) {
                my $t = $vl;
                $t =~ s/^\s+|\s+$//g;
                next if length($t) < 4;
                next if $t =~ /^(\/\/|\/\*|\*|\{|\}|package\s+)/;
                if (!$fixed_set{$t}) {
                    $has_diff = 1;
                    last;
                }
            }
            next unless $has_diff;

            # Prompt compilation
            my $pkg = $entry->{package} // "Golang Package";
            my $summary = $entry->{summary} // "Security Vulnerability";
            my $aliases_str = join(", ", @{$entry->{aliases} // []});
            my $truncated_diff = (length($diff_hunk) > 2500) ? (substr($diff_hunk, 0, 2500) . "\n... [diff truncated]") : $diff_hunk;

            my $prompt = <<"EOF";
You are an expert security engineer and Semgrep rule author specializing in Golang security.

### TASK:
Synthesize a precise, high-accuracy Semgrep YAML rule for the security vulnerability in package \x27$pkg\x27 ($id).

### VULNERABILITY CONTEXT:
- Advisory ID: $id
- Aliases: $aliases_str
- Package: $pkg
- Summary: $summary
- File Modified: $go_file_name

### SECURITY FIX DIFF:
\`\`\`diff
$truncated_diff
\`\`\`

### OBJECTIVE & REQUIREMENTS:
1. Identify the unsafe pattern in the pre-patch code.
2. Synthesize a Semgrep rule targeting Go (\`languages: [go]\`).
3. Match \x27$ws/vuln.go\x27 (True Positive >= 1).
4. Do NOT match \x27$ws/fixed.go\x27 (False Positive == 0).
5. Write rule to \x27$ws/rule.yaml\x27.
6. Validate inside subagent with: \`bash /src/rulegen/pipeline.sh validate $id /src/rulegen/workspaces/$id/rule.yaml\`
7. Iterate and refine until validation passes (True Positives >= 1, False Positives == 0).

### CONSTRAINTS:
- NO curl, wget, or network/API calls.
- NO ad-hoc scripts. Use ONLY \`bash /src/rulegen/pipeline.sh validate $id\`.

### REQUIRED YAML FORMAT:
\`\`\`yaml
rules:
  - id: $id
    languages: [go]
    severity: WARNING
    message: "$summary"
    metadata:
      cve: "$aliases_str"
      ghsa: "$id"
      confidence: HIGH
    pattern-either:
      - pattern: <UNSAFE_PATTERN>
\`\`\`
EOF
            open(my $prf, ">", "$ws/prompt.txt") or next;
            print $prf $prompt;
            close($prf);

            my %out = (
                status        => "READY_FOR_AGENT",
                advisory_id   => $id,
                package       => $pkg,
                aliases       => $entry->{aliases} // [],
                workspace     => $ws,
                prompt_file   => "$ws/prompt.txt",
                vuln_file     => "$ws/vuln.go",
                fixed_file    => "$ws/fixed.go"
            );
            print JSON::PP->new->utf8->canonical->encode(\%out) . "\n";
            $prepared++;
            last if ($target_count > 0 && $prepared >= $target_count);
        }
    ' "$DATASET_JSON" "$RULES_STORE_DIR" "$LEDGER_FILE" "$WORKSPACES_DIR" "$target_count" "$target_adv"
}

cmd_report() {
    perl -MJSON::PP -e '
        use strict;
        use warnings;

        my ($dataset_path, $rules_dir, $ledger_path) = @ARGV;
        my $dataset = {};
        if (-f $dataset_path && open(my $df, "<", $dataset_path)) {
            local $/;
            $dataset = eval { decode_json(<$df>) } // {};
            close($df);
        }

        my $ledger = {};
        if (-f $ledger_path && open(my $lf, "<", $ledger_path)) {
            local $/;
            $ledger = eval { decode_json(<$lf>) } // {};
            close($lf);
        }

        opendir(my $dh, $rules_dir) or die $!;
        my @rule_files = grep { /\.ya?ml$/ && -f "$rules_dir/$_" } readdir($dh);
        closedir($dh);

        my $total_adv = scalar keys %$dataset;
        my $verified_targets = 0;
        for my $id (keys %$dataset) {
            my $e = $dataset->{$id};
            $verified_targets++ if ($e->{validation}{has_commit_diff});
        }

        my $passed_cnt = 0;
        my $failed_cnt = 0;
        for my $id (keys %$ledger) {
            if ($ledger->{$id}{passed}) { $passed_cnt++; }
            else { $failed_cnt++; }
        }

        my $total_rules = scalar @rule_files;
        my $remaining = $verified_targets - $total_rules;
        $remaining = 0 if $remaining < 0;
        my $rate = ($passed_cnt + $failed_cnt > 0) ? ($passed_cnt / ($passed_cnt + $failed_cnt) * 100.0) : 100.0;

        print "======================================================================\n";
        print "        GOLANG SEMGREP RULE GENERATION PROGRESS & METRICS             \n";
        print "======================================================================\n";
        printf " Dataset Source              : %s\n", $dataset_path;
        printf " Total Advisories in Dataset : %d\n", $total_adv;
        printf " Verified Commit Diff Targets: %d\n", $verified_targets;
        print "----------------------------------------------------------------------\n";
        printf " Verified Rules in Store     : %d\n", $total_rules;
        printf " Ledger Verified (APPROVED)  : %d\n", $passed_cnt;
        printf " Ledger Rejected (FAILED)    : %d\n", $failed_cnt;
        printf " Validation Approval Rate    : %.1f%%\n", $rate;
        printf " Remaining Actionable Targets: %d\n", $remaining;
        print "======================================================================\n";
    ' "$DATASET_JSON" "$RULES_STORE_DIR" "$LEDGER_FILE"
}

# Main Command Dispatcher
case "${1:-}" in
    next)
        shift
        cmd_next "${1:-1}" "${2:-}"
        ;;
    prepare)
        shift
        if [ $# -lt 1 ]; then echo "Usage: $0 prepare <GHSA_ID>"; exit 1; fi
        cmd_next 1 "$1"
        ;;
    validate)
        shift
        if [ $# -lt 1 ]; then echo "Usage: $0 validate <GHSA_ID> [RULE_FILE]"; exit 1; fi
        cmd_validate "$@"
        ;;
    report|status)
        cmd_report
        ;;
    select|list)
        shift
        cmd_select "$@"
        ;;
    -h|--help|help)
        show_help
        ;;
    *)
        show_help
        exit 1
        ;;
esac
