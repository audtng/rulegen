package rules

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string {
	return "test"
}

func (c *GinContext) File(filepath string) {}

func testVulnerableHTTPPathJoinRemove(req *http.Request) {
	filename := req.URL.Query().Get("filename")
	targetPath := path.Join("/var/wiki/pages", filename+".md")
	// ruleid: path-traversal-arbitrary-file-deletion
	os.Remove(targetPath)
}

func testVulnerableHTTPPostFormFilepathJoinRemoveAll(req *http.Request) {
	dir := req.PostFormValue("dir")
	targetPath := filepath.Join("/var/uploads", dir)
	// ruleid: path-traversal-arbitrary-file-deletion
	os.RemoveAll(targetPath)
}

func testVulnerableGinQueryFile(c *GinContext) {
	page := c.Query("page")
	targetPath := fmt.Sprintf("/data/wiki/%s", page)
	// ruleid: path-traversal-arbitrary-file-deletion
	os.Remove(targetPath)
}

func testSafeFilePathBaseSanitizer(req *http.Request) {
	page := req.FormValue("page")
	cleanName := filepath.Base(page)
	targetPath := filepath.Join("/var/wiki/pages", cleanName)
	// ok: path-traversal-arbitrary-file-deletion
	os.Remove(targetPath)
}

func testSafeFilePathIsLocalValidation(req *http.Request) {
	relPath := req.FormValue("path")
	if !filepath.IsLocal(relPath) {
		return
	}
	targetPath := filepath.Join("/var/wiki/pages", relPath)
	// ok: path-traversal-arbitrary-file-deletion
	os.Remove(targetPath)
}

func testSafeConstantPath(req *http.Request) {
	const defaultPage = "Home.md"
	targetPath := filepath.Join("/var/wiki/pages", defaultPage)
	// ok: path-traversal-arbitrary-file-deletion
	os.Remove(targetPath)
}
