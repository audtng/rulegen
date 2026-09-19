package rules

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Mock GinContext for testing web sources without external dependencies
type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return "test"
}

func (c *GinContext) Param(key string) string {
	return "test"
}

func (c *GinContext) PostForm(key string) string {
	return "test"
}

// 1. Direct standard lib vulnerability
func testDirectVuln(req *http.Request) {
	path := req.URL.Query().Get("path")
	// ruleid: go-path-traversal
	os.Open(path)
}

// 2. Proper standard lib patch
func testProperPatch(req *http.Request) {
	path := req.URL.Query().Get("path")
	if !filepath.IsLocal(path) {
		return
	}
	// ok: go-path-traversal
	os.Open(path)
}

// 3. Cross-function taint (wrapper function bypass)
func runWrapper(fn func()) {
	fn()
}

func testCrossFunctionTaint(req *http.Request) {
	path := req.URL.Query().Get("path")
	runWrapper(func() {
		// ruleid: go-path-traversal
		os.Open(path)
	})
}

// 4. Interface abstraction bypass
func testInterfaceAbstractionBypass(req *http.Request, fsys fs.FS) {
	path := req.URL.Query().Get("path")
	// ruleid: go-path-traversal
	fsys.Open(path)
}

// 5. Fake sanitizer usage (must trigger alert)
func testFakeSanitizer(c *GinContext) {
	path := c.Query("path")
	// filepath.Join is a fake sanitizer that does not prevent path traversal
	fullPath := filepath.Join("/var/www/static", path)
	// ruleid: go-path-traversal
	os.Open(fullPath)
}

// 6. Real sanitizer usage (must not trigger alert)
func testRealSanitizer(req *http.Request) {
	path := req.URL.Query().Get("path")
	fullPath := filepath.Join("/var/www/static", path)
	if !strings.HasPrefix(fullPath, "/var/www/static/") {
		return
	}
	// ok: go-path-traversal
	os.Open(fullPath)
}
