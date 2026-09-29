package rules

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func testVulnSignEmptyByteSlice(t *jwt.Token) (string, error) {
	// ruleid: jwt-empty-secret-authentication-bypass
	return t.SignedString([]byte(""))
}

func testVulnSignEmptyString(t *jwt.Token) (string, error) {
	// ruleid: jwt-empty-secret-authentication-bypass
	return t.SignedString("")
}

func testVulnSignEmptyByteLiteral(t *jwt.Token) (string, error) {
	// ruleid: jwt-empty-secret-authentication-bypass
	return t.SignedString([]byte{})
}

func testVulnSignNilSecret(t *jwt.Token) (string, error) {
	// ruleid: jwt-empty-secret-authentication-bypass
	return t.SignedString(nil)
}

func testVulnInlineKeyfuncEmptyBytes(tokenStr string) (*jwt.Token, error) {
	// ruleid: jwt-empty-secret-authentication-bypass
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte(""), nil
	})
}

func testVulnStandaloneKeyfuncEmptyBytes(tokenStr string) (*jwt.Token, error) {
	// ruleid: jwt-empty-secret-authentication-bypass
	keyLookupFn := func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte{}, nil
	}
	return jwt.Parse(tokenStr, keyLookupFn)
}

// ruleid: jwt-empty-secret-authentication-bypass
func emptyKeyLookup(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected method")
	}
	return "", nil
}

// ruleid: jwt-empty-secret-authentication-bypass
func nilKeyLookup(t *jwt.Token) (interface{}, error) {
	return nil, nil
}

func testSafeSignSecureSecret(t *jwt.Token) (string, error) {
	// ok: jwt-empty-secret-authentication-bypass
	return t.SignedString([]byte("super-secure-production-secret-key-12345"))
}

func testSafeInlineParse(tokenStr string, secret []byte) (*jwt.Token, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("secret cannot be empty")
	}
	// ok: jwt-empty-secret-authentication-bypass
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return secret, nil
	})
}

func testSafeStandaloneKeyfunc(tokenStr string) (*jwt.Token, error) {
	// ok: jwt-empty-secret-authentication-bypass
	keyLookupFn := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("valid-shared-hmac-secret"), nil
	}
	return jwt.Parse(tokenStr, keyLookupFn)
}

// ok: jwt-empty-secret-authentication-bypass
func safeNamedKeyLookup(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected method")
	}
	return []byte("valid-shared-hmac-secret"), nil
}

func testSafeParseWithClaims(tokenStr string) (*jwt.Token, error) {
	var claims jwt.MapClaims
	// ok: jwt-empty-secret-authentication-bypass
	return jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("valid-shared-hmac-secret"), nil
	})
}
