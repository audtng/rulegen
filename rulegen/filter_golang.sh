#!/usr/bin/env bash
#
# Script to parse the combined GHSA JSON dataset and extract all
# Golang-related advisories into a separate consolidated JSON array file.
#
set -euo pipefail

# Help message
if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    echo "Usage: $0 [INPUT_COMBINED_JSON] [OUTPUT_GOLANG_JSON]"
    echo ""
    echo "Parses INPUT_COMBINED_JSON, filters all Golang-related advisories, and saves them to OUTPUT_GOLANG_JSON."
    echo ""
    echo "Arguments:"
    echo "  INPUT_COMBINED_JSON  Path to combined JSON master file (default: /src/ghsa_combined.json)"
    echo "  OUTPUT_GOLANG_JSON   Path for the filtered Golang JSON file (default: /src/ghsa_golang.json)"
    echo ""
    exit 0
fi

INPUT_FILE="${1:-/src/ghsa_combined.json}"
OUTPUT_FILE="${2:-/src/ghsa_golang.json}"

# Resolve relative path if default does not exist in absolute path
if [ ! -f "$INPUT_FILE" ]; then
    if [ -f "./ghsa_combined.json" ]; then
        INPUT_FILE="./ghsa_combined.json"
    else
        echo "Error: Input file '$INPUT_FILE' not found." >&2
        exit 1
    fi
fi

# Ensure destination directory exists
mkdir -p "$(dirname "$OUTPUT_FILE")"

echo "=================================================="
echo "GHSA Golang Advisory Filter & Parser"
echo "=================================================="
echo "Input Master File : $INPUT_FILE"
echo "Output Go File    : $OUTPUT_FILE"
echo "--------------------------------------------------"

# Execute streaming filter in Perl for high throughput, robust JSON validation, and O(1) memory
perl -MTime::HiRes=time -MJSON::PP -e '
    use strict;
    use warnings;

    my ($input_path, $output_path) = @ARGV;

    open(my $in, "<", $input_path) or die "Cannot open input file $input_path: $!\n";
    open(my $out, ">", $output_path) or die "Cannot create output file $output_path: $!\n";

    # Enable output buffering for high write performance
    select((select($out), $| = 0)[0]);

    print "Scanning and filtering records from $input_path...\n";
    my $t0 = time();

    print $out "[\n";

    my $total_scanned = 0;
    my $go_matched = 0;
    my $record_buffer = "";
    my $in_record = 0;
    my $report_interval = 50000;
    my $is_first_go_record = 1;

    while (my $line = <$in>) {
        if (!$in_record) {
            if ($line =~ /^\{/) {
                $in_record = 1;
                $record_buffer = $line;
            }
        } else {
            $record_buffer .= $line;
            # Record terminator: "}" or "}," at start of line
            if ($line =~ /^\}(,)?\s*$/) {
                $total_scanned++;

                # Clean trailing comma for valid JSON parsing of the single object
                my $json_text = $record_buffer;
                $json_text =~ s/,\s*$//;

                # Quick pre-filter check before full decode for speed
                my $candidate = 0;
                if ($json_text =~ /"ecosystem":\s*"Go"/i ||
                    $json_text =~ m{golang\.org|pkg\.go\.dev|github\.com/golang/}i ||
                    $json_text =~ /"name":\s*"(stdlib|toolchain|cmd\/go)"/i) {
                    $candidate = 1;
                }

                if ($candidate) {
                    my $data = eval { decode_json($json_text) };
                    if ($data) {
                        my $is_go = 0;

                        # 1. Primary check: affected ecosystem is Go
                        if (ref($data->{affected}) eq "ARRAY") {
                            for my $aff (@{$data->{affected}}) {
                                my $pkg = $aff->{package} // {};
                                my $eco = $pkg->{ecosystem} // "";
                                if ($eco =~ /^go(lang)?$/i) {
                                    $is_go = 1;
                                    last;
                                }
                                # Go stdlib packages or cmd/go
                                my $pkg_name = $pkg->{name} // "";
                                if ($pkg_name =~ /^(stdlib|toolchain|cmd\/go)$/i) {
                                    $is_go = 1;
                                    last;
                                }
                            }
                        }

                        # 2. Secondary check: unreviewed Go advisories referencing Go repositories
                        if (!$is_go && ref($data->{references}) eq "ARRAY") {
                            for my $ref (@{$data->{references}}) {
                                my $url = $ref->{url} // "";
                                if ($url =~ m{https?://(golang\.org|pkg\.go\.dev|vuln\.go\.dev|github\.com/golang/go)}i) {
                                    $is_go = 1;
                                    last;
                                }
                            }
                        }

                        if ($is_go) {
                            $go_matched++;
                            if (!$is_first_go_record) {
                                print $out ",\n";
                            } else {
                                $is_first_go_record = 0;
                            }
                            # Print trimmed record content
                            $json_text =~ s/^\s+//;
                            $json_text =~ s/\s+$//;
                            print $out $json_text;
                        }
                    } else {
                        warn "Warning: JSON decode error in record $total_scanned: $@\n";
                    }
                }

                $in_record = 0;
                $record_buffer = "";

                if ($total_scanned % $report_interval == 0) {
                    my $elapsed = time() - $t0;
                    printf "  Progress: %d records scanned, %d Go advisories matched (%.2fs)\n",
                        $total_scanned, $go_matched, $elapsed;
                }
            }
        }
    }

    print $out "\n]\n";
    close($in);
    close($out);

    my $total_elapsed = time() - $t0;
    my $size_bytes = -s $output_path;
    my $size_mb = $size_bytes / (1024 * 1024);

    print "--------------------------------------------------\n";
    printf "Scanned %d total records in %.2fs\n", $total_scanned, $total_elapsed;
    printf "Extracted %d Golang advisories into %s\n", $go_matched, $output_path;
    printf "Output file size: %.2f MB (%d bytes)\n", $size_mb, $size_bytes;
' "$INPUT_FILE" "$OUTPUT_FILE"

echo "--------------------------------------------------"
echo "Filtering complete!"
echo "=================================================="
