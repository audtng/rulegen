package rules

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string    { return "" }
func (c *GinContext) Param(key string) string    { return "" }
func (c *GinContext) PostForm(key string) string { return "" }
func (c *GinContext) File(path string)           {}

type EchoContext struct{}

func (c EchoContext) FormValue(key string) string  { return "" }
func (c EchoContext) QueryParam(key string) string { return "" }
func (c EchoContext) File(path string)             {}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string    { return "" }
func (c *FiberCtx) Params(key string) string   { return "" }
func (c *FiberCtx) SendFile(path string) error { return nil }

// Case 1: net/http query parameter directly flows into os.RemoveAll (Vulnerable)
func case1NetHttpRemoveAllVulnerable(w http.ResponseWriter, r *http.Request) {
	plugin := r.URL.Query().Get("plugin")
	pathToRemove := filepath.Join("/var/cache/plugins", plugin)
	// ruleid: path-traversal-arbitrary-file-deletion
	os.RemoveAll(pathToRemove)
}

// Case 2: net/http query parameter sanitized with filepath.Base before os.RemoveAll (Safe)
func case2NetHttpRemoveAllSafeBase(w http.ResponseWriter, r *http.Request) {
	plugin := r.URL.Query().Get("plugin")
	safeName := filepath.Base(plugin)
	pathToRemove := filepath.Join("/var/cache/plugins", safeName)
	// ok: path-traversal-arbitrary-file-deletion
	os.RemoveAll(pathToRemove)
}

// Case 3: Gin context Query parameter used directly in file sink (Vulnerable)
func case3GinQueryFileVulnerable(c *GinContext) {
	filePath := filepath.Join("/var/www/uploads", c.Query("file"))
	// ruleid: path-traversal-arbitrary-file-deletion
	c.File(filePath)
}

// Case 4: Gin context Query validated with filepath.IsLocal before deletion (Safe)
func case4GinIsLocalSafe(c *GinContext) {
	plugin := c.Query("plugin")
	if !filepath.IsLocal(plugin) {
		return
	}
	pathToRemove := filepath.Join("/var/cache/plugins", plugin)
	// ok: path-traversal-arbitrary-file-deletion
	os.RemoveAll(pathToRemove)
}

// Case 5: Echo FormValue flowing into os.Remove without validation (Vulnerable)
func case5EchoFormValueRemoveVulnerable(c EchoContext) {
	target := filepath.Join("/tmp/data", c.FormValue("target"))
	// ruleid: path-traversal-arbitrary-file-deletion
	os.Remove(target)
}

// Case 6: Fiber Query validated with strings.HasPrefix boundary check before SendFile (Safe)
func case6FiberSafePrefix(c *FiberCtx) {
	baseDir := "/var/data"
	filename := c.Query("doc")
	target := filepath.Join(baseDir, filename)
	if !strings.HasPrefix(target, baseDir) {
		return
	}
	// ok: path-traversal-arbitrary-file-deletion
	c.SendFile(target)
}
