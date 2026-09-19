package rules

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

// Direct standard lib vulnerability
func DirectStandardLibVuln(req *http.Request) {
	filename := req.URL.Query().Get("file")
	// ruleid: path-traversal-filesystem
	os.Open(filename)
}

// Proper standard lib patch
func ProperStandardLibPatch(req *http.Request) {
	filename := req.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		return
	}
	// ok: path-traversal-filesystem
	os.Open(filename)
}

func resolveUserPath(baseDir, untrustedInput string) string {
	return filepath.Join(baseDir, untrustedInput)
}

// Cross-function taint (wrapper function bypass)
func CrossFunctionWrapperBypass(req *http.Request) {
	filename := req.URL.Query().Get("file")
	targetPath := resolveUserPath("/tmp/storage", filename)
	// ruleid: path-traversal-filesystem
	os.RemoveAll(targetPath)
}

type PathTransformer interface {
	Transform(input string) string
}

type IdentityPathTransformer struct{}

func (t *IdentityPathTransformer) Transform(input string) string {
	return input
}

// Interface abstraction bypass
func InterfaceAbstractionBypass(req *http.Request, transformer PathTransformer) {
	filename := req.URL.Query().Get("file")
	resolved := transformer.Transform(filename)
	// ruleid: path-traversal-filesystem
	os.OpenFile(resolved, os.O_RDONLY, 0)
}

// Fake sanitizer usage (must trigger alert)
func FakeSanitizerVuln(req *http.Request) {
	filename := req.URL.Query().Get("file")
	cleaned := filepath.Clean(filename)
	joined := filepath.Join("/var/www/uploads", cleaned)
	// ruleid: path-traversal-filesystem
	os.ReadFile(joined)
}

// Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(req *http.Request) {
	filename := req.URL.Query().Get("file")
	if !fs.ValidPath(filename) {
		return
	}
	// ok: path-traversal-filesystem
	os.ReadFile(filename)
}
