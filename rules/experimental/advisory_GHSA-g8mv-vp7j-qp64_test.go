package rules

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

// Edge Case 1: Arbitrary File Write via PUT upload path (goshs archetype)
func testNetHttpPutUploadVuln(w http.ResponseWriter, req *http.Request) {
	upath := req.URL.Path
	filename := strings.Split(upath, "/")
	outName := filename[len(filename)-1]
	savepath := fmt.Sprintf("/var/uploads/%s", outName)
	// ruleid: go-path-traversal
	os.Create(savepath)
}

// Edge Case 2: Upload path sanitized using filepath.Base
func testNetHttpUploadSafeBase(w http.ResponseWriter, req *http.Request) {
	upath := req.URL.Path
	filename := filepath.Base(upath)
	savepath := filepath.Join("/var/uploads", filename)
	// ok: go-path-traversal
	os.Create(savepath)
}

// Edge Case 3: Guarded by filepath.IsLocal
func testNetHttpIsLocalSafe(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Query().Get("path")
	if !filepath.IsLocal(p) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	target := filepath.Join("/var/data", p)
	// ok: go-path-traversal
	os.OpenFile(target, os.O_RDWR, 0644)
}

// Edge Case 4: FormValue with filepath.Join into os.WriteFile
func testNetHttpFormValueJoinVuln(req *http.Request, data []byte) {
	userFile := req.FormValue("file")
	dest := filepath.Join("/var/storage", userFile)
	// ruleid: go-path-traversal
	os.WriteFile(dest, data, 0644)
}

// Edge Case 5: Gin framework context Query parameter directly into File sink
func testGinContextVuln(c *GinContext) {
	doc := c.Query("doc")
	target := filepath.Join("/static/docs", doc)
	// ruleid: go-path-traversal
	c.File(target)
}

// Edge Case 6: Traversal guard using strings.Contains check
func testNetHttpContainsSafe(req *http.Request) {
	name := req.URL.Query().Get("name")
	target := filepath.Join("/tmp/files", name)
	if strings.Contains(target, "..") {
		return
	}
	// ok: go-path-traversal
	os.Remove(target)
}
