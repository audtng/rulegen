package rules

import (
	"net/http"
)

// Case 1: Reverse proxy conditionally setting Remote-User without deleting incoming header
func vulnReverseProxyDirector(req *http.Request, authResp *http.Response) {
	user := authResp.Header.Get("Remote-User")
	// ruleid: http-reverse-proxy-unstripped-identity-headers
	if user != "" {
		req.Header.Set("Remote-User", user)
	}
}

// Case 2: Reverse proxy deleting client-supplied Remote-User prior to conditional set
func safeReverseProxyDirectorPreDel(req *http.Request, authResp *http.Response) {
	user := authResp.Header.Get("Remote-User")
	req.Header.Del("Remote-User")
	// ok: http-reverse-proxy-unstripped-identity-headers
	if user != "" {
		req.Header.Set("Remote-User", user)
	}
}

// Case 3: Reverse proxy deleting client-supplied header in else branch
func safeReverseProxyDirectorElseDel(req *http.Request, authResp *http.Response) {
	user := authResp.Header.Get("X-Remote-User")
	// ok: http-reverse-proxy-unstripped-identity-headers
	if user != "" {
		req.Header.Set("X-Remote-User", user)
	} else {
		req.Header.Del("X-Remote-User")
	}
}

// Case 4: Inline check setting X-Forwarded-User without stripping client-supplied header
func vulnForwardAuthInlineCheck(req *http.Request, authResp *http.Response) {
	// ruleid: http-reverse-proxy-unstripped-identity-headers
	if v := authResp.Header.Get("X-Forwarded-User"); len(v) > 0 {
		req.Header.Set("X-Forwarded-User", v)
	}
}

// Case 5: Setting a standard non-identity header conditionally
func safeStandardHeaderCopy(req *http.Request, authResp *http.Response) {
	// ok: http-reverse-proxy-unstripped-identity-headers
	if ct := authResp.Header.Get("Content-Type"); ct != "" {
		req.Header.Set("Content-Type", ct)
	}
}

// Case 6: Forward-auth handler setting identity header upon authentication without stripping
func vulnForwardAuthSessionCheck(req *http.Request, user string, authenticated bool) {
	// ruleid: http-reverse-proxy-unstripped-identity-headers
	if authenticated && user != "" {
		req.Header.Set("X-Auth-Request-User", user)
	}
}

// Case 7: Forward-auth handler deleting header prior to setting it
func safeForwardAuthSessionCheck(req *http.Request, user string, authenticated bool) {
	req.Header.Del("X-Auth-Request-User")
	// ok: http-reverse-proxy-unstripped-identity-headers
	if authenticated && user != "" {
		req.Header.Set("X-Auth-Request-User", user)
	}
}
