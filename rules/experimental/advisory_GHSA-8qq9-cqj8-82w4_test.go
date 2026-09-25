package rules

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (g *GinContext) Query(key string) string {
	return ""
}

func (g *GinContext) File(filepath string) {}

type EchoContext struct{}

func (e *EchoContext) FormValue(name string) string {
	return ""
}

func (e *EchoContext) File(file string) error {
	return nil
}

type SyncFileInfo struct {
	Name string `json:"name"`
}

func vulnHTTPQueryOpenFile(r *http.Request) {
	name := r.URL.Query().Get("file")
	target := filepath.Join("/data/uploads", name)
	// ruleid: go-path-traversal
	os.OpenFile(target, os.O_RDWR, 0644)
}

func safeHTTPQueryBaseSanitized(r *http.Request) {
	name := r.URL.Query().Get("file")
	safeName := filepath.Base(name)
	target := filepath.Join("/data/uploads", safeName)
	// ok: go-path-traversal
	os.OpenFile(target, os.O_RDWR, 0644)
}

func vulnJSONSyncPayload(data []byte) {
	var info SyncFileInfo
	json.Unmarshal(data, &info)
	target := filepath.Join("/data/uploads", info.Name)
	// ruleid: go-path-traversal
	os.Create(target)
}

func safeJSONSyncIsLocal(data []byte) {
	var info SyncFileInfo
	json.Unmarshal(data, &info)
	if !filepath.IsLocal(info.Name) {
		return
	}
	target := filepath.Join("/data/uploads", info.Name)
	// ok: go-path-traversal
	os.Create(target)
}

func vulnGinServeFile(c *GinContext) {
	param := c.Query("path")
	target := filepath.Join("/static", param)
	// ruleid: go-path-traversal
	c.File(target)
}

func safeEchoServeFile(c *EchoContext) {
	name := c.FormValue("name")
	safeName := filepath.Base(name)
	target := filepath.Join("/static", safeName)
	// ok: go-path-traversal
	c.File(target)
}
