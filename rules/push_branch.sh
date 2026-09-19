#This is a script to push branches to github

#!/bin/bash

set -e

USERNAME="audtng"

if [ -n "$GITHUB_TOKEN" ]; then
    TOKEN="$GITHUB_TOKEN"
else
    echo -n "Enter token: "
    read -s TOKEN
    echo "" 
fi

git checkout -b "$BRANCH_NAME"

git add "$CWE_NAME.yaml" "${CWE_NAME}_test.go"

git commit -m "feat(rules): Add highly-generalized Grade A rule for $CWE_NAME" || true

git -c credential.helper="!f() { echo \"username=${USERNAME}\"; echo \"password=${TOKEN}\"; }; f" push -f -u origin "$BRANCH_NAME"

echo "Creating Pull Request on GitHub..."

# Capture the response from GitHub instead of sending it to /dev/null
API_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST -H "Authorization: token $TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    -d "{
          \"title\": \"🤖 Auto-Generated Rule: $CWE_NAME\",
          \"head\": \"$BRANCH_NAME\",
          \"base\": \"main\",
          \"body\": \"This PR introduces a highly-validated, Grade-A Semgrep rule for **$CWE_NAME**.\\n\\nPlease review the diff to ensure logical accuracy before merging.\"
        }" \
    "https://api.github.com/repos/$USERNAME/rulegen/pulls")

# Extract the HTTP status code from the bottom of the response
HTTP_STATUS=$(echo "$API_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)

if [ "$HTTP_STATUS" -eq 201 ]; then
    echo "✅ Pull Request created successfully!"
else
    echo "❌ Failed to create PR. GitHub API Response:"
    # Print the error message from GitHub (removing our status code hack)
    echo "$API_RESPONSE" | sed '/HTTP_STATUS/d' | jq .
fi

git checkout main

echo ""

echo "Successful branch push"

