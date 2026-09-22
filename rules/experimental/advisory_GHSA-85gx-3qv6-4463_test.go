package rules

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string              { return "" }
func (c *GinContext) Param(key string) string              { return "" }
func (c *GinContext) File(filepath string)                 {}
func (c *GinContext) FileAttachment(filepath, name string) {}

type EchoContext struct{}

func (c EchoContext) FormValue(name string) string  { return "" }
func (c EchoContext) QueryParam(name string) string { return "" }
func (c EchoContext) File(file string) error        { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, args ...string) string      { return "" }
func (c *FiberCtx) Params(key string, args ...string) string     { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error { return nil }

// Edge case 1: Standard net/http request with Query parameter directly opened
func testNetHTTPQuery(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	// ruleid: go-path-traversal
	os.Open(filename)

	safe := filepath.Base(filename)
	// ok: go-path-traversal
	os.Open(safe)
}

// Edge case 2: Gin framework context Query parameter with IsLocal validation
func testGinQuery(c *GinContext) {
	path := c.Query("filepath")
	// ruleid: go-path-traversal
	c.File(path)

	if filepath.IsLocal(path) {
		// ok: go-path-traversal
		c.File(path)
	}
}

// Edge case 3: Echo framework FormValue propagated through filepath.Join
func testEchoFormValue(c EchoContext) {
	dir := c.FormValue("dir")
	joined := filepath.Join("/var/data", dir)
	// ruleid: go-path-traversal
	c.File(joined)

	safeName := filepath.Base(dir)
	safeJoined := filepath.Join("/var/data", safeName)
	// ok: go-path-traversal
	c.File(safeJoined)
}

// Edge case 4: Fiber framework route param formatted into path via fmt.Sprintf
func testFiberParams(c *FiberCtx) {
	path := c.Params("path")
	target := fmt.Sprintf("/uploads/%s", path)
	// ruleid: go-path-traversal
	c.SendFile(target)

	if filepath.IsLocal(path) {
		safeTarget := fmt.Sprintf("/uploads/%s", path)
		// ok: go-path-traversal
		c.SendFile(safeTarget)
	}
}

// Edge case 5: net/http POST form value constructed via strings.Builder
func testNetHTTPPostForm(w http.ResponseWriter, r *http.Request) {
	resource := r.PostFormValue("resource")
	var b strings.Builder
	b.WriteString("/srv/assets/")
	b.WriteString(resource)
	// ruleid: go-path-traversal
	os.ReadFile(b.String())

	safe := filepath.Base(resource)
	var b2 strings.Builder
	b2.WriteString("/srv/assets/")
	b2.WriteString(safe)
	// ok: go-path-traversal
	os.ReadFile(b2.String())
}

// Edge case 6: net/http URL.Path dispatched to http.ServeFile
func testNetHTTPRequestPath(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.Path
	// ruleid: go-path-traversal
	http.ServeFile(w, r, rawPath)

	safe := filepath.Base(rawPath)
	// ok: go-path-traversal
	http.ServeFile(w, r, safe)
}
