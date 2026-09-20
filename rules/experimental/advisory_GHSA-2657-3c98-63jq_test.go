package rules

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Test Case 1: Vulnerable Tar extraction using tar.Reader loop and filepath.Join
func ExtractTarVuln(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(destDir, hdr.Name)
		// ruleid: go-archive-path-traversal
		f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// Test Case 2: Vulnerable Zip extraction using path.Clean and path.Join (esm.sh pattern)
func ExtractZipVuln(zr *zip.Reader, destDir string) error {
	for _, f := range zr.File {
		cleaned := path.Clean(f.Name)
		target := path.Join(destDir, cleaned)
		// ruleid: go-archive-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// Test Case 3: Vulnerable helper function accepting *tar.Header directly with os.WriteFile
func WriteTarEntryVuln(hdr *tar.Header, targetDir string, data []byte) error {
	destPath := filepath.Join(targetDir, hdr.Name)
	// ruleid: go-archive-path-traversal
	return os.WriteFile(destPath, data, 0644)
}

// Test Case 4: Safe Tar extraction validated via filepath.Rel boundary check
func ExtractTarSafeRel(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(destDir, hdr.Name)
		rel, err := filepath.Rel(destDir, target)
		if err != nil || strings.HasPrefix(rel, "..") {
			return errors.New("path traversal detected")
		}
		// ok: go-archive-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// Test Case 5: Safe extraction using strings.HasPrefix boundary check on cleaned path
func ExtractTarSafePrefix(hdr *tar.Header, destDir string, data []byte) error {
	cleanDest := filepath.Clean(destDir)
	target := filepath.Clean(filepath.Join(cleanDest, hdr.Name))
	if !strings.HasPrefix(target, cleanDest+string(filepath.Separator)) {
		return fmt.Errorf("illegal file path: %s", target)
	}
	// ok: go-archive-path-traversal
	return os.WriteFile(target, data, 0644)
}

// Test Case 6: Safe Zip extraction using filepath.Base
func ExtractZipSafeBase(f *zip.File, destDir string) error {
	baseName := filepath.Base(f.Name)
	target := filepath.Join(destDir, baseName)
	// ok: go-archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	return out.Close()
}
