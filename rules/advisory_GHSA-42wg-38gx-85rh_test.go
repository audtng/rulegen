package rules

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FilePathResolver interface {
	ResolvePath() string
}

type RequestPathResolver struct {
	path string
}

func (r RequestPathResolver) ResolvePath() string {
	return r.path
}

func pathWrapper(p string) string {
	return filepath.Clean(p)
}

// 1. Direct standard lib vulnerability
func directVuln(req *http.Request) {
	filename := req.URL.Query().Get("file")
	// ruleid: go-path-traversal
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
}

// 2. Proper standard lib patch
func properPatch(req *http.Request) {
	filename := req.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		return
	}
	// ok: go-path-traversal
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
}

// 3. Cross-function taint (wrapper function bypass)
func crossFunctionTaint(req *http.Request) {
	filename := req.URL.Query().Get("file")
	target := pathWrapper(filename)
	// ruleid: go-path-traversal
	f, err := os.Open(target)
	if err != nil {
		return
	}
	defer f.Close()
}

// 4. Interface abstraction bypass
func interfaceAbstractionBypass(req *http.Request) {
	var resolver FilePathResolver = RequestPathResolver{
		path: req.URL.Query().Get("file"),
	}
	resolved := resolver.ResolvePath()
	// ruleid: go-path-traversal
	f, err := os.Open(resolved)
	if err != nil {
		return
	}
	defer f.Close()
}

// 5. Fake sanitizer usage (must trigger alert)
func fakeSanitizerUsage(req *http.Request) {
	baseDir := "/var/www/uploads"
	filename := req.URL.Query().Get("file")
	untrustedPath := filepath.Join(baseDir, filename)
	// ruleid: go-path-traversal
	f, err := os.Open(untrustedPath)
	if err != nil {
		return
	}
	defer f.Close()
}

// 6. Real sanitizer usage (must not trigger alert)
func realSanitizerUsage(req *http.Request) {
	filename := req.URL.Query().Get("file")
	if strings.ContainsAny(filename, "/\\") || strings.Contains(filename, "..") {
		return
	}
	// ok: go-path-traversal
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
}
