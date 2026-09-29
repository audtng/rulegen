package rules

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func testParseUnverified(tokenStr string) {
	parser := jwt.NewParser()
	var claims jwt.MapClaims
	// ruleid: jwt-improper-authentication
	token, parts, err := parser.ParseUnverified(tokenStr, claims)
	_ = token
	_ = parts
	_ = err
}

func testIgnoredError(tokenStr string) {
	// ruleid: jwt-improper-authentication
	token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
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
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	_ = token
}

func testUnhandledError(tokenStr string) {
	// ruleid: jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	fmt.Println("Proceeding without checking err:", err)
	_ = token
}

func testUnhandledErrorWithClaims(tokenStr string) {
	var claims jwt.MapClaims
	// ruleid: jwt-improper-authentication
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	fmt.Println("Proceeding without checking err:", err)
	_ = token
}

func testAlgConfusionParse(tokenStr string) {
	// ruleid: jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
	_ = token
}

func testAlgConfusionParseWithClaims(tokenStr string) {
	var claims jwt.MapClaims
	// ruleid: jwt-improper-authentication
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
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

func testSafeCheckNil(tokenStr string) {
	// ok: jwt-improper-authentication
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
	// ok: jwt-improper-authentication
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
