package rules

import (
	"fmt"
	"net/http"
)

func testVulnBasicAuthIgnoredOk(r *http.Request) {
	// ruleid: http-basic-auth-unverified
	user, pass, _ := r.BasicAuth()
	fmt.Println(user, pass)
}

func testVulnBasicAuthOnlyPasswordIgnoredOk(r *http.Request) {
	// ruleid: http-basic-auth-unverified
	_, pass, _ := r.BasicAuth()
	fmt.Println(pass)
}

func testVulnBasicAuthOnlyUserIgnoredOk(r *http.Request) {
	// ruleid: http-basic-auth-unverified
	user, _, _ := r.BasicAuth()
	fmt.Println(user)
}

func testVulnBasicAuthReassignIgnored(r *http.Request) {
	var user, pass string
	// ruleid: http-basic-auth-unverified
	user, pass, _ = r.BasicAuth()
	fmt.Println(user, pass)
}

func testVulnBasicAuthMissingReturn(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	// ruleid: http-basic-auth-unverified
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	fmt.Fprintf(w, "Welcome %s %s", user, pass)
}

func testVulnBasicAuthMissingReturnWriteHeader(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	// ruleid: http-basic-auth-unverified
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
	}
	fmt.Fprintf(w, "Welcome %s %s", user, pass)
}

func testSafeBasicAuth(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	// ok: http-basic-auth-unverified
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Welcome %s %s", user, pass)
}

func testSafeBasicAuthWriteHeader(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	// ok: http-basic-auth-unverified
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Welcome %s %s", user, pass)
}

func testSafeBasicAuthReturnValue(w http.ResponseWriter, r *http.Request) error {
	user, pass, ok := r.BasicAuth()
	// ok: http-basic-auth-unverified
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return fmt.Errorf("unauthorized: %s %s", user, pass)
	}
	return nil
}

func testSafeBasicAuthWriteHeaderReturn(w http.ResponseWriter, r *http.Request) error {
	_, _, ok := r.BasicAuth()
	// ok: http-basic-auth-unverified
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return fmt.Errorf("unauthorized")
	}
	return nil
}
