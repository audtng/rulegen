package rules

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) Param(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c EchoContext) FormValue(name string) string { return "" }
func (c EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key ...string) string                 { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error { return nil }

// Edge Case 1: Untrusted query parameter in net/http joined and written to file.
func testHTTPWriteVuln(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/safe/dir", filename)
	// ruleid: path-traversal
	os.WriteFile(target, []byte("data"), 0644)
}

// Edge Case 2: Untrusted parameter sanitized with filepath.Base.
func testHTTPSafeBase(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	baseName := filepath.Base(filename)
	target := filepath.Join("/safe/dir", baseName)
	// ok: path-traversal
	os.OpenFile(target, os.O_RDWR, 0644)
}

// Edge Case 3: Untrusted parameter guarded by filepath.IsLocal validation.
func testHTTPSafeIsLocal(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("file")
	if !filepath.IsLocal(filename) {
		return
	}
	target := filepath.Join("/safe/dir", filename)
	// ok: path-traversal
	os.Create(target)
}

// Edge Case 4: Gin context query parameter joined and served via File sink.
func testGinFileVuln(c *GinContext) {
	filename := c.Query("filepath")
	target := filepath.Join("/data/uploads", filename)
	// ruleid: path-traversal
	c.File(target)
}

// Edge Case 5: Echo context form parameter joined with path.Join and deleted via RemoveAll sink.
func testEchoDeleteVuln(c EchoContext) {
	docPath := c.FormValue("doc")
	target := path.Join("/var/docs", docPath)
	// ruleid: path-traversal
	os.RemoveAll(target)
}

// Edge Case 6: Fiber context query parameter joined and sent via SendFile sink.
func testFiberSendFileVuln(c *FiberCtx) {
	report := c.Query("report")
	target := filepath.Join("/reports", report)
	// ruleid: path-traversal
	c.SendFile(target)
}
