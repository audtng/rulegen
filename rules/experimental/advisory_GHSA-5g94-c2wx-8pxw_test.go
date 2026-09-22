package rules

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// VulnTarExtract demonstrates vulnerable extraction from tar.Reader Next().
func VulnTarExtract(tr *tar.Reader, targetDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		dest := filepath.Join(targetDir, hdr.Name)
		// ruleid: archive-path-traversal
		f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// VulnZipExtract demonstrates vulnerable extraction from zip.Reader File slice.
func VulnZipExtract(zr *zip.Reader, targetDir string) error {
	for _, f := range zr.File {
		dest := filepath.Join(targetDir, f.Name)
		// ruleid: archive-path-traversal
		out, err := os.Create(dest)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// VulnTarSymlink demonstrates vulnerable symlink creation from tar.Header.
func VulnTarSymlink(hdr *tar.Header, targetDir string) error {
	dest := filepath.Join(targetDir, hdr.Name)
	// ruleid: archive-path-traversal
	return os.Symlink(hdr.Linkname, dest)
}

// SafeTarBase demonstrates safe extraction using filepath.Base to prevent directory traversal.
func SafeTarBase(hdr *tar.Header, targetDir string, data []byte) error {
	safeName := filepath.Base(hdr.Name)
	dest := filepath.Join(targetDir, safeName)
	// ok: archive-path-traversal
	return os.WriteFile(dest, data, 0644)
}

// SafeZipIsLocal demonstrates safe extraction using filepath.IsLocal (Go 1.20+).
func SafeZipIsLocal(f *zip.File, targetDir string, data []byte) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return os.ErrInvalid
	}
	dest := filepath.Join(targetDir, name)
	// ok: archive-path-traversal
	return os.WriteFile(dest, data, 0644)
}

// SafeStaticPath demonstrates writing metadata or logs to a fixed, non-tainted file.
func SafeStaticPath(targetDir string, manifest []byte) error {
	dest := filepath.Join(targetDir, "manifest.json")
	// ok: archive-path-traversal
	return os.WriteFile(dest, manifest, 0644)
}
