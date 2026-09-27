package rules

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var yaml = struct {
	Unmarshal func([]byte, any) error
}{}

type Config struct {
	Path string `json:"path"`
}

// 1. Vulnerable: unmarshaled JSON path directly opened
func testVulnerableJSONDirect(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	// ruleid: path-traversal-unmarshaled-file
	_, _ = os.Open(cfg.Path)
}

// 2. Vulnerable: unmarshaled YAML path joined and read
func testVulnerableYAMLJoin(data []byte) {
	var cfg Config
	_ = yaml.Unmarshal(data, &cfg)
	target := filepath.Join("/base/dir", cfg.Path)
	// ruleid: path-traversal-unmarshaled-file
	_, _ = os.ReadFile(target)
}

// 3. Vulnerable: unmarshaled JSON path formatted and deleted
func testVulnerableCleanPropagator(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	cleaned := filepath.Clean(cfg.Path)
	target := fmt.Sprintf("/var/data/%s", cleaned)
	// ruleid: path-traversal-unmarshaled-file
	_ = os.RemoveAll(target)
}

// 4. Safe: sanitized via filepath.Base
func testSafeFilepathBase(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	safe := filepath.Base(cfg.Path)
	// ok: path-traversal-unmarshaled-file
	_, _ = os.Open(filepath.Join("/base/dir", safe))
}

// 5. Safe: validated via filepath.IsLocal
func testSafeFilepathIsLocal(data []byte) {
	var cfg Config
	_ = yaml.Unmarshal(data, &cfg)
	if !filepath.IsLocal(cfg.Path) {
		return
	}
	// ok: path-traversal-unmarshaled-file
	_, _ = os.ReadFile(filepath.Join("/base/dir", cfg.Path))
}

// 6. Safe: validated via fs.ValidPath
func testSafeFSValidPath(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	if !fs.ValidPath(cfg.Path) {
		return
	}
	// ok: path-traversal-unmarshaled-file
	_, _ = os.Create(filepath.Join("/base/dir", cfg.Path))
}
