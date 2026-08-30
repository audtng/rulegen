#!/usr/bin/env bash
#
# Master Pipeline Driver: Automates Commit Fetching, Prompt Generation, and Verification
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    echo "Usage: $0 [GHSA_ID] [MAX_ITEMS]"
    echo ""
    echo "Runs the automated pipeline loop across Medium, High, and Critical Go advisories."
    echo "If GHSA_ID is specified, runs for that specific advisory."
    echo "Otherwise, iterates over the first MAX_ITEMS (default: 5) advisories."
    echo ""
    exit 0
fi

DEFAULT_DATASET="${SCRIPT_DIR}/ghsa_golang_git_diffs_med_high_crit.json"
if [ ! -f "$DEFAULT_DATASET" ]; then
    DEFAULT_DATASET="${SCRIPT_DIR}/ghsa_golang_git_diffs.json"
fi
DATASET_JSON="$DEFAULT_DATASET"
TARGET_ID="${1:-}"
MAX_ITEMS="${2:-5}"

if [ ! -f "$DATASET_JSON" ]; then
    echo "Error: Dataset '$DATASET_JSON' not found." >&2
    exit 1
fi

echo "=================================================="
echo "Agentic Semgrep Rule Generation Pipeline Driver"
echo "=================================================="
echo "Dataset : $DATASET_JSON"
echo "=================================================="

process_advisory() {
    local adv_id="$1"
    echo ""
    echo ">>> Processing Advisory: $adv_id <<<"
    
    if ! bash "${SCRIPT_DIR}/01_fetch_commit.sh" "$adv_id" "$DATASET_JSON"; then
        echo "Warning: Step 1 failed for $adv_id, skipping..." >&2
        return 1
    fi

    if ! bash "${SCRIPT_DIR}/02_generate_prompt.sh" "$adv_id"; then
        echo "Warning: Step 2 failed for $adv_id, skipping..." >&2
        return 1
    fi

    echo "Workspace ready at: ${SCRIPT_DIR}/workspaces/$adv_id"
    echo "Prompt generated at: ${SCRIPT_DIR}/workspaces/$adv_id/prompt.txt"
    echo "Ready for Antigravity Subagent ingestion."
}

if [ -n "$TARGET_ID" ]; then
    process_advisory "$TARGET_ID"
else
    ADVISORY_LIST=$(perl -MJSON::PP -e '
        my ($dataset, $max) = @ARGV;
        open(my $fh, "<", $dataset) or die $!;
        local $/;
        my $data = decode_json(<$fh>);
        close($fh);
        my $count = 0;
        for my $id (sort keys %$data) {
            my $entry = $data->{$id};
            if ($entry->{validation}{has_commit_diff} && ref($entry->{fix_commits}) eq "ARRAY" && @{$entry->{fix_commits}}) {
                my $patch_url = $entry->{fix_commits}[0]{patch_url} // "";
                if ($patch_url =~ /github\.com/) {
                    print "$id\n";
                    $count++;
                    last if $count >= $max;
                }
            }
        }
    ' "$DATASET_JSON" "$MAX_ITEMS")

    for adv_id in $ADVISORY_LIST; do
        process_advisory "$adv_id" || true
    done
fi

echo ""
echo "=================================================="
echo "Pipeline driver preparation complete."
echo "=================================================="
