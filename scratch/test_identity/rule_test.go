package rules

import (
	"context"
	"net/http"
)

type Session struct {
	User      string
	UserAgent string
}

type AuthContext struct {
	Authorized bool
	User       string
}

func isTrustedProxy(addr string) bool {
	return addr == "127.0.0.1"
}

// Case 1: Vuln - Directly trusting X-Forwarded-User header into session/auth context
func vulnForwardedUser(w http.ResponseWriter, r *http.Request, s *Session) {
	// ruleid: http-identity-header-spoofing
	user := r.Header.Get("X-Forwarded-User")
	s.User = user
}

// Safe Case 1: Verifying the proxy via RemoteAddr or isTrustedProxy before extracting X-Forwarded-User
func safeForwardedUser(w http.ResponseWriter, r *http.Request, s *Session) {
	// ok: http-identity-header-spoofing
	if !isTrustedProxy(r.RemoteAddr) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	user := r.Header.Get("X-Forwarded-User")
	s.User = user
}

// Case 2: Vuln - Authenticated via shared token, but user identity is read from request header
func vulnSharedTokenHeaderUser(w http.ResponseWriter, r *http.Request, sharedToken string) *AuthContext {
	token := r.Header.Get("Authorization")
	if token == sharedToken {
		// ruleid: http-identity-header-spoofing
		return &AuthContext{
			Authorized: true,
			User:       r.Header.Get("X-User-Id"),
		}
	}
	return nil
}

// Safe Case 2: Using server-configured fixed owner for shared token
func safeSharedTokenFixedOwner(w http.ResponseWriter, r *http.Request, sharedToken, defaultOwner string) *AuthContext {
	token := r.Header.Get("Authorization")
	if token == sharedToken {
		// ok: http-identity-header-spoofing
		return &AuthContext{
			Authorized: true,
			User:       defaultOwner,
		}
	}
	return nil
}

// Case 3: Vuln - Storing unverified X-Remote-User in request context
func vulnRemoteUserContext(r *http.Request) *http.Request {
	// ruleid: http-identity-header-spoofing
	user := r.Header.Get("X-Remote-User")
	ctx := context.WithValue(r.Context(), "user", user)
	return r.WithContext(ctx)
}

// Safe Case 3: Verifying TLS client certificate before trusting X-Remote-User
func safeRemoteUserContext(r *http.Request) *http.Request {
	// ok: http-identity-header-spoofing
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		return r
	}
	user := r.Header.Get("X-Remote-User")
	ctx := context.WithValue(r.Context(), "user", user)
	return r.WithContext(ctx)
}

// Case 4: Vuln - Assigning identity header directly to session user
func vulnAssignIdentityHeader(r *http.Request, s *Session) {
	// ruleid: http-identity-header-spoofing
	s.User = r.Header.Get("X-Authenticated-User")
}

// Safe Case 4: Reading standard non-identity header
func safeStandardHeader(r *http.Request, s *Session) {
	// ok: http-identity-header-spoofing
	userAgent := r.Header.Get("User-Agent")
	s.UserAgent = userAgent
}
