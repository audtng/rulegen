package rules

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Case 1: Vulnerable zip extraction via os.Create
func VulnZipExtract(f *zip.File, destDir string) error {
	targetPath := filepath.Join(destDir, f.Name)
	// ruleid: archive-path-traversal
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 2: Vulnerable tar extraction via os.OpenFile
func VulnTarExtract(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		targetPath := filepath.Join(destDir, hdr.Name)
		// ruleid: archive-path-traversal
		outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		outFile.Close()
	}
	return nil
}

// Case 3: Vulnerable directory creation via os.MkdirAll
func VulnZipMkdirAll(f *zip.File, destDir string) error {
	targetPath := filepath.Join(destDir, f.Name)
	// ruleid: archive-path-traversal
	return os.MkdirAll(targetPath, 0755)
}

// Case 4: Safe zip extraction using filepath.Base sanitizer
func SafeZipBase(f *zip.File, destDir string) error {
	cleanName := filepath.Base(f.Name)
	targetPath := filepath.Join(destDir, cleanName)
	// ok: archive-path-traversal
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 5: Safe zip extraction using strings.HasPrefix path containment check
func SafeZipPrefixValidation(f *zip.File, destDir string, data []byte) error {
	targetPath := filepath.Join(destDir, f.Name)
	cleanDest := filepath.Clean(destDir) + string(os.PathSeparator)
	if !strings.HasPrefix(targetPath, cleanDest) {
		return os.ErrInvalid
	}
	// ok: archive-path-traversal
	return os.WriteFile(targetPath, data, 0644)
}

// Case 6: Safe zip extraction using filepath.IsLocal validation
func SafeZipIsLocal(f *zip.File, destDir string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return os.ErrInvalid
	}
	targetPath := filepath.Join(destDir, name)
	// ok: archive-path-traversal
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}
