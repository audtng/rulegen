package rules

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

// Mock framework context types to support compiling without external dependencies
type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return "mock"
}

func (c *GinContext) Param(key string) string {
	return "mock"
}

func (c *GinContext) PostForm(key string) string {
	return "mock"
}

// 1. Direct standard lib vulnerability
func TestDirectVuln(r *http.Request) {
	p := r.URL.Query().Get("path")
	// ruleid: go-path-traversal
	_, _ = os.Open(p)
}

// 2. Proper standard lib patch
func TestProperPatch(r *http.Request) {
	p := r.URL.Query().Get("path")
	if !fs.ValidPath(p) {
		return
	}
	// ok: go-path-traversal
	_, _ = os.Open(p)
}

// Helper wrapper for cross-function taint
func wrapPath(val string) string {
	return filepath.Clean(val)
}

// 3. Cross-function taint (wrapper function bypass)
func TestCrossFunctionVuln(r *http.Request) {
	p := r.URL.Query().Get("path")
	wrapped := wrapPath(p)
	// ruleid: go-path-traversal
	_, _ = os.Open(wrapped)
}

// Interface abstraction for path resolution
type PathResolver interface {
	Resolve(p string) string
}

type InsecureResolver struct{}

func (InsecureResolver) Resolve(p string) string {
	return filepath.Clean(p)
}

// 4. Interface abstraction bypass
func TestInterfaceAbstractionVuln(r *http.Request) {
	var resolver PathResolver = InsecureResolver{}
	p := r.URL.Query().Get("path")
	resolved := resolver.Resolve(p)
	// ruleid: go-path-traversal
	_, _ = os.Open(resolved)
}

// 5. Fake sanitizer usage (must trigger alert)
func TestFakeSanitizerVuln(r *http.Request) {
	p := r.URL.Query().Get("path")
	joined := filepath.Join("/base/directory", p)
	// ruleid: go-path-traversal
	_, _ = os.Open(joined)
}

// 6. Real sanitizer usage (must not trigger alert)
func TestRealSanitizerSafe(r *http.Request) {
	p := r.URL.Query().Get("path")
	if !filepath.IsLocal(p) {
		return
	}
	// ok: go-path-traversal
	_, _ = os.Open(p)
}
