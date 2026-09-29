package rules

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Edge case 1: Insecure ParseUnverified
func testVulnParseUnverified(tokenString string) {
	parser := jwt.NewParser()
	// ruleid: jwt-unverified-token-authentication
	token, _, _ := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	_ = token
}

// Edge case 2: Ignored error with blank identifier
func testVulnIgnoredError(tokenString string) {
	// ruleid: jwt-unverified-token-authentication
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	_ = token
}

// Edge case 3: Unchecked error variable
func testVulnUncheckedError(tokenString string) {
	// ruleid: jwt-unverified-token-authentication
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	_ = err
	_ = token
}

// Edge case 4: Missing return after http.Error on parse failure
func testVulnMissingReturnHttpError(w http.ResponseWriter, tokenString string) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	// ruleid: jwt-unverified-token-authentication
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
	}
	_ = token
}

// Edge case 5: Missing return after WriteHeader on !token.Valid
func testVulnMissingReturnValidCheck(w http.ResponseWriter, tokenString string) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
	// ruleid: jwt-unverified-token-authentication
	if !token.Valid {
		w.WriteHeader(http.StatusForbidden)
	}
	_ = token
}

// Safe case 1: Standard error check and validation
func testSafeParse(tokenString string) {
	// ok: jwt-unverified-token-authentication
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
	// ok: jwt-unverified-token-authentication
	if !token.Valid {
		return
	}
	_ = token
}

// Safe case 2: HTTP handler with proper return on error
func testSafeHttpHandler(w http.ResponseWriter, tokenString string) {
	// ok: jwt-unverified-token-authentication
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret"), nil
	})
	// ok: jwt-unverified-token-authentication
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	// ok: jwt-unverified-token-authentication
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_ = token
}

// Safe case 3: ParseWithClaims with proper error handling
func testSafeParseWithClaims(tokenString string) {
	var claims jwt.MapClaims
	// ok: jwt-unverified-token-authentication
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return
	}
	// ok: jwt-unverified-token-authentication
	if !token.Valid {
		return
	}
	_ = token
}
