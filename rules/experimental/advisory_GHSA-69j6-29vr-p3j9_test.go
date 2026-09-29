package rules

import (
	"crypto/subtle"
	"net/http"
)

var expectedAPIKey = "super-secret-key"

func vulnConstantTimeCompare(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-API-Key")
	// ruleid: http-empty-token-auth-bypass
	if subtle.ConstantTimeCompare([]byte(token), []byte(expectedAPIKey)) == 1 {
		w.Write([]byte("authenticated"))
	}
}

func vulnStringCompare(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-API-Key")
	// ruleid: http-empty-token-auth-bypass
	if token == expectedAPIKey {
		w.Write([]byte("authenticated"))
	}
}

func vulnMissingReturn(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-API-Key")
	// ruleid: http-empty-token-auth-bypass
	if subtle.ConstantTimeCompare([]byte(token), []byte(expectedAPIKey)) != 1 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	w.Write([]byte("authenticated"))
}

func vulnQueryTokenMissingReturn(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	// ruleid: http-empty-token-auth-bypass
	if token != expectedAPIKey {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.Write([]byte("authenticated"))
}

func safeConstantTimeCompare(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-API-Key")
	if token == "" || expectedAPIKey == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// ok: http-empty-token-auth-bypass
	if subtle.ConstantTimeCompare([]byte(token), []byte(expectedAPIKey)) == 1 {
		w.Write([]byte("authenticated"))
	}
}

func safeStringCompare(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-API-Key")
	if len(token) == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// ok: http-empty-token-auth-bypass
	if token == expectedAPIKey {
		w.Write([]byte("authenticated"))
	}
}

func safeWithReturn(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-API-Key")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// ok: http-empty-token-auth-bypass
	if subtle.ConstantTimeCompare([]byte(token), []byte(expectedAPIKey)) != 1 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Write([]byte("authenticated"))
}

func safeQueryTokenWithReturn(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// ok: http-empty-token-auth-bypass
	if token != expectedAPIKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte("authenticated"))
}
