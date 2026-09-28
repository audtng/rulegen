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

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string { return "" }
func (c *FiberCtx) SendFile(file string)    {}

// Edge Case 1: Standard net/http request with unvalidated filepath.Join flowing to os.Open
func VulnNetHTTP(w http.ResponseWriter, r *http.Request) {
	userPath := r.URL.Query().Get("path")
	target := filepath.Join("/var/www/uploads", userPath)
	// ruleid: go-path-traversal
	os.Open(target)
}

// Edge Case 2: Gin framework context with unvalidated path.Join flowing to Context.File sink
func VulnGinContext(c *GinContext) {
	relPath := c.Query("file")
	target := path.Join("/data/static", relPath)
	// ruleid: go-path-traversal
	c.File(target)
}

// Edge Case 3: Echo framework context form value flowing into os.ReadFile sink
func VulnEchoContext(c *EchoContext) {
	filename := c.FormValue("doc")
	target := filepath.Join("/tmp/documents", filename)
	// ruleid: go-path-traversal
	os.ReadFile(target)
}

// Edge Case 4: net/http request sanitized with filepath.Base before path join
func SafeFilepathBase(w http.ResponseWriter, r *http.Request) {
	userPath := r.URL.Query().Get("path")
	cleanName := filepath.Base(userPath)
	target := filepath.Join("/var/www/uploads", cleanName)
	// ok: go-path-traversal
	os.Open(target)
}

// Edge Case 5: net/http request sanitized with filepath.IsLocal boundary check
func SafeFilepathIsLocal(w http.ResponseWriter, r *http.Request) {
	userPath := r.URL.Path
	if !filepath.IsLocal(userPath) {
		return
	}
	target := filepath.Join("/var/www/data", userPath)
	// ok: go-path-traversal
	os.Open(target)
}

// Edge Case 6: Fiber framework context sanitized with path.Base before SendFile
func SafeFiberPathBase(c *FiberCtx) {
	relPath := c.Query("file")
	cleanName := path.Base(relPath)
	target := path.Join("/static", cleanName)
	// ok: go-path-traversal
	c.SendFile(target)
}
