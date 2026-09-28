package rules

import (
	"net/http"
	"os"
	"path/filepath"
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

func NetHTTPDirectWrite(r *http.Request) {
	userInput := r.URL.Query().Get("filename")
	// ruleid: path-traversal
	os.WriteFile(userInput, []byte("data"), 0644)
}

func NetHTTPSanitizedBase(r *http.Request) {
	userInput := r.FormValue("filename")
	safeName := filepath.Base(userInput)
	// ok: path-traversal
	os.OpenFile(safeName, os.O_RDWR, 0644)
}

func NetHTTPJoinPropagation(r *http.Request) {
	userInput := r.URL.Query().Get("path")
	fullPath := filepath.Join("/var/www/uploads", userInput)
	// ruleid: path-traversal
	os.RemoveAll(fullPath)
}

func EchoFileTraversal(c *EchoContext) {
	userInput := c.FormValue("file")
	// ruleid: path-traversal
	c.File(userInput)
}

func GinFileTraversalSanitized(c *GinContext) {
	userInput := c.Query("path")
	if !filepath.IsLocal(userInput) {
		return
	}
	// ok: path-traversal
	c.File(userInput)
}

func FiberSendFileSanitized(c *FiberCtx) {
	userInput := c.Query("file")
	safePath := filepath.Base(userInput)
	// ok: path-traversal
	c.SendFile(safePath)
}
