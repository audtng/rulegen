package rules

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return "untrusted"
}

func (c *GinContext) File(path string) {}

type EchoContext struct{}

func (c *EchoContext) FormValue(key string) string {
	return "untrusted"
}

func (c *EchoContext) File(path string) error {
	return nil
}

type DAGRequest struct {
	Name string `json:"name"`
}

// 1. Vulnerable: net/http URL query parameter joined to base directory and written directly.
func testVulnHTTPQueryFileWrite(r *http.Request) {
	filename := r.URL.Query().Get("file")
	dest := filepath.Join("/var/app/data", filename)
	// ruleid: path-traversal-arbitrary-file-write
	_ = os.WriteFile(dest, []byte("content"), 0644)
}

// 2. Safe: net/http parameter sanitized with filepath.Base.
func testSafeHTTPQueryBaseSanitized(r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	dest := filepath.Join("/var/app/data", safeName)
	// ok: path-traversal-arbitrary-file-write
	_ = os.WriteFile(dest, []byte("content"), 0644)
}

// 3. Vulnerable: JSON body unmarshaled file path joined and created without validation.
func testVulnJSONBodyFileCreate(body []byte) {
	var req DAGRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}
	dest := filepath.Join("/var/app/dags", req.Name)
	// ruleid: path-traversal-arbitrary-file-write
	f, err := os.Create(dest)
	if err != nil {
		return
	}
	defer f.Close()
}

// 4. Safe: JSON body checked with filepath.IsLocal before opening.
func testSafeJSONBodyIsLocalSanitized(body []byte) {
	var req DAGRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}
	if !filepath.IsLocal(req.Name) {
		return
	}
	dest := filepath.Join("/var/app/dags", req.Name)
	// ok: path-traversal-arbitrary-file-write
	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer f.Close()
}

// 5. Vulnerable: Gin context query parameter passed to framework file sink.
func testVulnGinContextFileServe(c *GinContext) {
	filename := c.Query("filename")
	dest := filepath.Join("/var/app/public", filename)
	// ruleid: path-traversal-arbitrary-file-write
	c.File(dest)
}

// 6. Safe: Echo context form parameter sanitized with filepath.Base before serving.
func testSafeEchoContextCleanedAndBase(c *EchoContext) {
	filename := c.FormValue("filename")
	safeName := filepath.Base(filename)
	dest := filepath.Join("/var/app/uploads", safeName)
	// ok: path-traversal-arbitrary-file-write
	_ = c.File(dest)
}
