package rules

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// Edge case 1: Untrusted URL Path joined and opened directly
func TestVulnDirectPath(w http.ResponseWriter, req *http.Request) {
	target := filepath.Join("/var/www/html", req.URL.Path)
	// ruleid: go-path-traversal
	f, err := os.Open(target)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	defer f.Close()
}

// Edge case 2: Query parameter formatted into system path and read
func TestVulnQueryParam(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	fullPath := fmt.Sprintf("/var/data/%s", filename)
	// ruleid: go-path-traversal
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return
	}
	w.Write(data)
}

// Edge case 3: FormValue joined using path.Join and passed to http.ServeFile
func TestVulnFormValueServeFile(w http.ResponseWriter, req *http.Request) {
	userDoc := req.FormValue("doc")
	docPath := path.Join("/var/docs", userDoc)
	// ruleid: go-path-traversal
	http.ServeFile(w, req, docPath)
}

// Edge case 4: Sanitized via filepath.Base extracting strictly the base filename
func TestSafeFilepathBase(w http.ResponseWriter, req *http.Request) {
	baseName := filepath.Base(req.URL.Path)
	safePath := filepath.Join("/var/www/html", baseName)
	// ok: go-path-traversal
	f, err := os.Open(safePath)
	if err != nil {
		return
	}
	defer f.Close()
}

// Edge case 5: Validated via filepath.IsLocal ensuring path remains local
func TestSafeFilepathIsLocal(w http.ResponseWriter, req *http.Request) {
	relPath := req.URL.Path
	if !filepath.IsLocal(relPath) {
		http.Error(w, "Forbidden path", http.StatusForbidden)
		return
	}
	safePath := filepath.Join("/var/www/html", relPath)
	// ok: go-path-traversal
	data, err := os.ReadFile(safePath)
	if err != nil {
		return
	}
	w.Write(data)
}

// Edge case 6: Safe fixed constant file path
func TestSafeConstantPath(w http.ResponseWriter, req *http.Request) {
	const safeConfigFile = "/etc/myapp/config.json"
	// ok: go-path-traversal
	f, err := os.Open(safeConfigFile)
	if err != nil {
		return
	}
	defer f.Close()
}
