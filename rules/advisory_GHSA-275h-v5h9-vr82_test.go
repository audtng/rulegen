package rules

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string    { return "" }
func (c *GinContext) Param(key string) string    { return "" }
func (c *GinContext) PostForm(key string) string { return "" }

type EchoContext interface {
	FormValue(name string) string
	QueryParam(name string) string
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string  { return "" }
func (c *FiberCtx) Params(key string, defaultValue ...string) string { return "" }

func openFileWrapper(filePath string) (*os.File, error) {
	return os.Open(filePath)
}

func testDirectVuln(req *http.Request) {
	filePath := req.URL.Query().Get("file")
	// ruleid: go-path-traversal
	os.Open(filePath)
}

func testProperPatch(req *http.Request) {
	filePath := req.URL.Query().Get("file")
	if !filepath.IsLocal(filePath) {
		return
	}
	// ok: go-path-traversal
	os.Open(filePath)
}

func testCrossFunctionTaint(req *http.Request) {
	filePath := req.URL.Query().Get("file")
	// ruleid: go-path-traversal
	openFileWrapper(filePath)
}

func testInterfaceAbstractionBypass(req *http.Request, fsys fs.FS) {
	filePath := req.URL.Query().Get("file")
	// ruleid: go-path-traversal
	fsys.Open(filePath)
}

func testFakeSanitizer(req *http.Request) {
	baseDir := "/var/www/uploads"
	userInput := req.URL.Query().Get("file")
	joinedPath := filepath.Join(baseDir, userInput)
	// ruleid: go-path-traversal
	os.ReadFile(joinedPath)
}

func testRealSanitizer(req *http.Request) {
	baseDir := "/var/www/uploads"
	userInput := req.URL.Query().Get("file")
	if !filepath.IsLocal(userInput) {
		return
	}
	safePath := filepath.Join(baseDir, userInput)
	// ok: go-path-traversal
	os.ReadFile(safePath)
}
