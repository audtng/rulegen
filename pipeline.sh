#!/bin/bash
set -e

# Prevent Go environment errors
export GOTMPDIR=/tmp

# Automatically configure Git if Docker restarted
git config --global user.email "pipeline@rulegen.local"
git config --global user.name "RuleGen Bot"

# Use absolute paths to permanently prevent directory state confusion
BASE_DIR=$(pwd)
DATABASE_DIR="$BASE_DIR/combined_output"
WORKSPACE="$BASE_DIR/rules"
CORPUS_DIR="/usr/local/go/src" # Strict False-Positive Baseline

mkdir -p "$WORKSPACE"

PROMPT=$(cat "$BASE_DIR/prompt_template.txt")

echo "Starting Grade-A Rule Generation Pipeline..."

for FILE in "$DATABASE_DIR"/*.json; do
    echo "======================================"
    echo "Processing $FILE..."
    
    ADVISORY=$(jq -c 'if type == "array" then .[0] else . end' "$FILE")
    CWE_NAME=$(basename "$FILE" .json)
    BRANCH_NAME="rule-staging/$CWE_NAME"
    
    # 1. Check if it is pending review
    if git show-ref --verify --quiet "refs/heads/$BRANCH_NAME"; then
        echo "⏭️  Skipping $CWE_NAME (Staging branch already exists)."
        continue
    fi
    
    # 2. Check if it was already merged
    if [ -f "$WORKSPACE/$CWE_NAME.yaml" ]; then
        echo "⏭️  Skipping $CWE_NAME (Rule already merged into main)."
        continue
    fi

    echo "Calling Antigravity CLI for Archetype Extraction (via File Reference)..."
    
    # 1. Combine the prompt and the massive JSON (diffs included) into a single file
    cat "$BASE_DIR/prompt_template.txt" > "$WORKSPACE/agent_task.txt"
    echo -e "\n\n=== THE ADVISORY (WITH DIFFS) ===\n" >> "$WORKSPACE/agent_task.txt"
    cat "$FILE" >> "$WORKSPACE/agent_task.txt"
    
    # 2. Give the agent a tiny CLI prompt telling it to read the file
    RAW_OUTPUT=$(agy --print-timeout 15m --dangerously-skip-permissions -p "I have placed your instructions and the full CVE JSON (including diffs) in the file '$WORKSPACE/agent_task.txt'. Read that file completely, then generate the YAML and Go blocks exactly as instructed.")
    
    # 3. Clean up the temporary file
    rm -f "$WORKSPACE/agent_task.txt"
    
    echo "$RAW_OUTPUT" | awk '/```yaml/{flag=1; next} /```/{flag=0} flag' > "$WORKSPACE/$CWE_NAME.yaml"
    echo "$RAW_OUTPUT" | awk '/```go/{flag=1; next} /```/{flag=0} flag' > "$WORKSPACE/${CWE_NAME}_test.go"
    
    # Move into the rules directory for isolated validation
    cd "$WORKSPACE"
    
    echo "Running Gate 1: Syntax and Semgrep Validation..."

    # 1. Ensure the LLM generated output and didn't timeout
    if [ ! -s "${CWE_NAME}_test.go" ]; then
        echo "❌ GATE 1 FAILED: LLM output was empty (likely timed out)."
        cd "$BASE_DIR"
        continue
    fi

    # 2. Ensure a Go module exists so the compiler works
    if [ ! -f "go.mod" ]; then
        go mod init ruletest >/dev/null 2>&1
    fi

    # 3. Safely copy to standard filename and rewrite 'package main' to 'package rules'
    cp "${CWE_NAME}_test.go" "temp_validate.go"
    sed -i 's/package main/package rules/g' "temp_validate.go"

    # 4. Compile natively to validate AST and imports
    if ! go build -o /dev/null "./temp_validate.go"; then
        echo "❌ GATE 1 FAILED: Invalid Go syntax in test file."
        rm -f "temp_validate.go"
        cd "$BASE_DIR"
        continue
    fi
    rm -f "temp_validate.go"

    # 5. Semgrep core validation (proves the rule catches the test file)
    if ! semgrep --validate --config "$CWE_NAME.yaml" || ! semgrep --test --config "$CWE_NAME.yaml" "${CWE_NAME}_test.go"; then
        echo "❌ GATE 1 FAILED: Rule failed Semgrep syntax/test validation."
        cd "$BASE_DIR"
        continue
    fi
    
    echo "Running Gate 2: Strict False-Positive Check against standard library..."
    semgrep --config "$CWE_NAME.yaml" "$CORPUS_DIR" --json -o "corpus_results.json" --quiet || true
    
    FP_COUNT=$(jq '.results | length' corpus_results.json)
    
    # ZERO TOLERANCE FOR FALSE POSITIVES IN THE STANDARD LIBRARY
    if [ "$FP_COUNT" -gt 0 ]; then
        echo "❌ GATE 2 FAILED: Rule caused $FP_COUNT false positives. Discarding."
        rm -f "$CWE_NAME.yaml" "${CWE_NAME}_test.go" corpus_results.json
        cd "$BASE_DIR"
        continue
    fi
    
    echo "✅ SUCCESS: Rule passed all QA gates."
    rm -f corpus_results.json
    

    echo "Executing push_branch.sh"
    bash "push_branch.sh"
    # Stage the highly-validated rule
   # git checkout main
   # git checkout -b "$BRANCH_NAME"
   # git add "$CWE_NAME.yaml" "${CWE_NAME}_test.go"
   # git commit -m "feat(rules): Add highly-generalized Grade A rule for $CWE_NAME"
   # git checkout main
    
    # Safely return to base directory
    cd "$BASE_DIR"
done
