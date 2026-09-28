package rules

import (
	"archive/tar"
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
)

// Edge case 1: Tar file extraction joining destination and Header.Name without validation
func extractTarVulnerable(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Edge case 2: Tar file extraction with filepath.Rel boundary validation (safe)
func extractTarSafeRel(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, filepath.Clean(hdr.Name))
	rel, err := filepath.Rel(dest, target)
	if err != nil || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return os.ErrInvalid
	}
	// ok: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Edge case 3: Tar symlink creation where Header.Linkname escapes workspace (GHSA-qxx2-7h4c-83f4 vector)
func extractTarSymlinkVulnerable(hdr *tar.Header, dest string) error {
	linkPath := filepath.Join(dest, "symlink")
	// ruleid: archive-path-traversal
	return os.Symlink(hdr.Linkname, linkPath)
}

// Edge case 4: Tar symlink creation sanitized using filepath.IsLocal (safe)
func extractTarSymlinkSafeLocal(hdr *tar.Header, dest string) error {
	linkTarget := hdr.Linkname
	if !filepath.IsLocal(linkTarget) {
		return os.ErrInvalid
	}
	linkPath := filepath.Join(dest, "symlink")
	// ok: archive-path-traversal
	return os.Symlink(linkTarget, linkPath)
}

// Edge case 5: Zip extraction joining destination and File.Name without sanitization
func extractZipVulnerable(f *zip.File, dest string) error {
	target := filepath.Join(dest, f.Name)
	// ruleid: archive-path-traversal
	return os.WriteFile(target, []byte("data"), 0644)
}

// Edge case 6: Zip extraction sanitized with filepath.Base (safe)
func extractZipSafeBase(f *zip.File, dest string) error {
	safeFilename := filepath.Base(f.Name)
	target := filepath.Join(dest, safeFilename)
	// ok: archive-path-traversal
	return os.WriteFile(target, []byte("data"), 0644)
}
