package rules

import (
	"fmt"
	"net/http"
	"net/url"
)

// Vulnerable: Discards error when parsing X-Original-URL inline
func testVulnDiscardErrorInline(w http.ResponseWriter, r *http.Request) {
	// ruleid: http-forward-auth-bypass
	target, _ := url.ParseRequestURI(r.Header.Get("X-Original-URL"))
	fmt.Println(target)
}

// Vulnerable: Discards error when parsing X-Forwarded-Uri from variable
func testVulnDiscardErrorVar(w http.ResponseWriter, r *http.Request) {
	orig := r.Header.Get("X-Forwarded-Uri")
	// ruleid: http-forward-auth-bypass
	target, _ := url.Parse(orig)
	fmt.Println(target)
}

// Vulnerable: Discards error when parsing target from query parameter
func testVulnDiscardErrorQuery(w http.ResponseWriter, r *http.Request) {
	// ruleid: http-forward-auth-bypass
	target, _ := url.Parse(r.URL.Query().Get("rd"))
	fmt.Println(target)
}

// Vulnerable: Error handled but returns without setting 401/403 (defaults to 200 OK)
func testVulnNo401ResponseInline(w http.ResponseWriter, r *http.Request) {
	// ruleid: http-forward-auth-bypass
	target, err := url.ParseRequestURI(r.Header.Get("X-Original-URL"))
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	fmt.Println(target)
}

// Vulnerable: Sets 500 Internal Server Error instead of 401/403
func testVuln500ResponseVar(w http.ResponseWriter, r *http.Request) {
	uri := r.Header.Get("X-Original-URI")
	// ruleid: http-forward-auth-bypass
	target, err := url.Parse(uri)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(target)
}

// Vulnerable: Query param parsed, returns on error without 401/403
func testVulnNo401QueryParam(w http.ResponseWriter, r *http.Request) {
	// ruleid: http-forward-auth-bypass
	target, err := url.Parse(r.URL.Query().Get("target"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	fmt.Println(target)
}

// Safe: Responds with http.StatusUnauthorized
func testSafeWriteHeaderUnauthorized(w http.ResponseWriter, r *http.Request) {
	// ok: http-forward-auth-bypass
	target, err := url.ParseRequestURI(r.Header.Get("X-Original-URL"))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(target)
}

// Safe: Uses http.Error with http.StatusUnauthorized
func testSafeHttpErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	uri := r.Header.Get("X-Forwarded-Uri")
	// ok: http-forward-auth-bypass
	target, err := url.Parse(uri)
	if err != nil {
		http.Error(w, "Unauthorized target URL", http.StatusUnauthorized)
		return
	}
	fmt.Println(target)
}

// Safe: Responds with http.StatusForbidden
func testSafeWriteHeaderForbidden(w http.ResponseWriter, r *http.Request) {
	uri := r.Header.Get("X-Original-URI")
	// ok: http-forward-auth-bypass
	target, err := url.ParseRequestURI(uri)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	fmt.Println(target)
}

// Safe: Responds with integer literal 401
func testSafeLiteral401(w http.ResponseWriter, r *http.Request) {
	// ok: http-forward-auth-bypass
	target, err := url.Parse(r.URL.Query().Get("rd"))
	if err != nil {
		w.WriteHeader(401)
		return
	}
	fmt.Println(target)
}

// Safe: Helper function that returns the error up the call stack
func testSafeHelperReturnError(r *http.Request) (*url.URL, error) {
	// ok: http-forward-auth-bypass
	target, err := url.ParseRequestURI(r.Header.Get("X-Original-URL"))
	if err != nil {
		return nil, err
	}
	return target, nil
}

// Safe: Unrelated standard URL parsing
func testUnrelatedURLParse(raw string) {
	// ok: http-forward-auth-bypass
	u, err := url.Parse(raw)
	if err != nil {
		return
	}
	fmt.Println(u)
}
