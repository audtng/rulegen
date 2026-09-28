package rules

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
)

type PackageManifest struct {
	Name    string `json:"name"`
	Browser string `json:"browser"`
}

type Config struct {
	FilePath string `json:"filePath"`
}

// 1. Vulnerable: json.Unmarshal of manifest struct browser field to os.ReadFile
func testVulnerableManifestUnmarshal(data []byte) {
	var pkg PackageManifest
	if err := json.Unmarshal(data, &pkg); err != nil {
		return
	}
	target := filepath.Join("/var/app/packages", pkg.Browser)
	// ruleid: go-unmarshal-path-traversal
	os.ReadFile(target)
}

// 2. Vulnerable: json.Unmarshal into map[string]interface{} to os.Open
func testVulnerableMapUnmarshal(data []byte) {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}
	if path, ok := m["browser"].(string); ok {
		target := filepath.Join("/var/app/static", path)
		// ruleid: go-unmarshal-path-traversal
		os.Open(target)
	}
}

// 3. Vulnerable: Stream-based json.NewDecoder with filepath.Clean to os.OpenFile
func testVulnerableDecoderClean(data []byte) {
	var cfg Config
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&cfg); err != nil {
		return
	}
	target := filepath.Clean(filepath.Join("/var/data", cfg.FilePath))
	// ruleid: go-unmarshal-path-traversal
	os.OpenFile(target, os.O_RDONLY, 0)
}

// 4. Safe: Sanitized with filepath.Base
func testSafeBaseSanitizer(data []byte) {
	var pkg PackageManifest
	if err := json.Unmarshal(data, &pkg); err != nil {
		return
	}
	safeName := filepath.Base(pkg.Browser)
	target := filepath.Join("/var/app/packages", safeName)
	// ok: go-unmarshal-path-traversal
	os.ReadFile(target)
}

// 5. Safe: Validated with filepath.IsLocal
func testSafeIsLocalValidation(data []byte) {
	var pkg PackageManifest
	if err := json.Unmarshal(data, &pkg); err != nil {
		return
	}
	if !filepath.IsLocal(pkg.Browser) {
		return
	}
	target := filepath.Join("/var/app/packages", pkg.Browser)
	// ok: go-unmarshal-path-traversal
	os.Open(target)
}

// 6. Safe: Hardcoded safe path, unmarshaled config value not used in file operation
func testSafeStaticPath(data []byte) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return
	}
	target := filepath.Join("/var/app/static", "default.json")
	// ok: go-unmarshal-path-traversal
	os.ReadFile(target)
}
