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

// Case 1: Unsanitized extraction from zip archive (vulnerable)
func vulnZipExtract(reader *zip.Reader, destDir string) error {
	for _, file := range reader.File {
		target := filepath.Join(destDir, file.Name)
		// ruleid: archive-zipslip-path-traversal
		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return nil
}

// Case 2: Extraction from zip archive sanitized via filepath.Base (safe)
func safeZipExtractBase(reader *zip.Reader, destDir string) error {
	for _, file := range reader.File {
		cleanName := filepath.Base(file.Name)
		target := filepath.Join(destDir, cleanName)
		// ok: archive-zipslip-path-traversal
		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return nil
}

// Case 3: Unsanitized extraction from tar archive (vulnerable)
func vulnTarExtract(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(destDir, hdr.Name)
		// ruleid: archive-zipslip-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return nil
}

// Case 4: Tar extraction sanitized using filepath.IsLocal (safe)
func safeTarExtractIsLocal(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		cleanName := filepath.Clean(hdr.Name)
		if !filepath.IsLocal(cleanName) {
			return fmt.Errorf("insecure path: %s", cleanName)
		}
		target := filepath.Join(destDir, cleanName)
		// ok: archive-zipslip-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return nil
}

// Case 5: Direct archive entry read via os.Open without path validation (vulnerable)
func vulnReadArchiveEntry(file *zip.File, rootDir string) (*os.File, error) {
	finalPath := filepath.Join(rootDir, file.Name)
	// ruleid: archive-zipslip-path-traversal
	return os.Open(finalPath)
}

// Case 6: Archive entry access sanitized with relative boundary validation (safe)
func safeReadArchiveEntryRel(file *zip.File, rootDir string) (*os.File, error) {
	cleanPath := filepath.Clean(file.Name)
	target := filepath.Join(rootDir, cleanPath)
	rel, err := filepath.Rel(rootDir, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("file path %s is outside root %s", target, rootDir)
	}
	// ok: archive-zipslip-path-traversal
	return os.Open(target)
}
