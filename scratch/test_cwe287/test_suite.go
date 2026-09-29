package rules

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type UserProfile struct {
	Email         string
	EmailVerified bool
}

func testParseUnverified(tokenStr string) {
	parser := jwt.NewParser()
	var claims jwt.MapClaims
	// ruleid: oauth-jwt-improper-authentication
	token, _, _ := parser.ParseUnverified(tokenStr, claims)
	_ = token
}

func testIgnoredParseError(tokenStr string) {
	// ruleid: oauth-jwt-improper-authentication
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
	// ruleid: oauth-jwt-improper-authentication
	token, _ := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("secret"), nil
	})
	_ = token
}

func testUnhandledParseError(tokenStr string) {
	// ruleid: oauth-jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte("secret"), nil
	})
	fmt.Println("Ignoring err:", err)
	_ = token
}

func testAlgorithmicConfusionAcceptEmptyKey(tokenStr string) {
	// ruleid: oauth-jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	if err != nil {
		return
	}
	_ = token
}

func testAlgorithmicConfusionNoneMethod(tokenStr string) {
	_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// ruleid: oauth-jwt-improper-authentication
		if t.Method == jwt.SigningMethodNone {
			return []byte(""), nil
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
}

func testOAuthEmailVerifiedInferredLiteral(email string) UserProfile {
	return UserProfile{
		Email: email,
		// ruleid: oauth-jwt-improper-authentication
		EmailVerified: email != "",
	}
}

func testOAuthEmailVerifiedInferredAssign(p *UserProfile, email string) {
	// ruleid: oauth-jwt-improper-authentication
	p.EmailVerified = email != ""
}

func testSafeParseWithAlgCheck(tokenStr string) {
	// ok: oauth-jwt-improper-authentication
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
	// ok: oauth-jwt-improper-authentication
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

func testSafeErrCheckEqualNil(tokenStr string) {
	// ok: oauth-jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err == nil {
		_ = token
	}
}

func testSafeEmailVerifiedExplicit(email string, providerVerified bool) UserProfile {
	// ok: oauth-jwt-improper-authentication
	return UserProfile{
		Email:         email,
		EmailVerified: providerVerified,
	}
}

func testSafeEmailVerifiedAssign(p *UserProfile, providerVerified bool) {
	// ok: oauth-jwt-improper-authentication
	p.EmailVerified = providerVerified
}
