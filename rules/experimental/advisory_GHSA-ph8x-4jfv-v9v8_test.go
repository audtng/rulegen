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

func (c EchoContext) FormValue(key string) string { return "" }
func (c EchoContext) File(file string) error      { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, args ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error              { return nil }

// Edge Case 1: Direct taint from net/http URL query to os.Open (Vulnerable)
func testNetHTTPDirect(r *http.Request) {
	p := r.URL.Query().Get("file")
	// ruleid: go-path-traversal
	os.Open(p)
}

// Edge Case 2: Taint propagated through filepath.Join to os.ReadFile (Vulnerable)
func testNetHTTPJoin(r *http.Request) {
	p := r.FormValue("path")
	joined := filepath.Join("/safe/dir", p)
	// ruleid: go-path-traversal
	os.ReadFile(joined)
}

// Edge Case 3: Web framework (Gin) query parameter flowing to File sink (Vulnerable)
func testGinContextSink(c *GinContext) {
	p := c.Query("path")
	// ruleid: go-path-traversal
	c.File(p)
}

// Edge Case 4: Web framework (Fiber) query parameter flowing to SendFile sink (Vulnerable)
func testFiberContextSink(c *FiberCtx) {
	p := c.Query("path")
	// ruleid: go-path-traversal
	c.SendFile(p)
}

// Edge Case 5: Sanitized using filepath.IsLocal guard (Safe)
func testSanitizedIsLocal(r *http.Request) {
	p := r.URL.Query().Get("file")
	if !filepath.IsLocal(p) {
		return
	}
	// ok: go-path-traversal
	os.ReadFile(p)
}

// Edge Case 6: Web framework (Echo) sanitized using filepath.Base (Safe)
func testEchoSanitized(c EchoContext) {
	name := c.FormValue("name")
	safeName := filepath.Base(name)
	// ok: go-path-traversal
	c.File(filepath.Join("/uploads", safeName))
}
