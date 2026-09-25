package rules

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string                             { return "" }
func (c *GinContext) Param(key string) string                             { return "" }
func (c *GinContext) PostForm(key string) string                          { return "" }
func (c *GinContext) File(filepath string)                                {}
func (c *GinContext) SaveUploadedFile(file interface{}, dst string) error { return nil }

type EchoContext struct{}

func (c *EchoContext) FormValue(key string) string  { return "" }
func (c *EchoContext) QueryParam(key string) string { return "" }
func (c *EchoContext) File(filepath string) error   { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string                  { return "" }
func (c *FiberCtx) Params(key string) string                 { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error { return nil }

func testNetHttpOpenFileVulnerable(r *http.Request) {
	relPath := r.URL.Query().Get("path")
	fullPath := filepath.Join("/var/www/uploads", relPath)
	// ruleid: path-traversal-arbitrary-file-write
	os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, 0644)
}

func testNetHttpWriteFileSanitizedBase(r *http.Request) {
	userInput := r.FormValue("filename")
	safeName := filepath.Base(userInput)
	fullPath := filepath.Join("/var/www/uploads", safeName)
	// ok: path-traversal-arbitrary-file-write
	os.WriteFile(fullPath, []byte("data"), 0644)
}

func testNetHttpMkdirAllSanitizedIsLocal(r *http.Request) {
	subDir := r.URL.Query().Get("dir")
	if !filepath.IsLocal(subDir) {
		return
	}
	targetDir := filepath.Join("/var/www/data", subDir)
	// ok: path-traversal-arbitrary-file-write
	os.MkdirAll(targetDir, 0755)
}

func testGinServeFileVulnerable(c *GinContext) {
	fileName := c.Query("file")
	target := filepath.Join("/static", fileName)
	// ruleid: path-traversal-arbitrary-file-write
	c.File(target)
}

func testEchoRemoveAllVulnerableCleanNotSafe(c *EchoContext) {
	folder := c.FormValue("folder")
	cleanFolder := filepath.Clean(folder)
	target := fmt.Sprintf("/var/cache/%s", cleanFolder)
	// ruleid: path-traversal-arbitrary-file-write
	os.RemoveAll(target)
}

func testFiberSendFileSafeConstant(c *FiberCtx) {
	const safePath = "/var/www/assets/index.html"
	// ok: path-traversal-arbitrary-file-write
	c.SendFile(safePath)
}
