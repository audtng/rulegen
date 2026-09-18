#!/bin/bash
set -e
export GOTMPDIR=/tmp

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
    
    BRANCH_NAME="rule-staging/$CWE_NAME"
    
    # 1. Check if it is pending review (branch exists)
    if git show-ref --verify --quiet "refs/heads/$BRANCH_NAME"; then
        echo "⏭️  Skipping $CWE_NAME (Staging branch already exists)."
        continue
    fi
    
    # 2. Check if it was already merged (file exists on main branch)
    # Because the script always returns to 'main', we can just check the filesystem.
    if [ -f "$WORKSPACE/$CWE_NAME.yaml" ]; then
        echo "⏭️  Skipping $CWE_NAME (Rule already merged into main)."
        continue
    fi

    echo "Calling Antigravity CLI for Archetype Extraction..."
    RAW_OUTPUT=$(agy --print-timeout 15m --dangerously-skip-permissions -p "$PROMPT The Advisory: $ADVISORY")
    
    echo "$RAW_OUTPUT" | awk '/```yaml/{flag=1; next} /```/{flag=0} flag' > "$WORKSPACE/$CWE_NAME.yaml"
    echo "$RAW_OUTPUT" | awk '/```go/{flag=1; next} /```/{flag=0} flag' > "$WORKSPACE/${CWE_NAME}_test.go"
    
    cd "$WORKSPACE"
    
    echo "Running Gate 1: Syntax and Edge Cases..."
    
    if [ ! -f "go.mod" ]; then
        go mod init ruletest >/dev/null 2>&1
    fi

    # Validate Go syntax natively to catch LLM stuttering/typos
    cp "${CWE_NAME}_test.go" ".temp_validate.go"
    if ! go build "temp_validate.go"; then
        echo "❌ GATE 1 FAILED: Invalid Go syntax in test file."
        rm -f ".temp_validate.go"
        cd - > /dev/null
        continue
    fi
    rm -f ".temp_validate.go"

    if ! semgrep --validate --config "$CWE_NAME.yaml" || ! semgrep --test --config "$CWE_NAME.yaml" "${CWE_NAME}_test.go"; then
        echo "❌ GATE 1 FAILED: Discarding $CWE_NAME."
#        rm -f "$CWE_NAME.yaml" "${CWE_NAME}_test.go"
        cd - > /dev/null
        continue
    fi
    
    echo "Running Gate 2: Scanning Go Standard Library ($CORPUS_DIR)..."
    semgrep --config "$CWE_NAME.yaml" "$CORPUS_DIR" --json -o "corpus_results.json" --quiet || true
    
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
