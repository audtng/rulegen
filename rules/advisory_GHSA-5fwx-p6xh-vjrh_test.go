package rules

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

// 1. Direct stdlib
func DirectStdlibTest(r *http.Request) {
	filename := r.URL.Query().Get("file")
	target := filepath.Join("/data/uploads", filename)
	// ruleid: path-traversal
	os.Open(target)
}

// 2. Proper patch
func ProperPatchTest(r *http.Request) {
	filename := r.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		return
	}
	target := filepath.Join("/data/uploads", filename)
	// ok: path-traversal
	os.Open(target)
}

func buildFilePath(dir, name string) string {
	return filepath.Join(dir, name)
}

// 3. Cross-function taint
func CrossFunctionTaintTest(r *http.Request) {
	filename := r.FormValue("file")
	target := buildFilePath("/data/uploads", filename)
	// ruleid: path-traversal
	os.Open(target)
}

// 4. Interface bypass
func InterfaceBypassTest(payload []byte) {
	var block map[string]interface{}
	if err := json.Unmarshal(payload, &block); err != nil {
		return
	}
	if fileID, ok := block["fileId"].(string); ok {
		// ruleid: path-traversal
		os.Open(fileID)
	}
}

// 5. Fake sanitizer
func FakeSanitizerTest(r *http.Request) {
	filename := r.URL.Query().Get("file")
	cleaned := filepath.Clean(filename)
	target := filepath.Join("/data/uploads", cleaned)
	// ruleid: path-traversal
	os.Open(target)
}

// 6. Real sanitizer
func RealSanitizerTest(r *http.Request) {
	filename := r.URL.Query().Get("file")
	safe := filepath.Base(filename)
	target := filepath.Join("/data/uploads", safe)
	// ok: path-traversal
	os.Open(target)
}
