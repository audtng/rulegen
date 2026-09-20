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
func (c *GinContext) File(filepath string)       {}

type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string  { return "" }
func (c *EchoContext) QueryParam(name string) string { return "" }
func (c *EchoContext) Param(name string) string      { return "" }
func (c *EchoContext) File(file string) error        { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string     { return "" }
func (c *FiberCtx) Params(key string) string    { return "" }
func (c *FiberCtx) FormValue(key string) string { return "" }
func (c *FiberCtx) SendFile(file string) error  { return nil }

func handleVulnNetHTTPPath(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Path
	targetPath := filepath.Join("/var/www/static", relPath)
	// ruleid: go-path-traversal
	http.ServeFile(w, r, targetPath)
}

func handleSafeNetHTTPBase(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("file")
	safeName := filepath.Base(userInput)
	targetPath := filepath.Join("/var/www/uploads", safeName)
	// ok: go-path-traversal
	os.Open(targetPath)
}

func handleVulnGinQuery(c *GinContext) {
	userInput := c.Query("file")
	targetPath := filepath.Join("/app/data", userInput)
	// ruleid: go-path-traversal
	c.File(targetPath)
}

func handleSafeGinIsLocal(c *GinContext) {
	userInput := c.Query("file")
	if !filepath.IsLocal(userInput) {
		return
	}
	targetPath := filepath.Join("/app/data", userInput)
	// ok: go-path-traversal
	c.File(targetPath)
}

func handleVulnEchoForm(c *EchoContext) {
	userInput := c.FormValue("doc")
	targetPath := filepath.Join("/var/documents", userInput)
	// ruleid: go-path-traversal
	os.ReadFile(targetPath)
}

func handleSafePrefixCheck(c *FiberCtx) {
	userInput := c.Params("path")
	baseDir := "/safe/storage"
	targetPath := filepath.Join(baseDir, userInput)
	if !strings.HasPrefix(targetPath, baseDir) {
		return
	}
	// ok: go-path-traversal
	c.SendFile(targetPath)
}
