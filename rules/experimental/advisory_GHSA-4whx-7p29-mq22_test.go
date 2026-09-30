package rules

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func testParseUnverified(tokenStr string) {
	parser := jwt.NewParser()
	var claims jwt.MapClaims
	// ruleid: jwt-authentication-bypass
	token, _, _ := parser.ParseUnverified(tokenStr, claims)
	_ = token
}

func testInsecureDirectUnverified(tokenStr string) {
	var claims jwt.MapClaims
	// ruleid: jwt-authentication-bypass
	token, _, _ := jwt.NewParser().ParseUnverified(tokenStr, claims)
	_ = token
}

func testIgnoredParseError(tokenStr string) {
	// ruleid: jwt-authentication-bypass
	token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	_ = token
}

func testIgnoredParseWithClaimsError(tokenStr string) {
	var claims jwt.MapClaims
	// ruleid: jwt-authentication-bypass
	token, _ := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	_ = token
}

func testUnhandledParseError(tokenStr string) {
	// ruleid: jwt-authentication-bypass
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	fmt.Println("Proceeding without checking err:", err)
	_ = token
}

func testAlgorithmicConfusionAcceptEmptyKey(tokenStr string) {
	// ruleid: jwt-authentication-bypass
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
		// ruleid: jwt-authentication-bypass
		if t.Method == jwt.SigningMethodNone {
			return []byte(""), nil
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
}

func testSafeParseWithAlgCheck(tokenStr string) {
	// ok: jwt-authentication-bypass
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
	// ok: jwt-authentication-bypass
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

func testSafeCheckNil(tokenStr string) {
	// ok: jwt-authentication-bypass
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

func testSafeIfInit(tokenStr string) {
	// ok: jwt-authentication-bypass
	if token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret"), nil
	}); err != nil {
		return
	} else {
		_ = token
	}
}
