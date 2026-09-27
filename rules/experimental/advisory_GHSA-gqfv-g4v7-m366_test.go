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

// Test 1: Vulnerable zip extraction using zip.Reader and os.OpenFile
func ExtractZipVulnerable(r *zip.Reader, targetDir string) error {
	for _, f := range r.File {
		path := filepath.Join(targetDir, f.Name)
		// ruleid: archive-zipslip-path-traversal
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer outFile.Close()
	}
	return nil
}

// Test 2: Vulnerable tar extraction using tar.Reader and os.Create
func ExtractTarVulnerable(tr *tar.Reader, targetDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		destPath := filepath.Join(targetDir, hdr.Name)
		// ruleid: archive-zipslip-path-traversal
		out, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return nil
}

// Test 3: Vulnerable directory creation from zip file header using os.MkdirAll
func ExtractZipDirVulnerable(f *zip.File, destDir string) error {
	destPath := filepath.Join(destDir, f.Name)
	// ruleid: archive-zipslip-path-traversal
	return os.MkdirAll(destPath, 0750)
}

// Test 4: Safe zip extraction validated via strings.HasPrefix and filepath.Clean
func ExtractZipSafePrefix(r *zip.Reader, targetDir string) error {
	cleanTarget := filepath.Clean(targetDir) + string(filepath.Separator)
	for _, f := range r.File {
		path := filepath.Join(targetDir, f.Name)
		cleanPath := filepath.Clean(path)
		if !strings.HasPrefix(cleanPath, cleanTarget) {
			return fmt.Errorf("path traversal detected: %s", path)
		}
		// ok: archive-zipslip-path-traversal
		outFile, err := os.OpenFile(cleanPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer outFile.Close()
	}
	return nil
}

// Test 5: Safe tar extraction using filepath.Base to prevent directory traversal
func ExtractTarSafeBase(tr *tar.Reader, targetDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		destPath := filepath.Join(targetDir, filepath.Base(hdr.Name))
		// ok: archive-zipslip-path-traversal
		out, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return nil
}

// Test 6: Safe extraction using filepath.IsLocal
func ExtractZipSafeIsLocal(f *zip.File, destDir string) error {
	fileName := f.Name
	if !filepath.IsLocal(fileName) {
		return fmt.Errorf("insecure file path: %s", fileName)
	}
	destPath := filepath.Join(destDir, fileName)
	// ok: archive-zipslip-path-traversal
	return os.WriteFile(destPath, []byte("content"), 0600)
}
