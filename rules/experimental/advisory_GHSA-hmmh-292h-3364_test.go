package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

// Mock third-party web frameworks for offline compilation
type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) Param(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string { return "" }
func (c *EchoContext) File(file string) error        { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string    { return "" }
func (c *FiberCtx) SendFile(file string) error { return nil }

func testNetHTTP_Vulnerable(r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/var/data", filename)
	// ruleid: go-path-traversal
	os.Open(target)
}

func testNetHTTP_Safe_Base(r *http.Request) {
	filename := r.FormValue("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/data", safeName)
	// ok: go-path-traversal
	os.ReadFile(target)
}

func testNetHTTP_Safe_IsLocal(r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filepath.IsLocal(filename) {
		target := filepath.Join("/var/data", filename)
		// ok: go-path-traversal
		os.Create(target)
	}
}

func testGin_Vulnerable(c *GinContext) {
	filename := c.Query("filepath")
	target := filepath.Join("/var/www/uploads", filename)
	// ruleid: go-path-traversal
	c.File(target)
}

func testEcho_Safe_Base(c *EchoContext) {
	filename := c.FormValue("download")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/www/downloads", safeName)
	// ok: go-path-traversal
	c.File(target)
}

func testFiber_Vulnerable(c *FiberCtx) {
	filename := c.Query("doc")
	target := filepath.Join("/var/www/docs", filename)
	// ruleid: go-path-traversal
	c.SendFile(target)
}
