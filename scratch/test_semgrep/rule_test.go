package rules

import (
	"github.com/golang-jwt/jwt/v5"
)

func testVulnerable(tokenString string) {
	// ruleid: jwt-improper-authentication
	jwt.ParseUnverified(tokenString, jwt.MapClaims{})
}

func testSafe(tokenString string) {
	// ok: jwt-improper-authentication
	jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
}
