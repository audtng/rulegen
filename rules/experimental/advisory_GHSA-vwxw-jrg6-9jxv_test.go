package rules

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
)

type GinContext struct{}

func (c *GinContext) Query(key string) string { return "" }
func (c *GinContext) File(filepath string)    {}

type EchoContext struct{}

func (c EchoContext) FormValue(key string) string { return "" }
func (c EchoContext) File(file string) error      { return nil }

type FiberCtx struct{}

func (c *FiberCtx) Query(key string, defaultValue ...string) string { return "" }
func (c *FiberCtx) SendFile(file string) error                      { return nil }

func vulnerableNetHTTP(r *http.Request) {
	userInput := r.URL.Query().Get("path")
	baseURL, _ := url.Parse("https://api.example.com/v1/")
	relURL, _ := url.Parse(userInput)
	// ruleid: url-path-traversal-resolve-reference
	target := baseURL.ResolveReference(relURL)
	_ = target
}

func vulnerableGin(c *GinContext) {
	param := c.Query("endpoint")
	baseURL, _ := url.Parse("https://api.example.com/service/")
	relURL, _ := url.Parse(param)
	// ruleid: url-path-traversal-resolve-reference
	resolved := baseURL.ResolveReference(relURL)
	_ = resolved
}

func vulnerableEcho(c EchoContext) {
	val := c.FormValue("path")
	baseURL, _ := url.Parse("https://api.example.com/base/")
	rel := &url.URL{Path: val}
	// ruleid: url-path-traversal-resolve-reference
	res := baseURL.ResolveReference(rel)
	_ = res
}

func vulnerableJSON(data []byte) {
	var payload struct {
		Target string `json:"target"`
	}
	_ = json.Unmarshal(data, &payload)
	baseURL, _ := url.Parse("https://api.example.com/base/")
	rel, _ := url.Parse(payload.Target)
	// ruleid: url-path-traversal-resolve-reference
	out := baseURL.ResolveReference(rel)
	_ = out
}

func safePathEscaped(r *http.Request) {
	userInput := r.URL.Query().Get("path")
	escaped := url.PathEscape(userInput)
	baseURL, _ := url.Parse("https://api.example.com/v1/")
	relURL, _ := url.Parse(escaped)
	// ok: url-path-traversal-resolve-reference
	target := baseURL.ResolveReference(relURL)
	_ = target
}

func safeIsLocal(r *http.Request) {
	userInput := r.URL.Query().Get("path")
	if !filepath.IsLocal(userInput) {
		return
	}
	baseURL, _ := url.Parse("https://api.example.com/v1/")
	relURL, _ := url.Parse(userInput)
	// ok: url-path-traversal-resolve-reference
	target := baseURL.ResolveReference(relURL)
	_ = target
}
