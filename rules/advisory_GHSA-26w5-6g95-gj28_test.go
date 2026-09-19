package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Workspace string `json:"workspace"`
}

// 1. Direct stdlib: Unmarshaled input joined directly into dangerous filesystem call
func testDirectStdlib(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	target := filepath.Join("/base/dir", cfg.Workspace)
	// ruleid: path-traversal-unmarshaled-data
	_ = os.RemoveAll(target)
}

// 2. Proper patch: Validate untrusted component using filepath.IsLocal
func testProperPatch(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	if filepath.IsLocal(cfg.Workspace) {
		target := filepath.Join("/base/dir", cfg.Workspace)
		// ok: path-traversal-unmarshaled-data
		_ = os.RemoveAll(target)
	}
}

// 3. Cross-function taint: Untrusted input flows across helper function boundary
func formatSubdir(sub string) string {
	return filepath.Join("subdir", sub)
}

func testCrossFunction(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	formatted := formatSubdir(cfg.Workspace)
	target := filepath.Join("/base/dir", formatted)
	// ruleid: path-traversal-unmarshaled-data
	_ = os.RemoveAll(target)
}

// 4. Interface bypass: Unmarshaled into generic map/interface
func testInterfaceBypass(data []byte) {
	var raw map[string]any
	_ = json.Unmarshal(data, &raw)
	val, ok := raw["path"].(string)
	if ok {
		target := filepath.Join("/base/dir", val)
		// ruleid: path-traversal-unmarshaled-data
		_ = os.RemoveAll(target)
	}
}

// 5. Fake sanitizer: filepath.Clean does not prevent traversal outside root
func testFakeSanitizer(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	cleaned := filepath.Clean(cfg.Workspace)
	target := filepath.Join("/base/dir", cleaned)
	// ruleid: path-traversal-unmarshaled-data
	_ = os.RemoveAll(target)
}

// 6. Real sanitizer: filepath.Base removes all directory traversal components
func testRealSanitizer(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	safe := filepath.Base(cfg.Workspace)
	target := filepath.Join("/base/dir", safe)
	// ok: path-traversal-unmarshaled-data
	_ = os.RemoveAll(target)
}
