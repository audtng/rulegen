package rules

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Mock web framework types to ensure test compiles without external internet dependencies
type GinContext struct{}

func (c *GinContext) Query(key string) string    { return "" }
func (c *GinContext) Param(key string) string    { return "" }
func (c *GinContext) PostForm(key string) string { return "" }

type EchoContext interface {
	FormValue(name string) string
	QueryParam(name string) string
}

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string  { return "" }
func (c *FiberCtx) Params(key string, defaultValue ...string) string { return "" }

// 1. Direct standard lib vulnerability
func DirectZipVuln(file *zip.File, dest string) error {
	target := filepath.Join(dest, file.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// 2. Proper standard lib patch
func ProperPatchedZip(file *zip.File, dest string) error {
	if !filepath.IsLocal(file.Name) {
		return errors.New("untrusted archive path escapes destination directory")
	}
	target := filepath.Join(dest, file.Name)
	// ok: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// 3. Cross-function taint (wrapper function bypass)
func resolveArchivePath(dest string, relativePath string) string {
	return fmt.Sprintf("%s/%s", dest, relativePath)
}

func WrapperBypass(file *zip.File, dest string) error {
	target := resolveArchivePath(dest, file.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// 4. Interface abstraction bypass
type PathModifier interface {
	Modify(path string) string
}

type IdentityModifier struct{}

func (IdentityModifier) Modify(p string) string {
	return p
}

func InterfaceBypass(req *http.Request, file *zip.File, dest string, modifier PathModifier) error {
	_ = req.Header.Get("X-Package-Name")
	modifiedName := modifier.Modify(file.Name)
	target := filepath.Join(dest, modifiedName)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// 5. Fake sanitizer usage (must trigger alert)
func FakeSanitizerUsage(file *zip.File, dest string) error {
	cleaned := path.Clean(file.Name)
	target := path.Join(dest, cleaned)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// 6. Real sanitizer usage (must not trigger alert)
func RealSanitizerUsage(file *zip.File, dest string) error {
	target := filepath.Join(dest, file.Name)
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

// Satisfy tar and framework imports
func UnusedMocks(c *GinContext, ec EchoContext, fc *FiberCtx, hdr *tar.Header) {
	_ = c.Query("q")
	_ = ec.FormValue("f")
	_ = fc.Query("q")
	_ = hdr.Name
}
