package rules

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func testWithoutClaimsValidation(tokenStr string) {
	// ruleid: jwt-policy-validation-bypass
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	var claims jwt.MapClaims
	token, err := parser.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
	_ = token
}

func testParseUnverifiedMethod(tokenStr string) {
	parser := jwt.NewParser()
	var claims jwt.MapClaims
	// ruleid: jwt-policy-validation-bypass
	token, parts, err := parser.ParseUnverified(tokenStr, claims)
	if err != nil {
		return
	}
	_ = token
	_ = parts
}

func testIgnoredParseError(tokenStr string) {
	// ruleid: jwt-policy-validation-bypass
	token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("secret"), nil
	})
	_ = token
}

func testIgnoredParseWithClaimsError(tokenStr string) {
	var claims jwt.MapClaims
	// ruleid: jwt-policy-validation-bypass
	token, _ := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("secret"), nil
	})
	_ = token
}

func testIgnoredParserInstanceError(tokenStr string) {
	parser := jwt.NewParser()
	// ruleid: jwt-policy-validation-bypass
	token, _ := parser.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("secret"), nil
	})
	_ = token
}

func testSafeParse(tokenStr string) {
	// ok: jwt-policy-validation-bypass
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %v", t.Header["alg"])
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
	// ok: jwt-policy-validation-bypass
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %v", t.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
	_ = token
}

func testSafeParserWithOptions(tokenStr string) {
	// ok: jwt-policy-validation-bypass
	parser := jwt.NewParser(jwt.WithExpirationRequired())
	var claims jwt.MapClaims
	token, err := parser.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
	_ = token
}
