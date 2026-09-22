package rules

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Test 1: Vulnerable Zip extraction loop directly writing to destination
func VulnZipExtract(zr *zip.Reader, destDir string) error {
	for _, f := range zr.File {
		targetPath := filepath.Join(destDir, f.Name)
		// ruleid: go-archive-path-traversal
		w, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		_ = w.Close()
	}
	return nil
}

// Test 2: Vulnerable Tar extraction reading header name and writing file
func VulnTarExtract(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err != nil {
			return err
		}
		targetPath := filepath.Join(destDir, hdr.Name)
		// ruleid: go-archive-path-traversal
		if err := os.WriteFile(targetPath, []byte("payload"), 0644); err != nil {
			return err
		}
	}
}

// Test 3: Vulnerable extraction using filepath.Clean without prefix check
func VulnCleanNotSufficient(f *zip.File, destDir string) error {
	joined := filepath.Join(destDir, f.Name)
	cleanPath := filepath.Clean(joined)
	// ruleid: go-archive-path-traversal
	out, err := os.Create(cleanPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Test 4: Safe extraction using filepath.IsLocal validation
func SafeZipIsLocal(f *zip.File, destDir string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return errors.New("insecure file path in archive")
	}
	targetPath := filepath.Join(destDir, name)
	// ok: go-archive-path-traversal
	return os.WriteFile(targetPath, []byte("data"), 0644)
}

// Test 5: Safe extraction using filepath.Base sanitization
func SafeZipBase(f *zip.File, destDir string) error {
	safeName := filepath.Base(f.Name)
	targetPath := filepath.Join(destDir, safeName)
	// ok: go-archive-path-traversal
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Test 6: Safe extraction using strings.HasPrefix containment check
func SafePrefixCheck(f *zip.File, destDir string) error {
	targetPath := filepath.Join(destDir, f.Name)
	cleanDest := filepath.Clean(destDir) + string(filepath.Separator)
	if !strings.HasPrefix(targetPath, cleanDest) {
		return errors.New("path escapes destination directory")
	}
	// ok: go-archive-path-traversal
	w, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer w.Close()
	return nil
}
