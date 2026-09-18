package main

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 1. Direct standard lib vulnerability
func DirectStandardLibVuln(hdr *tar.Header, targetDir string) error {
	dest := filepath.Join(targetDir, hdr.Name)
	// ruleid: go-archive-path-traversal
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// 2. Proper standard lib patch
func ProperStandardLibPatch(hdr *tar.Header, targetDir string) error {
	cleanName := filepath.Clean(hdr.Name)
	dest := filepath.Join(targetDir, cleanName)
	rel, err := filepath.Rel(targetDir, dest)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("path traversal attempt: %s", hdr.Name)
	}
	// ok: go-archive-path-traversal
	return os.WriteFile(dest, []byte("safe content"), 0644)
}

// Helper wrapper for cross-function taint propagation
func joinPathWrapper(base, part string) string {
	return filepath.Join(base, part)
}

// 3. Cross-function taint (wrapper function bypass)
func CrossFunctionVuln(hdr *tar.Header, targetDir string) error {
	dest := joinPathWrapper(targetDir, hdr.Name)
	// ruleid: go-archive-path-traversal
	return os.WriteFile(dest, []byte("untrusted payload"), 0644)
}

// ArchiveEntry abstracts archive header retrieval behind an interface
type ArchiveEntry interface {
	GetHeader() *tar.Header
}

// 4. Interface abstraction bypass
func InterfaceAbstractionVuln(entry ArchiveEntry, targetDir string) error {
	var hdr *tar.Header = entry.GetHeader()
	dest := filepath.Join(targetDir, hdr.Name)
	// ruleid: go-archive-path-traversal
	return os.WriteFile(dest, []byte("bypassed payload"), 0644)
}

// 5. Fake sanitizer usage (must trigger alert)
func FakeSanitizerVuln(hdr *tar.Header, targetDir string) error {
	// filepath.Clean and filepath.Join are fake sanitizers for relative traversals
	cleaned := filepath.Clean(hdr.Name)
	dest := filepath.Join(targetDir, cleaned)
	// ruleid: go-archive-path-traversal
	return os.WriteFile(dest, []byte("data"), 0644)
}

// 6. Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(hdr *tar.Header, targetDir string) error {
	name := hdr.Name
	if !filepath.IsLocal(name) {
		return fmt.Errorf("insecure path: %s", name)
	}
	dest := filepath.Join(targetDir, name)
	// ok: go-archive-path-traversal
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// HTTP handler helper utilizing net/http and archive/zip to satisfy package imports
func ServeHTTPArchiveStatus(w http.ResponseWriter, req *http.Request, zf *zip.File) {
	fmt.Fprintf(w, "Processing archive entry: %s from method %s", zf.Name, req.Method)
}
