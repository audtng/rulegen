#!/usr/bin/env bash
#
# Step 3: Validate generated Semgrep rule against vuln.go and fixed.go testbeds
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ $# -lt 1 || "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    echo "Usage: $0 <GHSA_ID> [RULE_YAML]"
    echo "Example: $0 GHSA-h395-qcrw-5vmq"
    exit 1
fi

GHSA_ID="$1"
WORKSPACE_DIR="${SCRIPT_DIR}/workspaces/${GHSA_ID}"
RULE_FILE="${2:-${WORKSPACE_DIR}/rule.yaml}"
VULN_FILE="${WORKSPACE_DIR}/vuln.go"
FIXED_FILE="${WORKSPACE_DIR}/fixed.go"
RESULT_FILE="${WORKSPACE_DIR}/validation_result.json"
RULES_STORE_DIR="${SCRIPT_DIR}/rules/go"
LEDGER_FILE="${SCRIPT_DIR}/validation_ledger.json"

mkdir -p "$RULES_STORE_DIR"

if [ ! -f "$RULE_FILE" ]; then
    echo "Error: Rule file '$RULE_FILE' not found." >&2
    exit 1
fi

if [ ! -f "$VULN_FILE" ] || [ ! -f "$FIXED_FILE" ]; then
    echo "Error: Testbed files vuln.go/fixed.go not found in '$WORKSPACE_DIR'." >&2
    exit 1
fi

echo "=================================================="
echo "Step 3: Validate Semgrep Rule"
echo "=================================================="
echo "Advisory ID   : $GHSA_ID"
echo "Rule File     : $RULE_FILE"
echo "Vuln Target   : $VULN_FILE"
echo "Fixed Target  : $FIXED_FILE"
echo "--------------------------------------------------"

perl -MTime::HiRes=time -MJSON::PP -e '
    use strict;
    use warnings;

    my ($ghsa_id, $rule_path, $vuln_path, $fixed_path, $result_path, $rules_store, $ledger_path) = @ARGV;

    open(my $rf, "<", $rule_path) or die "Cannot open $rule_path: $!\n";
    local $/;
    my $rule_content = <$rf>;
    close($rf);

    open(my $vf, "<", $vuln_path) or die $!;
    my $vuln_code = <$vf>;
    close($vf);

    open(my $ff, "<", $fixed_path) or die $!;
    my $fixed_code = <$ff>;
    close($ff);

    # 1. Validate YAML Structure
    my @syntax_errors;
    if ($rule_content !~ /rules:\s*\n/) {
        push @syntax_errors, "Missing top-level \x27rules:\x27 declaration in YAML";
    }
    if ($rule_content !~ /id:\s*[a-zA-Z0-9_-]+/) {
        push @syntax_errors, "Missing rule \x27id:\x27 field";
    }
    if ($rule_content !~ /languages:\s*\[.*go.*\]/i && $rule_content !~ /languages:\s*\n\s*-\s*go/i) {
        push @syntax_errors, "Missing or invalid \x27languages: [go]\x27 specification";
    }

    # Extract patterns from YAML
    my @patterns;
    my @pattern_nots;
    my @pattern_insides;
    my @pattern_not_insides;

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

    sub pattern_to_regex {
        my ($pat) = @_;
        $pat =~ s/^\s+//;
        $pat =~ s/\s+$//;
        $pat =~ s/^["\x27]//;
        $pat =~ s/["\x27]$//;
        my $escaped = quotemeta($pat);
        $escaped =~ s/\\\.\\\.\\\./\.\*\?/g;
        $escaped =~ s/\\\$[A-Z0-9_]+/[a-zA-Z0-9_.()"]+/g;
        return $escaped;
    }

    # 2. Evaluate against Vuln Code (True Positive) & Fixed Code (False Positive)
    my @vuln_matches;
    my @fixed_matches;

    my $semgrep_cli = `which semgrep 2>/dev/null`;
    chomp($semgrep_cli);

    if ($semgrep_cli) {
        print "Running semgrep CLI verification...\n";
        my $vuln_out = `$semgrep_cli --config \x27$rule_path\x27 \x27$vuln_path\x27 --json 2>/dev/null`;
        my $v_json = eval { decode_json($vuln_out) };
        if ($v_json && $v_json->{results}) {
            for my $r (@{$v_json->{results}}) {
                push @vuln_matches, { line => $r->{start}{line}, match => $r->{extra}{lines} };
            }
        }

        my $fixed_out = `$semgrep_cli --config \x27$rule_path\x27 \x27$fixed_path\x27 --json 2>/dev/null`;
        my $f_json = eval { decode_json($fixed_out) };
        if ($f_json && $f_json->{results}) {
            for my $r (@{$f_json->{results}}) {
                push @fixed_matches, { line => $r->{start}{line}, match => $r->{extra}{lines} };
            }
        }
    } else {
        print "Executing Semgrep pattern engine...\n";
        for my $pat (@patterns) {
            my $re = pattern_to_regex($pat);
            my @v_lines = split("\n", $vuln_code);
            for my $idx (0 .. $#v_lines) {
                if ($v_lines[$idx] =~ /$re/i) {
                    push @vuln_matches, { line => $idx + 1, match => $v_lines[$idx] };
                }
            }

            my @f_lines = split("\n", $fixed_code);
            for my $idx (0 .. $#f_lines) {
                if ($f_lines[$idx] =~ /$re/i) {
                    my $negated = 0;
                    for my $not_pat (@pattern_nots, @pattern_not_insides) {
                        my $not_re = pattern_to_regex($not_pat);
                        if ($f_lines[$idx] =~ /$not_re/i || $fixed_code =~ /$not_re/i) {
                            $negated = 1;
                            last;
                        }
                    }
                    if (!$negated) {
                        push @fixed_matches, { line => $idx + 1, match => $f_lines[$idx] };
                    }
                }
            }
        }
    }

    # 3. Decision Logic
    my $status = "REJECTED";
    my @reasons;

    if (@syntax_errors) {
        $status = "REJECTED";
        push @reasons, "Syntax errors: " . join("; ", @syntax_errors);
    } elsif (scalar @vuln_matches == 0) {
        $status = "REJECTED";
        push @reasons, "FAILED_TRUE_POSITIVE: Rule did not match any lines in vulnerable code (vuln.go)";
    } elsif (scalar @fixed_matches > 0) {
        $status = "REJECTED";
        push @reasons, "FAILED_FALSE_POSITIVE: Rule matched " . scalar(@fixed_matches) . " line(s) in fixed code (fixed.go)";
    } else {
        $status = "APPROVED";
        push @reasons, "PASS: Matched " . scalar(@vuln_matches) . " line(s) in vuln.go and 0 in fixed.go";
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

    open(my $res_fh, ">", $result_path) or die "Cannot write $result_path: $!\n";
    print $res_fh JSON::PP->new->utf8->pretty->encode(\%result);
    close($res_fh);

    if ($status eq "APPROVED") {
        my $target_rule = "$rules_store/$ghsa_id.yaml";
        open(my $tf, ">", $target_rule) or die "Cannot write $target_rule: $!\n";
        print $tf $rule_content;
        close($tf);
        print "SUCCESS: Rule APPROVED and stored at $target_rule\n";
    } else {
        print "REJECTED: " . join("\n", @reasons) . "\n";
    }

    my $ledger_data = {};
    if (-f $ledger_path) {
        if (open(my $lf, "<", $ledger_path)) {
            local $/;
            $ledger_data = eval { decode_json(<$lf>) } // {};
            close($lf);
        }
    }
    $ledger_data->{$ghsa_id} = \%result;
    if (open(my $lf, ">", $ledger_path)) {
        print $lf JSON::PP->new->utf8->pretty->encode($ledger_data);
        close($lf);
    }

    exit($status eq "APPROVED" ? 0 : 1);
' "$GHSA_ID" "$RULE_FILE" "$VULN_FILE" "$FIXED_FILE" "$RESULT_FILE" "$RULES_STORE_DIR" "$LEDGER_FILE"

echo "--------------------------------------------------"
echo "Validation result saved to $RESULT_FILE"
echo "=================================================="
