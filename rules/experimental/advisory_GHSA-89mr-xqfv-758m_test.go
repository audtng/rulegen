package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string { return "" }
func (c *EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string) string                      { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error { return nil }

func testVulnHTTPUpload(r *http.Request) error {
	filename := r.FormValue("filename")
	target := filepath.Join("/var/app/uploads", filename)
	if fi, err := os.Lstat(target); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	// ruleid: path-traversal
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func testSafeHTTPBase(r *http.Request) error {
	filename := r.FormValue("filename")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/app/uploads", safeName)
	// ok: path-traversal
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func testSafeHTTPIsLocal(r *http.Request) error {
	filename := r.URL.Query().Get("path")
	if !filepath.IsLocal(filename) {
		return nil
	}
	target := filepath.Join("/var/app/data", filename)
	// ok: path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func testVulnGin(c *GinContext) {
	filename := c.Query("name")
	target := filepath.Join("/public", filename)
	// ruleid: path-traversal
	c.File(target)
}

func testVulnEcho(c *EchoContext) error {
	filename := c.FormValue("file")
	target := filepath.Join("/tmp/storage", filename)
	// ruleid: path-traversal
	return os.RemoveAll(target)
}

func testSafeFiber(c *FiberCtx) error {
	filename := c.Query("file")
	safe := filepath.Base(filename)
	target := filepath.Join("/assets", safe)
	// ok: path-traversal
	return c.SendFile(target)
}
