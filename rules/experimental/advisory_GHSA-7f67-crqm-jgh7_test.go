package rules

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

var yaml = struct {
	Unmarshal func([]byte, any) error
}{
	Unmarshal: func(data []byte, v any) error { return nil },
}

type Config struct {
	Template string `json:"template" xml:"template"`
	Dest     string `json:"dest" xml:"dest"`
}

// 1. Vulnerable: unmarshaled JSON config path joined and read
func testVulnerableJSONReadFile(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	target := filepath.Join("/var/app/templates", cfg.Template)
	// ruleid: go-unmarshaled-path-traversal
	_, _ = os.ReadFile(target)
}

// 2. Vulnerable: unmarshaled YAML path manipulated with TrimLeft, joined, and created
func testVulnerableYAMLCreateFile(data []byte) {
	var cfg Config
	_ = yaml.Unmarshal(data, &cfg)
	relPath := strings.TrimLeft(cfg.Dest, "/")
	target := filepath.Join("/var/app/rootfs", relPath)
	// ruleid: go-unmarshaled-path-traversal
	_, _ = os.Create(target)
}

// 3. Vulnerable: unmarshaled XML path joined and removed
func testVulnerableXMLRemoveAll(data []byte) {
	var cfg Config
	_ = xml.Unmarshal(data, &cfg)
	target := filepath.Join("/tmp/scratch", cfg.Dest)
	// ruleid: go-unmarshaled-path-traversal
	_ = os.RemoveAll(target)
}

// 4. Safe: unmarshaled path sanitized using filepath.Base
func testSafeSanitizedWithBase(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	safeName := filepath.Base(cfg.Template)
	target := filepath.Join("/var/app/templates", safeName)
	// ok: go-unmarshaled-path-traversal
	_, _ = os.ReadFile(target)
}

// 5. Safe: unmarshaled path validated with filepath.IsLocal
func testSafeValidatedWithIsLocal(data []byte) {
	var cfg Config
	_ = yaml.Unmarshal(data, &cfg)
	if !filepath.IsLocal(cfg.Dest) {
		return
	}
	target := filepath.Join("/var/app/rootfs", cfg.Dest)
	// ok: go-unmarshaled-path-traversal
	_, _ = os.Create(target)
}

// 6. Safe: unmarshaled data used, but target path is hardcoded constant
func testSafeConstantPath(data []byte) {
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	target := filepath.Join("/var/app/templates", "default.tmpl")
	// ok: go-unmarshaled-path-traversal
	_, _ = os.ReadFile(target)
}
