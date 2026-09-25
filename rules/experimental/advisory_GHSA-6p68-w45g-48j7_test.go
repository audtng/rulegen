package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return ""
}

func (c *GinContext) File(filepath string) {}

type EchoContext struct{}

func (c EchoContext) FormValue(key string) string {
	return ""
}

func (c EchoContext) File(file string) error {
	return nil
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string {
	return ""
}

func (c *FiberCtx) SendFile(file string) error {
	return nil
}

// 1. Edge Case 1: Unsanitized net/http URL path flows directly into os.Open (Vulnerable)
func testNetHTTPDirectPath(w http.ResponseWriter, req *http.Request) {
	filePath := req.URL.Path
	// ruleid: go-http-path-traversal
	_, _ = os.Open(filePath)
}

// 2. Edge Case 2: net/http URL path sanitized using filepath.Base (Safe)
func testNetHTTPSanitizedBase(w http.ResponseWriter, req *http.Request) {
	filePath := req.URL.Path
	safePath := filepath.Base(filePath)
	// ok: go-http-path-traversal
	_, _ = os.Open(safePath)
}

// 3. Edge Case 3: net/http URL path validated using filepath.IsLocal (Safe)
func testNetHTTPSanitizedIsLocal(w http.ResponseWriter, req *http.Request) {
	filePath := req.URL.Path
	if !filepath.IsLocal(filePath) {
		return
	}
	// ok: go-http-path-traversal
	_, _ = os.Open(filePath)
}

// 4. Edge Case 4: Query parameter joined with root directory flows into http.ServeFile (Vulnerable)
func testNetHTTPServeFileJoin(w http.ResponseWriter, req *http.Request) {
	doc := req.URL.Query().Get("doc")
	target := filepath.Join("/var/www/uploads", doc)
	// ruleid: go-http-path-traversal
	http.ServeFile(w, req, target)
}

// 5. Edge Case 5: Gin framework context query parameter flows into c.File (Vulnerable)
func testGinFileDirect(c *GinContext) {
	filename := c.Query("filepath")
	// ruleid: go-http-path-traversal
	c.File(filename)
}

// 6. Edge Case 6: Echo framework context form value sanitized using filepath.Base before File (Safe)
func testEchoFileSanitized(c EchoContext) {
	filename := c.FormValue("doc")
	safeName := filepath.Base(filename)
	// ok: go-http-path-traversal
	_ = c.File(safeName)
}
