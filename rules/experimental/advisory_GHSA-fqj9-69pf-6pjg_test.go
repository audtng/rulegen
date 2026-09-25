package rules

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string { return "" }
func (c *EchoContext) File(file string) error       { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error     { return nil }

// Edge Case 1: net/http handler with path traversal via URL.Path and filepath.Join
func TestNetHttpJoinVulnerable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test/../secret", nil)
	urlpath := strings.Trim(req.URL.Path, "/")
	cleanPath := filepath.Clean(urlpath)
	fullPath := filepath.Join("/base/dir", cleanPath)
	// ruleid: go-path-traversal-http
	f, err := os.Open(fullPath)
	if err != nil {
		return
	}
	defer f.Close()
}

// Edge Case 2: net/http handler patched with canonical path validation (CVE archetype)
func TestNetHttpPatchedCanonical(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test/../secret", nil)
	urlpath := strings.Trim(req.URL.Path, "/")
	if path.Clean(urlpath) != urlpath {
		return
	}
	fullPath := filepath.Join("/base/dir", urlpath)
	// ok: go-path-traversal-http
	f, err := os.Open(fullPath)
	if err != nil {
		return
	}
	defer f.Close()
}

// Edge Case 3: net/http handler sanitized with filepath.Base
func TestNetHttpQuerySanitizedBase(t *testing.T) {
	req, _ := http.NewRequest("GET", "/download?file=../../etc/passwd", nil)
	filename := req.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	fullPath := filepath.Join("/safe/dir", safeName)
	// ok: go-path-traversal-http
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return
	}
	_ = data
}

// Edge Case 4: Gin framework handler vulnerable to path traversal
func TestGinFileVulnerable(t *testing.T) {
	c := &GinContext{}
	target := c.Query("file")
	fullPath := filepath.Join("/static/assets", target)
	// ruleid: go-path-traversal-http
	c.File(fullPath)
}

// Edge Case 5: Echo framework handler vulnerable to path traversal in RemoveAll
func TestEchoRemoveAllVulnerable(t *testing.T) {
	c := &EchoContext{}
	repo := c.FormValue("repo")
	targetPath := fmt.Sprintf("/var/data/%s", repo)
	// ruleid: go-path-traversal-http
	os.RemoveAll(targetPath)
}

// Edge Case 6: Fiber framework handler safe via filepath.IsLocal validation
func TestFiberSendFilePatchedIsLocal(t *testing.T) {
	c := &FiberCtx{}
	filename := c.Query("file")
	if !filepath.IsLocal(filename) {
		return
	}
	fullPath := filepath.Join("/safe/data", filename)
	// ok: go-path-traversal-http
	c.SendFile(fullPath)
}
