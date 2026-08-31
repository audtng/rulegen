---
name: semgrep-rulegen
description: >-
  Standard operational procedure for autonomous Golang Semgrep rule generation,
  commit fetching, testbed building, subagent synthesis, and dual-state validation
  targeting Medium, High, and Critical severity advisories.
---

# Golang Semgrep Rule Engineering Skill

This skill enforces the exact closed-loop workflow for synthesizing and validating Semgrep rules from the curated Medium/High/Critical Go advisory dataset (`/src/rulegen/rulegen/ghsa_golang_git_diffs_med_high_crit.json`).

## Orchestrator Protocol
1. **Target Setup**: `bash /src/rulegen/rulegen/pipeline.sh next <COUNT>`
2. **Subagent Spawning**: Spawn `semgrep_author` ($1\text{ subagent} : 1\text{ advisory}$) in waves of 5–7.
3. **Barrier Wait**: Passively await subagent completions.
4. **Audit**: Verify `validation_ledger.json` and auto-retry any failures.
5. **Report**: `bash /src/rulegen/rulegen/pipeline.sh report`

## Subagent Protocol (`semgrep_author`)
1. **Ingest (3 files only)**: `metadata.json`, `vuln.go`, `fixed.go`.
2. **Synthesize**: Write Semgrep YAML to `workspaces/<GHSA_ID>/rule.yaml`.
3. **Validate**: `bash /src/rulegen/rulegen/pipeline.sh validate <GHSA_ID> .../rule.yaml`.
4. **Decision Tree**: Refine on REJECTED (TP >= 1, FP == 0).
5. **Complete**: Report status and terminate.

## Strict Guardrails
- **Strict 1-to-1 Subagent Mapping**: Exactly 1 subagent per 1 advisory rule.
- **Isolated, Compact Testbeds**: Pre-filtered `vuln.go` and `fixed.go` (no 6000-line test fixtures).
- **Unified Engine Only**: Use ONLY `/src/rulegen/rulegen/pipeline.sh`.
- **NO External Network Calls**: All data is local.


