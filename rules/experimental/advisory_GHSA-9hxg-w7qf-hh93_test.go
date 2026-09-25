package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)     {}

type EchoContext struct{}

func (c EchoContext) FormValue(name string) string { return "" }
func (c EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error    { return nil }

func VulnHTTPJoin(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	target := filepath.Join("/var/www/uploads", name)
	// ruleid: path-traversal
	os.Open(target)
}

func VulnGin(c *GinContext) {
	path := c.Query("path")
	// ruleid: path-traversal
	c.File(path)
}

func VulnEcho(c EchoContext) {
	doc := c.FormValue("doc")
	target := filepath.Join("/data", doc)
	// ruleid: path-traversal
	c.File(target)
}

func VulnFiber(c *FiberCtx) {
	filename := c.Query("filename")
	target := filepath.Join("/static", filename)
	// ruleid: path-traversal
	c.SendFile(target)
}

func SafeHTTPBase(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	safe := filepath.Base(name)
	target := filepath.Join("/var/www/uploads", safe)
	// ok: path-traversal
	os.Open(target)
}

func SafeHTTPIsLocal(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	if !filepath.IsLocal(name) {
		return
	}
	target := filepath.Join("/var/www/uploads", name)
	// ok: path-traversal
	os.Open(target)
}
