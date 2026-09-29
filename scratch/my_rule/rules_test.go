package rules

import (
	"net/http"
)

type proxyHandler struct {
	token string
}

func (h proxyHandler) testVulnQueryCookieAuthBypass(w http.ResponseWriter, r *http.Request) {
	if h.token != "" {
		values := r.URL.Query()
		token := values.Get("auth_token")
		// ruleid: http-token-auth-bypass
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

func (h proxyHandler) testSafeQueryCookieAuthValidated(w http.ResponseWriter, r *http.Request) {
	if h.token != "" {
		values := r.URL.Query()
		token := values.Get("auth_token")
		// ok: http-token-auth-bypass
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

func testVulnIgnoredCookieError(w http.ResponseWriter, r *http.Request) {
	// ruleid: http-token-auth-bypass
	cookie, _ := r.Cookie("auth_token")
	if cookie != nil && cookie.Value == "admin" {
		w.Write([]byte("ok"))
	}
}

func testSafeCookieErrorCheck(w http.ResponseWriter, r *http.Request) {
	// ok: http-token-auth-bypass
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if cookie.Value == "admin" {
		w.Write([]byte("ok"))
	}
}

func testVulnMissingReturnCookieMismatch(w http.ResponseWriter, r *http.Request, expectedToken string) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// ruleid: http-token-auth-bypass
	if cookie.Value != expectedToken {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.Write([]byte("sensitive data"))
}

func testSafeCookieMismatchReturn(w http.ResponseWriter, r *http.Request, expectedToken string) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// ok: http-token-auth-bypass
	if cookie.Value != expectedToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte("sensitive data"))
}

func testVulnMissingReturnTokenMismatch(w http.ResponseWriter, r *http.Request, expectedToken string) {
	token := r.URL.Query().Get("token")
	// ruleid: http-token-auth-bypass
	if token != expectedToken {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
	w.Write([]byte("sensitive data"))
}

func testSafeTokenMismatchReturn(w http.ResponseWriter, r *http.Request, expectedToken string) {
	token := r.URL.Query().Get("token")
	// ok: http-token-auth-bypass
	if token != expectedToken {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	w.Write([]byte("sensitive data"))
}
