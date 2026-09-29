package rules

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func testParseUnverified(tokenStr string) {
	parser := jwt.NewParser()
	// ruleid: jwt-improper-authentication
	parser.ParseUnverified(tokenStr, jwt.MapClaims{})
}

func testIgnoredError(tokenStr string) {
	// ruleid: jwt-improper-authentication
	token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("secret"), nil
	})
	_ = token
}

func testIgnoredErrorWithClaims(tokenStr string) {
	var claims jwt.MapClaims
	// ruleid: jwt-improper-authentication
	token, _ := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("secret"), nil
	})
	_ = token
}

func testAlgConfusion(tokenStr string) {
	// ruleid: jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
	_ = token
}

func testSafeParse(tokenStr string) {
	// ok: jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
	_ = token
}

func testSafeParseWithClaims(tokenStr string) {
	var claims jwt.MapClaims
	// ok: jwt-improper-authentication
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
	_ = token
}
