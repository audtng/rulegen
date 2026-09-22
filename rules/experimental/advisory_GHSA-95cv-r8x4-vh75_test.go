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

type EchoContext interface {
	FormValue(name string) string
}

func testVulnerableNetHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	fullPath := filepath.Join("/var/data", filePath)
	// ruleid: path-traversal
	f, err := os.OpenFile(fullPath, os.O_RDONLY, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
}

func testSafeNetHTTPBase(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	safeName := filepath.Base(filePath)
	fullPath := filepath.Join("/var/data", safeName)
	// ok: path-traversal
	f, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
}

func testSafeNetHTTPRelCheck(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	fullPath := filepath.Join("/var/data", filePath)
	rel, err := filepath.Rel("/var/data", fullPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	// ok: path-traversal
	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func testVulnerableGin(c *GinContext) {
	target := filepath.Join("/var/www/static", c.Query("file"))
	// ruleid: path-traversal
	c.File(target)
}

func testSafeGinIsLocal(c *GinContext) {
	file := c.Query("file")
	if !filepath.IsLocal(file) {
		return
	}
	target := filepath.Join("/var/www/static", file)
	// ok: path-traversal
	c.File(target)
}

func testVulnerableEchoRemoveAll(c EchoContext) {
	folder := c.FormValue("dir")
	target := filepath.Join("/var/uploads", folder)
	// ruleid: path-traversal
	_ = os.RemoveAll(target)
}
