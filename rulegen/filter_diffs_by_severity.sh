#!/usr/bin/env bash
#
# Script to parse ghsa_golang_git_diffs.json and filter for Medium (Moderate),
# High, and Critical severity advisories, enriching each entry with severity metadata.
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Help message
if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    echo "Usage: $0 [OUTPUT_JSON] [ALLOWED_SEVERITIES]"
    echo ""
    echo "Filters ghsa_golang_git_diffs.json for specified severities (default: CRITICAL,HIGH,MODERATE,MEDIUM)."
    echo ""
    echo "Arguments:"
    echo "  OUTPUT_JSON         Path for filtered output JSON (default: /src/rulegen/ghsa_golang_git_diffs_med_high_crit.json)"
    echo "  ALLOWED_SEVERITIES  Comma-separated severities (default: CRITICAL,HIGH,MODERATE,MEDIUM)"
    echo ""
    exit 0
fi

DIFFS_JSON="${SCRIPT_DIR}/ghsa_golang_git_diffs.json"
ADVISORIES_JSON="${SCRIPT_DIR}/ghsa_golang.json"
OUTPUT_JSON="${1:-${SCRIPT_DIR}/ghsa_golang_git_diffs_med_high_crit.json}"
ALLOWED_SEV="${2:-CRITICAL,HIGH,MODERATE,MEDIUM}"

if [ ! -f "$DIFFS_JSON" ]; then
    echo "Error: Git diffs dataset '$DIFFS_JSON' not found." >&2
    exit 1
fi

if [ ! -f "$ADVISORIES_JSON" ]; then
    echo "Error: Advisory dataset '$ADVISORIES_JSON' not found." >&2
    exit 1
fi

echo "=================================================="
echo "GHSA Severity-Based Git Diff Filter"
echo "=================================================="
echo "Input Diffs  : $DIFFS_JSON"
echo "Advisories   : $ADVISORIES_JSON"
echo "Output File  : $OUTPUT_JSON"
echo "Severities   : $ALLOWED_SEV"
echo "--------------------------------------------------"

perl -MTime::HiRes=time -MJSON::PP -e '
    use strict;
    use warnings;

    my ($diffs_path, $adv_path, $out_path, $allowed_str) = @ARGV;

    my %allowed;
    for my $s (split(",", $allowed_str)) {
        $allowed{uc($s)} = 1;
    }
    # Treat MEDIUM and MODERATE as equivalent
    $allowed{MODERATE} = 1 if $allowed{MEDIUM};
    $allowed{MEDIUM} = 1 if $allowed{MODERATE};

    print "Loading advisory severity metadata from $adv_path...\n";
    my $t0 = time();

    my %severity_map;
    my %cvss_map;
    {
        open(my $afh, "<", $adv_path) or die "Cannot open $adv_path: $!\n";
        local $/;
        my $adv_list = decode_json(<$afh>);
        close($afh);

        for my $adv (@$adv_list) {
            my $id = $adv->{id};
            next unless $id;

            my $sev = $adv->{database_specific}{severity} // "UNSPECIFIED";
            $severity_map{$id} = uc($sev);

            if (ref($adv->{severity}) eq "ARRAY" && @{$adv->{severity}}) {
                $cvss_map{$id} = $adv->{severity}[0]{score} // $adv->{severity}[0]{type};
            }
        }
    }
    printf "Indexed metadata for %d advisories in %.2fs.\n", scalar keys %severity_map, time() - $t0;

    print "Filtering git diff dataset from $diffs_path...\n";
    open(my $dfh, "<", $diffs_path) or die "Cannot open $diffs_path: $!\n";
    local $/;
    my $diffs_data = decode_json(<$dfh>);
    close($dfh);

    my %filtered_diffs;
    my %stats = (
        CRITICAL => 0,
        HIGH     => 0,
        MODERATE => 0,
        LOW      => 0,
        OTHER    => 0,
        total_matched => 0,
        with_diff_sources => 0,
    );

    for my $id (sort keys %$diffs_data) {
        my $entry = $diffs_data->{$id};
        my $sev = $severity_map{$id} // "UNSPECIFIED";
        $sev = "MODERATE" if $sev eq "MEDIUM";

        if ($allowed{$sev}) {
            # Inject severity and CVSS info into entry
            $entry->{severity} = $sev;
            $entry->{cvss_score} = $cvss_map{$id} if $cvss_map{$id};

            $filtered_diffs{$id} = $entry;
            $stats{$sev}++;
            $stats{total_matched}++;
            $stats{with_diff_sources}++ if $entry->{validation}{has_diff_source};
        } else {
            if ($sev eq "LOW") {
                $stats{LOW}++;
            } else {
                $stats{OTHER}++;
            }
        }
    }

    print "Encoding and writing output to $out_path...\n";
    my $json_encoder = JSON::PP->new->utf8->pretty->canonical;
    my $encoded = $json_encoder->encode(\%filtered_diffs);

    open(my $out, ">", $out_path) or die "Cannot write $out_path: $!\n";
    print $out $encoded;
    close($out);

    my $size_bytes = -s $out_path;
    my $size_mb = $size_bytes / (1024 * 1024);
    my $elapsed = time() - $t0;

    print "--------------------------------------------------\n";
    print "SEVERITY FILTER AUDIT REPORT:\n";
    printf "  CRITICAL Severity Advisories : %d\n", $stats{CRITICAL};
    printf "  HIGH Severity Advisories     : %d\n", $stats{HIGH};
    printf "  MODERATE/MEDIUM Advisories   : %d\n", $stats{MODERATE};
    printf "  ──────────────────────────── : ────\n";
    printf "  TOTAL Filtered Advisories    : %d\n", $stats{total_matched};
    printf "  Advisories with Diff Sources : %d (%.1f%%)\n",
        $stats{with_diff_sources}, ($stats{with_diff_sources} / $stats{total_matched}) * 100;
    printf "  Excluded Low/Other Advisories: %d\n", $stats{LOW} + $stats{OTHER};
    printf "Output file size: %.2f MB (%d bytes) in %.2fs\n", $size_mb, $size_bytes, $elapsed;
' "$DIFFS_JSON" "$ADVISORIES_JSON" "$OUTPUT_JSON" "$ALLOWED_SEV"

echo "--------------------------------------------------"
echo "Severity filtering complete!"
echo "=================================================="
