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

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error   { return nil }

func vulnerableNetHTTP(r *http.Request) {
	file := r.URL.Query().Get("file")
	target := filepath.Join("/data/uploads", file)
	// ruleid: go-path-traversal
	f, _ := os.Open(target)
	_ = f
}

func vulnerableGin(c *GinContext) {
	doc := c.Query("doc")
	target := filepath.Join("/var/www/assets", doc)
	// ruleid: go-path-traversal
	c.File(target)
}

func vulnerableEcho(c EchoContext) {
	name := c.FormValue("name")
	target := filepath.Join("/var/reports", name)
	// ruleid: go-path-traversal
	_ = c.File(target)
}

func vulnerableFiber(c *FiberCtx) {
	item := c.Query("item")
	target := filepath.Join("/var/downloads", item)
	// ruleid: go-path-traversal
	_ = c.SendFile(target)
}

func safeFilepathBase(r *http.Request) {
	file := r.URL.Query().Get("file")
	base := filepath.Base(file)
	target := filepath.Join("/data/uploads", base)
	// ok: go-path-traversal
	f, _ := os.Open(target)
	_ = f
}

func safeIsLocal(r *http.Request) {
	file := r.FormValue("file")
	if !filepath.IsLocal(file) {
		return
	}
	target := filepath.Join("/data/uploads", file)
	// ok: go-path-traversal
	data, _ := os.ReadFile(target)
	_ = data
}
