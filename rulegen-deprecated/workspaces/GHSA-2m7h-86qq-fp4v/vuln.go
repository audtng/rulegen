package main

package rand

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandString generates, from a given charset, a cryptographically-secure pseudo-random string of a given length.
func RandString(n int) string {
	return RandStringCharset(n, letterBytes)
}

func RandStringCharset(n int, charset string) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(charset) {
			b[i] = charset[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// generateAppState creates an app state nonce
func (a *ClientApp) generateAppState(returnURL string, w http.ResponseWriter) (string, error) {
	randStr := rand.RandString(10)
	if returnURL == "" {
		returnURL = a.baseHRef
	}
	case GrantTypeAuthorizationCode:
		url = oauth2Config.AuthCodeURL(stateNonce, opts...)
	case GrantTypeImplicit:
		url = ImplicitFlowURL(oauth2Config, stateNonce, opts...)
	default:
		http.Error(w, fmt.Sprintf("Unsupported grant type: %v", grantType), http.StatusInternalServerError)
		return

// ImplicitFlowURL is an adaptation of oauth2.Config::AuthCodeURL() which returns a URL
// appropriate for an OAuth2 implicit login flow (as opposed to authorization code flow).
func ImplicitFlowURL(c *oauth2.Config, state string, opts ...oauth2.AuthCodeOption) string {
	opts = append(opts, oauth2.SetAuthURLParam("response_type", "id_token"))
	opts = append(opts, oauth2.SetAuthURLParam("nonce", rand.RandString(10)))
	return c.AuthCodeURL(state, opts...)
}

// OfflineAccess returns whether or not 'offline_access' is a supported scope
	// completionChan is to signal flow completed. Non-empty string indicates error
	completionChan := make(chan string)
	// stateNonce is an OAuth2 state nonce
	stateNonce := rand.RandString(10)
	var tokenString string
	var refreshToken string

	}

	// PKCE implementation of https://tools.ietf.org/html/rfc7636
	codeVerifier := rand.RandStringCharset(43, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~")
	codeChallengeHash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(codeChallengeHash[:])

		opts = append(opts, oauth2.SetAuthURLParam("code_challenge_method", "S256"))
		url = oauth2conf.AuthCodeURL(stateNonce, opts...)
	case oidcutil.GrantTypeImplicit:
		url = oidcutil.ImplicitFlowURL(oauth2conf, stateNonce, opts...)
	default:
		log.Fatalf("Unsupported grant type: %v", grantType)
	}
