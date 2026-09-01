package main

package markup

import (
	"regexp"
	"sync"

	"github.com/microcosm-cc/bluemonday"
)

	defaultPolicy     *bluemonday.Policy
	descriptionPolicy *bluemonday.Policy
	rendererPolicies  map[string]*bluemonday.Policy
	allowAllRegex     *regexp.Regexp
}

var (
	defaultSanitizer     *Sanitizer
	defaultSanitizerOnce sync.Once
)

func GetDefaultSanitizer() *Sanitizer {
	defaultSanitizerOnce.Do(func() {
		defaultSanitizer = &Sanitizer{
			rendererPolicies: map[string]*bluemonday.Policy{},
			allowAllRegex:    regexp.MustCompile(".+"),
		}
		for name, renderer := range renderers {
			sanitizerRules := renderer.SanitizerRules()
			if len(sanitizerRules) > 0 {
				policy := defaultSanitizer.createDefaultPolicy()
				defaultSanitizer.addSanitizerRules(policy, sanitizerRules)
				defaultSanitizer.rendererPolicies[name] = policy
			}
		}
		defaultSanitizer.defaultPolicy = defaultSanitizer.createDefaultPolicy()
		defaultSanitizer.descriptionPolicy = defaultSanitizer.createRepoDescriptionPolicy()
	})
	return defaultSanitizer
}

func ResetDefaultSanitizerForTesting() {
	defaultSanitizer = nil
	defaultSanitizerOnce = sync.Once{}
}
// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markup

import (
	"io"
	"net/url"
	"regexp"

	"code.gitea.io/gitea/modules/setting"

	"github.com/microcosm-cc/bluemonday"
)

func (st *Sanitizer) createDefaultPolicy() *bluemonday.Policy {
	policy := bluemonday.UGCPolicy()

	// For JS code copy and Mermaid loading state
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^code-block( is-loading)?$`)).OnElements("pre")

	// For code preview
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^code-preview-[-\w]+( file-content)?$`)).Globally()
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^lines-num$`)).OnElements("td")
	policy.AllowAttrs("data-line-number").OnElements("span")
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^lines-code chroma$`)).OnElements("td")
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^code-inner$`)).OnElements("div")

	// For code preview (unicode escape)
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^file-view( unicode-escaped)?$`)).OnElements("table")
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^lines-escape$`)).OnElements("td")
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^toggle-escape-button btn interact-bg$`)).OnElements("a") // don't use button, button might submit a form
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^(ambiguous-code-point|escaped-code-point|broken-code-point)$`)).OnElements("span")
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^char$`)).OnElements("span")
	policy.AllowAttrs("data-tooltip-content", "data-escaped").OnElements("span")

	// For color preview
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^color-preview$`)).OnElements("span")

	// For attention
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^attention-header attention-\w+$`)).OnElements("blockquote")
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^attention-\w+$`)).OnElements("strong")
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^attention-icon attention-\w+ svg octicon-[\w-]+$`)).OnElements("svg")
	policy.AllowAttrs("viewBox", "width", "height", "aria-hidden").OnElements("svg")
	policy.AllowAttrs("fill-rule", "d").OnElements("path")

	// For Chroma markdown plugin
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^(chroma )?language-[\w-]+( display)?( is-loading)?$`)).OnElements("code")

	// Checkboxes
	policy.AllowAttrs("type").Matching(regexp.MustCompile(`^checkbox$`)).OnElements("input")
	policy.AllowAttrs("checked", "disabled", "data-source-position").OnElements("input")

	// Custom URL-Schemes
	if len(setting.Markdown.CustomURLSchemes) > 0 {
		policy.AllowURLSchemes(setting.Markdown.CustomURLSchemes...)
	} else {
		policy.AllowURLSchemesMatching(st.allowAllRegex)

		// Even if every scheme is allowed, these three are blocked for security reasons
		disallowScheme := func(*url.URL) bool {
			return false
		}
		policy.AllowURLSchemeWithCustomPolicy("javascript", disallowScheme)
		policy.AllowURLSchemeWithCustomPolicy("vbscript", disallowScheme)
		policy.AllowURLSchemeWithCustomPolicy("data", disallowScheme)
	}

	// Allow classes for anchors
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`ref-issue( ref-external-issue)?`)).OnElements("a")

	// Allow classes for task lists
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`task-list-item`)).OnElements("li")

	// Allow classes for org mode list item status.
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^(unchecked|checked|indeterminate)$`)).OnElements("li")

	// Allow icons
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^icon(\s+[\p{L}\p{N}_-]+)+$`)).OnElements("i")

	// Allow classes for emojis
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`emoji`)).OnElements("img")

	// Allow icons, emojis, chroma syntax and keyword markup on span
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^((icon(\s+[\p{L}\p{N}_-]+)+)|(emoji)|(language-math display)|(language-math inline))$|^([a-z][a-z0-9]{0,2})$|^` + keywordClass + `$`)).OnElements("span")

	// Allow 'color' and 'background-color' properties for the style attribute on text elements.
	policy.AllowStyles("color", "background-color").OnElements("span", "p")

	// Allow generally safe attributes
	generalSafeAttrs := []string{
		"abbr", "accept", "accept-charset",
		"accesskey", "action", "align", "alt",
		"aria-describedby", "aria-hidden", "aria-label", "aria-labelledby",
		"axis", "border", "cellpadding", "cellspacing", "char",
		"charoff", "charset", "checked",
		"clear", "cols", "colspan", "color",
		"compact", "coords", "datetime", "dir",
		"disabled", "enctype", "for", "frame",
		"headers", "height", "hreflang",
		"hspace", "ismap", "label", "lang",
		"maxlength", "media", "method",
		"multiple", "name", "nohref", "noshade",
		"nowrap", "open", "prompt", "readonly", "rel", "rev",
		"rows", "rowspan", "rules", "scope",
		"selected", "shape", "size", "span",
		"start", "summary", "tabindex", "target",
		"title", "type", "usemap", "valign", "value",
		"vspace", "width", "itemprop",
	}

	generalSafeElements := []string{
		"h1", "h2", "h3", "h4", "h5", "h6", "h7", "h8", "br", "b", "i", "strong", "em", "a", "pre", "code", "img", "tt",
		"div", "ins", "del", "sup", "sub", "p", "ol", "ul", "table", "thead", "tbody", "tfoot", "blockquote", "label",
		"dl", "dt", "dd", "kbd", "q", "samp", "var", "hr", "ruby", "rt", "rp", "li", "tr", "td", "th", "s", "strike", "summary",
		"details", "caption", "figure", "figcaption",
		"abbr", "bdo", "cite", "dfn", "mark", "small", "span", "time", "video", "wbr",
	}

	policy.AllowAttrs(generalSafeAttrs...).OnElements(generalSafeElements...)

	policy.AllowAttrs("src", "autoplay", "controls").OnElements("video")

	policy.AllowAttrs("itemscope", "itemtype").OnElements("div")

	// FIXME: Need to handle longdesc in img but there is no easy way to do it

	// Custom keyword markup
	defaultSanitizer.addSanitizerRules(policy, setting.ExternalSanitizerRules)

	return policy
}

// Sanitize takes a string that contains a HTML fragment or document and applies policy whitelist.
func Sanitize(s string) string {
	return GetDefaultSanitizer().defaultPolicy.Sanitize(s)
}

// SanitizeReader sanitizes a Reader
func SanitizeReader(r io.Reader, renderer string, w io.Writer) error {
	policy, exist := GetDefaultSanitizer().rendererPolicies[renderer]
	if !exist {
		policy = GetDefaultSanitizer().defaultPolicy
	}
	return policy.SanitizeReaderToWriter(r, w)
}
// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markup

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

// createRepoDescriptionPolicy returns a minimal more strict policy that is used for
// repository descriptions.
func (st *Sanitizer) createRepoDescriptionPolicy() *bluemonday.Policy {
	policy := bluemonday.NewPolicy()
	policy.AllowStandardURLs()

	// Allow italics and bold.
	policy.AllowElements("i", "b", "em", "strong")

	// Allow code.
	policy.AllowElements("code")

	// Allow links
	policy.AllowAttrs("href", "target", "rel").OnElements("a")

	// Allow classes for emojis
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^emoji$`)).OnElements("img", "span")
	policy.AllowAttrs("aria-label").OnElements("span")

	return policy
}

// SanitizeDescription sanitizes the HTML generated for a repository description.
func SanitizeDescription(s string) string {
	return GetDefaultSanitizer().descriptionPolicy.Sanitize(s)
}
