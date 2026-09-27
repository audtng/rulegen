package rules

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GinContext struct {
	Request *http.Request
}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) Param(key string) string { return "" }
func (c *GinContext) File(path string)        {}

type EchoContext struct{}

func (c EchoContext) FormValue(key string) string { return "" }
func (c EchoContext) File(path string)            {}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string    { return "" }
func (c *FiberCtx) SendFile(path string) error { return nil }

// Edge Case 1: Direct net/http URL path traversal into os.Open (Vulnerable)
func case1NetHttpVulnerable(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("/var/www/uploads", r.URL.Path)
	// ruleid: path-traversal-arbitrary-file-read
	os.Open(path)
}

// Edge Case 2: net/http sanitized with filepath.Base (Safe)
func case2NetHttpSafeBase(w http.ResponseWriter, r *http.Request) {
	safePath := filepath.Join("/var/www/uploads", filepath.Base(r.URL.Path))
	// ok: path-traversal-arbitrary-file-read
	os.Open(safePath)
}

// Edge Case 3: Gin URL path traversal with strings.TrimPrefix matching CVE archetype (Vulnerable)
func case3GinTrimPrefixVulnerable(c *GinContext) {
	relPath := strings.TrimPrefix(c.Request.URL.Path, "/appearance/")
	filePath := filepath.Join("/data/themes", relPath)
	// ruleid: path-traversal-arbitrary-file-read
	c.File(filePath)
}

// Edge Case 4: Gin with filepath.IsLocal validation matching CVE patch archetype (Safe)
func case4GinIsLocalSafe(c *GinContext) {
	relPath := strings.TrimPrefix(c.Request.URL.Path, "/appearance/")
	if !filepath.IsLocal(relPath) {
		return
	}
	filePath := filepath.Join("/data/themes", relPath)
	// ok: path-traversal-arbitrary-file-read
	c.File(filePath)
}

// Edge Case 5: Echo FormValue into os.ReadFile without validation (Vulnerable)
func case5EchoVulnerable(c EchoContext) {
	filename := c.FormValue("file")
	target := filepath.Join("/tmp/docs", filename)
	// ruleid: path-traversal-arbitrary-file-read
	os.ReadFile(target)
}

// Edge Case 6: Fiber Query sanitized with strings.HasPrefix directory boundary check (Safe)
func case6FiberSafePrefix(c *FiberCtx) {
	baseDir := "/var/data"
	filename := c.Query("filename")
	target := filepath.Join(baseDir, filename)
	if !strings.HasPrefix(target, baseDir) {
		return
	}
	// ok: path-traversal-arbitrary-file-read
	c.SendFile(target)
}
