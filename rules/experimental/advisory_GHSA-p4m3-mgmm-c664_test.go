package rules

import (
	"fmt"
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

func (c *FiberCtx) Query(key ...string) string                   { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error { return nil }

func testNetHttpVuln(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/data/uploads", filename)
	// ruleid: go-path-traversal
	os.Open(target)
}

func testNetHttpSafeBase(r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/data/uploads", safeName)
	// ok: go-path-traversal
	os.ReadFile(target)
}

func testNetHttpSafeIsLocal(r *http.Request) {
	filename := r.FormValue("file")
	if !filepath.IsLocal(filename) {
		return
	}
	target := filepath.Join("/data/uploads", filename)
	// ok: go-path-traversal
	os.Create(target)
}

func testGinVuln(c *GinContext) {
	doc := c.Param("doc")
	target := fmt.Sprintf("/var/data/%s", doc)
	// ruleid: go-path-traversal
	c.File(target)
}

func testEchoVuln(c EchoContext) {
	docPath := c.FormValue("doc")
	target := path.Join("/var/docs", docPath)
	// ruleid: go-path-traversal
	os.RemoveAll(target)
}

func testFiberVuln(c *FiberCtx) {
	report := c.Query("report")
	target := filepath.Join("/reports", report)
	// ruleid: go-path-traversal
	c.SendFile(target)
}
