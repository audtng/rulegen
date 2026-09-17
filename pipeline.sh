#!/bin/bash
set -e

DATABASE_DIR="./combined_output"
WORKSPACE="./rules"
CORPUS_DIR="/usr/local/go/src" 
mkdir -p "$WORKSPACE"

PROMPT=$(cat prompt_template.txt)

echo "Starting Generalized Rule Generation Pipeline..."

for FILE in "$DATABASE_DIR"/*.json; do
    echo "======================================"
    echo "Processing $FILE..."
    
    ADVISORY=$(jq 'if type == "array" then .[0] else . end' "$FILE")
    CWE_NAME=$(basename "$FILE" .json)
    
    echo "Calling Antigravity CLI for Archetype Extraction..."
    RAW_OUTPUT=$(antigravity prompt "$PROMPT The Advisory: $ADVISORY")
    
    echo "$RAW_OUTPUT" | awk '/```yaml/{flag=1; next} /```/{flag=0} flag' > "$WORKSPACE/$CWE_NAME.yaml"
    echo "$RAW_OUTPUT" | awk '/```go/{flag=1; next} /```/{flag=0} flag' > "$WORKSPACE/${CWE_NAME}_test.go"
    
    cd "$WORKSPACE"
    
    echo "Running Gate 1: Syntax and Edge Cases..."
    if ! semgrep --validate --config "$CWE_NAME.yaml" || ! semgrep --test --config "$CWE_NAME.yaml" "${CWE_NAME}_test.go"; then
        echo "❌ GATE 1 FAILED: Discarding $CWE_NAME."
        rm -f "$CWE_NAME.yaml" "${CWE_NAME}_test.go"
        cd - > /dev/null
        continue
    fi
    
    echo "Running Gate 2: Scanning Go Standard Library ($CORPUS_DIR)..."
    semgrep --config "$CWE_NAME.yaml" "$CORPUS_DIR" --json -o "corpus_results.json" --quiet
    
    FP_COUNT=$(jq '.results | length' corpus_results.json)
    
    if [ "$FP_COUNT" -gt 0 ]; then
        echo "❌ GATE 2 FAILED: Rule caused $FP_COUNT false positives."
        rm -f "$CWE_NAME.yaml" "${CWE_NAME}_test.go" corpus_results.json
        cd - > /dev/null
        continue
    fi
    
    echo "✅ SUCCESS: Rule $CWE_NAME passed all OSS generalization tests."
    rm -f corpus_results.json
    
    BRANCH_NAME="rule-staging/$CWE_NAME"
    git checkout main
    
    git checkout -b "$BRANCH_NAME"
    git add "$CWE_NAME.yaml" "${CWE_NAME}_test.go"
    git commit -m "feat(rules): Add highly-generalized rule for $CWE_NAME"
    
    git checkout main
    
    cd - > /dev/null
done
