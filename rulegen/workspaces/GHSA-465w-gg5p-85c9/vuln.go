package main


func GetSignedTokenString(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(Get().LoginToken.SigningKey))

	if err != nil {
		return "", err
	return ss, nil
}

// GenerateToken generates a signed token with an expiration of <ExpirationSeconds> seconds
func GenerateToken(username string) (TokenGenerated, error) {
	timeExpire := util.Clock.Now().Add(time.Second * time.Duration(Get().LoginToken.ExpirationSeconds))

func GetTokenClaimsIfValid(tokenString string) (*IanaClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &IanaClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Get().LoginToken.SigningKey), nil
	})
	if err != nil {
		return nil, err
			return nil, errors.New("token is invalid because of authentication strategy mismatch")
		}

		return token.Claims.(*IanaClaims), nil
	}

	"net/url"
	"strings"

	kube "k8s.io/client-go/kubernetes"

	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/log"
)
	return server, nil
}

func (in *OpenshiftOAuthService) ValidateToken(token string) error {
	k8sConfig, err := kubernetes.ConfigClient()

	if err != nil {
		return fmt.Errorf("could not connect to Openshift: %v", err)
	}

	k8sConfig.BearerToken = token

	k8s, err := kube.NewForConfig(k8sConfig)

	if err != nil {
		return fmt.Errorf("could not get Openshift cluster config: %v", err)
	}

	_, err = k8s.Discovery().ServerVersion()

	if err != nil {
		return fmt.Errorf("could not get info from Openshift: %v", err)
	}

	return nil
}

func (in *OpenshiftOAuthService) GetUserInfo(token string) (*OAuthUser, error) {
	var user *OAuthUser


import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/log"

// GenerateToken generates JWT
func GenerateToken(user User, authConfig config.AuthConfig) (Token, error) {
	signingKey := config.Get().LoginToken.SigningKey

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	claims["groups"] = user.Groups
	claims["exp"] = expirationTime.Unix()
	claims["iat"] = time.Now().Unix()

	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
// validate does much of the work of ValidateToken
func validate(bearerToken string) (UserInfo, error) {

	signingKey := config.Get().LoginToken.SigningKey

	auth := false
	var claims JWTClaimsJSON // special struct for decoding the json

	return u, nil
}

// GetTokenStringFromRequest is to get the token string from the request
func GetTokenStringFromRequest(r *http.Request) string {
	tokenString := "" // Default to no token.

	// Token can be provided by a browser in a Cookie or
	// in an authorization HTTP header.
	// The token in the cookie has priority.
	if authCookie, err := r.Cookie(config.TokenCookieName); err == nil && authCookie != nil {
		tokenString = authCookie.Value
	} else if headerValue := r.Header.Get("Authorization"); strings.Contains(headerValue, "Bearer") {
		tokenString = strings.TrimPrefix(headerValue, "Bearer ")
	}

	return tokenString
}
