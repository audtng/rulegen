# Antigravity Workspace Guidelines: Golang Semgrep Rule Generation

## 1. Orchestrator Standard Operating Procedure (Step-by-Step)
- **Step O1 (Target Preparation)**: Execute `bash /src/rulegen/rulegen/pipeline.sh next [COUNT]` to construct isolated workspaces.
- **Step O2 (Strict 1-to-1 Subagent Spawning)**: Spawn exactly 1 child subagent (`semgrep_author`) per advisory ($1\text{ subagent} : 1\text{ advisory}$) in controlled waves of 5–7 workers.
- **Step O3 (Reactive Barrier Wait)**: Stop tool execution and passively await subagent completion messages (no polling / sleeping).
- **Step O4 (Audit Reconciliation & Auto-Retry)**: Verify `validation_ledger.json` and `rules/go/<GHSA_ID>.yaml`. Automatically retry any transient failures.
- **Step O5 (Reporting)**: Execute `bash /src/rulegen/rulegen/pipeline.sh report` and output results.

## 2. Child Subagent Standard Operating Procedure (`semgrep_author`)
- **Step S1 (Ingest Testbed)**: Read only `metadata.json`, `vuln.go`, and `fixed.go` in `/src/rulegen/rulegen/workspaces/<GHSA_ID>/`.
- **Step S2 (Synthesize Rule)**: Draft Semgrep YAML and write to `workspaces/<GHSA_ID>/rule.yaml`.
- **Step S3 (Validate)**: Run `bash /src/rulegen/rulegen/pipeline.sh validate <GHSA_ID> /src/rulegen/rulegen/workspaces/<GHSA_ID>/rule.yaml`.
- **Step S4 (Decision Tree)**: Refine `rule.yaml` on REJECTED (TP >= 1, FP == 0) up to 3 iterations.
- **Step S5 (Completion)**: Send structured completion report and terminate turn.

## 3. Core Prohibitions & System Guardrails
- **Strict 1-to-1 Subagent Mapping**: Exactly 1 subagent per 1 advisory rule. Never bundle multiple rules into a subagent.
- **Isolated, Compact Testbeds (Zero Code Bloat)**: Separate `vuln.go` and `fixed.go` per advisory, pre-filtered to remove huge test fixtures, vendor code, and generated files.
- **Unified Engine Only**: Use ONLY `/src/rulegen/rulegen/pipeline.sh`.
- **NO External Network Calls**: All data is local.
- **Zero Context Pollution**: Main orchestrator never drafts rules; all synthesis happens in isolated subagents.



