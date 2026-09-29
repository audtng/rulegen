package rules

import (
	"net/http"
	"net/url"
)

type handler struct {
	token string
}

func (h handler) testVulnQueryCookieAuthBypass(w http.ResponseWriter, r *http.Request) {
	if h.token != "" {
		values := r.URL.Query()
		// ruleid: http-token-auth-bypass
		token := values.Get("auth_token")
		if token != "" {
			tokenCookie := http.Cookie{
				Name:  "auth_token",
				Value: token,
			}
			http.SetCookie(w, &tokenCookie)
		} else {
			cookie, err := r.Cookie("auth_token")
			if err != nil || cookie.Value != h.token {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}
}

func (h handler) testSafeQueryCookieAuth(w http.ResponseWriter, r *http.Request) {
	if h.token != "" {
		values := r.URL.Query()
		// ok: http-token-auth-bypass
		token := values.Get("auth_token")
		if token != "" {
			tokenCookie := http.Cookie{
				Name:  "auth_token",
				Value: token,
			}
			http.SetCookie(w, &tokenCookie)
		} else {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token = cookie.Value
		}

		if token != h.token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}
