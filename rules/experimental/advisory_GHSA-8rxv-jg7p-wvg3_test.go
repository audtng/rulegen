package rules

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (g *GinContext) Query(key string) string {
	return ""
}

func (g *GinContext) File(filepath string) {}

type EchoContext interface {
	FormValue(key string) string
	File(file string) error
}

type FiberCtx struct{}

func (f *FiberCtx) Query(key string, defaultValue ...string) string {
	return ""
}

func (f *FiberCtx) SendFile(file string, compress ...bool) error {
	return nil
}

// 1. NetHTTP vulnerable handler: directly opens path constructed from URL.Path
func testNetHTTPVuln(req *http.Request) {
	p := req.URL.Path
	target := filepath.Join("/var/www/uploads", p)
	// ruleid: http-path-traversal
	os.Open(target)
}

// 2. NetHTTP safe handler: sanitized using filepath.Base
func testNetHTTPSafeBase(req *http.Request) {
	p := req.URL.Path
	filename := filepath.Base(p)
	target := filepath.Join("/var/www/uploads", filename)
	// ok: http-path-traversal
	os.Open(target)
}

// 3. NetHTTP safe handler: guarded using filepath.IsLocal
func testNetHTTPSafeIsLocal(req *http.Request) {
	p := req.URL.Path
	if !filepath.IsLocal(p) {
		return
	}
	target := filepath.Join("/var/www/uploads", p)
	// ok: http-path-traversal
	os.Open(target)
}

// 4. Gin vulnerable handler: serves file based on query parameter
func testGinVuln(c *GinContext) {
	doc := c.Query("doc")
	target := filepath.Join("/var/data", doc)
	// ruleid: http-path-traversal
	c.File(target)
}

// 5. Echo vulnerable handler: serves file based on form value
func testEchoVuln(c EchoContext) {
	name := c.FormValue("name")
	target := fmt.Sprintf("/var/data/%s", name)
	// ruleid: http-path-traversal
	c.File(target)
}

// 6. Fiber vulnerable handler: sends file based on query parameter
func testFiberVuln(c *FiberCtx) {
	asset := c.Query("asset")
	target := filepath.Join("/var/assets", asset)
	// ruleid: http-path-traversal
	c.SendFile(target)
}
