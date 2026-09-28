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

func (c EchoContext) FormValue(name string) string { return "" }
func (c EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string                      { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error { return nil }

// Test 1: Vulnerable net/http handler accessing arbitrary file via URL query param
func testNetHTTPVuln(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/var/www/uploads", filename)
	// ruleid: go-directory-traversal
	os.Open(target)
}

// Test 2: Vulnerable Gin handler passing query param directly to File sink
func testGinVuln(c *GinContext) {
	path := c.Query("path")
	// ruleid: go-directory-traversal
	c.File(path)
}

// Test 3: Vulnerable Echo handler reading directory contents without validation
func testEchoReadDirVuln(c EchoContext) {
	dir := c.FormValue("dir")
	target := filepath.Join("/data/workspace", dir)
	// ruleid: go-directory-traversal
	os.ReadDir(target)
}

// Test 4: Safe net/http handler sanitized with filepath.Base
func testNetHTTPSafe(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/www/uploads", safeName)
	// ok: go-directory-traversal
	os.Open(target)
}

// Test 5: Safe Gin handler protected by filepath.IsLocal check
func testGinSafeIsLocal(c *GinContext) {
	path := c.Query("path")
	if filepath.IsLocal(path) {
		// ok: go-directory-traversal
		c.File(path)
	}
}

// Test 6: Safe Fiber handler sanitized with filepath.Base
func testFiberSafe(c *FiberCtx) {
	filename := c.Query("file")
	safeName := filepath.Base(filename)
	// ok: go-directory-traversal
	c.SendFile(safeName)
}
