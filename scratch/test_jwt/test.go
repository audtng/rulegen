package rules

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func vulnEmptyKey(tokenStr string) {
	jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
}

func vulnNilKey(tokenStr string) {
	jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return nil, nil
	})
}

func vulnEmptyStringKey(tokenStr string) {
	jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return "", nil
	})
}

func vulnEmptyByteSlice(tokenStr string) {
	jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte{}, nil
	})
}

func vulnSignedStringEmpty(t *jwt.Token) {
	_, _ = t.SignedString([]byte(""))
}

func okKey(tokenStr string, secret string) {
	jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if secret == "" {
			return nil, fmt.Errorf("secret not configured")
		}
		return []byte(secret), nil
	})
}
