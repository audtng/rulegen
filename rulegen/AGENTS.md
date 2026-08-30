# Antigravity Workspace Guidelines: Golang Semgrep Rule Generation

## 1. Unified Pipeline Architecture & Subagent Self-Validation
All operations strictly follow the unified pipeline engine in `/src/rulegen/pipeline.sh`.
The orchestrator manages subagents and pipeline flow; the child subagent (`semgrep_author`) synthesizes and **self-validates** the rule within its own isolated workspace.

### The Unified Workflow:
1. **Prepare Actionable Target(s)**:
   ```bash
   bash /src/rulegen/pipeline.sh next [COUNT]
   ```
   - Automatically queries `/src/rulegen/ghsa_golang_git_diffs_med_high_crit.json`.
   - Reconstructs testbed (`vuln.go`, `fixed.go`, `metadata.json`).
   - Verifies code actionability (skipping non-Go and dependency bumps).
   - Generates the subagent prompt at `workspaces/<GHSA_ID>/prompt.txt`.
   - Returns structured JSON with workspace paths ready for direct subagent invocation.

2. **Delegate Rule Synthesis & Self-Validation to Subagent**:
   - Invoke `semgrep_author` subagent via `invoke_subagent` using the generated prompt.
   - Child subagent inspects testbed files, authors `workspaces/<GHSA_ID>/rule.yaml`, and directly executes:
     ```bash
     bash /src/rulegen/pipeline.sh validate <GHSA_ID> /src/rulegen/workspaces/<GHSA_ID>/rule.yaml
     ```
   - Subagent iterates and refines until validation passes ($TP \ge 1 \land FP = 0$).
   - When approved, `pipeline.sh validate` automatically commits the rule to `/src/rulegen/rules/go/<GHSA_ID>.yaml` and updates `/src/rulegen/validation_ledger.json`.

3. **Status & Progress Tracking**:
   ```bash
   bash /src/rulegen/pipeline.sh report
   ```

## 2. Core Prohibitions & Subagent Guardrails
- **Unified Engine Only**: Use ONLY `/src/rulegen/pipeline.sh`.
- **NO `curl` / `wget` / GitHub API Calls by Subagents**: All required data is pre-fetched locally in `/src/rulegen/workspaces/<GHSA_ID>/`.
- **NO Ad-hoc Script Generation**: Subagents must use ONLY `pipeline.sh validate`.
- **NO Context Pollution**: Every rule generation and validation iteration happens in isolated `semgrep_author` subagents.
- **NO Duplicate Reprocessing**: Automatically respects `/src/rulegen/validation_ledger.json` and `/src/rulegen/rules/go/`.
