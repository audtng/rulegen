package rules

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

// Direct standard lib vulnerability
func DirectVuln(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	// ruleid: go-path-traversal-file-read
	data, err := os.ReadFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Proper standard lib patch
func ProperPatch(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	if !fs.ValidPath(filename) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	// ok: go-path-traversal-file-read
	data, err := os.ReadFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Cross-function taint (wrapper function bypass)
func CrossFunctionVuln(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	readWrapper := func() ([]byte, error) {
		// ruleid: go-path-traversal-file-read
		return os.ReadFile(filename)
	}
	data, err := readWrapper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Interface abstraction bypass
func InterfaceAbstractionVuln(w http.ResponseWriter, req *http.Request) {
	dir := req.URL.Query().Get("dir")
	// ruleid: go-path-traversal-file-read
	var fsys http.FileSystem = http.Dir(dir)
	f, err := fsys.Open("config.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

// Fake sanitizer usage (must trigger alert)
func FakeSanitizerVuln(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	safePath := filepath.Join("/base/directory", filename)
	// ruleid: go-path-traversal-file-read
	data, err := os.ReadFile(safePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	safePath := filepath.Join("/base/directory", filename)
	// ok: go-path-traversal-file-read
	data, err := os.ReadFile(safePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
