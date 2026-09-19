package rules

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var validIDPattern = regexp.MustCompile(`^[a-fA-F0-9\-]+$`)

// 1. Direct stdlib: direct HTTP input into file sink
func directStdlibTest(r *http.Request) {
	filename := r.URL.Query().Get("filename")
	target := filepath.Join("/var/log", filename)
	// ruleid: path-traversal-file-write
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0600)
}

// 2. Proper patch: strict regex validation preventing traversal characters
func properPatchTest(r *http.Request) {
	trackingID := r.URL.Query().Get("tracking_id")
	if !validIDPattern.MatchString(trackingID) {
		return
	}
	target := filepath.Join("/var/log", trackingID)
	// ok: path-traversal-file-write
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0600)
}

// 3. Cross-function taint: input routed through helper function
func formatPath(id string) string {
	return filepath.Join("/var/log", id)
}

func crossFunctionTaintTest(r *http.Request) {
	fileID := r.URL.Query().Get("file_id")
	target := formatPath(fileID)
	// ruleid: path-traversal-file-write
	os.Create(target)
}

// 4. Interface bypass: input wrapped in interface/any before usage
type FileHolder struct {
	Path any
}

func interfaceBypassTest(r *http.Request) {
	name := r.URL.Query().Get("name")
	holder := FileHolder{Path: name}
	target := filepath.Join("/var/log", holder.Path.(string))
	// ruleid: path-traversal-file-write
	os.Create(target)
}

// 5. Fake sanitizer: naive string replacement that can be bypassed
func fakeSanitizerTest(r *http.Request) {
	filename := r.URL.Query().Get("filename")
	unsafeClean := strings.ReplaceAll(filename, "../", "")
	target := filepath.Join("/var/log", unsafeClean)
	// ruleid: path-traversal-file-write
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0600)
}

// 6. Real sanitizer: secure standard library sanitization using filepath.Base
func realSanitizerTest(r *http.Request) {
	filename := r.URL.Query().Get("filename")
	safeName := filepath.Base(filename)
	target := filepath.Join("/var/log", safeName)
	// ok: path-traversal-file-write
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0600)
}
