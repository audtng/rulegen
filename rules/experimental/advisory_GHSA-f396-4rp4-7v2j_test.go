package rules

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Case 1: Vulnerable Tar extraction writing regular file without path validation
func VulnTarExtract(tr *tar.Reader, dest string) error {
	hdr, err := tr.Next()
	if err != nil {
		return err
	}
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

// Case 2: Vulnerable Zip extraction writing file without validation
func VulnZipExtract(f *zip.File, dest string) error {
	target := filepath.Join(dest, f.Name)
	// ruleid: archive-path-traversal
	return os.WriteFile(target, []byte("payload"), 0644)
}

// Case 3: Vulnerable Symlink creation pointing to untrusted target linkname
func VulnTarSymlink(hdr *tar.Header, dest string) error {
	linkPath := filepath.Join(dest, "symlink_entry")
	// ruleid: archive-path-traversal
	return os.Symlink(hdr.Linkname, linkPath)
}

// Case 4: Patched extraction using filepath.Base
func SafeTarBase(hdr *tar.Header, dest string) error {
	base := filepath.Base(hdr.Name)
	target := filepath.Join(dest, base)
	// ok: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

// Case 5: Patched extraction validating entry path with filepath.IsLocal
func SafeTarIsLocal(hdr *tar.Header, dest string) error {
	name := hdr.Name
	if !filepath.IsLocal(name) {
		return errors.New("unsafe path")
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

// Case 6: Patched extraction validating destination prefix
func SafeTarPrefix(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	cleanDest := filepath.Clean(dest) + string(filepath.Separator)
	if !strings.HasPrefix(target, cleanDest) {
		return errors.New("unsafe path")
	}
	// ok: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}
