package rules

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func testUnsafeAllowNone(tokenStr string) {
	// ruleid: jwt-none-algorithm-confusion
	_, _ = jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwt.UnsafeAllowNoneSignatureType, nil
	})
}

func testUnsafeAllowNoneWithClaims(tokenStr string) {
	var claims jwt.MapClaims
	// ruleid: jwt-none-algorithm-confusion
	_, _ = jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwt.UnsafeAllowNoneSignatureType, nil
	})
}

func testExplicitNoneMethodCheck(tokenStr string) {
	_, _ = jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// ruleid: jwt-none-algorithm-confusion
		if t.Method == jwt.SigningMethodNone {
			return []byte(""), nil
		}
		return []byte("secret"), nil
	})
}

func testAcceptEmptyKeyNoneAlg(tokenStr string) {
	// ruleid: jwt-none-algorithm-confusion
	_, _ = jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
}

func testAcceptEmptyKeyWithClaims(tokenStr string) {
	var claims jwt.MapClaims
	// ruleid: jwt-none-algorithm-confusion
	_, _ = jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
}

func testValidMethodsIncludesNone() {
	// ruleid: jwt-none-algorithm-confusion
	_ = jwt.WithValidMethods([]string{"HS256", "none"})
}

func testSafeParse(tokenStr string) {
	// ok: jwt-none-algorithm-confusion
	_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
}

func testSafeParseWithClaims(tokenStr string) {
	var claims jwt.MapClaims
	// ok: jwt-none-algorithm-confusion
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("rsa-public-key"), nil
	})
	if err != nil {
		return
	}
}

func testSafeValidMethods() {
	// ok: jwt-none-algorithm-confusion
	_ = jwt.WithValidMethods([]string{"HS256", "RS256"})
}
