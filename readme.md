🚀 RuleGen: Autonomous Variant Analysis Engine

RuleGen is a zero-human-in-the-loop static analysis pipeline that transforms GitHub Security Advisory (GHSA) and CVE JSON data into weaponized, production-ready Semgrep rules.

Built for Go vulnerability research, the engine uses a heavily constrained LLM prompt to generate abstract syntax tree (AST) patterns, validates them against native compilers, enforces a strict false-positive baseline, deduplicates them semantically, and autonomously deploys them to a fleet scanner for live zero-day hunting.
🏗️ Architecture

The system operates as a multi-gate factory, ruthlessly discarding brittle, hallucinated, or noisy rules before they ever reach the scanner.

    Ingestion (combined_output/*.json): Parses raw vulnerability metadata and source code diffs, filtering out massive architectural rewrites that cannot be modeled by AST rules.

    Generation (Universal Omnibus Prompt): Uses the Antigravity CLI (agy) to prompt an LLM to extract the core vulnerability archetype (e.g., SSRF, CWE-22, CWE-287) and output strictly formatted Semgrep YAML and Go test cases.

    Gate 1: Native Compilation Check: Isolates the generated _test.go file, dynamically fetches missing dependencies via go mod init and goimports, and passes it to the native Go compiler (go build). If the LLM hallucinated packages or wrote invalid Go, the rule is instantly killed.

    Gate 1.5: Semantic Deduplication: Parses the surviving YAML using Python, extracts only the logical data flow (pattern-sources, pattern-propagators, pattern-sinks), sorts it deterministically, and hashes it. If the engine has already learned this AST pattern from a previous CVE, the duplicate is dropped.

    Gate 2: The FP Kill Switch: Executes the compiled rule against the entire Go Standard Library (/usr/local/go/src). If the rule generates even a single hit on known-safe baseline code, it is permanently discarded.

    Gate 3: Autonomous Deployment: Survivors are automatically committed and pushed to the rules/experimental/ directory on the main branch, ready for fleet-wide execution.

⚙️ Prerequisites

The pipeline relies on native OS and networking tools to validate syntax dynamically. Ensure the host environment has:

    Go (with goimports installed and $GOPATH/bin in your $PATH)

    Semgrep (pip install semgrep)

    Python 3 (with pip install pyyaml for semantic hashing)

    jq and awk (for JSON parsing and payload extraction)

    Antigravity CLI (agy)

📁 Repository Structure
Plaintext

.
├── rulegen.sh                 # The core autonomous pipeline engine
├── fleet_scan.sh              # Mass-execution script for live hunting
├── prompt_template.txt        # The Universal Omnibus Prompt instructing the LLM
├── hash_rule.py               # AST logic extractor and semantic hasher
├── combined_output/           # Drop your raw GHSA/CVE JSON files here
├── rules/
│   ├── experimental/          # Auto-merged, zero-FP rules ready for testing
│   └── stable/                # Manually verified rules that have caught zero-days
└── scan_results/              # Output directory for live fleet scanner hits

🚀 Usage
1. Generate Rules

Populate the combined_output/ directory with your target CVE/GHSA JSON files. Start the engine and walk away.
Bash

chmod +x pipeline.sh
./pipeline.sh

The script will output logs detailing which rules passed compilation, survived the standard library false-positive check, and were successfully merged into the experimental tier.


LLMs naturally overfit to the specific variable names of the CVE they are analyzing. RuleGen combats this through Universal Prompt Constraints (forcing the LLM to target abstract standard library interfaces like net/http or tar.Header) and Semantic Hashing (stripping out markdown, metadata, and YAML formatting before comparing rule logic to prevent identical AST flows from bloating the registry).
