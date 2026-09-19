package rules

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

// 1. Direct standard lib vulnerability
func DirectStdlibVuln(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("path")
	// ruleid: path-traversal-filesystem
	os.RemoveAll(name)
}

// 2. Proper standard lib patch
func ProperStdlibPatch(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("path")
	if !filepath.IsLocal(name) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	// ok: path-traversal-filesystem
	os.RemoveAll(name)
}

// 3. Cross-function taint (wrapper function bypass)
func buildWorkspacePath(base, rel string) string {
	return filepath.Join(base, rel)
}

func CrossFunctionWrapperBypass(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("workspace")
	target := buildWorkspacePath("/var/workspaces", name)
	// ruleid: path-traversal-filesystem
	os.RemoveAll(target)
}

// 4. Interface abstraction bypass
type WorkspaceResolver interface {
	Resolve(input string) string
}

type DefaultWorkspaceResolver struct {
	baseDir string
}

func (d *DefaultWorkspaceResolver) Resolve(input string) string {
	return filepath.Join(d.baseDir, input)
}

func InterfaceAbstractionBypass(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("dir")
	var resolver WorkspaceResolver = &DefaultWorkspaceResolver{baseDir: "/data/repos"}
	resolved := resolver.Resolve(name)
	// ruleid: path-traversal-filesystem
	os.RemoveAll(resolved)
}

// 5. Fake sanitizer usage (must trigger alert)
func FakeSanitizerUsage(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	// Developer erroneously assumes filepath.Clean and filepath.Join prevent path traversal
	cleaned := filepath.Clean(name)
	joined := filepath.Join("/safe/root", cleaned)
	// ruleid: path-traversal-filesystem
	os.RemoveAll(joined)
}

// 6. Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	if !fs.ValidPath(name) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	joined := filepath.Join("/safe/root", name)
	// ok: path-traversal-filesystem
	os.RemoveAll(joined)
}
