#!/usr/bin/env bash
#
# Step 2: Generate concise Semgrep rule authoring prompt for Antigravity subagent
#
set -euo pipefail

if [[ $# -lt 1 || "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    echo "Usage: $0 <GHSA_ID>"
    echo "Example: $0 GHSA-h395-qcrw-5vmq"
    exit 1
fi

GHSA_ID="$1"
WORKSPACE_DIR="/src/workspaces/${GHSA_ID}"
METADATA_FILE="${WORKSPACE_DIR}/metadata.json"
PROMPT_FILE="${WORKSPACE_DIR}/prompt.txt"

if [ ! -f "$METADATA_FILE" ]; then
    echo "Error: Metadata file '$METADATA_FILE' not found. Run 01_fetch_commit.sh first." >&2
    exit 1
fi

echo "=================================================="
echo "Step 2: Generate Antigravity Author Prompt"
echo "=================================================="
echo "Advisory ID   : $GHSA_ID"
echo "Prompt File   : $PROMPT_FILE"
echo "--------------------------------------------------"

perl -MJSON::PP -e '
    use strict;
    use warnings;

    my ($ghsa_id, $workspace, $meta_path, $prompt_path) = @ARGV;

    open(my $mf, "<", $meta_path) or die "Cannot open $meta_path: $!\n";
    local $/;
    my $meta = decode_json(<$mf>);
    close($mf);

    my $pkg = $meta->{package} // "Golang Package";
    my $summary = $meta->{summary} // "Security Vulnerability";
    my $aliases = join(", ", @{$meta->{aliases} // []});
    my $file_modified = $meta->{file_modified} // "code.go";
    my $diff_hunk = $meta->{diff_hunk} // "";

    # Limit diff hunk length if huge
    if (length($diff_hunk) > 2000) {
        $diff_hunk = substr($diff_hunk, 0, 2000) . "\n... [diff truncated]";
    }

    my $prompt = <<"END_PROMPT";
You are an expert security engineer and Semgrep rule author specializing in Golang security.

### TASK:
Synthesize a precise, high-accuracy Semgrep YAML rule for the security vulnerability in package \x27$pkg\x27 ($ghsa_id).

### VULNERABILITY CONTEXT:
- Advisory ID: $ghsa_id
- Aliases: $aliases
- Package: $pkg
- Summary: $summary
- File Modified: $file_modified

### SECURITY FIX DIFF:
\`\`\`diff
$diff_hunk
\`\`\`

### OBJECTIVE & REQUIREMENTS:
1. Identify the unsafe pattern in the pre-patch code (lines marked with \x27-\x27 in diff).
2. Synthesize a Semgrep rule targeting Go (\`languages: [go]\`).
3. The rule must match the vulnerable code in \x27$workspace/vuln.go\x27.
4. The rule MUST NOT match the patched code in \x27$workspace/fixed.go\x27.
5. Write the final standalone rule directly to: \x27$workspace/rule.yaml\x27.
6. Once written, run: \`bash /src/03_validate_rule.sh $ghsa_id\` to verify the rule.

### REQUIRED YAML FORMAT:
\`\`\`yaml
rules:
  - id: $ghsa_id
    languages: [go]
    severity: WARNING
    message: "$summary"
    metadata:
      cve: "$aliases"
      ghsa: "$ghsa_id"
      confidence: HIGH
    patterns:
      - pattern: <UNSAFE_PATTERN>
\`\`\`
END_PROMPT

    open(my $pf, ">", $prompt_path) or die "Cannot write $prompt_path: $!\n";
    print $pf $prompt;
    close($pf);

    print "Generated author prompt at $prompt_path (" . (-s $prompt_path) . " bytes).\n";
' "$GHSA_ID" "$WORKSPACE_DIR" "$METADATA_FILE" "$PROMPT_FILE"

echo "--------------------------------------------------"
echo "Step 2 complete!"
echo "=================================================="
