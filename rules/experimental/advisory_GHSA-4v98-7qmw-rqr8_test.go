package rules

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var yaml = struct {
	Unmarshal func([]byte, any) error
}{}

type Config struct {
	MountPath string   `json:"mount_path" xml:"mount_path"`
	Stubs     []string `json:"stubs"`
}

func testVulnerableJSONJoinRemove(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	target := filepath.Join(baseDir, cfg.MountPath)
	// ruleid: untrusted-config-path-traversal
	return os.RemoveAll(target)
}

func testSafeBaseSanitizer(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	safeName := filepath.Base(cfg.MountPath)
	target := filepath.Join(baseDir, safeName)
	// ok: untrusted-config-path-traversal
	return os.RemoveAll(target)
}

func testSafeIsLocalValidation(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	if !filepath.IsLocal(cfg.MountPath) {
		return fmt.Errorf("invalid path: %s", cfg.MountPath)
	}
	target := filepath.Join(baseDir, cfg.MountPath)
	// ok: untrusted-config-path-traversal
	return os.RemoveAll(target)
}

func testVulnerableMountStubsCleaner(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	for _, stub := range cfg.Stubs {
		p := filepath.Join(baseDir, stub)
		parent := filepath.Dir(p)
		now := time.Now()
		// ruleid: untrusted-config-path-traversal
		if err := os.Chtimes(parent, now, now); err != nil {
			return err
		}
	}
	return nil
}

func testVulnerableYAMLOpenFile(data []byte, baseDir string) error {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	cleaned := filepath.Clean(filepath.Join(baseDir, cfg.MountPath))
	// ruleid: untrusted-config-path-traversal
	f, err := os.OpenFile(cleaned, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	return f.Close()
}

func testSafeXMLWithIsLocal(data []byte, baseDir string) error {
	var cfg Config
	if err := xml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	if !filepath.IsLocal(cfg.MountPath) {
		return fmt.Errorf("insecure path: %s", cfg.MountPath)
	}
	target := filepath.Join(baseDir, cfg.MountPath)
	// ok: untrusted-config-path-traversal
	return os.WriteFile(target, []byte("safe"), 0600)
}
