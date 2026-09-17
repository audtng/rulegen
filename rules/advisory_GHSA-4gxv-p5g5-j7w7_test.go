package test

import (
	"net/http"
	"os"
	"path/filepath"
)

// Case 1: Direct standard lib vulnerability
func DirectVulnerability(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get("path")
	// ruleid: go-path-traversal-arbitrary-file-write
	os.WriteFile(path, []byte("playlist data"), 0644)
}

// Case 2: Proper standard lib patch
func ProperStandardLibPatch(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get("path")
	if !filepath.IsLocal(path) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	// ok: go-path-traversal-arbitrary-file-write
	os.WriteFile(path, []byte("playlist data"), 0644)
}

// Case 3: Cross-function taint (wrapper function bypass)
func buildPathWrapper(input string) string {
	return filepath.Join("/var/media/playlists", input)
}

func CrossFunctionBypass(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	target := buildPathWrapper(name)
	// ruleid: go-path-traversal-arbitrary-file-write
	os.WriteFile(target, []byte("playlist data"), 0644)
}

// Case 4: Interface abstraction bypass
type PathResolver interface {
	Resolve(input string) string
}

type DefaultPathResolver struct{}

func (d *DefaultPathResolver) Resolve(input string) string {
	return filepath.Join("/var/media/playlists", input)
}

func InterfaceAbstractionBypass(w http.ResponseWriter, req *http.Request) {
	name := req.PostFormValue("name")
	var resolver PathResolver = &DefaultPathResolver{}
	target := resolver.Resolve(name)
	// ruleid: go-path-traversal-arbitrary-file-write
	f, _ := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if f != nil {
		f.Close()
	}
}

// Case 5: Fake sanitizer usage (must trigger alert)
func FakeSanitizerUsage(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get("path")
	cleaned := filepath.Clean(path)
	fullPath := filepath.Join("/var/media/playlists", cleaned)
	// ruleid: go-path-traversal-arbitrary-file-write
	os.WriteFile(fullPath, []byte("playlist data"), 0644)
}

// Case 6: Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get("path")
	safeName := filepath.Base(path)
	fullPath := filepath.Join("/var/media/playlists", safeName)
	// ok: go-path-traversal-arbitrary-file-write
	os.WriteFile(fullPath, []byte("playlist data"), 0644)
}
