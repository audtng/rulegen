package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c EchoContext) FormValue(name string) string { return "" }
func (c EchoContext) File(file string) error       { return nil }

// Edge Case 1: Direct net/http user input passed to os.Open
func testVulnerableNetHTTPDirect(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	// ruleid: go-path-traversal
	f, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
}

// Edge Case 2: Untrusted path joined via filepath.Join without boundary validation
func testVulnerableNetHTTPJoin(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("path")
	fullPath := filepath.Join("/var/www/uploads", target)
	// ruleid: go-path-traversal
	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Edge Case 3: Gin framework query parameter flowing directly to File sink
func testVulnerableGinFile(c *GinContext) {
	doc := c.Query("doc")
	// ruleid: go-path-traversal
	c.File(doc)
}

// Edge Case 4: Patched using filepath.Base to strip directory traversal sequences
func testSafeFilepathBase(w http.ResponseWriter, r *http.Request) {
	rawName := r.FormValue("name")
	safeName := filepath.Base(rawName)
	fullPath := filepath.Join("/var/www/uploads", safeName)
	// ok: go-path-traversal
	f, err := os.OpenFile(fullPath, os.O_RDONLY, 0600)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
}

// Edge Case 5: Patched using filepath.IsLocal guard
func testSafeFilepathIsLocal(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Path
	if filepath.IsLocal(target) {
		fullPath := filepath.Join("/var/www/public", target)
		// ok: go-path-traversal
		os.Open(fullPath)
	}
}

// Edge Case 6: Echo framework form value sanitized with filepath.Base
func testSafeEchoFile(c EchoContext) {
	file := c.FormValue("file")
	safeFile := filepath.Base(file)
	// ok: go-path-traversal
	c.File(safeFile)
}
