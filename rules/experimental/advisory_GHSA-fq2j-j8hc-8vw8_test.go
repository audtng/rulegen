package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

// Framework mocks for offline compilation
type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c EchoContext) FormValue(name string) string { return "" }
func (c EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error    { return nil }

// Edge case 1: Vulnerable net/http handler with filepath.Join
func netHttpVulnerable(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/var/data", filename)
	// ruleid: go-arbitrary-file-read
	f, err := os.Open(target)
	if err != nil {
		return
	}
	defer f.Close()
}

// Edge case 2: Safe net/http handler sanitized with filepath.Base
func netHttpSafeWithBase(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/data", safeName)
	// ok: go-arbitrary-file-read
	_, _ = os.ReadFile(target)
}

// Edge case 3: Vulnerable Gin handler sending file
func ginFileVulnerable(c *GinContext) {
	doc := c.Query("path")
	fullPath := filepath.Join("/assets", doc)
	// ruleid: go-arbitrary-file-read
	c.File(fullPath)
}

// Edge case 4: Safe Gin handler sanitized with filepath.IsLocal
func ginFileSafeWithIsLocal(c *GinContext) {
	relPath := c.Query("path")
	if !filepath.IsLocal(relPath) {
		return
	}
	// ok: go-arbitrary-file-read
	c.File(relPath)
}

// Edge case 5: Vulnerable Echo handler with filepath.Clean (Clean does not prevent traversal)
func echoFileVulnerable(c EchoContext) {
	rawPath := c.FormValue("doc")
	cleanPath := filepath.Clean(rawPath)
	// ruleid: go-arbitrary-file-read
	_ = c.File(cleanPath)
}

// Edge case 6: Safe Fiber handler sanitized with filepath.Base
func fiberSendFileSafeBase(c *FiberCtx) {
	asset := c.Query("asset")
	safeAsset := filepath.Base(asset)
	safePath := filepath.Join("/public", safeAsset)
	// ok: go-arbitrary-file-read
	_ = c.SendFile(safePath)
}
