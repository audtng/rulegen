package rules

import (
	"github.com/golang-jwt/jwt/v5"
)

func testInline() {
	_, _ = jwt.Parse("token", func(t *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
}

func testStandalone() {
	keyLookupFn := func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	}
	_ = keyLookupFn
}
