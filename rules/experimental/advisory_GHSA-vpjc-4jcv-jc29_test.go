package rules

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

type ActionRequest struct {
	TargetFile string `json:"target_file"`
}

func testUnmarshalDirectVuln(data []byte) {
	var req ActionRequest
	_ = json.Unmarshal(data, &req)
	// ruleid: path-traversal
	os.ReadFile(req.TargetFile)
}

func testUnmarshalSanitizedBaseOk(data []byte) {
	var req ActionRequest
	_ = json.Unmarshal(data, &req)
	clean := filepath.Base(req.TargetFile)
	// ok: path-traversal
	os.ReadFile(clean)
}

func testUnmarshalJoinVuln(data []byte) {
	var req ActionRequest
	_ = json.Unmarshal(data, &req)
	fullPath := filepath.Join("/var/app/data", req.TargetFile)
	// ruleid: path-traversal
	os.WriteFile(fullPath, []byte("data"), 0600)
}

func testUnmarshalIsLocalOk(data []byte) {
	var req ActionRequest
	_ = json.Unmarshal(data, &req)
	if !filepath.IsLocal(req.TargetFile) {
		return
	}
	fullPath := filepath.Join("/var/app/data", req.TargetFile)
	// ok: path-traversal
	os.WriteFile(fullPath, []byte("data"), 0600)
}

func testBufioReaderJoinVuln(r *bufio.Reader) {
	entry, _ := r.ReadString('\n')
	targetPath := filepath.Join("/var/spool", entry)
	// ruleid: path-traversal
	os.RemoveAll(targetPath)
}

func testBufioReaderSanitizedBaseOk(r *bufio.Reader) {
	entry, _ := r.ReadString('\n')
	safeEntry := filepath.Base(entry)
	targetPath := filepath.Join("/var/spool", safeEntry)
	// ok: path-traversal
	os.RemoveAll(targetPath)
}
