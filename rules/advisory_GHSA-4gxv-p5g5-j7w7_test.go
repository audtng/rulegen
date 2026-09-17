package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// Helper wrapper function for cross-function taint analysis
func writeToFile(target string, content []byte) error {
	// ruleid: go-arbitrary-file-write
	return os.WriteFile(target, content, 0644)
}

// Helper interface abstraction
type FileWriter interface {
	Write(path string, data []byte) error
}

type LocalFileWriter struct{}

func (l *LocalFileWriter) Write(path string, data []byte) error {
	// ruleid: go-arbitrary-file-write
	return os.WriteFile(path, data, 0644)
}

// 1. Direct standard lib vulnerability
func DirectVuln(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	// ruleid: go-arbitrary-file-write
	f, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, _ = f.WriteString("playlist data")
}

// 2. Proper standard lib patch
func ProperPatch(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	// ok: go-arbitrary-file-write
	f, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, _ = f.WriteString("playlist data")
}

// 3. Cross-function taint (wrapper function bypass)
func CrossFunctionWrapperBypass(w http.ResponseWriter, req *http.Request) {
	filename := req.FormValue("file")
	if err := writeToFile(filename, []byte("playlist data")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 4. Interface abstraction bypass
func InterfaceAbstractionBypass(w http.ResponseWriter, req *http.Request, writer FileWriter) {
	filename := req.PostFormValue("file")
	if err := writer.Write(filename, []byte("playlist data")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 5. Fake sanitizer usage (must trigger alert)
func FakeSanitizerUsage(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	unsafePath := filepath.Join("/var/data/playlists", filepath.Clean(filename))
	// ruleid: go-arbitrary-file-write
	f, err := os.OpenFile(unsafePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, _ = fmt.Fprintln(f, "playlist data")
}

// 6. Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("file")
	safeFilename := filepath.Base(filename)
	safePath := filepath.Join("/var/data/playlists", safeFilename)
	// ok: go-arbitrary-file-write
	f, err := os.OpenFile(safePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, _ = fmt.Fprintln(f, "playlist data")
}
