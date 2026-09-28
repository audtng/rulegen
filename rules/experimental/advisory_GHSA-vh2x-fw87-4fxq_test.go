package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return ""
}

// Case 1: Direct net/http query parameter reaching os.Remove (vulnerable)
func VulnHttpQuery(req *http.Request) {
	p := req.URL.Query().Get("path")
	// ruleid: arbitrary-file-deletion
	os.Remove(p)
}

// Case 2: net/http query parameter validated with filepath.IsLocal guard (safe)
func SafeHttpIsLocal(req *http.Request) {
	p := req.URL.Query().Get("path")
	if !filepath.IsLocal(p) {
		return
	}
	// ok: arbitrary-file-deletion
	os.Remove(p)
}

// Case 3: net/http FormValue joined into a directory path reaching os.RemoveAll (vulnerable)
func VulnHttpJoinRemoveAll(req *http.Request) {
	p := req.FormValue("filepath")
	full := filepath.Join("/data/uploads", p)
	// ruleid: arbitrary-file-deletion
	os.RemoveAll(full)
}

// Case 4: net/http PostFormValue sanitized with filepath.Base (safe)
func SafeHttpBase(req *http.Request) {
	p := req.PostFormValue("filename")
	base := filepath.Base(p)
	full := filepath.Join("/data/uploads", base)
	// ok: arbitrary-file-deletion
	os.RemoveAll(full)
}

// Case 5: Gin query parameter with filepath.Clean but missing traversal check reaching os.Remove (vulnerable)
func VulnGinCleanOnly(c *GinContext) {
	p := c.Query("path")
	clean := filepath.Clean(p)
	// ruleid: arbitrary-file-deletion
	os.Remove(clean)
}

// Case 6: Gin query parameter with filepath.IsLocal validation and Clean (safe)
func SafeGinIsLocal(c *GinContext) {
	p := c.Query("path")
	if !filepath.IsLocal(p) {
		return
	}
	clean := filepath.Clean(p)
	// ok: arbitrary-file-deletion
	os.Remove(clean)
}
