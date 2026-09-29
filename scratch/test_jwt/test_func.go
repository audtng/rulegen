package rules

import (
	"github.com/golang-jwt/jwt/v5"
)

func testStandaloneKeyFunc() {
	keyLookupFn := func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	}
	_ = keyLookupFn
}
