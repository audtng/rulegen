package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return key
}

func (c *GinContext) ShouldBindJSON(obj any) error {
	return nil
}

type EchoContext struct{}

func (c *EchoContext) FormValue(key string) string {
	return key
}

func (c *EchoContext) File(file string) error {
	return nil
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string {
	return key
}

func (c *FiberCtx) SendFile(file string) error {
	return nil
}

type DeleteChunkRequest struct {
	FilePath string
}

// Edge case 1: Gin Query input directly reaches os.Remove (vulnerable)
func vulnerableGinQueryRemove(c *GinContext) {
	filename := c.Query("file")
	// ruleid: arbitrary-file-deletion
	os.Remove(filename)
}

// Edge case 2: Gin JSON unmarshaling into struct, joined via filepath.Join, reaches os.RemoveAll (vulnerable)
func vulnerableGinBindRemoveAll(c *GinContext) {
	var req DeleteChunkRequest
	_ = c.ShouldBindJSON(&req)
	target := filepath.Join("/var/app/chunks", req.FilePath)
	// ruleid: arbitrary-file-deletion
	os.RemoveAll(target)
}

// Edge case 3: Gin Query sanitized via filepath.Base before removal (safe)
func safeGinBaseClean(c *GinContext) {
	filename := c.Query("file")
	safeName := filepath.Base(filename)
	// ok: arbitrary-file-deletion
	os.Remove(filepath.Join("/tmp/chunks", safeName))
}

// Edge case 4: Echo FormValue reaches c.File framework sink (vulnerable)
func vulnerableEchoFile(c *EchoContext) {
	doc := c.FormValue("doc")
	// ruleid: arbitrary-file-deletion
	_ = c.File(doc)
}

// Edge case 5: Fiber Query reaches Fiber SendFile framework sink (vulnerable)
func vulnerableFiberSendFile(c *FiberCtx) {
	page := c.Query("page")
	// ruleid: arbitrary-file-deletion
	_ = c.SendFile(page)
}

// Edge case 6: Net/HTTP FormValue validated with filepath.IsLocal before removal (safe)
func safeHttpIsLocal(w http.ResponseWriter, req *http.Request) {
	target := req.FormValue("target")
	if !filepath.IsLocal(target) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	// ok: arbitrary-file-deletion
	os.Remove(filepath.Join("/data/store", target))
}
