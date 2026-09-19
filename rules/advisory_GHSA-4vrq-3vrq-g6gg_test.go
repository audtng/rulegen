package rules

import (
	"errors"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type PathExtractor interface {
	ExtractPath() string
}

type URLFragmentWrapper struct {
	fragment string
}

func (w URLFragmentWrapper) ExtractPath() string {
	return w.fragment
}

// 1. Direct stdlib: direct extraction from parsed URL into file operation
func testDirectStdlib(rawURL string, baseDir string) (*os.File, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	target := filepath.Join(baseDir, u.Fragment)
	// ruleid: url-path-traversal-file-access
	return os.Open(target)
}

// 2. Proper patch: validation using filepath.IsLocal
func testProperPatch(rawURL string, baseDir string) (*os.File, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	subdir := u.Fragment
	if !filepath.IsLocal(subdir) {
		return nil, errors.New("path traversal detected")
	}
	target := filepath.Join(baseDir, subdir)
	// ok: url-path-traversal-file-access
	return os.Open(target)
}

// 3. Cross-function taint: taint flowing across functions via URL object
func openFromURLHelper(u *url.URL, baseDir string) (*os.File, error) {
	target := filepath.Join(baseDir, u.Fragment)
	// ruleid: url-path-traversal-file-access
	return os.Open(target)
}

func testCrossFunctionTaint(rawURL string, baseDir string) (*os.File, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return openFromURLHelper(u, baseDir)
}

// 4. Interface bypass: taint passed through interface method
func testInterfaceBypass(rawURL string, baseDir string) (*os.File, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	var extractor PathExtractor = URLFragmentWrapper{fragment: u.Fragment}
	target := filepath.Join(baseDir, extractor.ExtractPath())
	// ruleid: url-path-traversal-file-access
	return os.Open(target)
}

// 5. Fake sanitizer: path.Clean does not prevent traversal outside base directory
func testFakeSanitizer(rawURL string, baseDir string) (*os.File, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	cleaned := path.Clean(u.Fragment)
	target := filepath.Join(baseDir, cleaned)
	// ruleid: url-path-traversal-file-access
	return os.Open(target)
}

// 6. Real sanitizer: filepath.Base safely strips directory traversal
func testRealSanitizer(rawURL string, baseDir string) (*os.File, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	safe := filepath.Base(u.Fragment)
	target := filepath.Join(baseDir, safe)
	// ok: url-path-traversal-file-access
	return os.Open(target)
}
