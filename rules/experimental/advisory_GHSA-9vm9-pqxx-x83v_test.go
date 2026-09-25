package rules

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// 1. Tar extraction with filepath.Join without validation (vulnerable)
func ExtractTarVulnerable(tr *tar.Reader, dest string) error {
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

// 2. Zip extraction with filepath.Join without validation (vulnerable)
func ExtractZipVulnerable(zr *zip.Reader, dest string) error {
	for _, f := range zr.File {
		target := filepath.Join(dest, f.Name)
		// ruleid: archive-path-traversal
		if err := os.WriteFile(target, []byte("content"), 0644); err != nil {
			return err
		}
	}
	return nil
}

// 3. Direct tar header extraction without validation (vulnerable)
func ExtractTarHeaderVulnerable(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	return out.Close()
}

// 4. Tar extraction sanitized with filepath.Base (safe)
func ExtractTarSafeBase(tr *tar.Reader, dest string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		baseName := filepath.Base(hdr.Name)
		target := filepath.Join(dest, baseName)
		// ok: archive-path-traversal
		f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// 5. Zip extraction sanitized with filepath.IsLocal (safe)
func ExtractZipSafeIsLocal(zr *zip.Reader, dest string) error {
	for _, f := range zr.File {
		name := f.Name
		if !filepath.IsLocal(name) {
			continue
		}
		target := filepath.Join(dest, name)
		// ok: archive-path-traversal
		if err := os.WriteFile(target, []byte("content"), 0644); err != nil {
			return err
		}
	}
	return nil
}

// 6. Archive extraction writing to a static destination path (safe)
func ExtractTarSafeFixedPath(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, "static_output.dat")
	// ok: archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	return out.Close()
}
