#!/bin/bash
set -e # Terminate on unhandled errors

# Prevent Go environment errors
export GOTMPDIR=/tmp

# Automatically configure Git
git config --global user.email "pipeline@rulegen.local"
git config --global user.name "RuleGen Bot"

# Use absolute paths
BASE_DIR=$(pwd)
DATABASE_DIR="$BASE_DIR/CWE-287"
WORKSPACE="$BASE_DIR/rules_workspace"
EXP_DIR="$BASE_DIR/rules/experimental"
CORPUS_DIR="/usr/local/go/src" # Strict False-Positive Baseline

mkdir -p "$WORKSPACE"
mkdir -p "$EXP_DIR"
mkdir -p "$BASE_DIR/rules/stable"

PROMPT=$(cat "$BASE_DIR/p287.txt")

# =========================================================================
# PRE-FLIGHT: Create Semantic Hashing Helper
# =========================================================================
# This Python script guarantees that YAML formatting differences do not 
# bypass the deduplication check. It hashes ONLY the core logic.
cat << 'EOF' > "$BASE_DIR/hash_rule.py"
import yaml, json, sys, hashlib
try:
with open(sys.argv[1], 'r') as f:
    data = yaml.safe_load(f)
    rule = data['rules'][0]
    # Extract only the AST logic, ignoring metadata, messages, and IDs
    core = {
        'sources': rule.get('pattern-sources', []),
        'propagators': rule.get('pattern-propagators', []),
        'sinks': rule.get('pattern-sinks', [])
    }
    # Serialize deterministically
    dump = json.dumps(core, sort_keys=True)
    print(hashlib.sha256(dump.encode()).hexdigest())
except Exception as e:
sys.exit(1)
EOF

echo "[*] Building semantic hash database of existing rules..."
HASH_DB="$BASE_DIR/rule_hashes.txt"
> "$HASH_DB"
find "$BASE_DIR/rules" -type f -name "*.yaml" | while read -r EXISTING_RULE; do
    python3 "$BASE_DIR/hash_rule.py" "$EXISTING_RULE" >> "$HASH_DB" 2>/dev/null || true
done
TOTAL_HASHES=$(wc -l < "$HASH_DB" | tr -d ' ')
echo "[+] Loaded $TOTAL_HASHES unique rule archetypes."
echo "======================================"

echo "Starting Rule Generation Pipeline..."
NEW_RULES_GENERATED=0

for FILE in "$DATABASE_DIR"/*.json; do
    echo "======================================"
    
    ADVISORY_ID=$(basename "$FILE" .json)
    echo "Processing $ADVISORY_ID..."
    
    EXISTING_RULE=$(find "$BASE_DIR/rules" -type f \( -name "${ADVISORY_ID}.yaml" -o -name "${ADVISORY_ID}.yaml.bak" \) -print -quit 2>/dev/null || true)

    if [ -n "$EXISTING_RULE" ]; then
        echo "⏭️  Skipping $ADVISORY_ID as Rule (or backup) already exists."
        continue
    fi

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
  git clone https://audtng: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@github.com/audtng/rulegen.git
  fi
    
    echo "Calling Antigravity CLI for Archetype Extraction..."
    RAW_OUTPUT=$(agy --print-timeout 15m --dangerously-skip-permissions -p "$PROMPT The Advisory: $ADVISORY" || true) 
    
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

    mkdir -p "build_$ADVISORY_ID"
    mv "${ADVISORY_ID}_test.go" "build_$ADVISORY_ID/main.go"
    cd "build_$ADVISORY_ID"
    
    go mod init ruletest >/dev/null 2>&1
    
    if command -v goimports &> /dev/null; then
        goimports -w "main.go"
    fi

    go get -d ./... >/dev/null 2>&1 || go mod tidy >/dev/null 2>&1 || true

    if ! go build -o /dev/null "./main.go"; then
        echo "❌ GATE 1 FAILED: Invalid Go syntax or missing dependencies in test file."
        cd ..
        rm -rf "build_$ADVISORY_ID" "$ADVISORY_ID.yaml"
        cd "$BASE_DIR"
        continue
    fi
    
    mv "main.go" "../${ADVISORY_ID}_test.go"
    cd ..
    rm -rf "build_$ADVISORY_ID"

    if ! semgrep --validate --config "$ADVISORY_ID.yaml" > /dev/null 2>&1 || ! semgrep --test --config "$ADVISORY_ID.yaml" "${ADVISORY_ID}_test.go"; then
        echo "❌ GATE 1 FAILED: Rule failed Semgrep schema syntax or test validation."
        rm -f "$ADVISORY_ID.yaml" "${ADVISORY_ID}_test.go"
        cd "$BASE_DIR"
        continue
    fi

    # =========================================================================
    # GATE 1.5: SEMANTIC DEDUPLICATION
    # =========================================================================
    echo "Running Gate 1.5: Semantic Deduplication Check..."
    RULE_HASH=$(python3 "$BASE_DIR/hash_rule.py" "$ADVISORY_ID.yaml" 2>/dev/null || echo "error")
    
    if [ "$RULE_HASH" = "error" ]; then
        echo "❌ GATE 1.5 FAILED: Python helper could not parse YAML."
        rm -f "$ADVISORY_ID.yaml" "${ADVISORY_ID}_test.go"
        cd "$BASE_DIR"
        continue
    fi

    if grep -q "$RULE_HASH" "$HASH_DB"; then
        echo "❌ GATE 1.5 FAILED: Discarded. Rule is a semantic duplicate of an existing archetype."
        rm -f "$ADVISORY_ID.yaml" "${ADVISORY_ID}_test.go"
        cd "$BASE_DIR"
        continue
    fi
    
    echo "Running Gate 2: Strict False-Positive Check against standard library..."
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
    
    # Add the new hash to the DB so we don't generate a duplicate of it in this same run
    echo "$RULE_HASH" >> "$HASH_DB"
    NEW_RULES_GENERATED=$((NEW_RULES_GENERATED + 1))
    
    cd "$BASE_DIR"

    bash "push.sh"
done

echo "======================================"
echo "Pipeline Execution Complete."
echo "Total new rules generated: $NEW_RULES_GENERATED"

if [ "$NEW_RULES_GENERATED" -gt 0 ]; then
    echo "Running Gate 3: Auto-Committing batch to Experimental Tier..."
    cd "$BASE_DIR"
    
    git pull origin main --rebase --quiet || true 
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
