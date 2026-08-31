package main

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
	}

	atomic.AddUint64(&syncIdPrefix, 1)
	syncId := fmt.Sprintf("%05d-%s", syncIdPrefix, rand.RandString(5))

	logEntry := log.WithFields(log.Fields{"application": app.Name, "syncId": syncId})
	initialResourcesRes := make([]common.ResourceSyncResult, 0)
}

func (c *client) startGRPCProxy() (*grpc.Server, net.Listener, error) {
	serverAddr := fmt.Sprintf("%s/argocd-%s.sock", os.TempDir(), rand.RandString(16))
	ln, err := net.Listen("unix", serverAddr)

	if err != nil {
	FailOnErr(Run("", "mkdir", "-p", TmpDir))

	// random id - unique across test runs
	postFix := "-" + strings.ToLower(rand.RandString(5))
	id = t.Name() + postFix
	name = DnsFriendly(t.Name(), "")
	deploymentNamespace = DnsFriendly(fmt.Sprintf("argocd-e2e-%s", t.Name()), postFix)

	"github.com/argoproj/gitops-engine/pkg/health"
	. "github.com/argoproj/gitops-engine/pkg/sync/common"

	. "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/test/e2e/fixture"
}

func getNewNamespace(t *testing.T) string {
	postFix := "-" + strings.ToLower(rand.RandString(5))
	name := fixture.DnsFriendly(t.Name(), "")
	return fixture.DnsFriendly(fmt.Sprintf("argocd-e2e-%s", name), postFix)
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

import (
	"testing"
)

func TestRandString(t *testing.T) {
	ss := RandStringCharset(10, "A")
	if ss != "AAAAAAAAAA" {
		t.Errorf("Expected 10 As, but got %q", ss)
	}
	ss = RandStringCharset(5, "ABC123")
	if len(ss) != 5 {
		t.Errorf("Expected random string of length 10, but got %q", ss)
	}
}
