package main


import (
	"net/url"
	pathpkg "path"
	"strings"
)

//nolint:gocognit,gocyclo,gocyclo,cyclop,funlen
func PathHasDotSegments(path string) bool {
	iDot := strings.IndexByte(path, '.')
	iPct := strings.IndexByte(path, '%')
	iBsl := strings.IndexByte(path, '\\')

	idx := iDot
	if idx == -1 || (iPct != -1 && iPct < idx) {
		idx = iPct
	}

	if idx == -1 || (iBsl != -1 && iBsl < idx) {
		idx = iBsl
	}

	if idx == -1 {
		return false
	}

	segLen := 0
	for i := idx - 1; i >= 0 && path[i] != '/'; i-- {
		segLen++
	}

	dotCount := 0

	for i := idx; i < len(path); {
		switch path[i] {
		case '/', '\\':
			if (segLen == 1 && dotCount == 1) || (segLen == 2 && dotCount == 2) {
				return true
			}

			segLen = 0
			dotCount = 0
			i++
		case '.':
			segLen++
			dotCount++
			i++
		case '%':
			if i+2 >= len(path) {
				segLen++
				i++

				continue
			}

			h1 := path[i+1]
			h2 := path[i+2] | 0x20 //nolint:mnd

			switch {
			case h1 == '2' && h2 == 'e':
				segLen++
				dotCount++
				i += 3
			case h1 == '2' && h2 == 'f':
				if (segLen == 1 && dotCount == 1) || (segLen == 2 && dotCount == 2) {
					return true
				}

				segLen = 0
				dotCount = 0
				i += 3
			case h1 == '5' && h2 == 'c':
				if (segLen == 1 && dotCount == 1) || (segLen == 2 && dotCount == 2) {
					return true
				}

				segLen = 0
				dotCount = 0
				i += 3
			default:
				segLen++
				i++
			}
		default:
			segLen++
			i++
		}
	}

	return (segLen == 1 && dotCount == 1) || (segLen == 2 && dotCount == 2)
}

func NormalizePath(path string) string {
	if path == "/" {
		return path
	}

	hasTrailingSlash := strings.HasSuffix(path, "/")
	path = pathpkg.Clean(path)

	if hasTrailingSlash && path != "/" {
		path += "/"
	}

	return path
}

// ContainsEncodedSlash reports whether path contains a URL-encoded slash
// sequence, case-insensitive, e.g. %2F or %2f.
func ContainsEncodedSlash(path string) bool {

	"github.com/dadrus/heimdall/internal/heimdall"
	"github.com/dadrus/heimdall/internal/rules/rule"
	"github.com/dadrus/heimdall/internal/x"
	"github.com/dadrus/heimdall/internal/x/errorchain"
	"github.com/dadrus/heimdall/internal/x/urlx"
)

type ruleExecutor struct {
		Str("_url", request.URL.String()).
		Msg("Analyzing request")

	if urlx.PathHasDotSegments(x.IfThenElse(len(request.URL.RawPath) != 0, request.URL.RawPath, request.URL.Path)) {
		return nil, errorchain.NewWithMessage(heimdall.ErrArgument,
			"path contains dot segments, which are not allowed")
	}

	rul, err := e.r.FindRule(ctx)
	if err != nil {
		return nil, err
	"strings"

	"github.com/dadrus/heimdall/internal/x"
	"github.com/dadrus/heimdall/internal/x/urlx"
)

func extractURL(req *http.Request) url.URL {
		query = req.URL.RawQuery
	}

	rawPath = urlx.NormalizePath(rawPath)
	path, _ = url.PathUnescape(rawPath)

	return url.URL{
