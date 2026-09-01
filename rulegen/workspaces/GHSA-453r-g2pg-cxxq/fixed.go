package main

type remoteProxyHandler struct {
	s         incus.InstanceServer
	transport http.RoundTripper
	url       string

	mu           *sync.RWMutex
	connections  *uint64
	*h.connections++
	h.mu.Unlock()

	// Don't allow cross-origin requests.
	origin := r.Header.Get("Origin")
	if origin != "" && origin != h.url {
		return
	}

	// Basic auth.
	if h.token != "" {
		// Parse query URL.

		token := values.Get("auth_token")
		if token != "" {
			// If a token was passed through the URL, persist it as a cookie.

			tokenCookie := http.Cookie{
				Name:     "auth_token",
				Value:    token,

			http.SetCookie(w, &tokenCookie)
		} else {
			// If not, attempt to pull it from the cookie.
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				// Fail authentication if no cookie can be found.
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token = cookie.Value
		}

		// Check the user token against the expected value.
		if token != h.token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

		transactions: &transactions,

		token: token,
		url:   "http://" + server.Addr().String(),
	}

	// Print address.
