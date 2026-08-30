# Antigravity Autonomous Semgrep Rule Engineering Pipeline (Golang)

An end-to-end, consolidated agentic pipeline designed for **Google Antigravity** to parse Go security advisories, extract fix commits, reconstruct dual-state testbeds, and autonomously synthesize and self-validate high-precision Semgrep security rules in isolated subagent contexts.

---

## 1. Architectural Overview & Workflow

The entire system is self-contained within `/src/rulegen/` and driven through the unified pipeline engine [`pipeline.sh`](file:///src/rulegen/pipeline.sh):

```
+---------------------------------------------------------------------------------------------------+
| 1. DATASET & AUDIT INFRASTRUCTURE                                                                |
|    ├── ghsa_golang_git_diffs_med_high_crit.json  ──► 4,404 Med/High/Crit Go security advisories     |
|    ├── rules/go/<GHSA_ID>.yaml                   ──► Verified Semgrep rules store                 |
|    └── validation_ledger.json                    ──► Central audit ledger of validation verdicts  |
+---------------------------------------------------------------------------------------------------+
                                                  │
                                                  ▼
+---------------------------------------------------------------------------------------------------+
| 2. UNIFIED AGENTIC WORKFLOW LOOP                                                                  |
|                                                                                                   |
|    [ Step 1: Target Preparation ]                                                                 |
|    Command: bash /src/rulegen/pipeline.sh next [COUNT]                                            |
|    • Scans dataset, skips already-approved/non-actionable diffs                                   |
|    • Downloads patch, reconstructs /workspaces/<ID>/vuln.go & fixed.go                           |
|    • Compiles /workspaces/<ID>/prompt.txt & returns JSON with workspace paths                     |
|                                                                                                   |
|    [ Step 2: Subagent Synthesis & Self-Validation ]                                               |
|    Tool: invoke_subagent (TypeName: "semgrep_author")                                             |
|    • Child subagent reads vuln.go / fixed.go in isolated context                                  |
|    • Subagent crafts and writes /workspaces/<ID>/rule.yaml                                        |
|    • Subagent runs: bash /src/rulegen/pipeline.sh validate <ID> /workspaces/<ID>/rule.yaml        |
|    • Subagent iterates and refines until validation passes                                        |
|                                                                                                   |
|    [ Step 3: Deterministic Dual-State Validation & Store Commit ]                                 |
|    Command: bash /src/rulegen/pipeline.sh validate <ID> [RULE_FILE]                               |
|    • True Positive Check : vuln.go (>= 1 match)                                                   |
|    • False Positive Check: fixed.go (== 0 matches)                                                |
|    • Verdict: APPROVED ──► Automatically commits to /rules/go/<ID>.yaml & records in ledger       |
|               REJECTED ──► Emits diagnostic line matches for subagent refinement                  |
|                                                                                                   |
|    [ Step 4: Progress Tracking & Metrics ]                                                        |
|    Command: bash /src/rulegen/pipeline.sh report                                                  |
|    • Displays current approved rules, approval rate, and remaining actionable targets             |
+---------------------------------------------------------------------------------------------------+
```

---

## 2. Directory Layout (All Inside `/src/rulegen`)

```
/src/rulegen/
├── pipeline.sh                                     # Unified pipeline engine (next, prepare, validate, report, select)
├── README.md                                       # Complete system documentation
├── AGENTS.md                                       # Workspace guidelines & execution rules
├── validation_ledger.json                          # Central audit ledger of validation verdicts
│
├── rules/
│   └── go/
│       └── <GHSA_ID>.yaml                          # Production store of approved, validated Semgrep rules
│
├── workspaces/
│   └── <GHSA_ID>/
│       ├── commit.patch                            # Raw fix commit diff
│       ├── vuln.go                                 # Pre-patch code (True Positive target)
│       ├── fixed.go                                # Post-patch code (False Positive target)
│       ├── metadata.json                           # Advisory metadata (CVE, package, summary, diff)
│       ├── prompt.txt                              # Tailored subagent authoring prompt
│       ├── rule.yaml                               # Drafted Semgrep YAML rule
│       └── validation_result.json                  # Output verdict and match diagnostics
│
└── Datasets:
    ├── ghsa_golang_git_diffs_med_high_crit.json   # PRIMARY: 4,404 Med/High/Crit advisories (7.92 MB)
    ├── ghsa_golang_git_diffs.json                  # All 4,484 Go advisories with Git diffs (8.13 MB)
    ├── ghsa_golang_vulnerable_versions.json        # 7,663 version constraints for 1,442 packages (5.07 MB)
    └── ghsa_golang.json                            # 4,777 raw Go OSV records (18.56 MB)
```

---

## 3. Command Reference: `/src/rulegen/pipeline.sh`

| Command | Usage | Description |
| :--- | :--- | :--- |
| **`next`** | `bash /src/rulegen/pipeline.sh next [COUNT]` | Single-pass scan to find and prepare the next unhandled, actionable target(s). Outputs structured JSON for subagent invocation. |
| **`prepare`** | `bash /src/rulegen/pipeline.sh prepare <GHSA_ID>` | Prepares the workspace, testbed (`vuln.go`/`fixed.go`), and prompt for a specific advisory ID. |
| **`validate`** | `bash /src/rulegen/pipeline.sh validate <GHSA_ID> [RULE]` | Evaluates Semgrep rule against `vuln.go` ($TP \ge 1$) and `fixed.go` ($FP = 0$). Auto-commits approved rules to `rules/go/` and updates `validation_ledger.json`. |
| **`report`** | `bash /src/rulegen/pipeline.sh report` | Prints real-time metrics, total store rules, approval rate, and remaining targets. |
| **`select`** | `bash /src/rulegen/pipeline.sh select [-n N] [-s SEV]` | Queries unhandled candidate advisory IDs from the dataset. |

---

## 4. Subagent Contract & System Prompt (`semgrep_author`)

The `semgrep_author` subagent is spawned via `invoke_subagent` to synthesize and self-validate rules in an isolated context:

### Subagent Spec:
- **`TypeName`**: `semgrep_author`
- **`Role`**: `Go Semgrep Author - <GHSA_ID>`
- **`Tools`**: File tools (`view_file`, `write_to_file`) + `run_command` (for `pipeline.sh validate`).

### Subagent Execution Cycle:
1. Subagent reads `workspaces/<GHSA_ID>/vuln.go`, `fixed.go`, and `metadata.json`.
2. Subagent drafts Semgrep rule and writes to `workspaces/<GHSA_ID>/rule.yaml`.
3. Subagent validates directly using:
   ```bash
   bash /src/rulegen/pipeline.sh validate <GHSA_ID> /src/rulegen/workspaces/<GHSA_ID>/rule.yaml
   ```
4. If validation fails, subagent inspects diagnostics and refines `rule.yaml`.
5. Upon approval, subagent summarizes the verified pattern and reports completion to the orchestrator.

---

## 5. Dual-State Validation Criteria & Ground Truth

| Test Criterion | Test Target | Pass Requirement |
| :--- | :--- | :--- |
| **True Positive (TP)** | `vuln.go` (Pre-patch state) | $\ge 1$ match at modified vulnerable lines |
| **False Positive (FP)** | `fixed.go` (Post-patch state) | Exactly $0$ matches |
| **Syntax Integrity** | Semgrep YAML Schema | Valid top-level `rules:`, `languages: [go]`, `id:` |
| **Store Commit** | `/src/rulegen/rules/go/<GHSA_ID>.yaml` | Automatically written on `APPROVED` |
| **Audit Ledger** | `/src/rulegen/validation_ledger.json` | Persistent record with timestamp and match details |

---

## 6. Prohibitions & Guardrails

- **Unified Engine Only**: Use ONLY `/src/rulegen/pipeline.sh`. Do not create fragmented or ad-hoc scripts.
- **NO External Network Calls**: All patch diffs and metadata are pre-fetched locally in `/src/rulegen/workspaces/<GHSA_ID>/`. Subagents must never call `curl`, `wget`, or query GitHub APIs.
- **NO Context Pollution**: Every rule synthesis iteration happens in an isolated child subagent context window.
- **NO Duplicate Reprocessing**: The pipeline engine automatically respects `validation_ledger.json` and existing rules in `rules/go/`.
