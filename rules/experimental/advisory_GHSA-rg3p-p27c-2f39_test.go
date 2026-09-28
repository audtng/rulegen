package rules

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return ""
}

func (c *GinContext) File(filepath string) {}

type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string {
	return ""
}

func (c *EchoContext) File(file string) error {
	return nil
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string {
	return ""
}

func (c *FiberCtx) SendFile(file string, compress ...bool) error {
	return nil
}

// 1. Vulnerable net/http: request URL path joined into filesystem path
func handlerVulnerableURLPath(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	fullPath := filepath.Join("/var/www/uploads", p)
	// ruleid: go-http-path-traversal
	f, err := os.Open(fullPath)
	if err != nil {
		return
	}
	defer f.Close()
}

// 2. Vulnerable net/http: RequestURI unescaped and read directly
func handlerVulnerableRequestURI(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	unescaped, err := url.QueryUnescape(uri)
	if err != nil {
		return
	}
	relPath := strings.TrimPrefix(unescaped, "/")
	fullPath := filepath.Join("/var/www/html", relPath)
	// ruleid: go-http-path-traversal
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return
	}
	_ = data
}

// 3. Vulnerable Gin framework: user query param flowing directly to File sink
func handlerVulnerableGin(c *GinContext) {
	p := c.Query("filepath")
	// ruleid: go-http-path-traversal
	c.File(p)
}

// 4. Safe net/http: user query parameter sanitized via filepath.Base
func handlerSafeBaseSanitized(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	fullPath := filepath.Join("/var/www/uploads", safeName)
	// ok: go-http-path-traversal
	f, err := os.Open(fullPath)
	if err != nil {
		return
	}
	defer f.Close()
}

// 5. Safe net/http: path validated via filepath.IsLocal
func handlerSafeCheckLocal(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	cleanPath := filepath.Clean(p)
	if !filepath.IsLocal(cleanPath) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	// ok: go-http-path-traversal
	f, err := os.Open(cleanPath)
	if err != nil {
		return
	}
	defer f.Close()
}

// 6. Safe Echo framework: user form value stripped down via filepath.Base
func handlerSafeEchoBase(c *EchoContext) {
	raw := c.FormValue("doc")
	safe := filepath.Base(raw)
	// ok: go-http-path-traversal
	_ = c.File(safe)
}
