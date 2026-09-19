package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

// Direct standard lib vulnerability
func DirectVuln(r *http.Request) {
	p := r.URL.Query().Get("file")
	// ruleid: go-path-traversal
	os.Open(p)
}

// Proper standard lib patch
func ProperPatch(r *http.Request) {
	p := r.URL.Query().Get("file")
	safe := filepath.Base(p)
	// ok: go-path-traversal
	os.Open(filepath.Join("/var/www/uploads", safe))
}

// Cross-function taint (wrapper function bypass)
func wrapPath(userPath string) string {
	return filepath.Join("/var/www/data", userPath)
}

func CrossFunctionTaint(r *http.Request) {
	p := r.URL.Query().Get("file")
	target := wrapPath(p)
	// ruleid: go-path-traversal
	os.Open(target)
}

// Interface abstraction bypass
type PathResolver interface {
	Resolve() string
}

type UserPath struct {
	path string
}

func (u UserPath) Resolve() string {
	return u.path
}

func InterfaceAbstractionBypass(r *http.Request) {
	p := r.URL.Query().Get("file")
	var resolver PathResolver = UserPath{path: p}
	// ruleid: go-path-traversal
	os.Open(resolver.Resolve())
}

// Fake sanitizer usage (must trigger alert)
func FakeSanitizerUsage(r *http.Request) {
	p := r.URL.Query().Get("file")
	// filepath.Join and filepath.Clean are fake sanitizers that do not prevent traversal
	fakeSanitized := filepath.Join("/var/www/uploads", filepath.Clean(p))
	// ruleid: go-path-traversal
	os.Open(fakeSanitized)
}

// Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(r *http.Request) {
	p := r.URL.Query().Get("file")
	if !filepath.IsLocal(p) {
		return
	}
	// ok: go-path-traversal
	os.Open(filepath.Join("/var/www/uploads", p))
}
