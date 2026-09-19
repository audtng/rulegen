package rules

import (
	"archive/tar"
	"os"
	"path/filepath"
)

// Direct stdlib: untrusted archive path passed directly to file sink
func testDirectStdlib(tr *tar.Reader) {
	hdr, err := tr.Next()
	if err != nil {
		return
	}
	target := filepath.Join("/tmp/base", hdr.Name)
	// ruleid: archive-path-traversal
	os.Create(target)
}

// Proper patch: validation using filepath.IsLocal to ensure the path stays within bounds
func testProperPatch(tr *tar.Reader) {
	hdr, err := tr.Next()
	if err != nil {
		return
	}
	if !filepath.IsLocal(hdr.Name) {
		return
	}
	target := filepath.Join("/tmp/base", hdr.Name)
	// ok: archive-path-traversal
	os.Create(target)
}

// Cross-function taint: taint propagated across a helper function to a sink
func helperJoin(base, name string) string {
	return filepath.Join(base, name)
}

func testCrossFunctionTaint(tr *tar.Reader) {
	hdr, err := tr.Next()
	if err != nil {
		return
	}
	target := helperJoin("/tmp/base", hdr.Name)
	// ruleid: archive-path-traversal
	os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

// Interface bypass: taint propagated through an interface method before reaching a sink
type EntryGetter interface {
	EntryPath() string
}

type TarEntryWrapper struct {
	Header *tar.Header
}

func (w *TarEntryWrapper) EntryPath() string {
	return w.Header.Name
}

func testInterfaceBypass(tr *tar.Reader) {
	hdr, err := tr.Next()
	if err != nil {
		return
	}
	var getter EntryGetter = &TarEntryWrapper{Header: hdr}
	target := filepath.Join("/tmp/base", getter.EntryPath())
	// ruleid: archive-path-traversal
	os.Create(target)
}

// Fake sanitizer: filepath.Clean does not prevent traversal outside of the target base
func testFakeSanitizer(tr *tar.Reader) {
	hdr, err := tr.Next()
	if err != nil {
		return
	}
	cleaned := filepath.Clean(hdr.Name)
	target := filepath.Join("/tmp/base", cleaned)
	// ruleid: archive-path-traversal
	os.Create(target)
}

// Real sanitizer: filepath.Base safely strips any directory traversal components
func testRealSanitizer(tr *tar.Reader) {
	hdr, err := tr.Next()
	if err != nil {
		return
	}
	safeName := filepath.Base(hdr.Name)
	target := filepath.Join("/tmp/base", safeName)
	// ok: archive-path-traversal
	os.Create(target)
}
