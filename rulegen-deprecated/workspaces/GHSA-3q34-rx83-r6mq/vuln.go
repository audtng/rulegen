package main


import (
	"net/url"
	"strings"
)

// ContainsEncodedSlash reports whether path contains a URL-encoded slash
// sequence, case-insensitive, e.g. %2F or %2f.
func ContainsEncodedSlash(path string) bool {

	"github.com/dadrus/heimdall/internal/heimdall"
	"github.com/dadrus/heimdall/internal/rules/rule"
)

type ruleExecutor struct {
		Str("_url", request.URL.String()).
		Msg("Analyzing request")

	rul, err := e.r.FindRule(ctx)
	if err != nil {
		return nil, err
	"strings"

	"github.com/dadrus/heimdall/internal/x"
)

func extractURL(req *http.Request) url.URL {
		query = req.URL.RawQuery
	}

	path, _ = url.PathUnescape(rawPath)

	return url.URL{
