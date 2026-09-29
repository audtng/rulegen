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
	Owner      string
}

func isTrusted(addr string) bool {
	return addr == "127.0.0.1"
}

// Case 1: Struct initialization with untrusted identity header
func vulnStructAuthContext(r *http.Request) *AuthContext {
	// ruleid: http-identity-header-spoofing
	return &AuthContext{
		Authorized: true,
		User:       r.Header.Get("X-User-Id"),
	}
}

// Safe Case 1: Struct initialization using server-configured fixed identity
func safeStructAuthContext(r *http.Request, fixedOwner string) *AuthContext {
	// ok: http-identity-header-spoofing
	return &AuthContext{
		Authorized: true,
		Owner:      fixedOwner,
	}
}

// Case 2: Assigning untrusted identity header directly to session struct
func vulnSessionUserAssign(r *http.Request, s *Session) {
	// ruleid: http-identity-header-spoofing
	s.User = r.Header.Get("X-Remote-User")
}

// Safe Case 2: Checking trusted proxy before assigning identity header
func safeSessionUserAssign(r *http.Request, s *Session) {
	// ok: http-identity-header-spoofing
	if !isTrusted(r.RemoteAddr) {
		return
	}
	s.User = r.Header.Get("X-Remote-User")
}

// Case 3: Storing unverified identity header into request context
func vulnContextUser(r *http.Request) *http.Request {
	// ruleid: http-identity-header-spoofing
	ctx := context.WithValue(r.Context(), "user", r.Header.Get("X-Forwarded-User"))
	return r.WithContext(ctx)
}

// Safe Case 3: Verifying TLS connection before storing identity in context
func safeContextUser(r *http.Request) *http.Request {
	// ok: http-identity-header-spoofing
	if r.TLS == nil {
		return r
	}
	ctx := context.WithValue(r.Context(), "user", r.Header.Get("X-Forwarded-User"))
	return r.WithContext(ctx)
}

// Case 4: Returning caller-supplied identity header directly
func vulnReturnIdentity(r *http.Request) string {
	// ruleid: http-identity-header-spoofing
	return r.Header.Get("X-Auth-User")
}

// Safe Case 4: Reading standard non-identity header
func safeStandardHeader(r *http.Request, s *Session) {
	// ok: http-identity-header-spoofing
	s.UserAgent = r.Header.Get("User-Agent")
}

// Safe Case 5: Extracting identity inside trusted proxy block
func safeInsideProxyBlock(r *http.Request, s *Session) {
	// ok: http-identity-header-spoofing
	if isTrusted(r.RemoteAddr) {
		s.User = r.Header.Get("X-Forwarded-User")
	}
}
