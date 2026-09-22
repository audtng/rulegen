package rules

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// 1. Tar extraction vulnerable to path traversal via Header.Name
func VulnTarExtract(tr *tar.Reader, dest string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
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
	}
	return nil
}

// 2. Tar extraction safe using filepath.Base
func SafeTarExtractBase(tr *tar.Reader, dest string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		base := filepath.Base(hdr.Name)
		target := filepath.Join(dest, base)
		// ok: archive-path-traversal
		f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// 3. Zip extraction vulnerable to path traversal via range loop
func VulnZipExtract(zr *zip.Reader, dest string) error {
	for _, f := range zr.File {
		target := filepath.Join(dest, f.Name)
		// ruleid: archive-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// 4. Zip extraction safe using filepath.IsLocal
func SafeZipExtractLocal(zr *zip.Reader, dest string) error {
	for _, f := range zr.File {
		name := f.Name
		if !filepath.IsLocal(name) {
			continue
		}
		target := filepath.Join(dest, name)
		// ok: archive-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// 5. Tar hardlink extraction vulnerable via Linkname
func VulnTarHardlink(hdr *tar.Header, dest string) error {
	linkTarget := filepath.Join(dest, hdr.Linkname)
	target := filepath.Join(dest, "safe.txt")
	// ruleid: archive-path-traversal
	return os.Link(linkTarget, target)
}

// 6. Tar directory creation safe using filepath.IsLocal
func SafeTarMkdirLocal(hdr *tar.Header, dest string) error {
	name := hdr.Name
	if !filepath.IsLocal(name) {
		return nil
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	return os.MkdirAll(target, 0755)
}
