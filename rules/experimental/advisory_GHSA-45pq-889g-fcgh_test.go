package rules

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string    { return "untrusted" }
func (c *GinContext) Param(key string) string    { return "untrusted" }
func (c *GinContext) PostForm(key string) string { return "untrusted" }
func (c *GinContext) File(filepath string)       {}

type EchoContext interface {
	FormValue(key string) string
	QueryParam(key string) string
	File(file string) error
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string    { return "untrusted" }
func (c *FiberCtx) Params(key string) string   { return "untrusted" }
func (c *FiberCtx) SendFile(file string) error { return nil }

// 1. Direct standard lib vulnerability
func DirectStdlibVuln(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	target := filepath.Join("/var/data", filePath)
	// ruleid: go-path-traversal
	os.Open(target)
}

// 2. Proper standard lib patch
func ProperStdlibPatch(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if !filepath.IsLocal(filePath) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	target := filepath.Join("/var/data", filePath)
	// ok: go-path-traversal
	os.Open(target)
}

// 3. Cross-function taint (wrapper function bypass)
func buildFilePath(base, rel string) string {
	return filepath.Join(base, rel)
}

func CrossFunctionTaint(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	target := buildFilePath("/var/data", name)
	// ruleid: go-path-traversal
	os.Open(target)
}

// 4. Interface abstraction bypass
func InterfaceBypass(w http.ResponseWriter, r *http.Request, fsys fs.FS) {
	filePath := r.URL.Query().Get("file")
	target := filepath.Join("/var/data", filePath)
	// ruleid: go-path-traversal
	fsys.Open(target)
}

// 5. Fake sanitizer usage (must trigger alert)
func FakeSanitizerVuln(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	// Fake sanitizer: filepath.Clean does not prevent traversal for relative paths
	cleaned := filepath.Clean(name)
	target := filepath.Join("/var/data", cleaned)
	// ruleid: go-path-traversal
	os.Open(target)
}

// 6. Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	if !fs.ValidPath(name) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	target := filepath.Join("/var/data", name)
	// ok: go-path-traversal
	os.Open(target)
}
