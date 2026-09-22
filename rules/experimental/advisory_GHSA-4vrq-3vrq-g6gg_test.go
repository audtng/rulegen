package rules

import (
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

// Case 1: URL fragment subdir component extracted and joined with repo root (vulnerable)
func TestVulnerableGitSubdirFragment(t *testing.T) {
	rawURL := "https://example.com/repo.git#main:../../etc/passwd"
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	_, subdir, _ := strings.Cut(u.Fragment, ":")
	target := filepath.Join("/repo/root", subdir)
	// ruleid: url-path-traversal
	_, _ = os.Open(target)
}

// Case 2: URL fragment subdir component validated with filepath.IsLocal (safe)
func TestSafeGitSubdirFragmentIsLocal(t *testing.T) {
	rawURL := "https://example.com/repo.git#main:sub/dir"
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	_, subdir, _ := strings.Cut(u.Fragment, ":")
	if !filepath.IsLocal(subdir) {
		t.Skip("path is not local")
	}
	target := filepath.Join("/repo/root", subdir)
	// ok: url-path-traversal
	_, _ = os.Open(target)
}

// Case 3: URL Path parameter directly joined into filesystem path (vulnerable)
func TestVulnerableURLPath(t *testing.T) {
	u, err := url.Parse("https://example.com/../../etc/shadow")
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join("/var/data", u.Path)
	// ruleid: url-path-traversal
	_, _ = os.ReadFile(target)
}

// Case 4: URL Path parameter sanitized with filepath.Base (safe)
func TestSafeURLPathBase(t *testing.T) {
	u, err := url.Parse("https://example.com/../../etc/shadow")
	if err != nil {
		t.Fatal(err)
	}
	safeName := filepath.Base(u.Path)
	target := filepath.Join("/var/data", safeName)
	// ok: url-path-traversal
	_, _ = os.ReadFile(target)
}

// Case 5: URL parsed with ParseRequestURI and cleaned with path.Clean (vulnerable: path.Clean does not prevent traversal for relative paths)
func TestVulnerablePathClean(t *testing.T) {
	rawURL := "/../../tmp/evil"
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	cleaned := path.Clean(u.Path)
	target := filepath.Join("/tmp/workspace", cleaned)
	// ruleid: url-path-traversal
	_ = os.RemoveAll(target)
}

// Case 6: URL Query parameter validated with fs.ValidPath (safe)
func TestSafeURLQueryValidPath(t *testing.T) {
	rawURL := "https://example.com/download?file=valid/path.txt"
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	fileName := u.Query().Get("file")
	if !fs.ValidPath(fileName) {
		t.Skip("invalid path")
	}
	target := filepath.Join("/tmp/workspace", fileName)
	// ok: url-path-traversal
	_ = os.WriteFile(target, []byte("contents"), 0644)
}
