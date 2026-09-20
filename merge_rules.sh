#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

# Configuration
INPUT_DIR="./rules"
OUTPUT_FILE="combined_rules.yaml"

# Check if yq is installed
if ! command -v yq &> /dev/null; then
    echo "Error: 'yq' is required but not installed."
    echo "Please install it from https://github.com"
    exit 1
fi

# Check if input directory exists
if [ ! -d "$INPUT_DIR" ]; then
    echo "Error: Input directory '$INPUT_DIR' does not exist."
    exit 1
fi

echo "Scanning '$INPUT_DIR' for YAML files..."

# Find all .yaml and .yml files and read them using yq
# This command automatically handles files whether they start with 'rules:' or are bare dictionaries/lists
yq eval-all '
  . as $item |
  select(has("rules")) |= .rules |
  [.] | flatten |
  {"rules": .}
' "$INPUT_DIR"/*.{yaml,yml} 2>/dev/null > "$OUTPUT_FILE" || true

# Check if the output file was successfully created and has content
if [ -s "$OUTPUT_FILE" ]; then
    # Count how many rules were processed by checking the array size
    RULE_COUNT=$(yq '.rules | length' "$OUTPUT_FILE")
    echo "Success! Merged $RULE_COUNT rules into '$OUTPUT_FILE'."
else
    echo "Error: No valid YAML rule files found or merge failed."
    rm -f "$OUTPUT_FILE"
    exit 1
fi
