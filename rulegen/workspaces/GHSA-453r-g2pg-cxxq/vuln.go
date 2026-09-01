package main

type remoteProxyHandler struct {
	s         incus.InstanceServer
	transport http.RoundTripper

	mu           *sync.RWMutex
	connections  *uint64
	*h.connections++
	h.mu.Unlock()

	// Basic auth.
	if h.token != "" {
		// Parse query URL.

		token := values.Get("auth_token")
		if token != "" {
			tokenCookie := http.Cookie{
				Name:     "auth_token",
				Value:    token,

			http.SetCookie(w, &tokenCookie)
		} else {
			cookie, err := r.Cookie("auth_token")
			if err != nil || cookie.Value != h.token {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}

		transactions: &transactions,

		token: token,
	}

	// Print address.
