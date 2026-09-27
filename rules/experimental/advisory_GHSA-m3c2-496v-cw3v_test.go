package rules

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// Mock third-party web frameworks for offline compilation
type GinContext struct{}

func (c *GinContext) Param(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext interface {
	FormValue(name string) string
	File(file string) error
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key ...string) string                   { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error { return nil }

// Edge Case 1: net/http handler with path traversal via filepath.Join and os.Open
func testNetHTTPJoinVuln(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	target := filepath.Join("/var/www/static", filename)
	// ruleid: go-web-path-traversal
	os.Open(target)
}

// Edge Case 2: net/http handler sanitized with filepath.Base
func testNetHTTPSafeBase(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/www/static", safeName)
	// ok: go-web-path-traversal
	os.Open(target)
}

// Edge Case 3: Fiber static file serving using path.Clean (vulnerable archetype from CVE)
func testFiberPathCleanVuln(c *FiberCtx) error {
	userPath := c.Query("path")
	cleaned := path.Clean("/" + userPath)
	// ruleid: go-web-path-traversal
	return c.SendFile(cleaned)
}

// Edge Case 4: Fiber static file serving safely guarded by filepath.IsLocal
func testFiberIsLocalSafe(c *FiberCtx) error {
	userPath := c.Query("path")
	if !filepath.IsLocal(userPath) {
		return nil
	}
	target := filepath.Join("./public", userPath)
	// ok: go-web-path-traversal
	return c.SendFile(target)
}

// Edge Case 5: Gin context route parameter propagated through fmt.Sprintf into File sink
func testGinSprintfVuln(c *GinContext) {
	doc := c.Param("doc")
	target := fmt.Sprintf("/var/data/%s", doc)
	// ruleid: go-web-path-traversal
	c.File(target)
}

// Edge Case 6: Echo context form value sanitized using filepath.Base into File sink
func testEchoSafeBase(c EchoContext) error {
	doc := c.FormValue("doc")
	safeDoc := filepath.Base(doc)
	// ok: go-web-path-traversal
	return c.File(safeDoc)
}
