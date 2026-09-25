package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return "test"
}

func (c *GinContext) File(filepath string) {}

type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string {
	return "test"
}

func (c *EchoContext) File(file string) error {
	return nil
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string {
	return "test"
}

func (c *FiberCtx) SendFile(file string) error {
	return nil
}

// 1. Gin Context query to file write (Vulnerable)
func testGinQueryToFileWriteVuln(c *GinContext) {
	dest := c.Query("dest")
	// ruleid: arbitrary-file-write-path-traversal
	os.WriteFile(dest, []byte("payload"), 0644)
}

// 2. Gin Context query sanitized with filepath.Base (Safe)
func testGinQueryToFileWriteSafe(c *GinContext) {
	dest := c.Query("dest")
	safe := filepath.Base(dest)
	// ok: arbitrary-file-write-path-traversal
	os.WriteFile(safe, []byte("payload"), 0644)
}

// 3. Echo Context form value to File sink (Vulnerable)
func testEchoFormValueToFileVuln(c *EchoContext) {
	filename := c.FormValue("file")
	// ruleid: arbitrary-file-write-path-traversal
	c.File(filename)
}

// 4. Echo Context form value sanitized with filepath.IsLocal (Safe)
func testEchoFormValueToFileSafe(c *EchoContext) {
	filename := c.FormValue("file")
	if !filepath.IsLocal(filename) {
		return
	}
	// ok: arbitrary-file-write-path-traversal
	c.File(filename)
}

// 5. Net/HTTP query param joined and passed to os.Create (Vulnerable)
func testNetHTTPQueryParamToCreateVuln(r *http.Request) {
	p := r.URL.Query().Get("file")
	target := filepath.Join("/var/data", p)
	// ruleid: arbitrary-file-write-path-traversal
	f, _ := os.Create(target)
	if f != nil {
		f.Close()
	}
}

// 6. Fiber Ctx query sanitized with filepath.Base to SendFile (Safe)
func testFiberQueryToSendFileSafe(c *FiberCtx) {
	p := c.Query("path")
	safe := filepath.Base(p)
	// ok: arbitrary-file-write-path-traversal
	c.SendFile(safe)
}
