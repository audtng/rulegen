package rules

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string    { return "" }
func (c *GinContext) Param(key string) string    { return "" }
func (c *GinContext) PostForm(key string) string { return "" }
func (c *GinContext) File(filepath string)       {}

type EchoContext struct{}

func (c *EchoContext) FormValue(key string) string  { return "" }
func (c *EchoContext) QueryParam(key string) string { return "" }
func (c *EchoContext) Param(key string) string      { return "" }
func (c *EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, args ...string) string             { return "" }
func (c *FiberCtx) Params(key string, defaultValue ...string) string    { return "" }
func (c *FiberCtx) FormValue(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error        { return nil }

func testNetHttpQueryOpenVuln(req *http.Request) {
	filename := req.URL.Query().Get("file")
	target := filepath.Join("/var/www/uploads", filename)
	// ruleid: path-traversal-arbitrary-file-read
	os.Open(target)
}

func testNetHttpFormValueBaseSafe(req *http.Request) {
	filename := req.FormValue("doc")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/www/uploads", safeName)
	// ok: path-traversal-arbitrary-file-read
	os.ReadFile(target)
}

func testGinQueryFileVuln(c *GinContext) {
	p := c.Query("path")
	fullPath := filepath.Join("/data/files", p)
	// ruleid: path-traversal-arbitrary-file-read
	c.File(fullPath)
}

func testGinParamIsLocalSafe(c *GinContext) {
	userPath := c.Param("filepath")
	if !filepath.IsLocal(userPath) {
		return
	}
	// ok: path-traversal-arbitrary-file-read
	os.OpenFile(userPath, os.O_RDONLY, 0400)
}

func testEchoFormValueServeFileVuln(w http.ResponseWriter, req *http.Request, c *EchoContext) {
	doc := c.FormValue("name")
	filePath := fmt.Sprintf("/var/data/%s", doc)
	// ruleid: path-traversal-arbitrary-file-read
	http.ServeFile(w, req, filePath)
}

func testFiberQuerySendFileSafe(c *FiberCtx) {
	input := c.Query("attachment")
	cleanName := path.Base(input)
	fullPath := filepath.Join("/storage/media", cleanName)
	// ok: path-traversal-arbitrary-file-read
	c.SendFile(fullPath)
}
