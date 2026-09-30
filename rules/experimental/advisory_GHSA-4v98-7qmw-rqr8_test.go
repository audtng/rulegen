package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	MountPath string `json:"mount_path"`
}

// 1. Direct stdlib
func testDirectStdlib(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	target := filepath.Join(baseDir, cfg.MountPath)
	// ruleid: path-traversal-untrusted-config
	return os.RemoveAll(target)
}

// 2. Proper patch
func testProperPatch(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	if !filepath.IsLocal(cfg.MountPath) {
		return nil
	}
	target := filepath.Join(baseDir, cfg.MountPath)
	// ok: path-traversal-untrusted-config
	return os.RemoveAll(target)
}

// 3. Cross-function taint
func getMountPath(cfg Config) string {
	return cfg.MountPath
}

func testCrossFunctionTaint(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	p := getMountPath(cfg)
	target := filepath.Join(baseDir, p)
	// ruleid: path-traversal-untrusted-config
	return os.RemoveAll(target)
}

// 4. Interface bypass
type PathGetter interface {
	GetPath() string
}

func (c Config) GetPath() string {
	return c.MountPath
}

func testInterfaceBypass(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	var getter PathGetter = cfg
	target := filepath.Join(baseDir, getter.GetPath())
	// ruleid: path-traversal-untrusted-config
	return os.RemoveAll(target)
}

// 5. Fake sanitizer
func testFakeSanitizer(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	// filepath.Clean does not prevent traversal outside baseDir for absolute or parent paths
	cleaned := filepath.Clean(cfg.MountPath)
	target := filepath.Join(baseDir, cleaned)
	// ruleid: path-traversal-untrusted-config
	return os.RemoveAll(target)
}

// 6. Real sanitizer
func testRealSanitizer(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	safeBase := filepath.Base(cfg.MountPath)
	target := filepath.Join(baseDir, safeBase)
	// ok: path-traversal-untrusted-config
	return os.RemoveAll(target)
}
