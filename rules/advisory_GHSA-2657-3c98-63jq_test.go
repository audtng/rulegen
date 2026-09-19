package rules

import (
	"archive/tar"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// 1. Direct standard lib vulnerability
func testDirectVuln(req *http.Request) {
	name := req.URL.Query().Get("file")
	target := filepath.Join("/safe/dir", name)
	// ruleid: go-path-traversal-file-write
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
}

// 2. Proper standard lib patch
func testProperPatch(req *http.Request) {
	name := req.URL.Query().Get("file")
	if !filepath.IsLocal(name) {
		return
	}
	target := filepath.Join("/safe/dir", name)
	// ok: go-path-traversal-file-write
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
}

// 3. Cross-function taint (wrapper function bypass)
func buildTarget(base, name string) string {
	return filepath.Join(base, name)
}

func testCrossFunctionTaint(req *http.Request) {
	name := req.URL.Query().Get("file")
	target := buildTarget("/safe/dir", name)
	// ruleid: go-path-traversal-file-write
	os.Create(target)
}

// 4. Interface abstraction bypass
type PathResolver interface {
	Resolve(base, input string) string
}

type DefaultResolver struct{}

func (DefaultResolver) Resolve(base, input string) string {
	return filepath.Join(base, input)
}

func testInterfaceAbstractionBypass(req *http.Request) {
	name := req.URL.Query().Get("file")
	var resolver PathResolver = DefaultResolver{}
	target := resolver.Resolve("/safe/dir", name)
	// ruleid: go-path-traversal-file-write
	os.WriteFile(target, []byte("data"), 0644)
}

// 5. Fake sanitizer usage (must trigger alert)
func testFakeSanitizer(h *tar.Header) {
	name := h.Name
	cleaned := path.Clean(name)
	target := path.Join("/safe/dir", cleaned)
	// ruleid: go-path-traversal-file-write
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
}

// 6. Real sanitizer usage (must not trigger alert)
func testRealSanitizer(req *http.Request) {
	name := req.URL.Query().Get("file")
	target := filepath.Join("/safe/dir", name)
	if !strings.HasPrefix(target, "/safe/dir/") {
		return
	}
	// ok: go-path-traversal-file-write
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
}
