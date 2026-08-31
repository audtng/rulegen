# Antigravity Workspace Guidelines: Golang Semgrep Rule Generation

## 1. Orchestrator Standard Operating Procedure (Step-by-Step)

The Orchestrator Agent manages pipeline execution, concurrency control, and validation audit. It MUST follow these exact steps without deviation:

### Step O1: Target Preparation
- Execute via `run_command`:
  ```bash
  bash /src/rulegen/rulegen/pipeline.sh next [COUNT]
  ```
- Parses structured JSON emitted per advisory.
- Verifies testbed files are generated: `workspaces/<GHSA_ID>/vuln.go`, `fixed.go`, `metadata.json`, and `prompt.txt`.

### Step O2: Strict 1-to-1 Subagent Spawning
- **Rule**: Exactly 1 child subagent (`semgrep_author`) is spawned for 1 advisory ($1\text{ subagent} : 1\text{ advisory}$).
- **Concurrency**: Invoke subagents in controlled waves (5 to 7 subagents max per `invoke_subagent` call).
- **Tool**: Call `invoke_subagent` with `TypeName: "semgrep_author"`, `Role: "Author - <GHSA_ID>"`, and the exact prompt from `workspaces/<GHSA_ID>/prompt.txt`.

### Step O3: Reactive Barrier Wait
- The Orchestrator immediately stops calling tools and enters a passive wait state.
- PROHIBITION: Do NOT execute sleep loops, polling scripts, or status checks. The Antigravity messaging runtime automatically wakes the agent upon subagent completion.

### Step O4: Audit Reconciliation & Auto-Retry
- Inspect incoming completion messages.
- Verify that `rules/go/<GHSA_ID>.yaml` is created and recorded in `validation_ledger.json`.
- If any subagent encountered a transient API failure (e.g., 503), immediately re-spawn a single replacement subagent for that specific `<GHSA_ID>`.

### Step O5: Metrics & Progress Reporting
- Execute via `run_command`:
  ```bash
  bash /src/rulegen/rulegen/pipeline.sh report
  ```
- Output final structured audit table.

---

## 2. Child Subagent Standard Operating Procedure (`semgrep_author`)

Each child subagent operates in strict isolation on **exactly ONE advisory** following this rigid 5-step loop:

### Step S1: Ingest Testbed (3 Files Only)
- Call `view_file` on ONLY:
  1. `/src/rulegen/rulegen/workspaces/<GHSA_ID>/metadata.json`
  2. `/src/rulegen/rulegen/workspaces/<GHSA_ID>/vuln.go`
  3. `/src/rulegen/rulegen/workspaces/<GHSA_ID>/fixed.go`
- PROHIBITION: Do NOT search unrelated folders, do NOT read raw dataset JSON, do NOT make network calls (`curl`/`wget`).

### Step S2: Synthesize Grade A Structural AST Rule
- Analyze AST diff between `vuln.go` and `fixed.go`.
- **Mandatory Grade A Requirements**:
  1. *Structural AST Depth*: Do NOT write bare keywords, literal constants, or isolated function calls (e.g. `exec.LookPath(...)` alone is prohibited).
  2. *Enclosing Scope*: MUST capture enclosing function signature/receiver type (`func ($R *$TYPE) $METHOD(...) ...`), control-flow block (`if ... { ... }`), struct initialization, or assignment.
  3. *Negative Filtering*: Use `pattern-not` or `pattern-not-inside` whenever the remediation adds a guard check, sanitizer, or defensive parameter.
  4. *Metavariable Generalization*: Use uppercase metavariables (`$PARAM`, `$CTX`, `$REQ`) for local variable flexibility.
- Construct Semgrep YAML matching schema:
  ```yaml
  rules:
    - id: <GHSA_ID>
      languages: [go]
      severity: WARNING
      message: "<Summary>"
      metadata:
        cve: "<CVE>"
        ghsa: "<GHSA_ID>"
        confidence: HIGH
      patterns:
        - pattern: <STRUCTURAL_AST_VULNERABLE_PATTERN>
  ```
- Call `write_to_file` to save to `/src/rulegen/rulegen/workspaces/<GHSA_ID>/rule.yaml`.

### Step S3: Deterministic Dual-State Validation
- Execute via `run_command`:
  ```bash
  bash /src/rulegen/rulegen/pipeline.sh validate <GHSA_ID> /src/rulegen/rulegen/workspaces/<GHSA_ID>/rule.yaml
  ```

### Step S4: Refinement Decision Tree
- **`SUCCESS: Rule APPROVED`** ($TP \ge 1 \land FP = 0$): Proceed directly to Step S5.
- **`REJECTED: 0 matches in vuln.go`**: Relax metavariables/pattern, update `rule.yaml`, re-run Step S3.
- **`REJECTED: N matches in fixed.go`**: Add `pattern-not` or restrict pattern scope, update `rule.yaml`, re-run Step S3.
- **`REJECTED: Syntax errors`**: Fix YAML schema / indentation, update `rule.yaml`, re-run Step S3.
- (Maximum 3 refinement iterations).

### Step S5: Completion & Immediate Termination
- Send structured completion summary to orchestrator with:
  - Advisory ID & CVE
  - Package
  - Pattern YAML snippet
  - Verdict: APPROVED (TP >= 1, FP == 0)
- Stop calling tools and terminate turn.

---

## 3. Core Prohibitions & System Guardrails

1. **Strict 1-to-1 Subagent Mapping**: Exactly one subagent parses and generates exactly one rule. Never bundle multiple rules into a single subagent.
2. **Isolated, Compact Testbeds (Zero Code Bloat)**: Every advisory has its own separate `workspaces/<GHSA_ID>/vuln.go` and `fixed.go` containing exclusively high-signal security hunks. Test fixtures, vendor code, and generated files are filtered out so no subagent parses redundant thousands of lines of code.
3. **Unified Engine Only**: Use ONLY `/src/rulegen/rulegen/pipeline.sh`. Never create ad-hoc scripts.
4. **NO External Network Calls**: All commit patches, pre-patch, and post-patch files are pre-staged locally.
5. **Zero Context Pollution**: Every rule synthesis and validation iteration occurs exclusively within isolated child subagents.
6. **Zero Duplicate Reprocessing**: The pipeline engine automatically respects `validation_ledger.json` and committed rules in `rules/go/`.



