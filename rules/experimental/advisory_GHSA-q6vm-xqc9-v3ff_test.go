package rules

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 1. Vulnerable: zip file entry extracted directly using filepath.Join
func VulnZipExtract(f *zip.File, destDir string) error {
	destPath := filepath.Join(destDir, f.Name)
	// ruleid: archive-zip-slip
	out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// 2. Safe: zip file entry sanitized with filepath.Base
func SafeZipBase(f *zip.File, destDir string) error {
	baseName := filepath.Base(f.Name)
	destPath := filepath.Join(destDir, baseName)
	// ok: archive-zip-slip
	out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// 3. Vulnerable: tar header entry extracted with filepath.Clean to os.Create
func VulnTarExtract(hdr *tar.Header, destDir string) error {
	cleanName := filepath.Clean(hdr.Name)
	target := filepath.Join(destDir, cleanName)
	// ruleid: archive-zip-slip
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// 4. Safe: tar header entry guarded with filepath.IsLocal
func SafeTarIsLocal(hdr *tar.Header, destDir string) error {
	name := hdr.Name
	if !filepath.IsLocal(name) {
		return fmt.Errorf("insecure path: %s", name)
	}
	target := filepath.Join(destDir, name)
	// ok: archive-zip-slip
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// 5. Vulnerable: zip.Reader loop creating directories with os.MkdirAll
func VulnZipReaderLoop(r *zip.Reader, destDir string) error {
	for _, f := range r.File {
		dirPath := filepath.Join(destDir, f.Name)
		// ruleid: archive-zip-slip
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}
	return nil
}

// 6. Safe: zip extraction with destination prefix containment check (strings.HasPrefix)
func SafeZipPrefixCheck(f *zip.File, destDir string) error {
	destPath := filepath.Join(destDir, f.Name)
	cleanDest := filepath.Clean(destDir) + string(filepath.Separator)
	if !strings.HasPrefix(destPath, cleanDest) {
		return fmt.Errorf("zip slip detected: %s", destPath)
	}
	// ok: archive-zip-slip
	return os.WriteFile(destPath, []byte("data"), 0644)
}
