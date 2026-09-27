package rules

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (g *GinContext) Query(key string) string { return "" }
func (g *GinContext) File(filepath string)    {}

type LogConfig struct {
	LogPath string `json:"log_path"`
}

func testVulnerableHTTP(req *http.Request) {
	filename := req.URL.Query().Get("file")
	target := filepath.Join("/var/log", filename)
	// ruleid: path-traversal-arbitrary-file-access
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
}

func testSafeHTTPBase(req *http.Request) {
	filename := req.URL.Query().Get("file")
	safe := filepath.Base(filename)
	target := filepath.Join("/var/log", safe)
	// ok: path-traversal-arbitrary-file-access
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
}

func testVulnerableJSONConfig(data []byte) {
	var cfg LogConfig
	_ = json.Unmarshal(data, &cfg)
	target := filepath.Join("/var/log", cfg.LogPath)
	// ruleid: path-traversal-arbitrary-file-access
	os.Open(target)
}

func testSafeHTTPIsLocal(req *http.Request) {
	p := req.FormValue("path")
	if filepath.IsLocal(p) {
		target := filepath.Join("/var/log", p)
		// ok: path-traversal-arbitrary-file-access
		os.Open(target)
	}
}

func testVulnerableGin(c *GinContext) {
	filename := c.Query("file")
	target := filepath.Join("/var/log", filename)
	// ruleid: path-traversal-arbitrary-file-access
	c.File(target)
}

func testSafeConstantPath() {
	target := filepath.Join("/var/log", "server.log")
	// ok: path-traversal-arbitrary-file-access
	os.Open(target)
}
