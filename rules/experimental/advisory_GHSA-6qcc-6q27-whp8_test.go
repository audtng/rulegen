package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

// Mock third-party web framework contexts for offline compilation
type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return ""
}

func (c *GinContext) File(filepath string) {}

type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string {
	return ""
}

func (c *EchoContext) File(file string) error {
	return nil
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string {
	return ""
}

func (c *FiberCtx) SendFile(file string) error {
	return nil
}

// 1. net/http handler with URL path traversal into os.Open
func testVulnURLPath(w http.ResponseWriter, r *http.Request) {
	userPath := r.URL.Path
	fullPath := filepath.Join("/var/www/static", userPath)
	// ruleid: go-path-traversal
	f, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()
}

// 2. net/http handler with filepath.Base sanitization
func testSafeBase(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)
	safePath := filepath.Join("/var/www/static", filename)
	// ok: go-path-traversal
	f, err := os.Open(safePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()
}

// 3. net/http handler with FormValue into os.RemoveAll (from goshs CVE archetype)
func testVulnFormDelete(w http.ResponseWriter, r *http.Request) {
	target := r.FormValue("target")
	deletePath := filepath.Join("/tmp/storage", target)
	// ruleid: go-path-traversal
	os.RemoveAll(deletePath)
}

// 4. net/http handler with filepath.IsLocal validation
func testSafeIsLocal(w http.ResponseWriter, r *http.Request) {
	userPath := r.URL.Path
	if !filepath.IsLocal(userPath) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	safePath := filepath.Join("/var/www/static", userPath)
	// ok: go-path-traversal
	os.ReadFile(safePath)
}

// 5. Framework handler (Gin) with vulnerable Query into File
func testVulnGin(c *GinContext) {
	filePath := c.Query("file")
	target := filepath.Join("/uploads", filePath)
	// ruleid: go-path-traversal
	c.File(target)
}

// 6. Framework handler (Echo) with sanitized Base before File
func testSafeEcho(c *EchoContext) {
	fileParam := c.FormValue("download")
	safeName := filepath.Base(fileParam)
	target := filepath.Join("/downloads", safeName)
	// ok: go-path-traversal
	c.File(target)
}
