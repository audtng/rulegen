package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

// Mock web framework contexts for offline compilation
type GinContext struct{}

func (g *GinContext) Query(key string) string {
	return ""
}

func (g *GinContext) File(filepath string) {}

type EchoContext struct{}

func (e *EchoContext) FormValue(name string) string {
	return ""
}

func (e *EchoContext) File(file string) error {
	return nil
}

type FiberCtx struct{}

func (f *FiberCtx) Query(key string, defaultValue ...string) string {
	return ""
}

func (f *FiberCtx) SendFile(file string, compress ...bool) error {
	return nil
}

// 1. Vulnerable standard net/http handler with path traversal via filepath.Join
func edgeCase1_NetHTTP_DirectTraversal(r *http.Request) {
	userInput := r.URL.Query().Get("file")
	target := filepath.Join("/var/www/uploads", userInput)
	// ruleid: go-path-traversal
	_, _ = os.Open(target)
}

// 2. Safe standard net/http handler sanitized with filepath.Base
func edgeCase2_NetHTTP_SanitizedBase(r *http.Request) {
	userInput := r.URL.Query().Get("file")
	safeName := filepath.Base(userInput)
	target := filepath.Join("/var/www/uploads", safeName)
	// ok: go-path-traversal
	_, _ = os.Open(target)
}

// 3. Safe standard net/http handler validated with filepath.IsLocal
func edgeCase3_NetHTTP_SanitizedIsLocal(r *http.Request) {
	userInput := r.URL.Query().Get("file")
	if !filepath.IsLocal(userInput) {
		return
	}
	target := filepath.Join("/var/www/uploads", userInput)
	// ok: go-path-traversal
	_, _ = os.ReadFile(target)
}

// 4. Vulnerable Gin handler passing tainted query input to File sink
func edgeCase4_Gin_VulnerableFileServing(c *GinContext) {
	userInput := c.Query("filename")
	target := filepath.Join("/var/www/static", userInput)
	// ruleid: go-path-traversal
	c.File(target)
}

// 5. Vulnerable Echo handler passing tainted form input to File sink
func edgeCase5_Echo_VulnerableFileServing(c *EchoContext) {
	userInput := c.FormValue("path")
	target := filepath.Join("/var/www/data", userInput)
	// ruleid: go-path-traversal
	_ = c.File(target)
}

// 6. Vulnerable Fiber handler passing tainted query parameter to SendFile sink
func edgeCase6_Fiber_VulnerableSendFile(c *FiberCtx) {
	userInput := c.Query("doc")
	target := filepath.Join("/var/www/docs", userInput)
	// ruleid: go-path-traversal
	_ = c.SendFile(target)
}
