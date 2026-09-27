package rules

import (
	"archive/zip"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) Param(key string) string { return "" }
func (c *GinContext) File(path string)        {}

type EchoContext struct{}

func (c EchoContext) FormValue(key string) string { return "" }
func (c EchoContext) File(path string)            {}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string    { return "" }
func (c *FiberCtx) SendFile(path string) error { return nil }

// Edge Case 1: net/http upload filename used directly in zip archive entry creation (Vulnerable)
func case1ZipCreationVulnerable(w http.ResponseWriter, r *http.Request, zw *zip.Writer) {
	filename := r.URL.Query().Get("filename")
	// ruleid: tainted-path-traversal
	zw.Create(filename)
}

// Edge Case 2: net/http filename sanitized with filepath.Base before zip creation (Safe)
func case2ZipCreationSafeBase(w http.ResponseWriter, r *http.Request, zw *zip.Writer) {
	filename := r.URL.Query().Get("filename")
	safeName := filepath.Base(filename)
	// ok: tainted-path-traversal
	zw.Create(safeName)
}

// Edge Case 3: Gin context Query parameter used directly in File sink (Vulnerable)
func case3GinVulnerable(c *GinContext) {
	filePath := filepath.Join("/var/www/uploads", c.Query("file"))
	// ruleid: tainted-path-traversal
	c.File(filePath)
}

// Edge Case 4: Gin context Query validated with filepath.IsLocal (Safe)
func case4GinIsLocalSafe(c *GinContext) {
	filename := c.Query("file")
	if !filepath.IsLocal(filename) {
		return
	}
	filePath := filepath.Join("/var/www/uploads", filename)
	// ok: tainted-path-traversal
	c.File(filePath)
}

// Edge Case 5: Echo FormValue flowing into os.Create (Vulnerable)
func case5EchoVulnerable(c EchoContext) {
	target := filepath.Join("/tmp/output", c.FormValue("name"))
	// ruleid: tainted-path-traversal
	os.Create(target)
}

// Edge Case 6: Fiber Query validated with strings.HasPrefix boundary check (Safe)
func case6FiberSafePrefix(c *FiberCtx) {
	baseDir := "/var/data"
	filename := c.Query("doc")
	target := filepath.Join(baseDir, filename)
	if !strings.HasPrefix(target, baseDir) {
		return
	}
	// ok: tainted-path-traversal
	c.SendFile(target)
}
