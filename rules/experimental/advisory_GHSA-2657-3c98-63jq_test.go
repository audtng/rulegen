package rules

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 1. Direct standard lib vulnerability
func DirectTarVuln(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(destDir, hdr.Name)
		// ruleid: go-archive-path-traversal-file-write
		f, err := os.Create(target)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// 2. Proper standard lib patch
func ProperTarPatch(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		name := hdr.Name
		if !filepath.IsLocal(name) {
			return fmt.Errorf("illegal path traversal in tar: %s", name)
		}
		target := filepath.Join(destDir, name)
		// ok: go-archive-path-traversal-file-write
		f, err := os.Create(target)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// 3. Cross-function taint (wrapper function bypass)
func writeTarEntry(hdr *tar.Header, destDir string) error {
	target := filepath.Join(destDir, hdr.Name)
	// ruleid: go-archive-path-traversal-file-write
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

func CrossFunctionWrapper(tr *tar.Reader, destDir string) error {
	hdr, err := tr.Next()
	if err != nil {
		return err
	}
	return writeTarEntry(hdr, destDir)
}

// 4. Interface abstraction bypass
type TarEntrySource interface {
	GetEntry() *tar.Header
}

type DefaultTarEntrySource struct {
	Header *tar.Header
}

func (s *DefaultTarEntrySource) GetEntry() *tar.Header {
	return s.Header
}

func InterfaceAbstractionVuln(source TarEntrySource, destDir string) error {
	var hdr *tar.Header = source.GetEntry()
	target := filepath.Join(destDir, hdr.Name)
	// ruleid: go-archive-path-traversal-file-write
	return os.WriteFile(target, []byte("payload"), 0644)
}

// 5. Fake sanitizer usage (must trigger alert)
// Uses filepath.Clean/path.Clean or Join without containment checking, exactly like GHSA-2657-3c98-63jq
func FakeSanitizerVuln(hdr *tar.Header, destDir string) error {
	cleaned := filepath.Clean(hdr.Name)
	target := filepath.Join(destDir, cleaned)
	// ruleid: go-archive-path-traversal-file-write
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	return f.Close()
}

// 6. Real sanitizer usage (must not trigger alert)
// Uses path prefix containment verification to ensure path does not escape destination
func RealSanitizerPrefixCheck(zf *zip.File, destDir string) error {
	target := filepath.Join(destDir, zf.Name)
	cleanDest := filepath.Clean(destDir) + string(filepath.Separator)
	if !strings.HasPrefix(target, cleanDest) {
		return fmt.Errorf("path traversal attempt: %s", zf.Name)
	}
	// ok: go-archive-path-traversal-file-write
	return os.WriteFile(target, []byte("data"), 0600)
}
