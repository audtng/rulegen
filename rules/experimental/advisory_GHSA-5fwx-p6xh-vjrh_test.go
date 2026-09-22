package rules

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c EchoContext) FormValue(key string) string { return "" }
func (c EchoContext) File(filepath string) error  { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error                      { return nil }

type BoardBlock struct {
	FileID string `json:"fileId"`
}

// Edge case 1: Deserialized JSON payload (Mattermost board import vulnerability archetype)
func testJSONImportEdgeCase(data []byte) {
	var block BoardBlock
	_ = json.Unmarshal(data, &block)

	badPath := filepath.Join("/data/uploads", block.FileID)
	// ruleid: go-path-traversal
	_, _ = os.ReadFile(badPath)

	safeID := filepath.Base(block.FileID)
	safePath := filepath.Join("/data/uploads", safeID)
	// ok: go-path-traversal
	_, _ = os.ReadFile(safePath)
}

// Edge case 2: net/http URL query parameter with os.Open
func testNetHTTPQueryEdgeCase(r *http.Request) {
	filename := r.URL.Query().Get("file")

	vulnPath := filepath.Join("/static", filename)
	// ruleid: go-path-traversal
	_, _ = os.Open(vulnPath)

	if filepath.IsLocal(filename) {
		safePath := filepath.Join("/static", filename)
		// ok: go-path-traversal
		_, _ = os.Open(safePath)
	}
}

// Edge case 3: net/http POST form parameter with os.RemoveAll
func testNetHTTPPostFormEdgeCase(w http.ResponseWriter, r *http.Request) {
	targetFile := r.PostFormValue("filename")

	vulnPath := filepath.Join("/tmp/app", targetFile)
	// ruleid: go-path-traversal
	_ = os.RemoveAll(vulnPath)

	cleanFile := filepath.Base(targetFile)
	safePath := filepath.Join("/tmp/app", cleanFile)
	// ok: go-path-traversal
	_ = os.RemoveAll(safePath)
}

// Edge case 4: Gin Framework context query source and c.File sink
func testGinFrameworkEdgeCase(c *GinContext) {
	userInput := c.Query("path")

	vulnPath := filepath.Join("/assets", userInput)
	// ruleid: go-path-traversal
	c.File(vulnPath)

	safeInput := filepath.Base(userInput)
	safePath := filepath.Join("/assets", safeInput)
	// ok: go-path-traversal
	c.File(safePath)
}

// Edge case 5: Echo Framework context FormValue source and c.File sink
func testEchoFrameworkEdgeCase(c EchoContext) {
	docName := c.FormValue("doc")

	vulnPath := filepath.Join("/documents", docName)
	// ruleid: go-path-traversal
	_ = c.File(vulnPath)

	safeDoc := filepath.Base(docName)
	safePath := filepath.Join("/documents", safeDoc)
	// ok: go-path-traversal
	_ = c.File(safePath)
}

// Edge case 6: Fiber Framework context query source and c.SendFile sink
func testFiberFrameworkEdgeCase(c *FiberCtx) {
	attachment := c.Query("attachment")

	vulnPath := filepath.Join("/storage", attachment)
	// ruleid: go-path-traversal
	_ = c.SendFile(vulnPath)

	if filepath.IsLocal(attachment) {
		safePath := filepath.Join("/storage", attachment)
		// ok: go-path-traversal
		_ = c.SendFile(safePath)
	}
}
