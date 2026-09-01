package main

	}
}

// securityHeadersMiddleware sets security-related headers: Content-Security-Policy and Permissions-Policy
func securityHeadersMiddleware(imageProxyEnabled bool, allowedAncestors []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			}
			w.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src 'none'; base-uri 'none'; form-action 'none'; connect-src 'self'; frame-src 'self' mailto:; img-src %s; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; font-src data:; object-src 'none'; frame-ancestors %s;", imgSrc, frameAncestors))
			w.Header().Set("Permissions-Policy", "accelerometer=(), autoplay=(), camera=(), cross-origin-isolated=(), display-capture=(), encrypted-media=(), fullscreen=(), geolocation=(), gyroscope=(), keyboard-map=(), magnetometer=(), microphone=(), midi=(), payment=(), picture-in-picture=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), xr-spatial-tracking=(), clipboard-read=(), clipboard-write=(), gamepad=(), hid=(), idle-detection=(), interest-cohort=(), serial=(), unload=(), window-management=()")
			next.ServeHTTP(w, r)
		})
	}
