package main


func GetSignedTokenString(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(GetSigningKey()))

	if err != nil {
		return "", err
	return ss, nil
}

func GetSigningKey() string {
	cfg := Get()
	signKey := cfg.LoginToken.SigningKey

	if len(signKey) == 0 || signKey == "kiali" {
		// "kiali" is a well-known signing key reported in a CVE. We ban it's usage.
		// An empty key is also just not allowed.
		panic("signing key for login tokens is invalid")
	}

	if cfg.Auth.Strategy == AuthStrategyLogin {
		// If we are using "login" strategy, let's combine the login passphrase
		// and the token signing key to form a new signing key. This way, if
		// either the login passphrase or the signing key is changed, active
		// sessions will be invalidated.
		signKey = fmt.Sprintf("%s+%s", signKey, cfg.Server.Credentials.Passphrase)
	}

	return signKey
}

// GenerateToken generates a signed token with an expiration of <ExpirationSeconds> seconds
func GenerateToken(username string) (TokenGenerated, error) {
	timeExpire := util.Clock.Now().Add(time.Second * time.Duration(Get().LoginToken.ExpirationSeconds))

func GetTokenClaimsIfValid(tokenString string) (*IanaClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &IanaClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetSigningKey()), nil
	})
	if err != nil {
		return nil, err
			return nil, errors.New("token is invalid because of authentication strategy mismatch")
		}

		// A token with no expiration claim is invalid for Kiali
		if claims.ExpiresAt == 0 {
			return nil, errors.New("token is invalid because expiration claim is missing")
		}

		// If auth strategy is login and the subject claim does not match the username in the Kiali secret,
		// the token is invalid.
		if cfg.Auth.Strategy == AuthStrategyLogin && claims.Subject != cfg.Server.Credentials.Username {
			return nil, errors.New("username has changed")
		}

		return token.Claims.(*IanaClaims), nil
	}

	"net/url"
	"strings"

	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/log"
)
	return server, nil
}

func (in *OpenshiftOAuthService) GetUserInfo(token string) (*OAuthUser, error) {
	var user *OAuthUser


import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/log"

// GenerateToken generates JWT
func GenerateToken(user User, authConfig config.AuthConfig) (Token, error) {
	signingKey := config.GetSigningKey()

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	claims["groups"] = user.Groups
	claims["exp"] = expirationTime.Unix()
	claims["iat"] = time.Now().Unix()
	claims["iss"] = config.AuthStrategyLDAPIssuer

	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
// validate does much of the work of ValidateToken
func validate(bearerToken string) (UserInfo, error) {

	signingKey := config.GetSigningKey()

	auth := false
	var claims JWTClaimsJSON // special struct for decoding the json

	return u, nil
}
