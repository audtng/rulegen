package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (g *GinContext) Query(key string) string { return "" }
func (g *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (e EchoContext) FormValue(key string) string { return "" }
func (e EchoContext) File(filepath string) error  { return nil }

type FiberCtx struct{}

func (f *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (f *FiberCtx) SendFile(filepath string, compress ...bool) error { return nil }

// Case 1: Standard net/http request with URL path joined and removed directly
func testNetHTTPRemoveAllTraversal(req *http.Request) {
	p := req.URL.Path
	target := filepath.Join("/var/www/uploads", p)
	// ruleid: go-path-traversal
	os.RemoveAll(target)
}

// Case 2: Gin framework context with query parameter passed directly to File sink
func testGinQueryFileTraversal(c *GinContext) {
	doc := c.Query("doc")
	// ruleid: go-path-traversal
	c.File(doc)
}

// Case 3: Fiber framework context with parameter joined to upload directory
func testFiberQuerySendFileTraversal(c *FiberCtx) {
	filename := c.Query("name")
	target := filepath.Join("/app/uploads", filename)
	// ruleid: go-path-traversal
	c.SendFile(target)
}

// Case 4: Standard net/http request with FormValue sanitized using filepath.Base
func testNetHTTPBaseSanitized(req *http.Request) {
	userPath := req.FormValue("file")
	safeName := filepath.Base(userPath)
	target := filepath.Join("/data/safe", safeName)
	// ok: go-path-traversal
	os.RemoveAll(target)
}

// Case 5: Echo framework request validated with filepath.IsLocal before serving
func testEchoIsLocalSanitized(c EchoContext) {
	doc := c.FormValue("doc")
	if !filepath.IsLocal(doc) {
		return
	}
	// ok: go-path-traversal
	c.File(doc)
}

// Case 6: Fiber framework request sanitized using filepath.Base before sending
func testFiberBaseSanitized(c *FiberCtx) {
	report := c.Query("report")
	safe := filepath.Base(report)
	target := filepath.Join("/reports", safe)
	// ok: go-path-traversal
	c.SendFile(target)
}
