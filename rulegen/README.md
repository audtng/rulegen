# Antigravity Autonomous Semgrep Rule Engineering Pipeline (Golang)

An end-to-end, script-first agentic ecosystem designed for **Google Antigravity** to parse the GitHub Advisory Database, filter **Medium, High, and Critical severity** Go advisories, extract vulnerable dependencies and Git diffs, and autonomously synthesize and validate high-precision Semgrep security rules in isolated subagent chat contexts.

---

## 1. Architectural Overview & Design Philosophy

The system transforms **364,715 raw OSV advisory records** into verified, production-ready Semgrep security rules through a **hybrid script-agent architecture**:

```
+---------------------------------------------------------------------------------------------------+
| 1. DATA EXTRACTION, SEVERITY FILTERING & NORMALIZATION (Deterministic Core)                       |
|    GitHub Advisory DB (364k JSONs)                                                               |
|    ├── move_to_ghsa.sh                 ──► Flatten all advisories into /ghsa/                     |
|    ├── combine_ghsa.sh                 ──► Consolidate into ghsa_combined.json (541 MB)           |
|    ├── filter_golang.sh                ──► Filter 4,777 Go advisories (ghsa_golang.json)          |
|    ├── extract_vulnerable_versions.sh ──► Package-Centric Index (1,442 Go packages)               |
|    ├── extract_git_diffs.sh            ──► Validated Git Diffs & Commits (ghsa_golang_git_diffs)  |
|    └── filter_diffs_by_severity.sh     ──► Filter Med/High/Crit Dataset (4,404 advisories)        |
|                                            (ghsa_golang_git_diffs_med_high_crit.json)             |
+---------------------------------------------------------------------------------------------------+
                                                  │
                                                  ▼
+---------------------------------------------------------------------------------------------------+
| 2. AUTONOMOUS SYNTHESIS & VALIDATION LOOP (Antigravity Agentic Ecosystem)                        |
|                                                                                                   |
|    [ 01_fetch_commit.sh <GHSA_ID> ]                                                               |
|    └── Fetches remote commit patch & reconstructs vuln.go (pre-patch) & fixed.go (post-patch)     |
|                                                                                                   |
|    [ 02_generate_prompt.sh <GHSA_ID> ]                                                            |
|    └── Builds compact ~300-token prompt in /workspaces/<GHSA_ID>/prompt.txt                       |
|                                                                                                   |
|    [ Antigravity Subagent: invoke_subagent ]                                                      |
|    └── Fresh, isolated chat context per advisory; synthesizes rule.yaml                           |
|                                                                                                   |
|    [ 03_validate_rule.sh <GHSA_ID> ]                                                              |
|    └── Dual-State Ground-Truth Oracle:                                                            |
|        • True Positive Check : vuln.go (>= 1 match)                                               |
|        • False Positive Check: fixed.go (== 0 matches)                                            |
|        • Verdict: APPROVED ──► /rules/go/<GHSA_ID>.yaml & records to validation_ledger.json       |
|                   REJECTED ──► Diagnostic feedback emitted for in-context refinement              |
+---------------------------------------------------------------------------------------------------+
```

---

## 2. Google Antigravity (AGY) Integration & Runtime Context

This pipeline is built to leverage the native multi-agent architecture and execution primitives of the **Google Antigravity Agent Runtime**:

### A. Subagents & Isolated Chat Contexts (`invoke_subagent`)
- **Conversation Isolation**: Calling `invoke_subagent` spawns an independent child conversation with its own unique `conversationId`.
- **Zero Token Drift**: Each advisory is authored in a pristine, compact context ($< 1,000$ tokens total context) containing only the pre-digested Go diff hunk and metadata.
- **Workspace Modes**:
  - `Workspace: "inherit"`: Shares the `/src` workspace for read/write access to `/workspaces/<GHSA_ID>/` and `/rules/go/`.
  - `Workspace: "branch"` or `"share"`: Optional isolated worktree checkouts.

### B. Reactive Message Passing (Zero Polling Overhead)
- **Automatic Wakeups**: When a subagent finishes or a background task completes, the runtime automatically delivers the completion message and resumes the orchestrator turn without polling loops.

### C. Transcripts & Auditability
- Every subagent execution produces a structured JSONL transcript stored at:
  ```
  <appDataDir>/brain/<conversation-id>/.system_generated/logs/transcript.jsonl
  ```

### D. Native Tools Used in the Loop
| Tool | Role in Pipeline |
| :--- | :--- |
| `invoke_subagent` | Spawns fresh subagent contexts for advisory rule authoring |
| `run_command` | Executes deterministic shell scripts (`01_fetch`, `02_generate`, `03_validate`, `04_run`) |
| `view_file` | Inspects testbeds, generated prompts, and validation ledgers |
| `write_to_file` / `replace_file_content` | Writes and edits YAML rules and metadata |
| `manage_task` | Monitors background downloads and long-running batch jobs |

### E. Script-First Strategy: Zero External LLM Keys
- **0% LLM on Deterministic Tasks**: Git patch retrieval, unified diff parsing, Go testbed generation, and AST evaluation are 100% offloaded to local scripts.
- **No Third-Party SaaS / Keys Required**: Runs entirely on Antigravity's internal model.

---

## 3. Project Directory Structure

```
/src/rulegen/
├── README.md                                       # Complete system & Antigravity documentation
│
├── Data Pipeline Scripts:
│   ├── move_to_ghsa.sh                             # Flattens nested advisory folders into /ghsa/
│   ├── combine_ghsa.sh                             # Aggregates 364k JSON files into master JSON array
│   ├── filter_golang.sh                            # Filters Golang-specific advisories (4,777 records)
│   ├── extract_vulnerable_versions.sh              # Generates package-centric vulnerable version mappings
│   ├── extract_git_diffs.sh                        # Extracts validated fix commits & diff comparison URLs
│   └── filter_diffs_by_severity.sh                 # Extracts Med/High/Crit dataset (4,404 records)
│
├── Autonomous Rule Generation Suite:
│   ├── 01_fetch_commit.sh                          # Fetches commit patch and builds vuln.go & fixed.go
│   ├── 02_generate_prompt.sh                       # Generates compact author prompt for Antigravity
│   ├── 03_validate_rule.sh                         # Dual-state verification oracle (vuln vs fixed)
│   └── 04_run_pipeline.sh                          # Master driver for sequential & batch execution
│
├── Generated Datasets:
│   ├── ghsa_golang_git_diffs_med_high_crit.json   # PRIMARY DATASET: 4,404 Med/High/Crit advisories (7.92 MB)
│   ├── ghsa_golang_git_diffs.json                  # All 4,484 Go advisories with Git diffs (8.13 MB)
│   ├── ghsa_golang_vulnerable_versions.json        # 7,663 version constraints for 1,442 packages (5.07 MB)
│   └── ghsa_golang.json                            # 4,777 raw Go OSV records (18.56 MB)
│
├── Workspaces & Artifacts:
│   ├── workspaces/<GHSA_ID>/                       # Isolated testbed per advisory (vuln.go, fixed.go, rule.yaml)
│   ├── rules/go/<GHSA_ID>.yaml                     # Store of approved, validated Semgrep rules
│   └── validation_ledger.json                      # Central audit ledger of all verification verdicts
```

---

## 4. Dataset Specifications & Schemas

### Primary Dataset: `ghsa_golang_git_diffs_med_high_crit.json`
Contains **4,404 Medium, High, and Critical severity** Go security advisories (93.8% diff coverage), enriched with explicit severity and CVSS scores:

```json
{
  "GHSA-h395-qcrw-5vmq": {
    "severity": "HIGH",
    "cvss_score": "CVSS:3.1/AV:N/AC:L/PR:N/UI:R/S:U/C:H/I:L/A:N",
    "aliases": ["CVE-2020-28483"],
    "package": "github.com/gin-gonic/gin",
    "repository_url": "https://github.com/gin-gonic/gin",
    "summary": "Inconsistent Interpretation of HTTP Requests in github.com/gin-gonic/gin",
    "validation": {
      "has_diff_source": true,
      "has_commit_diff": true,
      "has_version_diff": true,
      "has_pr_diff": false,
      "warnings": []
    },
    "fix_commits": [
      {
        "commit_sha": "03e5e05ae089bc989f1ca41841f05504d29e3fd9",
        "commit_url": "https://github.com/gin-gonic/gin/commit/03e5e05ae089bc989f1ca41841f05504d29e3fd9",
        "diff_url": "https://github.com/gin-gonic/gin/commit/03e5e05ae089bc989f1ca41841f05504d29e3fd9.diff",
        "patch_url": "https://github.com/gin-gonic/gin/commit/03e5e05ae089bc989f1ca41841f05504d29e3fd9.patch",
        "git_show_cmd": "git show 03e5e05ae089bc989f1ca41841f05504d29e3fd9"
      }
    ],
    "version_diffs": [
      {
        "introduced_version": "0",
        "fixed_version": "1.6.0",
        "last_vulnerable_version": null,
        "fixed_tag": "v1.6.0",
        "compare_url": "https://github.com/gin-gonic/gin/compare/v1.5.0...v1.6.0",
        "diff_url": "https://github.com/gin-gonic/gin/compare/v1.5.0...v1.6.0.diff",
        "git_diff_cmd": "git diff v1.5.0..v1.6.0"
      }
    ]
  }
}
```

---

## 5. Step-by-Step Execution Guide

### Step 1: Prepare Testbed for an Advisory
Fetches the commit patch from the remote Git forge and generates `vuln.go` and `fixed.go` in `/src/rulegen/workspaces/<GHSA_ID>/`:
```bash
bash /src/rulegen/01_fetch_commit.sh GHSA-2286-hxv5-cmp2
```

### Step 2: Generate Antigravity Author Prompt
Extracts minimal diff hunks and builds `/src/rulegen/workspaces/<GHSA_ID>/prompt.txt`:
```bash
bash /src/rulegen/02_generate_prompt.sh GHSA-2286-hxv5-cmp2
```

### Step 3: Author Rule via Isolated Antigravity Subagent Context
In Antigravity, invoke a fresh subagent:
```json
{
  "TypeName": "self",
  "Role": "Go Semgrep Author - GHSA-2286-hxv5-cmp2",
  "Prompt": "Read /src/rulegen/workspaces/GHSA-2286-hxv5-cmp2/prompt.txt, inspect vuln.go vs fixed.go, synthesize rule.yaml, and run bash /src/rulegen/03_validate_rule.sh GHSA-2286-hxv5-cmp2."
}
```

### Step 4: Validate Rule & Record Verdict
Runs dual-state ground-truth verification:
```bash
bash /src/rulegen/03_validate_rule.sh GHSA-2286-hxv5-cmp2
```

### Step 5: Batch Pipeline Driver
To prepare multiple advisories in sequence from the Medium/High/Critical dataset:
```bash
# Process next 5 advisories
bash /src/rulegen/04_run_pipeline.sh "" 5

# Or target a specific advisory
bash /src/rulegen/04_run_pipeline.sh GHSA-2286-hxv5-cmp2
```

---

## 6. Verification Criteria & Decision Matrix

| Test Criterion | Method | Pass Requirement |
| :--- | :--- | :--- |
| **True Positive (TP)** | Match against `vuln.go` (Pre-patch state) | $\ge 1$ match at modified lines |
| **False Positive (FP)** | Match against `fixed.go` (Post-patch state) | Exactly $0$ matches |
| **Syntax Integrity** | YAML & Semgrep Schema Validation | Valid top-level `rules:`, `languages: [go]`, `id:` |
| **Action on PASS** | Store approved YAML in `/rules/go/<GHSA_ID>.yaml` | Ledger marked `APPROVED` |
| **Action on FAIL** | Emit match lines & error reasons | Diagnostic fed back for in-context refinement |

---

## 7. Validated Rules Catalog (Sample)

| Advisory ID | Affected Package | Severity | CVE | Verified Pattern |
| :--- | :--- | :---: | :--- | :--- |
| `GHSA-227x-7mh8-3cf6` | `gardener-extension-provider-aws` | `HIGH` | CVE-2025-59823 | `featurevalidation.ValidateFeatureGates(...)` |
| `GHSA-2286-hxv5-cmp2` | `github.com/bishopfox/sliver` | `HIGH` | CVE-2026-25760 | `os.ReadFile(filepath.Join($DIR, $W.Path))` |
| `GHSA-22qq-3xwm-r5x4` | `github.com/cometbft/cometbft` | `HIGH` | CVE-2025-24371 | `$POOL.Logger.Debug("Ignoring banned peer", $PEER)` |
| `GHSA-239w-m3h6-ch8v` | `github.com/filebrowser/filebrowser/v2` | `HIGH` | CVE-2026-54094 | `WithinScope($FS, $PATH)` |
| `GHSA-2464-8j7c-4cjm` | `github.com/go-viper/mapstructure/v2` | `MODERATE` | CVE-2025-11065 | `return time.ParseDuration(...)` |
| `GHSA-h395-qcrw-5vmq` | `github.com/gin-gonic/gin` | `HIGH` | CVE-2020-28483 | `$CIDRS, _ := $C.engine.prepareTrustedCIDRs()` |

---

## 8. License & Provenance

- **Advisory Source Data**: [GitHub Advisory Database](https://github.com/github/advisory-database) (CC-BY 4.0).
- **Semgrep Rules**: Synthesized and verified autonomously by the Antigravity Multi-Agent Pipeline.
