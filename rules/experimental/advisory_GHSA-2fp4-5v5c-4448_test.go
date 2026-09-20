package rules

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c *EchoContext) FormValue(key string) string { return "" }
func (c *EchoContext) File(filepath string) error  { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string    { return "" }
func (c *FiberCtx) SendFile(path string) error { return nil }

// Edge Case 1: Vulnerable standard net/http handler joining unvalidated query parameter to base path
func VulnNetHttpJoin(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	absPath := filepath.Join("/var/app/playlists", id)
	// ruleid: go-path-traversal
	f, err := os.Open(absPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
}

// Edge Case 2: Safe standard net/http handler using filepath.Base to strip directory traversal segments
func SafeNetHttpBase(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	cleanName := filepath.Base(filename)
	absPath := filepath.Join("/var/app/uploads", cleanName)
	// ok: go-path-traversal
	data, err := os.ReadFile(absPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Edge Case 3: Vulnerable handler with flawed sanitization using filepath.Clean which does not prevent escaping base path
func VulnCleanBypass(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.Query().Get("path")
	cleanPath := filepath.Clean(rawPath)
	target := filepath.Join("/var/app/data", cleanPath)
	// ruleid: go-path-traversal
	if err := os.Remove(target); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Edge Case 4: Safe handler using filepath.IsLocal to verify relative path does not escape base directory
func SafeIsLocalCheck(w http.ResponseWriter, r *http.Request) {
	relPath := r.FormValue("path")
	if !filepath.IsLocal(relPath) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	target := filepath.Join("/var/app/storage", relPath)
	// ok: go-path-traversal
	f, err := os.OpenFile(target, os.O_RDWR, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
}

// Edge Case 5: Vulnerable Gin framework handler serving file from user query parameter
func VulnGinFileServing(c *GinContext) {
	file := c.Query("file")
	path := filepath.Join("/var/www/static", file)
	// ruleid: go-path-traversal
	c.File(path)
}

// Edge Case 6: Safe Echo framework handler verifying resolved path prefix matches base directory
func SafeEchoPrefixCheck(c *EchoContext) {
	file := c.FormValue("file")
	absPath := filepath.Join("/var/www/static", file)
	if !strings.HasPrefix(absPath, "/var/www/static/") {
		return
	}
	// ok: go-path-traversal
	c.File(absPath)
}
