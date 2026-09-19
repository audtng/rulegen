package rules

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// EchoContext mocks the Echo framework Context interface for testing interface-based taint sources and sinks.
type EchoContext interface {
	FormValue(name string) string
	File(file string) error
}

// Direct stdlib: HTTP request query parameter directly used in file operations.
func DirectStdlibTest(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/data/uploads", filename)
	// ruleid: go-path-traversal
	os.Open(target)
}

// Proper patch: validates that the resolved path does not escape the base directory.
func ProperPatchTest(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/data/uploads", filename)
	if !strings.HasPrefix(target, "/data/uploads/") {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	// ok: go-path-traversal
	os.Open(target)
}

func resolveUserPath(filename string) string {
	return filepath.Join("/var/www/docs", filename)
}

// Cross-function taint: tainted input passed through a helper function.
func CrossFunctionTaintTest(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("doc")
	target := resolveUserPath(filename)
	// ruleid: go-path-traversal
	os.Open(target)
}

// Interface bypass: untrusted input passed through an interface abstraction and sunk via framework interface.
func InterfaceBypassTest(c EchoContext) {
	filename := c.FormValue("file")
	var iface interface{} = filename
	target := filepath.Join("/var/app/static", iface.(string))
	// ruleid: go-path-traversal
	c.File(target)
}

// Fake sanitizer: filepath.Clean does NOT prevent path traversal when joined with a base path.
func FakeSanitizerTest(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	cleanName := filepath.Clean(filename)
	target := filepath.Join("/data/uploads", cleanName)
	// ruleid: go-path-traversal
	os.Open(target)
}

// Real sanitizer: filepath.Base removes all directory traversal components.
func RealSanitizerTest(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	baseName := filepath.Base(filename)
	target := filepath.Join("/data/uploads", baseName)
	// ok: go-path-traversal
	os.Open(target)
}
