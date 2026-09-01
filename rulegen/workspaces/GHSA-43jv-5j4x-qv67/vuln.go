package main

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/rs/zerolog"

	"github.com/dadrus/heimdall/internal/rules/config"
	"github.com/dadrus/heimdall/internal/rules/rule"
	"github.com/dadrus/heimdall/internal/x/errorchain"
)

type ruleImpl struct {
		// unescape path
		request.URL.RawPath = ""
	case config.EncodedSlashesOff:
		if strings.Contains(request.URL.RawPath, "%2F") {
			return nil, errorchain.NewWithMessage(heimdall.ErrArgument,
				"path contains encoded slash, which is not allowed")
		}
	// unescape captures
	captures := request.URL.Captures
	for k, v := range captures {
		captures[k] = unescape(v, r.slashesHandling)
	}

	// authenticators
func (b backend) URL() *url.URL { return b.targetURL }

func (b backend) ForwardHostHeader() bool { return b.forwardHostHeader }

func unescape(value string, handling config.EncodedSlashesHandling) string {
	if handling == config.EncodedSlashesOn {
		unescaped, _ := url.PathUnescape(value)

		return unescaped
	}

	unescaped, _ := url.PathUnescape(strings.ReplaceAll(value, "%2F", "$$$escaped-slash$$$"))

	return strings.ReplaceAll(unescaped, "$$$escaped-slash$$$", "%2F")
}
import (
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/dadrus/heimdall/internal/rules/config"
	"github.com/dadrus/heimdall/internal/x/errorchain"
	"github.com/dadrus/heimdall/internal/x/slicex"
)

var (
	if len(request.URL.RawPath) != 0 {
		switch m.slashHandling {
		case config.EncodedSlashesOff:
			if strings.Contains(request.URL.RawPath, "%2F") {
				return errorchain.NewWithMessage(ErrRequestPathMismatch,
					"request path contains encoded slashes which are not allowed")
			}
		case config.EncodedSlashesOn:
			value, _ = url.PathUnescape(value)
		default:
			unescaped, _ := url.PathUnescape(strings.ReplaceAll(value, "%2F", "$$$escaped-slash$$$"))
			value = strings.ReplaceAll(unescaped, "$$$escaped-slash$$$", "%2F")
		}
	}

