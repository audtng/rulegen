package rules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type GinContext struct{}

func (g *GinContext) Query(key string) string { return "" }
func (g *GinContext) File(filepath string)     {}

type FiberCtx struct{}

func (f *FiberCtx) Query(key string) string  { return "" }
func (f *FiberCtx) SendFile(filepath string) {}

type ExecutionRequest struct {
	TrackingID string `json:"uniqueTrackingId"`
}

// Case 1: Untrusted HTTP query parameter interpolated into path and written to disk
func testCase1VulnerableHttpWrite(req *http.Request) {
	trackingID := req.URL.Query().Get("tracking_id")
	filename := fmt.Sprintf("exec.%s.log", trackingID)
	logPath := path.Join("/var/log/app", filename)
	// ruleid: go-path-traversal
	os.WriteFile(logPath, []byte("log data"), 0600)
}

// Case 2: Untrusted HTTP query parameter sanitized using filepath.Base before write
func testCase2SafeHttpBase(req *http.Request) {
	trackingID := req.URL.Query().Get("tracking_id")
	safeName := filepath.Base(trackingID)
	logPath := filepath.Join("/var/log/app", safeName)
	// ok: go-path-traversal
	os.WriteFile(logPath, []byte("log data"), 0600)
}

// Case 3: Untrusted HTTP query parameter sanitized with filepath.IsLocal validation
func testCase3SafeHttpIsLocal(req *http.Request) {
	target := req.URL.Query().Get("file")
	if !filepath.IsLocal(target) {
		return
	}
	resolvedPath := filepath.Join("/data/files", target)
	// ok: go-path-traversal
	os.OpenFile(resolvedPath, os.O_RDWR, 0600)
}

// Case 4: Untrusted JSON unmarshaled payload tracking ID joined and created on filesystem
func testCase4VulnerableJsonUnmarshal(data []byte) {
	var execReq ExecutionRequest
	if err := json.Unmarshal(data, &execReq); err != nil {
		return
	}
	targetPath := filepath.Join("/data/reports", execReq.TrackingID)
	// ruleid: go-path-traversal
	os.Create(targetPath)
}

// Case 5: Untrusted web framework query parameter passed directly to File serving sink
func testCase5VulnerableGinFile(ctx *GinContext) {
	doc := ctx.Query("doc")
	target := filepath.Join("/static/docs", doc)
	// ruleid: go-path-traversal
	ctx.File(target)
}

// Case 6: Untrusted web framework query parameter sanitized with filepath.Base before SendFile
func testCase6SafeFiberSendFile(ctx *FiberCtx) {
	doc := ctx.Query("doc")
	safeDoc := filepath.Base(doc)
	target := filepath.Join("/static/docs", safeDoc)
	// ok: go-path-traversal
	ctx.SendFile(target)
}
