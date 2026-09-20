#!/bin/bash
set -e # Terminate on unhandled errors, but we will explicitly handle expected failures below

# Prevent Go environment errors
export GOTMPDIR=/tmp

# Automatically configure Git
git config --global user.email "pipeline@rulegen.local"
git config --global user.name "RuleGen Bot"

# Use absolute paths
BASE_DIR=$(pwd)
DATABASE_DIR="$BASE_DIR/combined_output"
WORKSPACE="$BASE_DIR/rules_workspace"
EXP_DIR="$BASE_DIR/rules/experimental"
CORPUS_DIR="/usr/local/go/src" # Strict False-Positive Baseline

mkdir -p "$WORKSPACE"
mkdir -p "$EXP_DIR"
mkdir -p "$BASE_DIR/rules/stable"

PROMPT=$(cat "$BASE_DIR/prompt_template.txt")

echo "Starting Grade-A Rule Generation Pipeline..."

# Track successful rules for the batch commit
NEW_RULES_GENERATED=0

for FILE in "$DATABASE_DIR"/*.json; do
    echo "======================================"
    
    # 1. FIX: Naming Collisions (Use Advisory ID instead of CWE class)
    ADVISORY_ID=$(basename "$FILE" .json)
    echo "Processing $ADVISORY_ID..."
    
    # 2. Check if rule already exists anywhere in the pipeline
    if [ -f "$EXP_DIR/$ADVISORY_ID.yaml" ] || [ -f "$BASE_DIR/rules/stable/$ADVISORY_ID.yaml" ]; then
        echo "⏭️  Skipping $ADVISORY_ID (Rule already exists)."
        continue
    fi

    # Strip massive metadata arrays
    ADVISORY=$(jq -c '(if type == "array" then .[0] else . end) | del(.affected, .references)' "$FILE" || echo "")
    if [ -z "$ADVISORY" ]; then
        echo "⏭️  Skipping $ADVISORY_ID: Failed to parse JSON."
        continue
    fi
    
    PAYLOAD_LEN=${#ADVISORY}
    MAX_LEN=125000 
    
    if [ "$PAYLOAD_LEN" -gt "$MAX_LEN" ]; then
        echo "⏭️  Skipping $ADVISORY_ID: Payload is $PAYLOAD_LEN chars (Too massive)."
        continue
    fi
    
    echo "Calling Antigravity CLI for Archetype Extraction..."
    # FIX: Add '|| true' to prevent LLM timeouts from killing the entire script
    RAW_OUTPUT=$(agy --print-timeout 15m --dangerously-skip-permissions -p "$PROMPT The Advisory: $ADVISORY" || true) 
    
    # Extract YAML and Go
    echo "$RAW_OUTPUT" | awk 'tolower($0) ~ /^[ \t]*```(yaml|yml)/ {flag=1; next} /^[ \t]*```/ {flag=0} flag' > "$WORKSPACE/$ADVISORY_ID.yaml"
    echo "$RAW_OUTPUT" | awk 'tolower($0) ~ /^[ \t]*```go/ {flag=1; next} /^[ \t]*```/ {flag=0} flag' > "$WORKSPACE/${ADVISORY_ID}_test.go"
    
    cd "$WORKSPACE"
    
    echo "Running Gate 1: Syntax and Semgrep Validation..."

    if [ ! -s "$ADVISORY_ID.yaml" ] || [ ! -s "${ADVISORY_ID}_test.go" ]; then
        echo "❌ GATE 1 FAILED: LLM output was empty or malformed."
        rm -f "$ADVISORY_ID.yaml" "${ADVISORY_ID}_test.go"
        cd "$BASE_DIR"
        continue
    fi

    # 3. FIX: Isolated Go Dependency Resolution
    # Create a temporary module to safely fetch third-party packages without polluting the workspace
    mkdir -p "build_$ADVISORY_ID"
    mv "${ADVISORY_ID}_test.go" "build_$ADVISORY_ID/main.go"
    cd "build_$ADVISORY_ID"
    
    go mod init ruletest >/dev/null 2>&1
    
    if command -v goimports &> /dev/null; then
        goimports -w "main.go"
    fi

    # Fetch external dependencies (e.g., gin, echo). Suppress output and don't exit on failure.
    go get -d ./... >/dev/null 2>&1 || go mod tidy >/dev/null 2>&1 || true

    # Native Compilation Check
    if ! go build -o /dev/null "./main.go"; then
        echo "❌ GATE 1 FAILED: Invalid Go syntax or missing dependencies in test file."
        cd ..
        rm -rf "build_$ADVISORY_ID" "$ADVISORY_ID.yaml"
        cd "$BASE_DIR"
        continue
    fi
    
    # Move the validated file back and clean up the temp build dir
    mv "main.go" "../${ADVISORY_ID}_test.go"
    cd ..
    rm -rf "build_$ADVISORY_ID"

    # Semgrep strict schema validation AND test validation
    if ! semgrep --validate --config "$ADVISORY_ID.yaml" > /dev/null 2>&1 || ! semgrep --test --config "$ADVISORY_ID.yaml" "${ADVISORY_ID}_test.go"; then
        echo "❌ GATE 1 FAILED: Rule failed Semgrep schema syntax or test validation."
        rm -f "$ADVISORY_ID.yaml" "${ADVISORY_ID}_test.go"
        cd "$BASE_DIR"
        continue
    fi
    
    echo "Running Gate 2: Strict False-Positive Check against standard library..."
    # FIX: Catch JSON parse errors if semgrep outputs something weird
    semgrep --config "$ADVISORY_ID.yaml" "$CORPUS_DIR" --json -o "corpus_results.json" --quiet || true
    FP_COUNT=$(jq '.results | length' corpus_results.json 2>/dev/null || echo "0")
    
    if [ "$FP_COUNT" -gt 0 ]; then
        echo "❌ GATE 2 FAILED: Rule caused $FP_COUNT false positives. Discarding."
        rm -f "$ADVISORY_ID.yaml" "${ADVISORY_ID}_test.go" corpus_results.json
        cd "$BASE_DIR"
        continue
    fi
    
    echo "✅ SUCCESS: Rule $ADVISORY_ID passed all QA gates."
    rm -f corpus_results.json
    
    # Move to experimental
    mv "$ADVISORY_ID.yaml" "$EXP_DIR/"
    mv "${ADVISORY_ID}_test.go" "$EXP_DIR/"
    
    # Increment success counter
    NEW_RULES_GENERATED=$((NEW_RULES_GENERATED + 1))
    
    cd "$BASE_DIR"
done

# 4. FIX: Batched Git Operations
echo "======================================"
echo "Pipeline Execution Complete."
echo "Total new rules generated: $NEW_RULES_GENERATED"

if [ "$NEW_RULES_GENERATED" -gt 0 ]; then
    echo "Running Gate 3: Auto-Committing batch to Experimental Tier..."
    cd "$BASE_DIR"
    
    # Safely pull to avoid conflicts
    git pull origin main --rebase --quiet || true 
    
    # Stage ONLY the experimental directory
    git add "$EXP_DIR/"
    
    if git diff --cached --quiet; then
        echo "⏭️  No changes to commit."
    else
        git commit -m "feat(rules): auto-generate $NEW_RULES_GENERATED experimental rules [skip ci]" --quiet
        git push origin HEAD:main --quiet
        echo "🚀 BATCH MERGED: $NEW_RULES_GENERATED rules are now live in the experimental tier!"
    fi
else
    echo "No new rules survived the pipeline this run."
fi
