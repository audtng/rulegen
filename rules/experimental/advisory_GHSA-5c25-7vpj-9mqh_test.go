package rules

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Mock framework contexts for environment-independent compilation
type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)     {}

type EchoContext struct{}

func (c *EchoContext) FormValue(key string) string { return "" }
func (c *EchoContext) File(file string) error      { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error                      { return nil }

// 1. Direct stdlib
func DirectStdlib(r *http.Request) {
	p := r.URL.Path
	// ruleid: go-path-traversal
	os.Open(p)
}

// 2. Proper patch
func ProperPatch(r *http.Request) {
	p := r.URL.Path
	if !filepath.IsLocal(p) {
		return
	}
	// ok: go-path-traversal
	os.Open(p)
}

func sanitizeHelper(p string) string {
	return p
}

// 3. Cross-function taint
func CrossFunctionTaint(r *http.Request) {
	p := sanitizeHelper(r.URL.Path)
	// ruleid: go-path-traversal
	os.Open(p)
}

// 4. Interface bypass
func InterfaceBypass(r *http.Request) {
	var val any = r.URL.Path
	if str, ok := val.(string); ok {
		// ruleid: go-path-traversal
		os.Open(str)
	}
}

// 5. Fake sanitizer
func FakeSanitizer(r *http.Request) {
	p := r.URL.Path
	if strings.HasPrefix(p, "/dashboard") {
		strip := strings.TrimPrefix(p, "/dashboard")
		local := path.Join("/var/www", strip)
		// ruleid: go-path-traversal
		os.Open(local)
	}
}

// 6. Real sanitizer
func RealSanitizer(r *http.Request) {
	p := filepath.Base(r.URL.Path)
	// ok: go-path-traversal
	os.Open(p)
}
