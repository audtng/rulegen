package rules

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type GinContext struct{}

func (g *GinContext) Query(key string) string { return "" }
func (g *GinContext) Param(key string) string { return "" }
func (g *GinContext) File(filepath string)    {}

type PluginConfig struct {
	Executable string `json:"executable"`
}

func testVulnerableHTTP(req *http.Request) {
	filename := req.URL.Query().Get("file")
	target := filepath.Join("/opt/plugins", filename)
	// ruleid: tainted-path-traversal
	_, _ = os.Open(target)
}

func testSafeBaseSanitizer(req *http.Request) {
	filename := req.URL.Query().Get("file")
	safe := filepath.Base(filename)
	target := filepath.Join("/opt/plugins", safe)
	// ok: tainted-path-traversal
	_, _ = os.Open(target)
}

func testSafeIsLocalSanitizer(req *http.Request) {
	filename := req.URL.Query().Get("file")
	if filepath.IsLocal(filename) {
		target := filepath.Join("/opt/plugins", filename)
		// ok: tainted-path-traversal
		_, _ = os.Open(target)
	}
}

func testSafeOpenRoot(req *http.Request) {
	filename := req.URL.Query().Get("file")
	root, err := os.OpenRoot("/opt/plugins")
	if err != nil {
		return
	}
	defer root.Close()
	// ok: tainted-path-traversal
	_, _ = root.Stat(filename)
}

func testVulnerableJSON(data []byte) {
	var cfg PluginConfig
	_ = json.Unmarshal(data, &cfg)
	target := filepath.Join("/opt/plugins", cfg.Executable)
	// ruleid: tainted-path-traversal
	_, _ = os.Stat(target)
}

func testVulnerableGin(c *GinContext) {
	filename := c.Query("file")
	target := filepath.Join("/opt/plugins", filename)
	// ruleid: tainted-path-traversal
	c.File(target)
}
