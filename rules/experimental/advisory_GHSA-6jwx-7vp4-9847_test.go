package rules

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Case 1: Vulnerable - Trimming prefix from RawPath using prefix extracted from decoded Path
func TestVulnPrefixTrim(req *http.Request) {
	prefix := "/api"
	if strings.HasPrefix(req.URL.Path, prefix) {
		p := strings.TrimPrefix(req.URL.Path, prefix)
		// ruleid: go-url-path-rawpath-desync
		req.URL.RawPath = strings.TrimPrefix(req.URL.RawPath, p)
	}
}

// Case 2: Vulnerable - Using regex match against Path to replace in RawPath
func TestVulnRegexReplace(req *http.Request) {
	re := regexp.MustCompile("^/api/v[0-9]+")
	matched := re.FindString(req.URL.Path)
	if matched != "" {
		// ruleid: go-url-path-rawpath-desync
		req.URL.RawPath = strings.Replace(req.URL.RawPath, matched, "", 1)
	}
}

// Case 3: Vulnerable - Suffix trim directly from Path assigned to RawPath
func TestVulnSuffixTrim(req *http.Request) {
	clean := strings.TrimSuffix(req.URL.Path, "/")
	// ruleid: go-url-path-rawpath-desync
	req.URL.RawPath = clean
}

// Case 4: Patched - URL synchronized and canonicalized using JoinPath
func TestSafeJoinPath(req *http.Request) {
	prefix := "/api"
	if strings.HasPrefix(req.URL.Path, prefix) {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		// ok: go-url-path-rawpath-desync
		req.URL = req.URL.JoinPath()
	}
}

// Case 5: Patched - Clearing RawPath so EscapedPath computes canonical representation
func TestSafeClearRawPath(req *http.Request) {
	prefix := "/api"
	if strings.HasPrefix(req.URL.Path, prefix) {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		// ok: go-url-path-rawpath-desync
		req.URL.RawPath = ""
	}
}

// Case 6: Patched - Explicitly escaping path before assignment to RawPath
func TestSafePathEscape(req *http.Request) {
	trimmed := strings.TrimPrefix(req.URL.Path, "/api")
	// ok: go-url-path-rawpath-desync
	req.URL.RawPath = url.PathEscape(trimmed)
}
