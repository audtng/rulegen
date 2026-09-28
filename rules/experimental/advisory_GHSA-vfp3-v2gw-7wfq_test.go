package rules

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext interface {
	FormValue(key string) string
	File(file string) error
}

func testNetHTTPDirectOpenVuln(req *http.Request) {
	p := req.URL.Query().Get("file")
	fullPath := filepath.Join("/var/data", p)
	// ruleid: path-traversal
	os.Open(fullPath)
}

func testNetHTTPSanitizedWithBaseOk(req *http.Request) {
	p := req.URL.Query().Get("file")
	safeName := filepath.Base(p)
	safePath := filepath.Join("/var/data", safeName)
	// ok: path-traversal
	os.Open(safePath)
}

func testNetHTTPSanitizedWithIsLocalOk(req *http.Request) {
	p := req.URL.Query().Get("file")
	if !filepath.IsLocal(p) {
		return
	}
	// ok: path-traversal
	os.Open(p)
}

func testURLUnescapeServeFileVuln(w http.ResponseWriter, req *http.Request) {
	raw := req.URL.Path
	unescaped, err := url.PathUnescape(raw)
	if err != nil {
		return
	}
	cleaned := filepath.Clean(unescaped)
	// ruleid: path-traversal
	http.ServeFile(w, req, cleaned)
}

func testGinContextFileVuln(c *GinContext) {
	p := c.Query("doc")
	// ruleid: path-traversal
	c.File(p)
}

func testEchoContextSanitizedPrefixOk(c EchoContext) {
	p := c.FormValue("doc")
	target := filepath.Join("/static", p)
	if !strings.HasPrefix(target, "/static/") {
		return
	}
	// ok: path-traversal
	c.File(target)
}
