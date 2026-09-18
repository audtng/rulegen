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

git -c credential.helper="!f() { echo \"username=${USERNAME}\"; echo \"password=${TOKEN}\"; }; f" push -u origin "$BRANCH_NAME"

git checkout main

echo ""

echo "Successful branch push"

