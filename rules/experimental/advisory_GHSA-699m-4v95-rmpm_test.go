package rules

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(s string) string { return "" }
func (c *GinContext) File(s string)         {}

type EchoContext struct{}

func (c EchoContext) FormValue(s string) string { return "" }
func (c EchoContext) File(s string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(s string) string          { return "" }
func (c *FiberCtx) SendFile(s string) error        { return nil }

func VulnHttpHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/data", filename)
	// ruleid: go-path-traversal
	os.Open(target)
}

func VulnGinHandler(c *GinContext) {
	filename := c.Query("file")
	target := filepath.Join("/var/www", filename)
	// ruleid: go-path-traversal
	c.File(target)
}

func VulnEchoHandler(c EchoContext) {
	filename := c.FormValue("doc")
	target := fmt.Sprintf("/docs/%s", filename)
	// ruleid: go-path-traversal
	c.File(target)
}

func VulnFiberHandler(c *FiberCtx) {
	filename := c.Query("name")
	target := filepath.Join("/uploads", filename)
	// ruleid: go-path-traversal
	c.SendFile(target)
}

func SafeBaseHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/safe/path", safeName)
	// ok: go-path-traversal
	os.ReadFile(target)
}

func SafeIsLocalHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		return
	}
	target := filepath.Join("/safe/path", filename)
	// ok: go-path-traversal
	os.Open(target)
}
