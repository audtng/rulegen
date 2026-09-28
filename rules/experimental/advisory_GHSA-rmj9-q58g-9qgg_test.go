package rules

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Case 1: Vulnerable zip extraction using os.OpenFile directly
func testCase1VulnZipOpenFile(f *zip.File, dest string) error {
	path := filepath.Join(dest, f.Name)
	// ruleid: archive-path-traversal
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 2: Safe zip extraction with strings.HasPrefix boundary check (matches CVE fix)
func testCase2SafeZipPrefix(f *zip.File, dest string) error {
	path := filepath.Join(dest, f.Name)
	cleanDest := filepath.Clean(dest) + string(os.PathSeparator)
	if !strings.HasPrefix(path, cleanDest) {
		return fmt.Errorf("illegal file path: %s", path)
	}
	// ok: archive-path-traversal
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 3: Vulnerable tar directory creation using os.MkdirAll
func testCase3VulnTarMkdir(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	return os.MkdirAll(target, 0755)
}

// Case 4: Safe tar file extraction using filepath.Base sanitization
func testCase4SafeTarBase(hdr *tar.Header, dest string) error {
	safeName := filepath.Base(hdr.Name)
	target := filepath.Join(dest, safeName)
	// ok: archive-path-traversal
	return os.WriteFile(target, []byte("content"), 0644)
}

// Case 5: Vulnerable zip extraction using os.Create
func testCase5VulnZipCreate(f *zip.File, dest string) error {
	target := filepath.Join(dest, f.Name)
	// ruleid: archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 6: Safe zip extraction using filepath.IsLocal validation
func testCase6SafeZipIsLocal(f *zip.File, dest string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return fmt.Errorf("insecure path: %s", name)
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	return os.WriteFile(target, []byte("content"), 0644)
}
