package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string { return "" }
func (c *EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, args ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error              { return nil }

type FileRequest struct {
	FilePath string `json:"filePath"`
}

// Case 1: JSON body unmarshaled path joined to base dir and read via os.ReadFile (Vulnerable - Goploy archetype)
func testVulnerableJSONBodyRead(data []byte) {
	var req FileRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return
	}
	targetPath := filepath.Join("/var/data/projects", req.FilePath)
	// ruleid: path-traversal-arbitrary-file-read
	os.ReadFile(targetPath)
}

// Case 2: JSON body unmarshaled path sanitized with filepath.Base (Safe)
func testSafeJSONBodyWithBase(data []byte) {
	var req FileRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return
	}
	safeName := filepath.Base(req.FilePath)
	targetPath := filepath.Join("/var/data/projects", safeName)
	// ok: path-traversal-arbitrary-file-read
	os.ReadFile(targetPath)
}

// Case 3: Gin query parameter directly used in Context.File sink (Vulnerable)
func testVulnerableGinQueryFile(c *GinContext) {
	userPath := c.Query("filepath")
	fullPath := filepath.Join("/static/assets", userPath)
	// ruleid: path-traversal-arbitrary-file-read
	c.File(fullPath)
}

// Case 4: Gin query parameter validated with filepath.IsLocal (Safe)
func testSafeGinQueryIsLocal(c *GinContext) {
	userPath := c.Query("filepath")
	if !filepath.IsLocal(userPath) {
		return
	}
	fullPath := filepath.Join("/static/assets", userPath)
	// ok: path-traversal-arbitrary-file-read
	c.File(fullPath)
}

// Case 5: Echo FormValue parameter cleaned with filepath.Clean and passed to Context.File (Vulnerable)
func testVulnerableEchoFormFile(c *EchoContext) {
	filename := c.FormValue("file")
	targetPath := filepath.Clean(filepath.Join("/var/www/uploads", filename))
	// ruleid: path-traversal-arbitrary-file-read
	c.File(targetPath)
}

// Case 6: Fiber query parameter sanitized with filepath.Base before SendFile (Safe)
func testSafeFiberQueryBase(c *FiberCtx) {
	filename := c.Query("doc")
	safeName := filepath.Base(filename)
	docPath := filepath.Join("/docs", safeName)
	// ok: path-traversal-arbitrary-file-read
	c.SendFile(docPath)
}
