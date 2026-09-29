package rules

import (
	"context"
	"net/http"
)

type oauth2Config struct{}

func (c *oauth2Config) Exchange(ctx context.Context, code string) (string, error) {
	return "token", nil
}

type oidcVerifier struct{}

func (v *oidcVerifier) Verify(ctx context.Context, rawIDToken string) (string, error) {
	return "id_token", nil
}

// 1. Vulnerable: Ignoring OAuth2 token exchange error
func vulnExchangeIgnoredError(ctx context.Context, conf *oauth2Config, code string) string {
	// ruleid: oauth-improper-authentication
	token, _ := conf.Exchange(ctx, code)
	return token
}

// 1. Safe: Properly handling OAuth2 token exchange error
func safeExchangeErrorCheck(ctx context.Context, conf *oauth2Config, code string) (string, error) {
	// ok: oauth-improper-authentication
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return "", err
	}
	return token, nil
}

// 2. Vulnerable: Ignoring OIDC token verification error
func vulnOIDCVerifyIgnoredError(ctx context.Context, verifier *oidcVerifier, rawIDToken string) string {
	// ruleid: oauth-improper-authentication
	idToken, _ := verifier.Verify(ctx, rawIDToken)
	return idToken
}

// 2. Safe: Properly checking OIDC verification error
func safeOIDCVerifyErrorCheck(ctx context.Context, verifier *oidcVerifier, rawIDToken string) (string, error) {
	// ok: oauth-improper-authentication
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", err
	}
	return idToken, nil
}

// 3. Vulnerable: OAuth state parameter mismatch does not return
func vulnOAuthStateNoReturn(w http.ResponseWriter, r *http.Request, conf *oauth2Config) {
	state := r.URL.Query().Get("state")
	expectedState := "expected_random_state"
	// ruleid: oauth-improper-authentication
	if state != expectedState {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
	}

	code := r.URL.Query().Get("code")
	token, err := conf.Exchange(r.Context(), code)
	if err != nil {
		return
	}
	_ = token
}

// 3. Safe: OAuth state parameter mismatch aborts and returns immediately
func safeOAuthStateWithReturn(w http.ResponseWriter, r *http.Request, conf *oauth2Config) {
	state := r.URL.Query().Get("state")
	expectedState := "expected_random_state"
	// ok: oauth-improper-authentication
	if state != expectedState {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := conf.Exchange(r.Context(), code)
	if err != nil {
		return
	}
	_ = token
}

// 4. Vulnerable: State token validation error with WriteHeader but no return
func vulnOAuthTokenWriteHeaderNoReturn(w http.ResponseWriter, r *http.Request) {
	oauthToken := r.Header.Get("X-OAuth-State")
	expectedToken := "secret_token"
	// ruleid: oauth-improper-authentication
	if oauthToken != expectedToken {
		w.WriteHeader(http.StatusUnauthorized)
	}

	w.Write([]byte("authenticated"))
}

// 4. Safe: State token validation error with WriteHeader and returns
func safeOAuthTokenWriteHeaderWithReturn(w http.ResponseWriter, r *http.Request) {
	oauthToken := r.Header.Get("X-OAuth-State")
	expectedToken := "secret_token"
	// ok: oauth-improper-authentication
	if oauthToken != expectedToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte("authenticated"))
}
