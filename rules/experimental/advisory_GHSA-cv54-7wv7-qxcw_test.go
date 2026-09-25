package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c EchoContext) FormValue(name string) string { return "" }
func (c EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error   { return nil }

func testNetHTTPRequestUnsafe(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	// ruleid: arbitrary-file-read-path-traversal
	os.Open(filename)
}

func testNetHTTPBaseSanitized(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	// ok: arbitrary-file-read-path-traversal
	os.ReadFile(safeName)
}

func testNetHTTPIsLocalSanitized(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	// ok: arbitrary-file-read-path-traversal
	http.ServeFile(w, r, filename)
}

func testGinJoinPropagationUnsafe(c *GinContext) {
	userInput := c.Query("path")
	fullPath := filepath.Join("/var/data", userInput)
	// ruleid: arbitrary-file-read-path-traversal
	c.File(fullPath)
}

func testEchoBaseSanitized(c EchoContext) {
	filename := c.FormValue("doc")
	clean := filepath.Base(filename)
	// ok: arbitrary-file-read-path-traversal
	c.File(clean)
}

func testFiberAbsPropagationUnsafe(c *FiberCtx) {
	userInput := c.Query("doc")
	absPath, _ := filepath.Abs(userInput)
	// ruleid: arbitrary-file-read-path-traversal
	c.SendFile(absPath)
}
