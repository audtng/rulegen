---
name: semgrep-rulegen
description: >-
  Standard operational procedure for autonomous Golang Semgrep rule generation,
  commit fetching, testbed building, subagent synthesis, and dual-state validation
  targeting Medium, High, and Critical severity advisories.
---

# Golang Semgrep Rule Engineering Skill

This skill enforces the exact closed-loop workflow for synthesizing and validating Semgrep rules from the curated Medium/High/Critical Go advisory dataset (`/src/rulegen/ghsa_golang_git_diffs_med_high_crit.json`).

## Standard Execution Sequence

1. **Testbed Setup**:
   `bash /src/rulegen/01_fetch_commit.sh <GHSA_ID>`

2. **Prompt Compilation**:
   `bash /src/rulegen/02_generate_prompt.sh <GHSA_ID>`

3. **Isolated Subagent Authoring**:
   Call `invoke_subagent` with `TypeName: "self"`, `Role: "Go Semgrep Author - <GHSA_ID>"`, and the prompt from `/src/rulegen/workspaces/<GHSA_ID>/prompt.txt`.

4. **Deterministic Validation**:
   `bash /src/rulegen/03_validate_rule.sh <GHSA_ID>`

5. **Batch Processing**:
   `bash /src/rulegen/04_run_pipeline.sh "" <N>`

## Rules & Constraints
- Always use the standalone scripts in `/src/rulegen/`.
- Default dataset is `/src/rulegen/ghsa_golang_git_diffs_med_high_crit.json`.
- Never create ad-hoc scripts.
- Never author rules in the main conversation context.
- Always require $TP \ge 1 \land FP == 0$ in validation.
