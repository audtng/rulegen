package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

// Mock framework types to ensure test compiles without external internet dependencies.
type GinContext struct{}

func (c *GinContext) Query(key string) string    { return "untrusted" }
func (c *GinContext) Param(key string) string    { return "untrusted" }
func (c *GinContext) PostForm(key string) string { return "untrusted" }

type EchoContext interface {
	FormValue(key string) string
	QueryParam(key string) string
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string  { return "untrusted" }
func (c *FiberCtx) Params(key string) string { return "untrusted" }

// 1. Direct standard lib vulnerability
func DirectStandardLibVuln(req *http.Request) {
	path := req.URL.Query().Get("file")
	// ruleid: go-path-traversal
	_, _ = os.Open(path)
}

// 2. Proper standard lib patch
func ProperStandardLibPatch(req *http.Request) {
	path := req.URL.Query().Get("file")
	if !filepath.IsLocal(path) {
		return
	}
	// ok: go-path-traversal
	_, _ = os.Open(path)
}

// 3. Cross-function taint (wrapper function bypass)
func openFileWrapper(target string) (*os.File, error) {
	// ruleid: go-path-traversal
	return os.Open(target)
}

func CrossFunctionBypass(req *http.Request) {
	path := req.URL.Query().Get("file")
	_, _ = openFileWrapper(path)
}

// 4. Interface abstraction bypass
type FileOpener interface {
	Open(name string) (*os.File, error)
}

type LocalFileOpener struct{}

func (l *LocalFileOpener) Open(name string) (*os.File, error) {
	// ruleid: go-path-traversal
	return os.Open(name)
}

func InterfaceAbstractionBypass(req *http.Request, opener FileOpener) {
	path := req.URL.Query().Get("file")
	_, _ = opener.Open(path)
}

// 5. Fake sanitizer usage (must trigger alert)
func FakeSanitizerUsage(req *http.Request) {
	relPath := req.URL.Query().Get("file")
	// filepath.Join cleans path lexically but does not contain traversal out of the root directory
	joinedPath := filepath.Join("/var/www/uploads", relPath)
	// ruleid: go-path-traversal
	_, _ = os.Open(joinedPath)
}

// 6. Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(req *http.Request) {
	relPath := req.URL.Query().Get("file")
	if filepath.IsLocal(relPath) {
		// ok: go-path-traversal
		_, _ = os.ReadFile(relPath)
	}
}
