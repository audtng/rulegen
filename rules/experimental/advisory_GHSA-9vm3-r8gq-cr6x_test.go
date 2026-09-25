package rules

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c *EchoContext) FormValue(key string) string { return "" }
func (c *EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error                       { return nil }

// 1. net/http handler vulnerable to arbitrary file write via path traversal
func vulnNetHTTP(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	target := filepath.Join("/data", filename)
	// ruleid: tainted-path-traversal
	os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0644)
}

// 2. net/http handler sanitized with filepath.Base
func safeNetHTTP(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/data", safeName)
	// ok: tainted-path-traversal
	os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0644)
}

// 3. Gin handler vulnerable to path traversal
func vulnGin(c *GinContext) {
	doc := c.Query("doc")
	target := filepath.Join("/static", doc)
	// ruleid: tainted-path-traversal
	c.File(target)
}

// 4. Echo handler sanitized with filepath.IsLocal
func safeEcho(c *EchoContext) error {
	doc := c.FormValue("doc")
	if !filepath.IsLocal(doc) {
		return nil
	}
	target := filepath.Join("/docs", doc)
	// ok: tainted-path-traversal
	return c.File(target)
}

// 5. Fiber handler vulnerable to path traversal
func vulnFiber(c *FiberCtx) error {
	doc := c.Query("doc")
	target := filepath.Join("/uploads", doc)
	// ruleid: tainted-path-traversal
	return c.SendFile(target)
}

// 6. net/http upload handler sanitized with traversal check (Casdoor archetype)
func safeCasdoorUpload(w http.ResponseWriter, req *http.Request) error {
	fullFilePath := req.FormValue("fullFilePath")
	if strings.Contains(fullFilePath, "..") {
		return errors.New("the fullFilePath is not allowed")
	}
	target := filepath.Join("/uploads", fullFilePath)
	// ok: tainted-path-traversal
	return os.WriteFile(target, []byte("data"), 0644)
}
