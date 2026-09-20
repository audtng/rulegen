#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

# --- Configuration ---
SOURCE_FILE="semgrep.yaml"
DEST_FOLDER="../go-scanner"

# --- Validation & Copying ---

# 1. Check if the source file exists
if [ ! -f "$SOURCE_FILE" ]; then
    echo "Error: Source file '$SOURCE_FILE' does not exist."
    exit 1
fi

# 2. Create destination folder if it doesn't exist yet
if [ ! -d "$DEST_FOLDER" ]; then
    echo "Destination folder missing. Creating '$DEST_FOLDER'..."
    mkdir -p "$DEST_FOLDER"
fi

# 3. Perform the copy operation
cp "$SOURCE_FILE" "$DEST_FOLDER"

# Confirm success
FILENAME=$(basename "$SOURCE_FILE")
echo "Success! CoPIED '$FILENAME' to '$DEST_FOLDER'."
