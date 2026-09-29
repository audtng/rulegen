package rules

import (
	"github.com/golang-jwt/jwt/v5"
)

func testInline() {
	// ruleid: test-rule
	_, _ = jwt.Parse("token", func(t *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
}

func testStandalone() {
	// ruleid: test-rule
	keyLookupFn := func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	}
	_ = keyLookupFn
}

func testSafe() {
	// ok: test-rule
	keyLookupFn := func(token *jwt.Token) (interface{}, error) {
		return []byte("secure-key"), nil
	}
	_ = keyLookupFn
}
