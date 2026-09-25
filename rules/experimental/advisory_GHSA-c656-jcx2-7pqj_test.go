package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) Param(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c EchoContext) FormValue(key string) string  { return "" }
func (c EchoContext) QueryParam(key string) string { return "" }
func (c EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string  { return "" }
func (c *FiberCtx) Params(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error     { return nil }

func testNetHTTPDirectVulnerable(w http.ResponseWriter, req *http.Request) {
	file := req.URL.Query().Get("file")
	target := filepath.Join("/var/www/uploads", file)
	// ruleid: go-path-traversal
	os.Open(target)
}

func testNetHTTPSanitizedIsLocal(w http.ResponseWriter, req *http.Request) {
	file := req.FormValue("file")
	if !filepath.IsLocal(file) {
		return
	}
	target := filepath.Join("/var/www/uploads", file)
	// ok: go-path-traversal
	os.ReadFile(target)
}

func testGinVulnerable(c *GinContext) {
	doc := c.Query("doc")
	fullPath := filepath.Join("/data/docs", doc)
	// ruleid: go-path-traversal
	c.File(fullPath)
}

func testEchoSanitizedBase(c EchoContext) {
	filename := c.FormValue("filename")
	safeName := filepath.Base(filename)
	dest := filepath.Join("/safe/dir", safeName)
	// ok: go-path-traversal
	c.File(dest)
}

func testFiberVulnerable(c *FiberCtx) {
	page := c.Query("page")
	target := filepath.Join("/srv/static", page)
	// ruleid: go-path-traversal
	c.SendFile(target)
}

func testNetHTTPSanitizedBase(w http.ResponseWriter, req *http.Request) {
	cleanName := filepath.Base(req.URL.Path)
	target := filepath.Join("/tmp/storage", cleanName)
	// ok: go-path-traversal
	os.Create(target)
}
