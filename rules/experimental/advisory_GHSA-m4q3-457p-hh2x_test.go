package rules

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

// GinContext mocks the gin framework context
type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)     {}

// EchoContext mocks the echo framework context
type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string { return "" }
func (c *EchoContext) File(file string) error       { return nil }

// FiberCtx mocks the fiber framework context
type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error                       { return nil }

type DAGRunRequest struct {
	DagRunID string `json:"dagRunId"`
}

func testHttpQueryParamRemoveAll_Vuln(r *http.Request) {
	runID := r.URL.Query().Get("dagRunId")
	dir := filepath.Join("/tmp/dags", runID)
	// ruleid: go-path-traversal
	os.RemoveAll(dir)
}

func testHttpQueryParamBase_Safe(r *http.Request) {
	runID := r.URL.Query().Get("dagRunId")
	safeID := filepath.Base(runID)
	dir := filepath.Join("/tmp/dags", safeID)
	// ok: go-path-traversal
	os.RemoveAll(dir)
}

func testHttpFormValueIsLocal_Safe(r *http.Request) {
	filePath := r.FormValue("path")
	if !filepath.IsLocal(filePath) {
		return
	}
	target := filepath.Join("/var/data", filePath)
	// ok: go-path-traversal
	os.OpenFile(target, os.O_RDONLY, 0)
}

func testJsonUnmarshalOpenFile_Vuln(data []byte) {
	var req DAGRunRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return
	}
	target := filepath.Join("/var/runs", req.DagRunID)
	// ruleid: go-path-traversal
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0600)
}

func testGinFrameworkFile_Vuln(c *GinContext) {
	runID := c.Query("dagRunId")
	target := filepath.Join("/tmp/executions", runID)
	// ruleid: go-path-traversal
	c.File(target)
}

func testEchoFrameworkFile_Safe(c *EchoContext) {
	runID := c.FormValue("dagRunId")
	safeID := filepath.Base(runID)
	target := filepath.Join("/tmp/executions", safeID)
	// ok: go-path-traversal
	_ = c.File(target)
}
