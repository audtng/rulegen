package rules

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c EchoContext) FormValue(key string) string { return "" }
func (c EchoContext) File(file string) error      { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error   { return nil }

type ExportPayload struct {
	Slug string `json:"slug"`
}

// 1. Vulnerable net/http handler with query param flowing into os.OpenFile
func VulnerableHTTPQuery(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/data/uploads", filename)
	// ruleid: go-path-traversal
	f, err := os.OpenFile(target, os.O_RDONLY, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
}

// 2. Safe net/http handler with filepath.IsLocal validation
func SafeHTTPQueryIsLocal(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	target := filepath.Join("/data/uploads", filename)
	// ok: go-path-traversal
	f, err := os.Open(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
}

// 3. Vulnerable JSON unmarshaling payload flowing into os.RemoveAll
func VulnerableJSONBody(body []byte) error {
	var payload ExportPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}
	target := filepath.Join("/base/export", payload.Slug)
	// ruleid: go-path-traversal
	return os.RemoveAll(target)
}

// 4. Safe JSON unmarshaling sanitized with filepath.Base
func SafeJSONBodyWithBase(body []byte) (*os.File, error) {
	var payload ExportPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	safeName := filepath.Base(payload.Slug)
	target := filepath.Join("/base/export", safeName)
	// ok: go-path-traversal
	return os.Create(target)
}

// 5. Vulnerable Gin framework context query flowing into c.File
func VulnerableGinContext(c *GinContext) {
	filename := c.Query("path")
	target := filepath.Join("/var/www/static", filename)
	// ruleid: go-path-traversal
	c.File(target)
}

// 6. Safe Fiber framework context query validated with fs.ValidPath
func SafeFiberContextValidPath(c *FiberCtx) error {
	filename := c.Query("path")
	if !fs.ValidPath(filename) {
		return nil
	}
	// ok: go-path-traversal
	return c.SendFile(filename)
}
