import sys, json, hashlib

try:
    import yaml
except ImportError:
    print("FATAL: Missing pyyaml dependency. Run: pip3 install pyyaml", file=sys.stderr)
    sys.exit(1)

try:
    with open(sys.argv[1], 'r') as f:
        data = yaml.safe_load(f)
        rule = data['rules'][0]
        
        # Universal AST Extraction (Supports both Taint and Search modes)
        core = {
            'sources': rule.get('pattern-sources', []),
            'sinks': rule.get('pattern-sinks', []),
            'patterns': rule.get('patterns', []),
            'pattern': rule.get('pattern', None)
        }
        
        dump = json.dumps(core, sort_keys=True)
        print(hashlib.sha256(dump.encode()).hexdigest())
except Exception as e:
    print(f"FATAL: Python parse error: {e}", file=sys.stderr)
    sys.exit(1)
