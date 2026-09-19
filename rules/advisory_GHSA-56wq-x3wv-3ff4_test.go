package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

// Direct stdlib: tainted HTTP request input flows directly to os.Open
func testDirectStdlib(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Query().Get("path")
	full := filepath.Join("/var/www/uploads", p)
	// ruleid: go-path-traversal
	os.Open(full)
}

// Proper patch: input is validated with filepath.IsLocal before use
func testProperPatch(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Query().Get("path")
	if !filepath.IsLocal(p) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	full := filepath.Join("/var/www/uploads", p)
	// ok: go-path-traversal
	os.Open(full)
}

func loadUserFile(p string) string {
	return filepath.Join("/var/www/uploads", p)
}

// Cross-function taint: taint flows across function call boundary
func testCrossFunction(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	full := loadUserFile(p)
	// ruleid: go-path-traversal
	os.Open(full)
}

// Interface bypass: boxing tainted string into an interface and asserting back
func testInterfaceBypass(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	var i any = p
	pathStr := i.(string)
	full := filepath.Join("/var/www/uploads", pathStr)
	// ruleid: go-path-traversal
	os.Open(full)
}

// Fake sanitizer: filepath.Clean alone is insufficient to prevent traversal
func testFakeSanitizer(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	cleaned := filepath.Clean(p)
	full := filepath.Join("/var/www/uploads", cleaned)
	// ruleid: go-path-traversal
	os.Open(full)
}

// Real sanitizer: filepath.Base safely strips any directory traversal prefix
func testRealSanitizer(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	base := filepath.Base(p)
	full := filepath.Join("/var/www/uploads", base)
	// ok: go-path-traversal
	os.Open(full)
}
