package rules

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext interface {
	FormValue(name string) string
	File(file string) error
}

func CVEVariantVulnerable(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path
	urlpath = strings.Trim(urlpath, "/")
	if urlpath != "" && path.Clean(urlpath) != urlpath {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	target := filepath.Join("/var/data", urlpath)
	// ruleid: go-path-traversal
	os.Open(target)
}

func CVEVariantPatched(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path
	urlpath = strings.Trim(urlpath, "/")
	if urlpath == "." {
		return
	}
	if !fs.ValidPath(urlpath) {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	target := filepath.Join("/var/data", urlpath)
	// ok: go-path-traversal
	os.Open(target)
}

func QueryParamSanitizedWithBase(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	safeName := filepath.Base(filename)
	target := filepath.Join("/uploads", safeName)
	// ok: go-path-traversal
	os.ReadFile(target)
}

func QueryParamVulnerableServeFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/public", filename)
	// ruleid: go-path-traversal
	http.ServeFile(w, r, target)
}

func GinVulnerableFile(c *GinContext) {
	filename := c.Query("file")
	target := filepath.Join("/static", filename)
	// ruleid: go-path-traversal
	c.File(target)
}

func EchoPatchedFile(c EchoContext) {
	filename := c.FormValue("doc")
	if !filepath.IsLocal(filename) {
		return
	}
	target := filepath.Join("/documents", filename)
	// ok: go-path-traversal
	c.File(target)
}
