#!/usr/bin/env bash
#
# Step 1: Fetch commit patch using local JSON metadata and reconstruct vuln.go & fixed.go
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ $# -lt 1 || "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    echo "Usage: $0 <GHSA_ID> [DATASET_JSON]"
    echo "Example: $0 GHSA-h395-qcrw-5vmq"
    exit 1
fi

GHSA_ID="$1"
DEFAULT_DATASET="${SCRIPT_DIR}/ghsa_golang_git_diffs_med_high_crit.json"
if [ ! -f "$DEFAULT_DATASET" ]; then
    DEFAULT_DATASET="${SCRIPT_DIR}/ghsa_golang_git_diffs.json"
fi
DATASET_JSON="${2:-$DEFAULT_DATASET}"
WORKSPACE_DIR="${SCRIPT_DIR}/workspaces/${GHSA_ID}"

if [ ! -f "$DATASET_JSON" ]; then
    echo "Error: Dataset JSON '$DATASET_JSON' not found." >&2
    exit 1
fi

mkdir -p "$WORKSPACE_DIR"

echo "=================================================="
echo "Step 1: Fetch Commit & Build Testbed"
echo "=================================================="
echo "Advisory ID   : $GHSA_ID"
echo "Dataset       : $DATASET_JSON"
echo "Workspace     : $WORKSPACE_DIR"
echo "--------------------------------------------------"

perl -MJSON::PP -e '
    use strict;
    use warnings;

    my ($ghsa_id, $dataset_path, $workspace) = @ARGV;

    my $entry;
    {
        open(my $fh, "<", $dataset_path) or die "Cannot read $dataset_path: $!\n";
        local $/;
        my $data = decode_json(<$fh>);
        close($fh);
        $entry = $data->{$ghsa_id};
    }
    die "Error: Advisory $ghsa_id not found in $dataset_path\n" unless $entry;

    my $patch_url = "";
    my $commit_sha = "";
    if (ref($entry->{fix_commits}) eq "ARRAY" && @{$entry->{fix_commits}}) {
        $patch_url = $entry->{fix_commits}[0]{patch_url} // $entry->{fix_commits}[0]{diff_url};
        $commit_sha = $entry->{fix_commits}[0]{commit_sha} // "";
    }

    if (!$patch_url && ref($entry->{pull_requests}) eq "ARRAY" && @{$entry->{pull_requests}}) {
        $patch_url = $entry->{pull_requests}[0]{patch_url} // $entry->{pull_requests}[0]{diff_url};
    }

    if (!$patch_url && ref($entry->{version_diffs}) eq "ARRAY" && @{$entry->{version_diffs}}) {
        $patch_url = $entry->{version_diffs}[0]{diff_url};
    }

    die "Error: No patch or diff URL available for $ghsa_id\n" unless $patch_url;

    print "Fetching patch from: $patch_url\n";
    my $patch_file = "$workspace/commit.patch";
    my $curl_cmd = "curl -sL --max-time 15 -o \x27$patch_file\x27 \x27$patch_url\x27";
    my $ret = system($curl_cmd);
    die "Failed to download patch: $!\n" if $ret != 0 || ! -s $patch_file;

    print "Patch downloaded successfully (" . (-s $patch_file) . " bytes).\n";

    open(my $pf, "<", $patch_file) or die "Cannot open $patch_file: $!\n";
    my @lines = <$pf>;
    close($pf);

    my @vuln;
    my @fixed;
    my $in_go = 0;
    my $go_file_name = "";
    my @hunk_diff;

    for my $l (@lines) {
        if ($l =~ /^diff --git/i) {
            $in_go = ($l =~ /\.go/i) ? 1 : 0;
            if ($l =~ /b\/(.+\.go)/i && !$go_file_name) {
                $go_file_name = $1;
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

    open(my $vf, ">", "$workspace/vuln.go") or die $!;
    print $vf "package main\n\n" . join("", @vuln);
    close($vf);

    open(my $ff, ">", "$workspace/fixed.go") or die $!;
    print $ff "package main\n\n" . join("", @fixed);
    close($ff);

    my %meta = (
        advisory_id   => $ghsa_id,
        package       => $entry->{package} // "",
        summary       => $entry->{summary} // "",
        severity      => $entry->{severity} // "MODERATE",
        cvss_score    => $entry->{cvss_score} // "",
        aliases       => $entry->{aliases} // [],
        commit_sha    => $commit_sha,
        patch_url     => $patch_url,
        file_modified => $go_file_name || "target.go",
        diff_hunk     => join("", @hunk_diff),
    );

    open(my $mf, ">", "$workspace/metadata.json") or die $!;
    print $mf JSON::PP->new->utf8->pretty->encode(\%meta);
    close($mf);

    print "Reconstructed testbed:\n";
    printf "  - %s/vuln.go (%d bytes)\n", $workspace, -s "$workspace/vuln.go";
    printf "  - %s/fixed.go (%d bytes)\n", $workspace, -s "$workspace/fixed.go";
    print "  - $workspace/metadata.json\n";
' "$GHSA_ID" "$DATASET_JSON" "$WORKSPACE_DIR"

echo "--------------------------------------------------"
echo "Step 1 complete!"
echo "=================================================="
