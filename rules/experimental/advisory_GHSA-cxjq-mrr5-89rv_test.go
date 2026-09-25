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

func (c *EchoContext) File(filepath string) error {
	return nil
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string {
	return ""
}

func (c *FiberCtx) SendFile(file string) error {
	return nil
}

// 1. net/http: Raw URL path reaches file system sink without sanitization
func testHttpPathVuln(req *http.Request) {
	path := req.URL.Path
	target := filepath.Join("/var/www/uploads", path)
	// ruleid: go-path-traversal
	os.Open(target)
}

// 2. net/http: Query parameter sanitized using filepath.Base
func testHttpQuerySafe(req *http.Request) {
	param := req.URL.Query().Get("file")
	safeName := filepath.Base(param)
	target := filepath.Join("/var/www/uploads", safeName)
	// ok: go-path-traversal
	os.OpenFile(target, os.O_RDONLY, 0400)
}

// 3. Gin: Query parameter passed directly to File sink
func testGinQueryVuln(c *GinContext) {
	param := c.Query("filepath")
	target := filepath.Join("/assets", param)
	// ruleid: go-path-traversal
	c.File(target)
}

// 4. Gin: Query parameter sanitized via filepath.IsLocal
func testGinQuerySafe(c *GinContext) {
	param := c.Query("filepath")
	if !filepath.IsLocal(param) {
		return
	}
	target := filepath.Join("/assets", param)
	// ok: go-path-traversal
	c.File(target)
}

// 5. Echo: FormValue parameter passed directly to File sink
func testEchoFormVuln(c *EchoContext) {
	param := c.FormValue("doc")
	// ruleid: go-path-traversal
	c.File(param)
}

// 6. Fiber: Query parameter sanitized via filepath.Base before SendFile
func testFiberQuerySafe(c *FiberCtx) {
	param := c.Query("filename")
	safe := filepath.Base(param)
	target := filepath.Join("/static", safe)
	// ok: go-path-traversal
	c.SendFile(target)
}
