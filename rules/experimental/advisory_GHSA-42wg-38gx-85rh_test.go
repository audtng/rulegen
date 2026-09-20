package rules

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// 1. Vulnerable: Classic Zip Slip extracting files using filepath.Join without validation.
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

// 2. Safe: Zip extraction with filepath.IsLocal validation.
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

// 3. Vulnerable: Tar extraction using tar.Reader.Next() and writing to constructed path.
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

// 4. Safe: Tar extraction sanitized using filepath.Base.
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

// 5. Vulnerable: Archive entry with prefix manipulation (CVE pattern) without escape check.
func RestoreDatabaseEntryVulnerable(f *zip.File, destDir string, data []byte) error {
	if strings.HasPrefix(f.Name, "database/") {
		rel := strings.TrimPrefix(f.Name, "database/")
		outPath := filepath.Join(destDir, rel)
		// ruleid: go-archive-path-traversal
		return os.WriteFile(outPath, data, 0600)
	}
	return nil
}

// 6. Safe: Archive entry sanitized with path.Base before file operations.
func RestoreDatabaseEntrySafe(f *zip.File, destDir string, data []byte) error {
	if strings.HasPrefix(f.Name, "database/") {
		rel := strings.TrimPrefix(f.Name, "database/")
		safeName := path.Base(rel)
		outPath := filepath.Join(destDir, safeName)
		// ok: go-archive-path-traversal
		return os.WriteFile(outPath, data, 0600)
	}
	return nil
}
