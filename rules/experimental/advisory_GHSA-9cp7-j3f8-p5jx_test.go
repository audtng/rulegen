package rules

import (
	"archive/zip"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
func testGinContextFile(c *GinContext) {
	path := c.Query("filepath")
	// ruleid: go-path-traversal
	c.File(path)

	if filepath.IsLocal(path) {
		// ok: go-path-traversal
		c.File(path)
	}
}

// Edge case 3: Echo framework FormValue propagated through filepath.Join
func testEchoContextFile(c EchoContext) {
	name := c.FormValue("name")
	joined := filepath.Join("/tmp", name)
	// ruleid: go-path-traversal
	c.File(joined)

	safe := filepath.Base(name)
	safeJoined := filepath.Join("/tmp", safe)
	// ok: go-path-traversal
	c.File(safeJoined)
}

// Edge case 4: Fiber framework route param formatted into path via fmt.Sprintf
func testFiberContextSendFile(c *FiberCtx) {
	p := c.Params("path")
	target := fmt.Sprintf("/var/www/%s", p)
	// ruleid: go-path-traversal
	c.SendFile(target)

	if filepath.IsLocal(p) {
		safeTarget := fmt.Sprintf("/var/www/%s", p)
		// ok: go-path-traversal
		c.SendFile(safeTarget)
	}
}

// Edge case 5: Zip archive extraction (Zip Slip) from zip.File entry Name
func testZipArchiveExtraction(f *zip.File, targetDir string) {
	dest := filepath.Join(targetDir, f.Name)
	// ruleid: go-path-traversal
	out, _ := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if out != nil {
		out.Close()
	}

	safe := filepath.Base(f.Name)
	safeDest := filepath.Join(targetDir, safe)
	// ok: go-path-traversal
	safeOut, _ := os.OpenFile(safeDest, os.O_CREATE|os.O_WRONLY, 0644)
	if safeOut != nil {
		safeOut.Close()
	}
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
