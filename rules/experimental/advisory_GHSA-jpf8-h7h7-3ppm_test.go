package rules

import (
	"archive/tar"
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
)

// Case 1: Direct tar extraction with hdr.Name joined to dest (Vulnerable)
func testVulnTarDirect(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Case 2: Zip extraction with f.Name via os.OpenFile (Vulnerable)
func testVulnZipOpenFile(f *zip.File, dest string) error {
	target := filepath.Join(dest, f.Name)
	// ruleid: archive-path-traversal
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 3: Flawed sanitization with filepath.Clean which does not prevent directory traversal (Vulnerable)
func testVulnCleanFalseSanitization(hdr *tar.Header, dest string) error {
	cleaned := filepath.Clean(filepath.Join(dest, hdr.Name))
	// ruleid: archive-path-traversal
	return os.WriteFile(cleaned, []byte("data"), 0600)
}

// Case 4: Safe extraction using filepath.Base (Safe)
func testSafeBase(hdr *tar.Header, dest string) error {
	base := filepath.Base(hdr.Name)
	target := filepath.Join(dest, base)
	// ok: archive-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Case 5: Safe extraction using filepath.IsLocal (Safe)
func testSafeIsLocal(f *zip.File, dest string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return os.ErrInvalid
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	return os.WriteFile(target, []byte("data"), 0600)
}

// Case 6: Safe extraction using strings.HasPrefix path containment check (Safe)
func testSafePrefixCheck(hdr *tar.Header, dest string) error {
	cleanDest := filepath.Clean(dest) + string(filepath.Separator)
	target := filepath.Join(dest, hdr.Name)
	if !strings.HasPrefix(target, cleanDest) {
		return os.ErrInvalid
	}
	// ok: archive-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
