package rules

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// Case 1: Direct link from tar.Header.Linkname (vulnerable)
func DirectTarHardlink(hdr *tar.Header, linkPath string) error {
	// ruleid: archive-hardlink-path-traversal
	return os.Link(hdr.Linkname, linkPath)
}

// Case 2: Tar linkname propagated through filepath.Join (vulnerable)
func JoinedTarHardlink(hdr *tar.Header, extractDir, linkPath string) error {
	target := filepath.Join(extractDir, hdr.Linkname)
	// ruleid: archive-hardlink-path-traversal
	return os.Link(target, linkPath)
}

// Case 3: Untyped tar.Reader iteration with syscall.Link (vulnerable)
func ReaderTarHardlink(r io.Reader, linkPath string) error {
	tr := tar.NewReader(r)
	hdr, err := tr.Next()
	if err != nil {
		return err
	}
	target := filepath.Clean(hdr.Linkname)
	// ruleid: archive-hardlink-path-traversal
	return syscall.Link(target, linkPath)
}

// Case 4: Zip file entry name used as hardlink source (vulnerable)
func ZipEntryHardlink(zf *zip.File, linkPath string) error {
	target := filepath.Join("/tmp/extract", zf.Name)
	// ruleid: archive-hardlink-path-traversal
	return os.Link(target, linkPath)
}

// Case 5: Tar linkname sanitized with filepath.Base (safe)
func SafeBaseTarHardlink(hdr *tar.Header, extractDir, linkPath string) error {
	safeName := filepath.Base(hdr.Linkname)
	target := filepath.Join(extractDir, safeName)
	// ok: archive-hardlink-path-traversal
	return os.Link(target, linkPath)
}

// Case 6: Tar linkname sanitized with filepath.IsLocal (safe)
func SafeIsLocalTarHardlink(hdr *tar.Header, linkPath string) error {
	target := hdr.Linkname
	if !filepath.IsLocal(target) {
		return os.ErrInvalid
	}
	// ok: archive-hardlink-path-traversal
	return os.Link(target, linkPath)
}
