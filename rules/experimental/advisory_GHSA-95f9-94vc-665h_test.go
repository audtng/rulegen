package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)     {}

type EchoContext struct{}

func (c EchoContext) FormValue(name string) string { return "" }
func (c EchoContext) File(file string) error       { return nil }

// Edge Case 1: Vulnerable - Direct HTTP URL Path in os.Open
func handleHttpURLPathVuln(req *http.Request) {
	p := req.URL.Path
	// ruleid: path-traversal
	os.Open(p)
}

// Edge Case 2: Safe - HTTP URL Path sanitized with filepath.Base
func handleHttpURLPathSanitized(req *http.Request) {
	p := req.URL.Path
	safe := filepath.Base(p)
	// ok: path-traversal
	os.Open(safe)
}

// Edge Case 3: Vulnerable - Query parameter joined to base directory in http.ServeFile
func handleHttpQueryParamVuln(w http.ResponseWriter, req *http.Request) {
	file := req.URL.Query().Get("file")
	target := filepath.Join("/var/www/static", file)
	// ruleid: path-traversal
	http.ServeFile(w, req, target)
}

// Edge Case 4: Safe - Form value validated with filepath.IsLocal
func handleHttpFormValueValidated(req *http.Request) {
	filename := req.FormValue("doc")
	if filepath.IsLocal(filename) {
		// ok: path-traversal
		os.ReadFile(filename)
	}
}

// Edge Case 5: Vulnerable - Gin Context Query used directly in File sink
func handleGinQueryVuln(c *GinContext) {
	doc := c.Query("doc")
	// ruleid: path-traversal
	c.File(doc)
}

// Edge Case 6: Safe - Echo Context FormValue sanitized with filepath.Base
func handleEchoFormValueSanitized(c EchoContext) {
	doc := c.FormValue("doc")
	clean := filepath.Base(doc)
	// ok: path-traversal
	c.File(clean)
}
