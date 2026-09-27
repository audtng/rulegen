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

func (c *EchoContext) FormValue(name string) string { return "" }
func (c *EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error                      { return nil }

// Edge Case 1: Direct net/http request param flowing to os.Open (Vulnerable)
func netHttpDirectVulnerable(r *http.Request) {
	p := r.URL.Query().Get("path")
	// ruleid: go-path-traversal
	_, _ = os.Open(p)
}

// Edge Case 2: net/http with filepath.Join and filepath.Base sanitization (Safe)
func netHttpSafeBase(r *http.Request) {
	p := r.FormValue("file")
	safe := filepath.Base(p)
	target := filepath.Join("/var/www/uploads", safe)
	// ok: go-path-traversal
	_, _ = os.ReadFile(target)
}

// Edge Case 3: Gin context Query tainted path to Gin File sink (Vulnerable)
func ginFileVulnerable(c *GinContext) {
	p := c.Query("doc")
	target := filepath.Join("/static", p)
	// ruleid: go-path-traversal
	c.File(target)
}

// Edge Case 4: Gin context with filepath.IsLocal validation guard (Safe)
func ginSafeIsLocal(c *GinContext) {
	p := c.Query("doc")
	if !filepath.IsLocal(p) {
		return
	}
	target := filepath.Join("/static", p)
	// ok: go-path-traversal
	c.File(target)
}

// Edge Case 5: Echo FormValue flowing to os.OpenFile sink (Vulnerable)
func echoOpenFileVulnerable(c *EchoContext) {
	name := c.FormValue("name")
	// ruleid: go-path-traversal
	_, _ = os.OpenFile(name, os.O_RDONLY, 0600)
}

// Edge Case 6: Fiber Query sanitized with filepath.IsLocal before SendFile (Safe)
func fiberSendFileSafe(c *FiberCtx) {
	doc := c.Query("file")
	if filepath.IsLocal(doc) {
		// ok: go-path-traversal
		_ = c.SendFile(doc)
	}
}
