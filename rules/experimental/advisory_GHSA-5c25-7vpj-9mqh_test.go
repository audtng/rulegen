package rules

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Mock third-party web frameworks to allow native offline compilation
type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return ""
}

func (c *GinContext) File(filepath string) {}

type EchoContext interface {
	FormValue(name string) string
	File(file string) error
}

type FiberCtx struct{}

func (f *FiberCtx) Query(key string, defaultValue ...string) string {
	return ""
}

func (f *FiberCtx) SendFile(file string) error {
	return nil
}

// Edge case 1: Prefix confusion / trimming with path.Join to http.ServeFile (Nezha archetype)
func testNetHTTPPrefixConfusion(w http.ResponseWriter, r *http.Request) {
	strip := strings.TrimPrefix(r.URL.Path, "/dashboard")
	target := path.Join("/var/www/dashboard", strip)
	// ruleid: go-http-path-traversal
	http.ServeFile(w, r, target)
}

// Edge case 2: Gin framework handler reading query parameter into c.File
func testGinContextQueryFile(c *GinContext) {
	filename := c.Query("file")
	target := filepath.Join("/var/www/uploads", filename)
	// ruleid: go-http-path-traversal
	c.File(target)
}

// Edge case 3: Echo framework handler reading form value into c.File
func testEchoContextFormValueFile(c EchoContext) {
	doc := c.FormValue("doc")
	target := filepath.Join("/data/reports", doc)
	// ruleid: go-http-path-traversal
	_ = c.File(target)
}

// Edge case 4: Fiber framework handler reading query into c.SendFile
func testFiberCtxQuerySendFile(c *FiberCtx) {
	asset := c.Query("asset")
	target := filepath.Join("/static/assets", asset)
	// ruleid: go-http-path-traversal
	_ = c.SendFile(target)
}

// Edge case 5: Safe path retrieval sanitized using filepath.Base
func testSafeFilepathBase(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)
	target := filepath.Join("/var/www/static", filename)
	// ok: go-http-path-traversal
	http.ServeFile(w, r, target)
}

// Edge case 6: Safe path retrieval sanitized using filepath.IsLocal validation
func testSafeFilepathIsLocal(r *http.Request) {
	userPath := r.URL.Path
	if !filepath.IsLocal(userPath) {
		return
	}
	target := filepath.Join("/data/public", userPath)
	// ok: go-http-path-traversal
	_, _ = os.Open(target)
}
