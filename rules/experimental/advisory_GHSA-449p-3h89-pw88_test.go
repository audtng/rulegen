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

// 1. Vulnerable: Zip extraction using filepath.Join without boundary validation.
func ExtractZipVulnerable(r *zip.Reader, destDir string) error {
	for _, f := range r.File {
		targetPath := filepath.Join(destDir, f.Name)
		// ruleid: go-archive-path-traversal
		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		outFile.Close()
	}
	return nil
}

// 2. Safe: Zip extraction validating path with filepath.IsLocal.
func ExtractZipSafe(r *zip.Reader, destDir string) error {
	for _, f := range r.File {
		entryName := f.Name
		if !filepath.IsLocal(entryName) {
			return errors.New("invalid path in archive")
		}
		targetPath := filepath.Join(destDir, entryName)
		// ok: go-archive-path-traversal
		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		outFile.Close()
	}
	return nil
}

// 3. Vulnerable: Tar extraction reading next entry and creating file directly.
func ExtractTarVulnerable(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(destDir, hdr.Name)
		// ruleid: go-archive-path-traversal
		f, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// 4. Safe: Tar extraction sanitizing entry name using filepath.Base.
func ExtractTarSafe(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		safeName := filepath.Base(hdr.Name)
		targetPath := filepath.Join(destDir, safeName)
		// ok: go-archive-path-traversal
		f, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// 5. Vulnerable: Insufficient sanitization using filepath.Clean which still allows traversal.
func ExtractCleanBypassVulnerable(f *zip.File, destDir string, data []byte) error {
	cleaned := filepath.Clean(f.Name)
	targetPath := filepath.Join(destDir, cleaned)
	// ruleid: go-archive-path-traversal
	return os.WriteFile(targetPath, data, 0644)
}

// 6. Safe: Zip entry validated by checking destination prefix boundary.
func ExtractZipPrefixCheckSafe(f *zip.File, destDir string, data []byte) error {
	targetPath := filepath.Join(destDir, f.Name)
	cleanDest := filepath.Clean(destDir) + string(filepath.Separator)
	if !strings.HasPrefix(targetPath, cleanDest) {
		return errors.New("path traversal detected")
	}
	// ok: go-archive-path-traversal
	return os.WriteFile(targetPath, data, 0644)
}
