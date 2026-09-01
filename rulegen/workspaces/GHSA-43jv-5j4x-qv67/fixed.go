package main

// Copyright 2026 Dimitrij Drus <dadrus@gmx.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package urlx

import (
	"net/url"
	"strings"
)

// ContainsEncodedSlash reports whether path contains a URL-encoded slash
// sequence, case-insensitive, e.g. %2F or %2f.
func ContainsEncodedSlash(path string) bool {
	for i := strings.IndexByte(path, '%'); i != -1; {
		if i+2 < len(path) && path[i+1] == '2' && (path[i+2]|0x20) == 'f' { //nolint:mnd
			return true
		}

		next := strings.IndexByte(path[i+1:], '%')
		if next == -1 {
			break
		}

		i += next + 1
	}

	return false
}

// Unescape decodes URL-escaped path value.
// If decodeEncodedSlash is false, encoded slashes (%2F / %2f) are preserved.
func Unescape(value string, decodeEncodedSlash bool) string { //nolint:cyclop
	start := strings.IndexByte(value, '%')
	if start == -1 {
		return value
	}

	if decodeEncodedSlash {
		unescaped, _ := url.PathUnescape(value)

		return unescaped
	}

	var builder strings.Builder
	builder.Grow(len(value))
	builder.WriteString(value[:start])

	for i := start; i < len(value); {
		j := i
		for j < len(value) && value[j] != '%' {
			j++
		}

		if j > i {
			builder.WriteString(value[i:j])
			i = j
		}

		if i >= len(value) {
			break
		}

		if i+2 >= len(value) {
			builder.WriteByte(value[i])
			i++

			continue
		}

		hi := hexValue(value[i+1])
		lo := hexValue(value[i+2])

		if hi == 0xFF || lo == 0xFF { //nolint:mnd
			builder.WriteByte(value[i])
			i++

			continue
		}

		decoded := (hi << 4) | lo //nolint:mnd
		if decoded == '/' {
			builder.WriteByte(value[i])
			builder.WriteByte(value[i+1])
			builder.WriteByte(value[i+2])
			i += 3

			continue
		}

		builder.WriteByte(decoded)
		i += 3
	}

	return builder.String()
}

func hexValue(ch byte) byte {
	switch {
	case ch >= '0' && ch <= '9':
		return ch - '0'
	case ch >= 'a' && ch <= 'f':
		return ch - 'a' + 10 //nolint:mnd
	case ch >= 'A' && ch <= 'F':
		return ch - 'A' + 10 //nolint:mnd
	default:
		return 0xFF //nolint:mnd
	}
}
import (
	"bytes"
	"net/url"

	"github.com/rs/zerolog"

	"github.com/dadrus/heimdall/internal/rules/config"
	"github.com/dadrus/heimdall/internal/rules/rule"
	"github.com/dadrus/heimdall/internal/x/errorchain"
	"github.com/dadrus/heimdall/internal/x/urlx"
)

type ruleImpl struct {
		// unescape path
		request.URL.RawPath = ""
	case config.EncodedSlashesOff:
		if urlx.ContainsEncodedSlash(request.URL.RawPath) {
			return nil, errorchain.NewWithMessage(heimdall.ErrArgument,
				"path contains encoded slash, which is not allowed")
		}
	// unescape captures
	captures := request.URL.Captures
	for k, v := range captures {
		captures[k] = urlx.Unescape(v, r.slashesHandling == config.EncodedSlashesOn)
	}

	// authenticators
func (b backend) URL() *url.URL { return b.targetURL }

func (b backend) ForwardHostHeader() bool { return b.forwardHostHeader }
import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/dadrus/heimdall/internal/rules/config"
	"github.com/dadrus/heimdall/internal/x/errorchain"
	"github.com/dadrus/heimdall/internal/x/slicex"
	"github.com/dadrus/heimdall/internal/x/urlx"
)

var (
	if len(request.URL.RawPath) != 0 {
		switch m.slashHandling {
		case config.EncodedSlashesOff:
			if urlx.ContainsEncodedSlash(request.URL.RawPath) {
				return errorchain.NewWithMessage(ErrRequestPathMismatch,
					"request path contains encoded slashes which are not allowed")
			}
		default:
			value = urlx.Unescape(value, m.slashHandling == config.EncodedSlashesOn)
		}
	}

