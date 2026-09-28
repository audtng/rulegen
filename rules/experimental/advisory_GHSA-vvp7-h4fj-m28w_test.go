package rules

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c *EchoContext) FormValue(key string) string { return "" }
func (c *EchoContext) File(file string) error      { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error                       { return nil }

// Edge case 1: Standard library net/http query parameter with filepath.Join (Vulnerable)
func vulnerableHTTPQueryJoinOpen(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	target := filepath.Join("/data/media", name)
	// ruleid: go-path-traversal
	os.Open(target)
}

// Edge case 2: Standard library net/http sanitized with filepath.Base (Safe / Patched)
func safeHTTPQueryBaseJoinOpen(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	safeName := filepath.Base(name)
	target := filepath.Join("/data/media", safeName)
	// ok: go-path-traversal
	os.Open(target)
}

// Edge case 3: Gin framework query parameter passed to c.File (Vulnerable)
func vulnerableGinQueryFile(c *GinContext) {
	filename := c.Query("file")
	target := filepath.Join("/static", filename)
	// ruleid: go-path-traversal
	c.File(target)
}

// Edge case 4: Gin framework query parameter guarded by filepath.IsLocal (Safe)
func safeGinQueryIsLocal(c *GinContext) {
	filename := c.Query("file")
	if !filepath.IsLocal(filename) {
		return
	}
	target := filepath.Join("/static", filename)
	// ok: go-path-traversal
	c.File(target)
}

// Edge case 5: Echo framework form value joined via path.Join to os.ReadFile (Vulnerable)
func vulnerableEchoFormValueReadFile(c *EchoContext) {
	doc := c.FormValue("doc")
	target := path.Join("/docs", doc)
	// ruleid: go-path-traversal
	os.ReadFile(target)
}

// Edge case 6: Fiber framework query parameter sanitized with filepath.Base to SendFile (Safe)
func safeFiberQueryBaseSendFile(c *FiberCtx) {
	asset := c.Query("asset")
	safeAsset := filepath.Base(asset)
	target := filepath.Join("/assets", safeAsset)
	// ok: go-path-traversal
	c.SendFile(target)
}
