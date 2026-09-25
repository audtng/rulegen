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

// 1. Tar Extraction - Vulnerable (Direct join without validation into file creation)
func testVulnerableTarExtract(tr *tar.Reader, dest string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dest, hdr.Name)
		// ruleid: go-archive-path-traversal
		f, err := os.Create(target)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// 2. Tar Extraction - Safe (Using filepath.Base to prevent directory traversal)
func testSafeTarExtractBase(tr *tar.Reader, dest string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		cleanName := filepath.Base(hdr.Name)
		target := filepath.Join(dest, cleanName)
		// ok: go-archive-path-traversal
		f, err := os.Create(target)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// 3. Zip Extraction - Vulnerable (Direct join into OpenFile without validation)
func testVulnerableZipExtract(zr *zip.Reader, dest string) error {
	for _, f := range zr.File {
		target := filepath.Join(dest, f.Name)
		// ruleid: go-archive-path-traversal
		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// 4. Incomplete Sanitization - Vulnerable (Missing path separator allows prefix match bypass)
func testVulnerablePrefixSanitization(hdr *tar.Header, dest string) error {
	destPath := filepath.Join(dest, hdr.Name)
	if !strings.HasPrefix(destPath, filepath.Clean(dest)) {
		return fmt.Errorf("path traversal attempt: %s", destPath)
	}
	// ruleid: go-archive-path-traversal
	return os.WriteFile(destPath, []byte("payload"), 0600)
}

// 5. Patched Sanitization - Safe (Clean destination verified with trailing path separator)
func testPatchedPrefixSanitization(hdr *tar.Header, dest string) error {
	destPath := filepath.Join(dest, hdr.Name)
	if !strings.HasPrefix(destPath, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("path traversal attempt: %s", destPath)
	}
	// ok: go-archive-path-traversal
	return os.WriteFile(destPath, []byte("payload"), 0600)
}

// 6. Modern Sanitization - Safe (filepath.IsLocal check before file write)
func testSafeIsLocalSanitization(f *zip.File, dest string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return fmt.Errorf("untrusted path: %s", name)
	}
	target := filepath.Join(dest, name)
	// ok: go-archive-path-traversal
	return os.WriteFile(target, []byte("safe"), 0600)
}
