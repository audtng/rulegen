package rules

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) Param(key string) string { return "" }
func (c *GinContext) PostForm(key string) string { return "" }
func (c *GinContext) File(filepath string) {}

type EchoContext struct{}

func (c *EchoContext) FormValue(key string) string { return "" }
func (c *EchoContext) QueryParam(key string) string { return "" }
func (c *EchoContext) Param(key string) string { return "" }
func (c *EchoContext) File(file string) error { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string { return "" }
func (c *FiberCtx) Params(key string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error { return nil }

func testNetHttpVulnerable(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("filename")
	targetPath := filepath.Join("/var/www/uploads", filename)
	// ruleid: go-path-traversal
	os.Open(targetPath)
}

func testNetHttpSafeBase(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("filename")
	safeName := filepath.Base(filename)
	targetPath := filepath.Join("/var/www/uploads", safeName)
	// ok: go-path-traversal
	os.Open(targetPath)
}

func testEchoVulnerable(ctx *EchoContext) {
	userPath := ctx.FormValue("path")
	fullPath := filepath.Join("/data/reports", userPath)
	// ruleid: go-path-traversal
	ctx.File(fullPath)
}

func testEchoSafeIsLocal(ctx *EchoContext) {
	userPath := ctx.FormValue("path")
	if filepath.IsLocal(userPath) {
		fullPath := filepath.Join("/data/reports", userPath)
		// ok: go-path-traversal
		ctx.File(fullPath)
	}
}

func testGinVulnerableSprintf(ctx *GinContext) {
	doc := ctx.Query("doc")
	path := fmt.Sprintf("/var/docs/%s.pdf", doc)
	// ruleid: go-path-traversal
	os.ReadFile(path)
}

func testFiberSafeValidPath(ctx *FiberCtx) {
	filename := ctx.Query("file")
	if fs.ValidPath(filename) {
		// ok: go-path-traversal
		ctx.SendFile(filename)
	}
}
