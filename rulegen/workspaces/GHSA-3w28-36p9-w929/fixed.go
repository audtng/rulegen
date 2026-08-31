package main


	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/macaron.v1"

	"gogs.io/gogs/internal/markup"
)

func ipynbSanitizer() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class", "data-prompt-number").OnElements("div")
	p.AllowAttrs("class").OnElements("img")
	// Only allow data URIs with safe image MIME types to prevent XSS via
	// "data:text/html" payloads.
	p.AllowURLSchemeWithCustomPolicy("data", markup.IsSafeDataURI)
	return p
}


		// Only allow data URIs with safe image MIME types to prevent XSS via
		// "data:text/html" payloads.
		sanitizer.policy.AllowURLSchemeWithCustomPolicy("data", IsSafeDataURI)

		// Custom URL-Schemes
		sanitizer.policy.AllowURLSchemes(conf.Markdown.CustomURLSchemes...)
	})
}

// IsSafeDataURI returns whether the given data URI uses a safe image MIME type.
func IsSafeDataURI(u *url.URL) bool {
	// The opaque data of a data URI has the form "mediatype;base64,data" or
	// "mediatype,data". We only allow common image MIME types.
	mediatype, _, _ := strings.Cut(u.Opaque, ";")
