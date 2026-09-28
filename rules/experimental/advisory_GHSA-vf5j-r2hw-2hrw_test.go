package rules

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c EchoContext) FormValue(name string) string { return "" }
func (c EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error                      { return nil }

func testNetHttpVulnerable(r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/data", filename)
	// ruleid: go-path-traversal
	os.Open(target)
}

func testNetHttpSanitizedBase(r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/data", safeName)
	// ok: go-path-traversal
	os.Open(target)
}

func testGinPrefixCheckVulnerable(c *GinContext) {
	param := c.Query("path")
	baseDir := "/var/data"
	fullPath := filepath.Join(baseDir, param)
	// Insecure prefix check does not prevent path traversal / directory boundary confusion
	if strings.HasPrefix(fullPath, baseDir) {
		// ruleid: go-path-traversal
		c.File(fullPath)
	}
}

func testGinIsLocalSafe(c *GinContext) {
	param := c.Query("path")
	if !filepath.IsLocal(param) {
		return
	}
	target := filepath.Join("/var/data", param)
	// ok: go-path-traversal
	c.File(target)
}

func testEchoFormValueVulnerable(c EchoContext) {
	doc := c.FormValue("doc")
	target := filepath.Join("/docs", doc)
	// ruleid: go-path-traversal
	c.File(target)
}

func testFiberSanitizedSafe(c *FiberCtx) {
	report := c.Query("report")
	safeReport := filepath.Base(report)
	target := filepath.Join("/reports", safeReport)
	// ok: go-path-traversal
	c.SendFile(target)
}
