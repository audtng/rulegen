package rules

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	Path string `json:"path"`
}

// 1. Direct stdlib: untrusted unmarshaled data flows directly into a stdlib filesystem sink
func DirectStdlib(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	target := filepath.Join("/tmp/base", cfg.Path)
	// ruleid: path-traversal-unmarshaled-data
	return os.RemoveAll(target)
}

// 2. Proper patch: validated using filepath.IsLocal to verify path remains within directory bounds
func ProperPatch(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	if !filepath.IsLocal(cfg.Path) {
		return errors.New("path escapes root directory")
	}
	target := filepath.Join("/tmp/base", cfg.Path)
	// ok: path-traversal-unmarshaled-data
	return os.RemoveAll(target)
}

// 3. Cross-function taint: taint propagates across a helper function into a filesystem sink
func extractPath(cfg Config) string {
	return cfg.Path
}

func CrossFunctionTaint(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	extracted := extractPath(cfg)
	target := filepath.Join("/tmp/base", extracted)
	// ruleid: path-traversal-unmarshaled-data
	return os.RemoveAll(target)
}

// 4. Interface bypass: tainted value passed through an interface{} / any type assertion
func InterfaceBypass(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	var val any = cfg.Path
	target := filepath.Join("/tmp/base", val.(string))
	// ruleid: path-traversal-unmarshaled-data
	return os.RemoveAll(target)
}

// 5. Fake sanitizer: filepath.Clean does not prevent traversal sequences escaping base directory
func FakeSanitizer(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	cleaned := filepath.Clean(cfg.Path)
	target := filepath.Join("/tmp/base", cleaned)
	// ruleid: path-traversal-unmarshaled-data
	return os.RemoveAll(target)
}

// 6. Real sanitizer: filepath.Base safely strips any directory traversal segments
func RealSanitizer(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	base := filepath.Base(cfg.Path)
	target := filepath.Join("/tmp/base", base)
	// ok: path-traversal-unmarshaled-data
	return os.RemoveAll(target)
}
