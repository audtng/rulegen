package rules

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) Param(key string) string { return "" }
func (c *GinContext) PostForm(key string) string { return "" }
func (c *GinContext) File(filepath string) {}
func (c *GinContext) BindJSON(obj any) error { return nil }

type EchoContext interface {
	FormValue(name string) string
	QueryParam(name string) string
	File(file string) error
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) Params(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string, compress ...bool) error { return nil }

type FileActionReq struct {
	Dir      string   `json:"dir"`
	Filename string   `json:"filename"`
	Names    []string `json:"names"`
}

// 1. Edge Case 1: Gin JSON unmarshaling into request struct followed by filepath.Join into os.RemoveAll
func testGinJsonUnmarshalPathTraversal(data []byte) {
	var req FileActionReq
	if err := json.Unmarshal(data, &req); err != nil {
		return
	}

	// Vulnerable: untrusted path field joined and passed to file deletion
	vulnPath := filepath.Join("/var/data", req.Filename)
	// ruleid: go-path-traversal
	os.RemoveAll(vulnPath)

	// Safe: sanitized using filepath.Base
	safeBase := filepath.Base(req.Filename)
	safePath := filepath.Join("/var/data", safeBase)
	// ok: go-path-traversal
	os.RemoveAll(safePath)
}

// 2. Edge Case 2: Standard net/http request with query param and filepath.Clean (false sanitizer) to os.ReadFile
func testNetHttpQueryReadFile(req *http.Request) {
	fileParam := req.URL.Query().Get("path")

	// Vulnerable: filepath.Clean does not restrict directory traversal outside base dir
	cleaned := filepath.Clean(filepath.Join("/app/public", fileParam))
	// ruleid: go-path-traversal
	_, _ = os.ReadFile(cleaned)

	// Safe: filepath.IsLocal ensures the path does not escape
	if !filepath.IsLocal(fileParam) {
		return
	}
	safePath := filepath.Join("/app/public", fileParam)
	// ok: go-path-traversal
	_, _ = os.ReadFile(safePath)
}

// 3. Edge Case 3: Echo framework FormValue source propagating to Echo File sink
func testEchoFormValueFileSink(c EchoContext) {
	doc := c.FormValue("doc")

	// Vulnerable: tainted form value passed directly into Echo File sink
	vulnTarget := filepath.Join("/srv/storage", doc)
	// ruleid: go-path-traversal
	_ = c.File(vulnTarget)

	// Safe: sanitized with filepath.Base
	safeTarget := filepath.Join("/srv/storage", filepath.Base(doc))
	// ok: go-path-traversal
	_ = c.File(safeTarget)
}

// 4. Edge Case 4: Fiber framework Query source to Fiber SendFile sink
func testFiberQuerySendFileSink(c *FiberCtx) {
	item := c.Query("item")

	// Vulnerable: Fiber query directly fed into SendFile sink
	vulnItem := filepath.Join("/assets", item)
	// ruleid: go-path-traversal
	_ = c.SendFile(vulnItem)

	// Safe: validated with filepath.IsLocal
	if !filepath.IsLocal(item) {
		return
	}
	safeItem := filepath.Join("/assets", item)
	// ok: go-path-traversal
	_ = c.SendFile(safeItem)
}

// 5. Edge Case 5: Gin query parameter in file move/rename (OpenList CVE archetype)
func testGinQueryRenameMove(c *GinContext) {
	targetName := c.Query("name")

	// Vulnerable: untrusted name joined to directory without validation before os.Rename
	vulnSrc := filepath.Join("/repo/uploads", targetName)
	// ruleid: go-path-traversal
	_ = os.Rename(vulnSrc, "/repo/archive/file.dat")

	// Safe: Base name extracted to prevent directory climbing
	safeSrc := filepath.Join("/repo/uploads", filepath.Base(targetName))
	// ok: go-path-traversal
	_ = os.Rename(safeSrc, "/repo/archive/file.dat")
}

// 6. Edge Case 6: Standard net/http POST form value used in os.OpenFile
func testNetHttpPostOpenFile(req *http.Request) {
	name := req.PostFormValue("filename")

	// Vulnerable: user input joined to directory and opened for write
	vulnPath := filepath.Join("/var/log/app", name)
	// ruleid: go-path-traversal
	f, _ := os.OpenFile(vulnPath, os.O_CREATE|os.O_WRONLY, 0600)
	if f != nil {
		f.Close()
	}

	// Safe: sanitized via filepath.IsLocal check
	if !filepath.IsLocal(name) {
		return
	}
	safePath := filepath.Join("/var/log/app", name)
	// ok: go-path-traversal
	f2, _ := os.OpenFile(safePath, os.O_CREATE|os.O_WRONLY, 0600)
	if f2 != nil {
		f2.Close()
	}
}
