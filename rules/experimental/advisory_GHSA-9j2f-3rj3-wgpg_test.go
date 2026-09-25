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

type EchoContext interface {
	FormValue(key string) string
	File(file string) error
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error   { return nil }

// Edge Case 1: net/http with insecure prefix check (vulnerable archetype)
func testNetHttpInsecurePrefix(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("path")
	fullPath := filepath.Join("/data/user", userInput)
	// Insecure prefix check vulnerable to traversal and prefix collision
	if strings.HasPrefix(fullPath, "/data/user") {
		// ruleid: go-path-traversal
		os.Open(fullPath)
	}
}

// Edge Case 2: Gin framework handler passing tainted input to File sink
func testGinContextVulnerable(c *GinContext) {
	fileParam := c.Query("file")
	target := filepath.Join("/static/uploads", fileParam)
	// ruleid: go-path-traversal
	c.File(target)
}

// Edge Case 3: Echo framework handler passing tainted form value to OpenFile
func testEchoContextVulnerable(c EchoContext) {
	doc := c.FormValue("doc")
	target := filepath.Join("/docs", doc)
	// ruleid: go-path-traversal
	os.OpenFile(target, os.O_RDONLY, 0600)
}

// Edge Case 4: net/http patched with filepath.Rel and negative prefix check
func testNetHttpPatchedRel(w http.ResponseWriter, r *http.Request) {
	userInput := r.FormValue("path")
	rel, err := filepath.Rel("/data/user", userInput)
	if err != nil {
		return
	}
	if !strings.HasPrefix(rel, "..") {
		safePath := filepath.Join("/data/user", rel)
		// ok: go-path-traversal
		os.Open(safePath)
	}
}

// Edge Case 5: Gin framework handler sanitized with filepath.Base
func testGinSanitizedBase(c *GinContext) {
	fileParam := c.Query("file")
	safeName := filepath.Base(fileParam)
	target := filepath.Join("/static/uploads", safeName)
	// ok: go-path-traversal
	c.File(target)
}

// Edge Case 6: Fiber framework handler sanitized with filepath.IsLocal
func testFiberSanitizedIsLocal(c *FiberCtx) {
	subPath := c.Query("page")
	if filepath.IsLocal(subPath) {
		target := filepath.Join("/pages", subPath)
		// ok: go-path-traversal
		c.SendFile(target)
	}
}
