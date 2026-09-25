package rules

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type DeleteRequest struct {
	Filename string `json:"filename"`
}

// Case 1: Vulnerable net/http handler with query parameter passed to filepath.Join and os.RemoveAll
func testNetHTTPQueryDeleteVuln(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	target := filepath.Join("/var/data", path)
	// ruleid: path-traversal-arbitrary-file-deletion
	os.RemoveAll(target)
}

// Case 2: Safe net/http handler sanitized with filepath.IsLocal
func testNetHTTPQueryDeleteSanitized(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if !filepath.IsLocal(path) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	target := filepath.Join("/var/data", path)
	// ok: path-traversal-arbitrary-file-deletion
	os.Remove(target)
}

// Case 3: Vulnerable JSON body unmarshaling passed to os.Remove
func testJSONUnmarshalDeleteVuln(body []byte) {
	var req DeleteRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}
	fullPath := filepath.Join("/uploads", req.Filename)
	// ruleid: path-traversal-arbitrary-file-deletion
	os.Remove(fullPath)
}

// Case 4: Safe JSON body unmarshaling sanitized with filepath.Base
func testJSONUnmarshalSafeBase(body []byte) {
	var req DeleteRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}
	safeName := filepath.Base(req.Filename)
	fullPath := filepath.Join("/uploads", safeName)
	// ok: path-traversal-arbitrary-file-deletion
	os.Remove(fullPath)
}

// Case 5: Vulnerable Gin framework handler with query parameter passed to c.File
func testGinContextVuln(c *GinContext) {
	userInput := c.Query("file")
	target := filepath.Join("/public", userInput)
	// ruleid: path-traversal-arbitrary-file-deletion
	c.File(target)
}

// Case 6: Safe static hardcoded path operation
func testStaticSafePath() {
	safePath := filepath.Join("/tmp", "cleanup.log")
	// ok: path-traversal-arbitrary-file-deletion
	os.RemoveAll(safePath)
}
