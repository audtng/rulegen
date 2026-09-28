package rules

import (
	"archive/tar"
	"archive/zip"
	"os"
	"path/filepath"
)

// 1. Vulnerable: tar header name directly joined and opened
func VulnTarExtract(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// 2. Safe: tar header name sanitized with filepath.Base before write
func SafeTarExtractBase(hdr *tar.Header, dest string, data []byte) error {
	cleanName := filepath.Base(hdr.Name)
	target := filepath.Join(dest, cleanName)
	// ok: archive-path-traversal
	return os.WriteFile(target, data, 0644)
}

// 3. Vulnerable: zip file entry name directly joined and created
func VulnZipExtract(f *zip.File, dest string) error {
	target := filepath.Join(dest, f.Name)
	// ruleid: archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// 4. Safe: zip file entry name sanitized with filepath.IsLocal
func SafeZipExtractIsLocal(f *zip.File, dest string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return os.ErrInvalid
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// 5. Vulnerable: tar symlink linkname directly written via os.Symlink
func VulnTarSymlink(hdr *tar.Header, dest string) error {
	linkTarget := hdr.Linkname
	target := filepath.Join(dest, filepath.Base(hdr.Name))
	// ruleid: archive-path-traversal
	return os.Symlink(linkTarget, target)
}

// 6. Safe: tar symlink linkname sanitized with filepath.IsLocal
func SafeTarSymlink(hdr *tar.Header, dest string) error {
	linkTarget := hdr.Linkname
	if !filepath.IsLocal(linkTarget) {
		return os.ErrInvalid
	}
	target := filepath.Join(dest, filepath.Base(hdr.Name))
	// ok: archive-path-traversal
	return os.Symlink(linkTarget, target)
}
