package rules

import (
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(s string) string { return "" }
func (c *GinContext) File(filepath string)  {}

type EchoContext interface {
	FormValue(string) string
	File(string) error
}

type FiberCtx struct{}

func (f *FiberCtx) Query(s string) string      { return "" }
func (f *FiberCtx) SendFile(file string) error { return nil }

// Edge Case 1: Untrusted multipart file upload filename used in file write without validation
func sampleMultipartUploadVuln(fh *multipart.FileHeader, baseDir string) {
	target := filepath.Join(baseDir, fh.Filename)
	// ruleid: go-path-traversal-file-write
	os.WriteFile(target, []byte("data"), 0644)
}

// Edge Case 2: Untrusted multipart file upload filename sanitized with filepath.Base
func sampleMultipartUploadSafeBase(fh *multipart.FileHeader, baseDir string) {
	cleanName := filepath.Base(fh.Filename)
	target := filepath.Join(baseDir, cleanName)
	// ok: go-path-traversal-file-write
	os.WriteFile(target, []byte("data"), 0644)
}

// Edge Case 3: Untrusted HTTP form value sanitized via filepath.IsLocal check
func sampleHttpRequestSafeIsLocal(r *http.Request, baseDir string) {
	param := r.FormValue("path")
	if !filepath.IsLocal(param) {
		return
	}
	target := filepath.Join(baseDir, param)
	// ok: go-path-traversal-file-write
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
}

// Edge Case 4: Gin query parameter joined with base directory and served via c.File
func sampleGinFileVuln(c *GinContext, baseDir string) {
	userPath := c.Query("path")
	fullPath := filepath.Join(baseDir, userPath)
	// ruleid: go-path-traversal-file-write
	c.File(fullPath)
}

// Edge Case 5: Echo form value joined with base directory and deleted via os.RemoveAll
func sampleEchoFileVuln(c EchoContext, baseDir string) {
	fileName := c.FormValue("name")
	target := filepath.Join(baseDir, fileName)
	// ruleid: go-path-traversal-file-write
	os.RemoveAll(target)
}

// Edge Case 6: Fiber query parameter sanitized via filepath.Base before sending file
func sampleFiberSendFileSafe(c *FiberCtx, baseDir string) {
	cleanName := filepath.Base(c.Query("name"))
	target := filepath.Join(baseDir, cleanName)
	// ok: go-path-traversal-file-write
	_ = c.SendFile(target)
}
