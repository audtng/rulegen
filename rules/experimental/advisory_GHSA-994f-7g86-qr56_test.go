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

func (c *EchoContext) FormValue(key string) string {
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

// 1. Vulnerable net/http handler with URL query parameter
func VulnNetHTTPHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/var/www/uploads", filename)
	// ruleid: path-traversal
	os.Open(target)
}

// 2. Safe net/http handler sanitized with filepath.Base
func SafeNetHTTPBaseHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/www/uploads", safeName)
	// ok: path-traversal
	os.Open(target)
}

// 3. Safe net/http handler validated with filepath.IsLocal
func SafeNetHTTPIsLocalHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("file")
	if !filepath.IsLocal(filename) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	target := filepath.Join("/var/www/uploads", filename)
	// ok: path-traversal
	os.ReadFile(target)
}

// 4. Vulnerable Gin handler with query param
func VulnGinHandler(c *GinContext) {
	filename := c.Query("file")
	target := filepath.Join("/var/www/uploads", filename)
	// ruleid: path-traversal
	c.File(target)
}

// 5. Vulnerable Fiber handler with query param
func VulnFiberHandler(c *FiberCtx) {
	filename := c.Query("file")
	target := filepath.Join("/var/www/uploads", filename)
	// ruleid: path-traversal
	c.SendFile(target)
}

// 6. Safe Echo handler sanitized with filepath.Base
func SafeEchoHandler(c *EchoContext) {
	filename := c.FormValue("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/www/uploads", safeName)
	// ok: path-traversal
	c.File(target)
}
