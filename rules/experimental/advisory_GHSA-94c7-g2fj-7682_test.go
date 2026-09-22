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

func (c EchoContext) FormValue(key string) string { return "" }
func (c EchoContext) File(filepath string) error  { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error   { return nil }

// Case 1: Vulnerable - net/http query param directly read via os.ReadFile
func testNetHttpDirectFileRead_Vuln(req *http.Request) {
	filename := req.URL.Query().Get("file")
	// ruleid: go-arbitrary-file-read
	os.ReadFile(filename)
}

// Case 2: Vulnerable - Gin context query param served directly via c.File
func testGinFileServe_Vuln(c *GinContext) {
	filePath := c.Query("path")
	// ruleid: go-arbitrary-file-read
	c.File(filePath)
}

// Case 3: Vulnerable - Echo FormValue concatenated via filepath.Join into os.Open
func testEchoJoinFileOpen_Vuln(c EchoContext) {
	userPath := c.FormValue("doc")
	target := filepath.Join("/data/uploads", userPath)
	// ruleid: go-arbitrary-file-read
	os.Open(target)
}

// Case 4: Vulnerable - Fiber context query param sent via c.SendFile
func testFiberSendFile_Vuln(c *FiberCtx) {
	name := c.Query("name")
	// ruleid: go-arbitrary-file-read
	c.SendFile(name)
}

// Case 5: Safe - net/http query sanitized using filepath.Base
func testNetHttpFileRead_BaseSanitized_Safe(req *http.Request) {
	rawFile := req.URL.Query().Get("file")
	safeName := filepath.Base(rawFile)
	safePath := filepath.Join("/safe/directory", safeName)
	// ok: go-arbitrary-file-read
	os.ReadFile(safePath)
}

// Case 6: Safe - Gin context input validated via filepath.IsLocal
func testGinFileOpen_IsLocalSanitized_Safe(c *GinContext) {
	target := c.Query("path")
	if !filepath.IsLocal(target) {
		return
	}
	// ok: go-arbitrary-file-read
	os.Open(target)
}
