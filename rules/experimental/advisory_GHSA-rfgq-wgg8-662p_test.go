package rules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"testing"
)

type Glob interface {
	Match(string) bool
}

type globPkg struct{}

func (globPkg) Compile(pattern string, separators ...rune) (Glob, error) { return nil, nil }
func (globPkg) MustCompile(pattern string, separators ...rune) Glob      { return nil }

var glob globPkg

type yamlPkg struct{}

func (yamlPkg) Unmarshal(in []byte, out any) error { return nil }

var yaml yamlPkg

type GinContext struct{}

func (*GinContext) Query(key string) string { return "" }

type EchoContext struct{}

func (EchoContext) FormValue(key string) string { return "" }

type FiberCtx struct{}

func (*FiberCtx) Query(key string, args ...string) string { return "" }

type Config struct {
	ResourcePath string `json:"resource_path" yaml:"resource_path"`
}

// Edge Case 1: YAML configuration unmarshaling into resource path
func TestYamlConfigUnmarshaling(t *testing.T) {
	data := []byte("resource_path: /upload/*/open")
	var cfg Config
	_ = yaml.Unmarshal(data, &cfg)

	// ruleid: missing-path-separator-glob-compile
	_, _ = glob.Compile(cfg.ResourcePath)

	// ok: missing-path-separator-glob-compile
	_, _ = glob.Compile(cfg.ResourcePath, '/')
}

// Edge Case 2: JSON configuration unmarshaling into struct with MustCompile
func TestJsonConfigMustCompile(t *testing.T) {
	data := []byte(`{"resource_path": "/api/*/public"}`)
	var cfg Config
	_ = json.Unmarshal(data, &cfg)

	// ruleid: missing-path-separator-glob-compile
	_ = glob.MustCompile(cfg.ResourcePath)

	// ok: missing-path-separator-glob-compile
	_ = glob.MustCompile(cfg.ResourcePath, '/')
}

// Edge Case 3: net/http Request URL Path directly passed to glob compilation
func TestNetHttpUrlPath(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test/path", nil)
	pattern := req.URL.Path

	// ruleid: missing-path-separator-glob-compile
	_, _ = glob.Compile(pattern)

	// ok: missing-path-separator-glob-compile
	_, _ = glob.Compile(pattern, '/')
}

// Edge Case 4: Gin framework query parameter with path.Join propagation
func TestGinContextWithPathJoin(t *testing.T) {
	c := &GinContext{}
	raw := c.Query("resource")
	p := path.Join("/base", raw)

	// ruleid: missing-path-separator-glob-compile
	_, _ = glob.Compile(p)

	// ok: missing-path-separator-glob-compile
	_, _ = glob.Compile(p, '/')
}

// Edge Case 5: Echo framework FormValue parameter with fmt.Sprintf propagation
func TestEchoContextWithSprintf(t *testing.T) {
	c := EchoContext{}
	sub := c.FormValue("subpath")
	pattern := fmt.Sprintf("/api/v1/%s/*", sub)

	// ruleid: missing-path-separator-glob-compile
	_ = glob.MustCompile(pattern)

	// ok: missing-path-separator-glob-compile
	_ = glob.MustCompile(pattern, '/')
}

// Edge Case 6: Fiber framework query parameter with filepath.Base sanitizer and filepath.IsLocal guard
func TestFiberContextWithSanitizer(t *testing.T) {
	c := &FiberCtx{}
	raw := c.Query("pattern")

	// ruleid: missing-path-separator-glob-compile
	_, _ = glob.Compile(raw)

	// Sanitized via filepath.Base
	safeBase := filepath.Base(raw)
	// ok: missing-path-separator-glob-compile
	_, _ = glob.Compile(safeBase)

	// Sanitized via filepath.IsLocal check
	if filepath.IsLocal(raw) {
		// ok: missing-path-separator-glob-compile
		_, _ = glob.Compile(raw)
	}
}
