package rules

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type Config struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

// 1. Direct stdlib
func TestDirectStdlib(t *testing.T) {
	data := []byte(`{"path": "../evil"}`)
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join("/storage/root", cfg.Path)
	// ruleid: path-traversal-unmarshaled-data
	_ = os.Mkdir(target, 0o755)
}

// 2. Proper patch
func TestProperPatch(t *testing.T) {
	data := []byte(`{"path": "safe/sub"}`)
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if !filepath.IsLocal(cfg.Path) {
		t.Fatal(errors.New("invalid path"))
	}
	target := filepath.Join("/storage/root", cfg.Path)
	// ok: path-traversal-unmarshaled-data
	_ = os.Mkdir(target, 0o755)
}

func resolvePath(root, sub string) string {
	return filepath.Join(root, sub)
}

// 3. Cross-function taint
func TestCrossFunctionTaint(t *testing.T) {
	data := []byte(`{"path": "../../escape"}`)
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	target := resolvePath("/storage/root", cfg.Path)
	// ruleid: path-traversal-unmarshaled-data
	_ = os.WriteFile(target, []byte("data"), 0o644)
}

// 4. Interface bypass
func TestInterfaceBypass(t *testing.T) {
	data := []byte(`{"path": "../../../etc/passwd"}`)
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatal(err)
	}
	val, ok := payload["path"].(string)
	if !ok {
		return
	}
	target := filepath.Join("/storage/root", val)
	// ruleid: path-traversal-unmarshaled-data
	_ = os.RemoveAll(target)
}

// 5. Fake sanitizer
func TestFakeSanitizer(t *testing.T) {
	data := []byte(`{"path": "../../../etc/shadow"}`)
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	// filepath.Clean does NOT prevent path traversal escapes
	cleaned := filepath.Clean(cfg.Path)
	target := filepath.Join("/storage/root", cleaned)
	// ruleid: path-traversal-unmarshaled-data
	f, _ := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0o644)
	if f != nil {
		_ = f.Close()
	}
}

// 6. Real sanitizer
func TestRealSanitizer(t *testing.T) {
	data := []byte(`{"path": "../../traversal"}`)
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	base := filepath.Base(cfg.Path)
	target := filepath.Join("/storage/root", base)
	// ok: path-traversal-unmarshaled-data
	f, _ := os.Create(target)
	if f != nil {
		_ = f.Close()
	}
}
