package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

// Mock third-party web framework structures for testing without external dependencies
type GinContext struct{}

func (c *GinContext) Query(key string) string    { return "untrusted" }
func (c *GinContext) Param(key string) string    { return "untrusted" }
func (c *GinContext) PostForm(key string) string { return "untrusted" }

type EchoContext interface {
	FormValue(key string) string
	QueryParam(key string) string
}

type FiberCtx struct{}

func (f *FiberCtx) Query(key string) string  { return "untrusted" }
func (f *FiberCtx) Params(key string) string { return "untrusted" }

// Direct standard lib vulnerability
func DirectStdlibVuln(req *http.Request) {
	filename := req.URL.Query().Get("filename")
	// ruleid: go-path-traversal
	os.Open(filename)
}

// Proper standard lib patch
func ProperStdlibPatch(req *http.Request) {
	filename := req.URL.Query().Get("filename")
	if filepath.IsLocal(filename) {
		// ok: go-path-traversal
		os.Open(filename)
	}
}

// Cross-function taint (wrapper function bypass)
func pathWrapper(input string) string {
	return input
}

func CrossFunctionTaint(req *http.Request) {
	filename := req.URL.Query().Get("filename")
	wrapped := pathWrapper(filename)
	// ruleid: go-path-traversal
	os.OpenFile(wrapped, os.O_RDONLY, 0)
}

// Interface abstraction bypass
type PathProvider interface {
	Resolve(p string) string
}

type DefaultPathProvider struct{}

func (d DefaultPathProvider) Resolve(p string) string {
	return p
}

func InterfaceAbstractionBypass(req *http.Request) {
	filename := req.URL.Query().Get("filename")
	var provider PathProvider = DefaultPathProvider{}
	resolved := provider.Resolve(filename)
	// ruleid: go-path-traversal
	os.ReadFile(resolved)
}

// Fake sanitizer usage (must trigger alert)
func FakeSanitizerUsage(req *http.Request) {
	filename := req.URL.Query().Get("filename")
	joinedPath := filepath.Join("/var/www/uploads", filename)
	// ruleid: go-path-traversal
	os.Open(joinedPath)
}

// Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(req *http.Request) {
	filename := req.URL.Query().Get("filename")
	if filepath.IsLocal(filename) {
		target := filepath.Join("/var/www/uploads", filename)
		// ok: go-path-traversal
		os.Open(target)
	}
}
