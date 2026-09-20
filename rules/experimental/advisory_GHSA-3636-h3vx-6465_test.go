package rules

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return ""
}

func (c *GinContext) File(filepath string) {}

// Case 1: Vulnerable - HTTP URL path concatenated into file creation (Arbitrary File Write)
func handleUploadVuln(w http.ResponseWriter, req *http.Request) {
	userPath := req.URL.Path
	targetPath := filepath.Join("/var/www/uploads", userPath)
	// ruleid: go-http-path-traversal
	file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
}

// Case 2: Safe - HTTP URL path sanitized with filepath.Base
func handleUploadSafeBase(w http.ResponseWriter, req *http.Request) {
	userPath := req.URL.Path
	safeFilename := filepath.Base(userPath)
	targetPath := filepath.Join("/var/www/uploads", safeFilename)
	// ok: go-http-path-traversal
	file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
}

// Case 3: Vulnerable - filepath.Clean does NOT prevent directory traversal outside root
func handleDownloadCleanVuln(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	cleanedPath := filepath.Clean(filepath.Join("/var/data", filename))
	// ruleid: go-http-path-traversal
	data, err := os.ReadFile(cleanedPath)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	w.Write(data)
}

// Case 4: Safe - Path validated with filepath.IsLocal
func handleDownloadIsLocalSafe(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	targetPath := filepath.Join("/var/data", filename)
	// ok: go-http-path-traversal
	data, err := os.ReadFile(targetPath)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	w.Write(data)
}

// Case 5: Vulnerable - Gin web framework query parameter flows into file sink
func handleGinFileVuln(c *GinContext) {
	filePath := c.Query("path")
	fullPath := fmt.Sprintf("/static/%s", filePath)
	// ruleid: go-http-path-traversal
	c.File(fullPath)
}

// Case 6: Safe - Path validated with fs.ValidPath before access
func handleFSValidPathSafe(w http.ResponseWriter, req *http.Request) {
	filename := req.FormValue("name")
	if !fs.ValidPath(filename) {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}
	fullPath := filepath.Join("/safe/store", filename)
	// ok: go-http-path-traversal
	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	w.Write(data)
}
