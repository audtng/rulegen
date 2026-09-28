package rules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext interface {
	FormValue(key string) string
	File(file string) error
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error    { return nil }

type Payload struct {
	Filename string `json:"filename"`
}

// 1. Direct HTTP Query to os.Open (Vulnerable)
func EdgeCase1_HTTPDirect(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	// ruleid: path-traversal-arbitrary-file-read
	os.Open(filename)
}

// 2. HTTP FormValue with filepath.Join and filepath.Base (Safe)
func EdgeCase2_HTTPSafeBase(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("file")
	safeName := filepath.Base(filename)
	fullPath := filepath.Join("/var/www/uploads", safeName)
	// ok: path-traversal-arbitrary-file-read
	os.Open(fullPath)
}

// 3. Gin Query to Gin File Sink with Sprintf propagation (Vulnerable)
func EdgeCase3_GinVuln(c *GinContext) {
	rawFile := c.Query("path")
	target := fmt.Sprintf("/data/%s", rawFile)
	// ruleid: path-traversal-arbitrary-file-read
	c.File(target)
}

// 4. Echo FormValue with filepath.IsLocal sanitizer (Safe)
func EdgeCase4_EchoSafeIsLocal(c EchoContext) {
	path := c.FormValue("doc")
	if !filepath.IsLocal(path) {
		return
	}
	// ok: path-traversal-arbitrary-file-read
	c.File(path)
}

// 5. JSON Unmarshal from request data to os.ReadFile (Vulnerable)
func EdgeCase5_JSONUnmarshalVuln(data []byte) {
	var payload Payload
	_ = json.Unmarshal(data, &payload)
	target := filepath.Join("/storage", payload.Filename)
	// ruleid: path-traversal-arbitrary-file-read
	os.ReadFile(target)
}

// 6. Fiber Query to Fiber SendFile with filepath.Base (Safe)
func EdgeCase6_FiberSafe(c *FiberCtx) {
	file := c.Query("file")
	safeFile := filepath.Base(file)
	// ok: path-traversal-arbitrary-file-read
	c.SendFile(safeFile)
}
