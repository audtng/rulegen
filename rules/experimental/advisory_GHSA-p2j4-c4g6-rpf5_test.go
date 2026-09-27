package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return ""
}

func (c *GinContext) File(filepath string) {}

type EchoContext interface {
	FormValue(key string) string
	File(file string) error
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string {
	return ""
}

func (c *FiberCtx) SendFile(file string, compress ...bool) error {
	return nil
}

func testNetHttpDirect(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	// ruleid: go-arbitrary-file-read
	os.Open(filename)
}

func testNetHttpSafeBase(r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	// ok: go-arbitrary-file-read
	os.ReadFile(safeName)
}

func testNetHttpSafeIsLocal(r *http.Request) {
	target := r.FormValue("path")
	if !filepath.IsLocal(target) {
		return
	}
	// ok: go-arbitrary-file-read
	os.OpenFile(target, os.O_RDONLY, 0o600)
}

func testGinVulnerable(c *GinContext) {
	doc := c.Query("doc")
	// ruleid: go-arbitrary-file-read
	c.File(doc)
}

func testEchoVulnerable(c EchoContext) {
	target := c.FormValue("file")
	// ruleid: go-arbitrary-file-read
	c.File(target)
}

func testFiberVulnerable(c *FiberCtx) {
	attachment := c.Query("attachment")
	// ruleid: go-arbitrary-file-read
	c.SendFile(attachment)
}
