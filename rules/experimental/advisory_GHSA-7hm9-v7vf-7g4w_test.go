package rules

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Mock web framework contexts for offline compilation
type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string { return "" }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error                      { return nil }

func handleNetHTTPVulnerable(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/var/data", filename)
	// ruleid: go-path-traversal
	os.Open(target)
}

func handleNetHTTPSanitizedBase(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/data", safeName)
	// ok: go-path-traversal
	os.Open(target)
}

func handleGinContextVulnerable(c *GinContext) {
	doc := c.Query("doc")
	fullPath := filepath.Join("/app/storage", doc)
	// ruleid: go-path-traversal
	c.File(fullPath)
}

func handleGinContextSanitizedIsLocal(c *GinContext) {
	doc := c.Query("doc")
	if !filepath.IsLocal(doc) {
		return
	}
	fullPath := filepath.Join("/app/storage", doc)
	// ok: go-path-traversal
	c.File(fullPath)
}

func handleEchoVulnerableReadFile(c *EchoContext) {
	reportID := c.FormValue("id")
	target := fmt.Sprintf("/var/reports/%s.pdf", reportID)
	// ruleid: go-path-traversal
	os.ReadFile(target)
}

func handleFiberSanitizedHasPrefix(c *FiberCtx) {
	inputPath := c.Query("path")
	baseDir := "/var/www/uploads"
	cleaned := filepath.Clean(filepath.Join(baseDir, inputPath))
	if !strings.HasPrefix(cleaned, baseDir+"/") {
		return
	}
	// ok: go-path-traversal
	c.SendFile(cleaned)
}
