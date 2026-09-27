package rules

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Case 1: Vulnerable zip extraction using zip.Reader and os.OpenFile
func ExtractZipVulnerable(r *zip.Reader, targetDir string) error {
	for _, f := range r.File {
		path := filepath.Join(targetDir, f.Name)
		// ruleid: archive-path-traversal
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer outFile.Close()
	}
	return nil
}

// Case 2: Vulnerable tar extraction using tar.Reader and os.Create
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
		// ruleid: archive-path-traversal
		out, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return nil
}

// Case 3: Vulnerable extraction with flawed prefix check (missing path separator)
func ExtractZipFlawedPrefix(f *zip.File, destDir string) error {
	destPath := filepath.Join(destDir, f.Name)
	if strings.HasPrefix(destPath, filepath.Clean(destDir)) {
		// ruleid: archive-path-traversal
		return os.WriteFile(destPath, []byte("payload"), 0600)
	}
	return errors.New("invalid path")
}

// Case 4: Safe extraction with proper separator prefix validation
func ExtractZipSafePrefix(f *zip.File, destDir string) error {
	destPath := filepath.Join(destDir, f.Name)
	if !strings.HasPrefix(destPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
		return errors.New("path traversal detected")
	}
	// ok: archive-path-traversal
	return os.WriteFile(destPath, []byte("payload"), 0600)
}

// Case 5: Safe extraction sanitized via filepath.Base
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
		// ok: archive-path-traversal
		out, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return nil
}

// Case 6: Safe extraction validated via filepath.IsLocal
func ExtractZipSafeIsLocal(f *zip.File, destDir string) error {
	entryName := f.Name
	if !filepath.IsLocal(entryName) {
		return errors.New("unsafe entry path")
	}
	destPath := filepath.Join(destDir, entryName)
	// ok: archive-path-traversal
	return os.MkdirAll(destPath, 0750)
}
