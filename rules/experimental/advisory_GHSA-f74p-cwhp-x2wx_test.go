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

// Edge Case 1: Standard net/http request with unvalidated query parameter directly opened
func testHTTPDirectTraversal(req *http.Request) {
	file := req.URL.Query().Get("file")
	// ruleid: go-path-traversal
	os.Open(file)
}

// Edge Case 2: Gin framework context with query parameter passed directly to File serving sink
func testGinFileTraversal(c *GinContext) {
	p := c.Query("path")
	// ruleid: go-path-traversal
	c.File(p)
}

// Edge Case 3: Fiber framework context with query parameter propagated via filepath.Join into SendFile
func testFiberSendFileTraversal(c *FiberCtx) {
	relPath := c.Query("rel")
	target := filepath.Join("/data/uploads", relPath)
	// ruleid: go-path-traversal
	c.SendFile(target)
}

// Edge Case 4: Standard net/http request with FormValue sanitized using filepath.Base before os.ReadFile
func testHTTPBaseSanitized(req *http.Request) {
	raw := req.FormValue("doc")
	safeName := filepath.Base(raw)
	target := filepath.Join("/safe/storage", safeName)
	// ok: go-path-traversal
	os.ReadFile(target)
}

// Edge Case 5: Echo framework request validated with filepath.IsLocal before serving
func testEchoIsLocalSanitized(c EchoContext) {
	target := c.FormValue("path")
	if !filepath.IsLocal(target) {
		return
	}
	// ok: go-path-traversal
	c.File(target)
}

// Edge Case 6: Fiber framework request sanitized using filepath.Base before SendFile
func testFiberBaseSanitized(c *FiberCtx) {
	p := c.Query("file")
	safe := filepath.Base(p)
	target := filepath.Join("/data/public", safe)
	// ok: go-path-traversal
	c.SendFile(target)
}
