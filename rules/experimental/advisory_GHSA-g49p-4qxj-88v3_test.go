package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(s string) string { return "" }
func (c *GinContext) File(filepath string)  {}

type EchoContext struct{}

func (c EchoContext) FormValue(s string) string { return "" }
func (c EchoContext) File(file string) error    { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error                      { return nil }

// Case 1: Standard library net/http with Header.Get (original Note Mark asset vector)
func handleNetHTTPVuln(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("X-Name")
	target := filepath.Join("/data/assets", name)
	// ruleid: go-path-traversal
	os.Create(target)
}

// Case 2: Standard library net/http sanitized with filepath.Base
func handleNetHTTPSafeBase(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("X-Name")
	safeName := filepath.Base(name)
	target := filepath.Join("/data/assets", safeName)
	// ok: go-path-traversal
	os.Create(target)
}

// Case 3: Standard library net/http guarded with filepath.IsLocal
func handleNetHTTPSafeIsLocal(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	if !filepath.IsLocal(name) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	target := filepath.Join("/data/uploads", name)
	// ok: go-path-traversal
	os.WriteFile(target, []byte("data"), 0644)
}

// Case 4: Gin framework query parameter to c.File sink
func handleGinVuln(c *GinContext) {
	filename := c.Query("file")
	target := filepath.Join("/var/www/uploads", filename)
	// ruleid: go-path-traversal
	c.File(target)
}

// Case 5: Echo framework form value to c.File sink
func handleEchoVuln(c EchoContext) {
	asset := c.FormValue("asset")
	target := filepath.Join("/storage/notes", asset)
	// ruleid: go-path-traversal
	c.File(target)
}

// Case 6: Fiber framework query parameter sanitized with filepath.Base
func handleFiberSafe(c *FiberCtx) {
	filename := c.Query("filename")
	safeFilename := filepath.Base(filename)
	target := filepath.Join("/static/assets", safeFilename)
	// ok: go-path-traversal
	c.SendFile(target)
}
