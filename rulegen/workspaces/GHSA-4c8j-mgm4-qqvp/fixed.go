package main

	}
}

// securityHeadersMiddleware sets security-related headers:
// - Content-Security-Policy: controls which resources the browser is allowed to load
// - Permissions-Policy: disables browser features (camera, mic, etc.) not needed by a comment widget
// - X-Content-Type-Options: prevents browsers from MIME-sniffing responses away from the declared type,
//   stopping e.g. a user-uploaded image from being reinterpreted as executable HTML/JS
// - Referrer-Policy: controls how much URL information leaks in the Referer header on cross-origin
//   requests; "strict-origin-when-cross-origin" sends only the origin (no path) to other domains
//   and nothing at all on HTTPS→HTTP downgrades
func securityHeadersMiddleware(imageProxyEnabled bool, allowedAncestors []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			}
			w.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src 'none'; base-uri 'none'; form-action 'none'; connect-src 'self'; frame-src 'self' mailto:; img-src %s; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; font-src data:; object-src 'none'; frame-ancestors %s;", imgSrc, frameAncestors))
			w.Header().Set("Permissions-Policy", "accelerometer=(), autoplay=(), camera=(), cross-origin-isolated=(), display-capture=(), encrypted-media=(), fullscreen=(), geolocation=(), gyroscope=(), keyboard-map=(), magnetometer=(), microphone=(), midi=(), payment=(), picture-in-picture=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), xr-spatial-tracking=(), clipboard-read=(), clipboard-write=(), gamepad=(), hid=(), idle-detection=(), interest-cohort=(), serial=(), unload=(), window-management=()")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			next.ServeHTTP(w, r)
		})
	}
