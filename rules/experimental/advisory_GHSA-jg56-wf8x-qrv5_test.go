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

func (c EchoContext) FormValue(name string) string {
	return ""
}

func (c EchoContext) File(file string) error {
	return nil
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string {
	return ""
}

func (c *FiberCtx) SendFile(file string, compress ...bool) error {
	return nil
}

func VulnMultipartUpload(r *http.Request) {
	_, hdr, err := r.FormFile("upload")
	if err != nil {
		return
	}
	target := filepath.Join("/var/uploads", hdr.Filename)
	// ruleid: go-path-traversal
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer f.Close()
}

func SafeMultipartUpload(r *http.Request) {
	_, hdr, err := r.FormFile("upload")
	if err != nil {
		return
	}
	cleanName := filepath.Base(hdr.Filename)
	target := filepath.Join("/var/uploads", cleanName)
	// ok: go-path-traversal
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer f.Close()
}

func VulnURLPathTraversal(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	fullPath := filepath.Join("/data/site", upath)
	// ruleid: go-path-traversal
	err := os.RemoveAll(fullPath)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}
}

func SafeURLPathIsLocal(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Path
	if !filepath.IsLocal(relPath) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	fullPath := filepath.Join("/data/site", relPath)
	// ok: go-path-traversal
	f, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	defer f.Close()
}

func VulnGinServeFile(c *GinContext) {
	filePath := c.Query("file")
	dest := filepath.Join("/assets", filePath)
	// ruleid: go-path-traversal
	c.File(dest)
}

func SafeFiberSendFile(c *FiberCtx) {
	userInput := c.Query("doc")
	safeName := filepath.Base(userInput)
	fullPath := filepath.Join("/docs", safeName)
	// ok: go-path-traversal
	_ = c.SendFile(fullPath)
}
