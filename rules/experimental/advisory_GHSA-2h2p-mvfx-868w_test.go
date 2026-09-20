package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return ""
}

func (c *GinContext) File(filepath string) {}

type EchoContext struct{}

func (c EchoContext) FormValue(name string) string {
	return ""
}

func (c EchoContext) File(file string) error {
	return nil
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string {
	return ""
}

func (c *FiberCtx) SendFile(file string) error {
	return nil
}

// Case 1: Vulnerable - net/http query param used directly in os.Open
func HandlerNetHttpQueryVuln(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	targetPath := filepath.Join("/var/data/exports", filename)
	// ruleid: path-traversal-arbitrary-file-read
	f, err := os.Open(targetPath)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()
}

// Case 2: Safe - net/http query param sanitized using filepath.Base
func HandlerNetHttpBaseSafe(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	targetPath := filepath.Join("/var/data/exports", safeName)
	// ok: path-traversal-arbitrary-file-read
	f, err := os.Open(targetPath)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()
}

// Case 3: Vulnerable - Echo FormValue concatenated and passed to os.ReadFile
func HandlerEchoFormReadFileVuln(c EchoContext) {
	userPath := c.FormValue("path")
	targetPath := filepath.Join("/var/data/exports", userPath)
	// ruleid: path-traversal-arbitrary-file-read
	_, _ = os.ReadFile(targetPath)
}

// Case 4: Safe - net/http query param guarded with filepath.IsLocal before http.ServeFile
func HandlerNetHttpIsLocalServeFileSafe(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	targetPath := filepath.Join("/var/data/exports", filename)
	// ok: path-traversal-arbitrary-file-read
	http.ServeFile(w, r, targetPath)
}

// Case 5: Vulnerable - Gin Query param passed directly to c.File sink
func HandlerGinQueryFileVuln(c *GinContext) {
	target := c.Query("filepath")
	fullPath := filepath.Join("/var/data/exports", target)
	// ruleid: path-traversal-arbitrary-file-read
	c.File(fullPath)
}

// Case 6: Safe - Fiber Query param sanitized with filepath.Base before c.SendFile
func HandlerFiberQueryBaseSendFileSafe(c *FiberCtx) {
	doc := c.Query("doc")
	safeDoc := filepath.Base(doc)
	fullPath := filepath.Join("/var/data/exports", safeDoc)
	// ok: path-traversal-arbitrary-file-read
	_ = c.SendFile(fullPath)
}
