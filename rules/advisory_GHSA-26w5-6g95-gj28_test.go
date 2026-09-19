package rules

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Mock web framework types to ensure test compiles without external internet dependencies
type GinContext struct{}

func (c *GinContext) Query(key string) string    { return "" }
func (c *GinContext) Param(key string) string    { return "" }
func (c *GinContext) PostForm(key string) string { return "" }

type EchoContext interface {
	FormValue(name string) string
	QueryParam(name string) string
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string  { return "" }
func (c *FiberCtx) Params(key string, defaultValue ...string) string { return "" }

// 1. Direct standard lib vulnerability
func DirectStandardLibVuln(req *http.Request) error {
	path := req.URL.Query().Get("path")
	// ruleid: go-path-traversal
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return f.Close()
}

// 2. Proper standard lib patch
func ProperStandardLibPatch(req *http.Request) error {
	path := req.URL.Query().Get("path")
	if !filepath.IsLocal(path) {
		return errors.New("invalid path: path must be local")
	}
	// ok: go-path-traversal
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return f.Close()
}

// 3. Cross-function taint (wrapper function bypass)
func resolveUserPath(baseDir, userPath string) string {
	return fmt.Sprintf("%s/%s", baseDir, userPath)
}

func CrossFunctionTaintBypass(req *http.Request, baseDir string) error {
	path := req.URL.Query().Get("path")
	target := resolveUserPath(baseDir, path)
	// ruleid: go-path-traversal
	f, err := os.Open(target)
	if err != nil {
		return err
	}
	return f.Close()
}

// 4. Interface abstraction bypass
type PathTransformer interface {
	Transform(path string) string
}

type DefaultTransformer struct{}

func (DefaultTransformer) Transform(p string) string {
	return p
}

func InterfaceAbstractionBypass(req *http.Request, baseDir string, transformer PathTransformer) error {
	path := req.URL.Query().Get("path")
	transformed := transformer.Transform(path)
	target := filepath.Join(baseDir, transformed)
	// ruleid: go-path-traversal
	f, err := os.Open(target)
	if err != nil {
		return err
	}
	return f.Close()
}

// 5. Fake sanitizer usage (must trigger alert)
func FakeSanitizerUsage(req *http.Request, baseDir string) error {
	path := req.URL.Query().Get("path")
	target := filepath.Join(baseDir, path)
	// ruleid: go-path-traversal
	f, err := os.Open(target)
	if err != nil {
		return err
	}
	return f.Close()
}

// 6. Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(req *http.Request, baseDir string) error {
	path := req.URL.Query().Get("path")
	target := filepath.Join(baseDir, path)
	cleanBase := filepath.Clean(baseDir) + string(filepath.Separator)
	if !strings.HasPrefix(target, cleanBase) {
		return errors.New("path traversal attempt detected")
	}
	// ok: go-path-traversal
	f, err := os.Open(target)
	if err != nil {
		return err
	}
	return f.Close()
}

// UnusedMockHelpers ensures mock interfaces and fs import compile without warnings
func UnusedMockHelpers(c *GinContext, ec EchoContext, fc *FiberCtx, fsys fs.FS) {
	_ = c.Query("q")
	_ = ec.FormValue("f")
	_ = fc.Query("q")
	_, _ = fsys.Open("file")
}
