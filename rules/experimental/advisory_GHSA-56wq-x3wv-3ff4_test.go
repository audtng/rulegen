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

func (c *GinContext) File(filepath string) {}

type EchoContext struct{}

func (c *EchoContext) FormValue(name string) string {
	return ""
}

// 1. Vulnerable HTTP header path traversal (CVE archetype: X-Amz-Copy-Source)
func VulnHttpCopySource(w http.ResponseWriter, r *http.Request) {
	source := r.Header.Get("X-Amz-Copy-Source")
	target := filepath.Join("/var/data", source)
	// ruleid: go-path-traversal
	os.Open(target)
}

// 2. Safe HTTP header path traversal with IsLocal boundary check
func SafeHttpCopySource(w http.ResponseWriter, r *http.Request) {
	source := r.Header.Get("X-Amz-Copy-Source")
	target := filepath.Join("/var/data", source)
	if !filepath.IsLocal(target) {
		return
	}
	// ok: go-path-traversal
	os.Open(target)
}

// 3. Vulnerable Gin framework file download
func VulnGinDownload(c *GinContext) {
	filename := c.Query("file")
	target := filepath.Join("/var/www/uploads", filename)
	// ruleid: go-path-traversal
	c.File(target)
}

// 4. Safe Gin framework file download using filepath.Base
func SafeGinDownload(c *GinContext) {
	filename := c.Query("file")
	clean := filepath.Base(filename)
	target := filepath.Join("/var/www/uploads", clean)
	// ok: go-path-traversal
	c.File(target)
}

// 5. Vulnerable Echo framework read file
func VulnEchoReadFile(c *EchoContext) {
	filename := c.FormValue("doc")
	target := filepath.Join("/srv/storage", filename)
	// ruleid: go-path-traversal
	os.ReadFile(target)
}

// 6. Safe Echo framework read file using filepath.Base
func SafeEchoReadFile(c *EchoContext) {
	filename := c.FormValue("doc")
	clean := filepath.Base(filename)
	target := filepath.Join("/srv/storage", clean)
	// ok: go-path-traversal
	os.ReadFile(target)
}
