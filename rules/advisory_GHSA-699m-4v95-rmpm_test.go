package rules

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Direct stdlib
func DirectStdlib(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/data/uploads", filename)
	// ruleid: go-path-traversal
	os.Open(target)
}

// Proper patch
func ProperPatch(w http.ResponseWriter, r *http.Request) {
	baseDir := "/data/uploads"
	filename := r.URL.Query().Get("file")
	target := filepath.Clean(filepath.Join(baseDir, filename))
	if !strings.HasPrefix(target, baseDir+string(filepath.Separator)) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	// ok: go-path-traversal
	os.Open(target)
}

func normalizePath(name string) string {
	return filepath.Clean(name)
}

// Cross-function taint
func CrossFunctionTaint(w http.ResponseWriter, r *http.Request) {
	userInput := r.FormValue("filepath")
	processed := normalizePath(userInput)
	fullPath := filepath.Join("/safe/path", processed)
	// ruleid: go-path-traversal
	os.Open(fullPath)
}

type PathProvider interface {
	Path() string
}

type UserPath struct {
	raw string
}

func (u UserPath) Path() string {
	return u.raw
}

// Interface bypass
func InterfaceBypass(w http.ResponseWriter, r *http.Request) {
	var provider PathProvider = UserPath{raw: r.URL.Query().Get("file")}
	target := filepath.Join("/var/www/uploads", provider.Path())
	// ruleid: go-path-traversal
	os.Open(target)
}

// Fake sanitizer
func FakeSanitizer(w http.ResponseWriter, r *http.Request) {
	baseDir := "/data/uploads"
	filename := r.URL.Query().Get("file")
	target := filepath.Clean(filepath.Join(baseDir, filename))
	// Incomplete validation: missing path separator allows prefix bypass / sibling directory traversal
	if !strings.HasPrefix(target, baseDir) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	// ruleid: go-path-traversal
	os.Open(target)
}

// Real sanitizer
func RealSanitizer(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/data/uploads", safeName)
	// ok: go-path-traversal
	os.Open(target)
}
