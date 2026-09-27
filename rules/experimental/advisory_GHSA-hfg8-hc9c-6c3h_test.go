package rules

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
)

// Case 1: Vulnerable Tar extraction using os.OpenFile directly
func VulnTarExtract(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// Case 2: Vulnerable Zip extraction using os.WriteFile directly
func VulnZipExtract(f *zip.File, dest string, data []byte) error {
	target := filepath.Join(dest, f.Name)
	// ruleid: archive-path-traversal
	return os.WriteFile(target, data, 0644)
}

// Case 3: Vulnerable Tar symlink creation using os.Symlink with Linkname
func VulnTarSymlink(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	return os.Symlink(hdr.Linkname, target)
}

// Case 4: Safe Tar extraction using filepath.Base
func SafeTarBase(hdr *tar.Header, dest string) error {
	safeName := filepath.Base(hdr.Name)
	target := filepath.Join(dest, safeName)
	// ok: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// Case 5: Safe Zip extraction using filepath.IsLocal validation
func SafeZipIsLocal(f *zip.File, dest string, data []byte) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return errors.New("unsafe entry name")
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	return os.WriteFile(target, data, 0644)
}

// Case 6: Safe Tar directory creation using filepath.IsLocal validation
func SafeTarDirIsLocal(hdr *tar.Header, dest string) error {
	name := hdr.Name
	if !filepath.IsLocal(name) {
		return errors.New("unsafe directory name")
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	return os.MkdirAll(target, 0755)
}
