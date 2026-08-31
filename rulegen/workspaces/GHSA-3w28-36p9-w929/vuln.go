package main


	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/macaron.v1"
)

func ipynbSanitizer() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class", "data-prompt-number").OnElements("div")
	p.AllowAttrs("class").OnElements("img")
	p.AllowURLSchemes("data")
	return p
}


		// Only allow data URIs with safe image MIME types to prevent XSS via
		// "data:text/html" payloads.
		sanitizer.policy.AllowURLSchemeWithCustomPolicy("data", isSafeDataURI)

		// Custom URL-Schemes
		sanitizer.policy.AllowURLSchemes(conf.Markdown.CustomURLSchemes...)
	})
}

// isSafeDataURI returns whether the given data URI uses a safe image MIME type.
func isSafeDataURI(u *url.URL) bool {
	// The opaque data of a data URI has the form "mediatype;base64,data" or
	// "mediatype,data". We only allow common image MIME types.
	mediatype, _, _ := strings.Cut(u.Opaque, ";")
