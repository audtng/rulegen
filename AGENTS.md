# Antigravity Workspace Guidelines: Golang Semgrep Rule Generation

## 1. Mandatory Workflow Enforcement
All future Antigravity sessions working in this repository MUST strictly follow the established 4-step pipeline and use ONLY the existing standalone scripts in `/src/rulegen/`.

DO NOT create ad-hoc scripts or bypass the pipeline.

### The 4-Step Standard Pipeline:
1. **Fetch Commit & Reconstruct Testbed**:
   ```bash
   bash /src/rulegen/01_fetch_commit.sh <GHSA_ID>
   ```
   - Uses local dataset `/src/rulegen/ghsa_golang_git_diffs_med_high_crit.json`.
   - Generates `/src/rulegen/workspaces/<GHSA_ID>/vuln.go` and `fixed.go`.

2. **Generate Subagent Prompt**:
   ```bash
   bash /src/rulegen/02_generate_prompt.sh <GHSA_ID>
   ```
   - Creates compact prompt in `/src/rulegen/workspaces/<GHSA_ID>/prompt.txt`.

3. **Synthesize Rule in Isolated Subagent Context**:
   - MUST invoke a fresh child subagent via `invoke_subagent`.
   - Never author rules directly in the main orchestrator conversation to prevent context window bloat and token accumulation.

4. **Deterministic Dual-State Validation**:
   ```bash
   bash /src/rulegen/03_validate_rule.sh <GHSA_ID>
   ```
   - Tests rule against `vuln.go` (asserts matches >= 1) and `fixed.go` (asserts matches == 0).
   - Only approved rules are committed to `/src/rulegen/rules/go/<GHSA_ID>.yaml`.
   - Central audit ledger is recorded in `/src/rulegen/validation_ledger.json`.

5. **Batch Driver**:
   ```bash
   bash /src/rulegen/04_run_pipeline.sh [GHSA_ID] [MAX_ITEMS]
   ```

## 2. Core Prohibitions
- **NO Ad-hoc Script Generation**: Use only the established standalone scripts.
- **NO External API Keys**: Operates 100% on the internal Antigravity runtime.
- **NO Context Pollution**: Every rule generation task MUST happen in an isolated subagent chat context.
- **NO Low Severity Distractions**: Targets the curated Medium, High, and Critical severity dataset (`/src/rulegen/ghsa_golang_git_diffs_med_high_crit.json`).
