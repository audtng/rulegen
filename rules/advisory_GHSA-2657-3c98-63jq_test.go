package rules

import (
	"archive/tar"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Direct standard lib vulnerability
func DirectTarVuln(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// Proper standard lib patch
func ProperPatchedTar(hdr *tar.Header, dest string) error {
	if !filepath.IsLocal(hdr.Name) {
		return errors.New("untrusted archive path escapes destination directory")
	}
	target := filepath.Join(dest, hdr.Name)
	// ok: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// Cross-function taint (wrapper function bypass)
func resolveArchivePath(dest string, relativePath string) string {
	return fmt.Sprintf("%s/%s", dest, relativePath)
}

func WrapperBypass(hdr *tar.Header, dest string) error {
	target := resolveArchivePath(dest, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// Interface abstraction bypass
type PathModifier interface {
	Modify(path string) string
}

type IdentityModifier struct{}

func (IdentityModifier) Modify(p string) string {
	return p
}

func InterfaceBypass(req *http.Request, hdr *tar.Header, dest string, modifier PathModifier) error {
	_ = req.Header.Get("X-Package-Name")
	modifiedName := modifier.Modify(hdr.Name)
	target := filepath.Join(dest, modifiedName)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// Fake sanitizer usage (must trigger alert)
func FakeSanitizerUsage(hdr *tar.Header, dest string) error {
	cleaned := path.Clean(hdr.Name)
	target := path.Join(dest, cleaned)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	cleanDest := filepath.Clean(dest) + string(filepath.Separator)
	if !strings.HasPrefix(target, cleanDest) {
		return errors.New("illegal file path outside destination")
	}
	// ok: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}
