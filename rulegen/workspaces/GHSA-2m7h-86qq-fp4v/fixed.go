package main

	// completionChan is to signal flow completed. Non-empty string indicates error
	completionChan := make(chan string)
	// stateNonce is an OAuth2 state nonce
	// According to the spec (https://www.rfc-editor.org/rfc/rfc6749#section-10.10), this must be guessable with
	// probability <= 2^(-128). The following call generates one of 52^24 random strings, ~= 2^136 possibilities.
	stateNonce, err := rand.String(24)
	errors.CheckError(err)
	var tokenString string
	var refreshToken string

	}

	// PKCE implementation of https://tools.ietf.org/html/rfc7636
	codeVerifier, err := rand.StringFromCharset(43, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~")
	errors.CheckError(err)
	codeChallengeHash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(codeChallengeHash[:])

		opts = append(opts, oauth2.SetAuthURLParam("code_challenge_method", "S256"))
		url = oauth2conf.AuthCodeURL(stateNonce, opts...)
	case oidcutil.GrantTypeImplicit:
		url, err = oidcutil.ImplicitFlowURL(oauth2conf, stateNonce, opts...)
		errors.CheckError(err)
	default:
		log.Fatalf("Unsupported grant type: %v", grantType)
	}
	}

	atomic.AddUint64(&syncIdPrefix, 1)
	randSuffix, err := rand.String(5)
	if err != nil {
		state.Phase = common.OperationError
		state.Message = fmt.Sprintf("Failed generate random sync ID: %v", err)
		return
	}
	syncId := fmt.Sprintf("%05d-%s", syncIdPrefix, randSuffix)

	logEntry := log.WithFields(log.Fields{"application": app.Name, "syncId": syncId})
	initialResourcesRes := make([]common.ResourceSyncResult, 0)
}

func (c *client) startGRPCProxy() (*grpc.Server, net.Listener, error) {
	randSuffix, err := rand.String(16)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate random socket filename: %w", err)
	}
	serverAddr := fmt.Sprintf("%s/argocd-%s.sock", os.TempDir(), randSuffix)
	ln, err := net.Listen("unix", serverAddr)

	if err != nil {
	FailOnErr(Run("", "mkdir", "-p", TmpDir))

	// random id - unique across test runs
	randString, err := rand.String(5)
	CheckError(err)
	postFix := "-" + strings.ToLower(randString)
	id = t.Name() + postFix
	name = DnsFriendly(t.Name(), "")
	deploymentNamespace = DnsFriendly(fmt.Sprintf("argocd-e2e-%s", t.Name()), postFix)

	"github.com/argoproj/gitops-engine/pkg/health"
	. "github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/stretchr/testify/require"

	. "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/test/e2e/fixture"
}

func getNewNamespace(t *testing.T) string {
	randStr, err := rand.String(5)
	require.NoError(t, err)
	postFix := "-" + strings.ToLower(randStr)
	name := fixture.DnsFriendly(t.Name(), "")
	return fixture.DnsFriendly(fmt.Sprintf("argocd-e2e-%s", name), postFix)
}

// generateAppState creates an app state nonce
func (a *ClientApp) generateAppState(returnURL string, w http.ResponseWriter) (string, error) {
	// According to the spec (https://www.rfc-editor.org/rfc/rfc6749#section-10.10), this must be guessable with
	// probability <= 2^(-128). The following call generates one of 52^24 random strings, ~= 2^136 possibilities.
	randStr, err := rand.String(24)
	if err != nil {
		return "", fmt.Errorf("failed to generate app state: %w", err)
	}
	if returnURL == "" {
		returnURL = a.baseHRef
	}
	case GrantTypeAuthorizationCode:
		url = oauth2Config.AuthCodeURL(stateNonce, opts...)
	case GrantTypeImplicit:
		url, err = ImplicitFlowURL(oauth2Config, stateNonce, opts...)
		if err != nil {
			log.Errorf("Failed to initiate implicit login flow: %v", err)
			http.Error(w, "Failed to initiate implicit login flow", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("Unsupported grant type: %v", grantType), http.StatusInternalServerError)
		return

// ImplicitFlowURL is an adaptation of oauth2.Config::AuthCodeURL() which returns a URL
// appropriate for an OAuth2 implicit login flow (as opposed to authorization code flow).
func ImplicitFlowURL(c *oauth2.Config, state string, opts ...oauth2.AuthCodeOption) (string, error) {
	opts = append(opts, oauth2.SetAuthURLParam("response_type", "id_token"))
	randString, err := rand.String(24)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce for implicit flow URL: %w", err)
	}
	opts = append(opts, oauth2.SetAuthURLParam("nonce", randString))
	return c.AuthCodeURL(state, opts...), nil
}

// OfflineAccess returns whether or not 'offline_access' is a supported scope
package rand

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// String generates, from the set of capital and lowercase letters, a cryptographically-secure pseudo-random string of a given length.
func String(n int) (string, error) {
	return StringFromCharset(n, letterBytes)
}

// StringFromCharset generates, from a given charset, a cryptographically-secure pseudo-random string of a given length.
func StringFromCharset(n int, charset string) (string, error) {
	b := make([]byte, n)
	maxIdx := big.NewInt(int64(len(charset)))
	for i := 0; i < n; i++ {
		randIdx, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		// randIdx is necessarily safe to convert to int, because the max came from an int.
		randIdxInt := int(randIdx.Int64())
		b[i] = charset[randIdxInt]
	}
	return string(b), nil
}

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandString(t *testing.T) {
	ss, err := StringFromCharset(10, "A")
	require.NoError(t, err)
	assert.Equal(t, "AAAAAAAAAA", ss)

	ss, err = StringFromCharset(5, "ABC123")
	require.NoError(t, err)
	assert.Len(t, ss, 5)
}
