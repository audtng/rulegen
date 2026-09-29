package rules

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func vulnParserParseUnverified(tokenStr string) (*jwt.Token, error) {
	parser := &jwt.Parser{}
	claims := jwt.MapClaims{}
	// ruleid: jwt-algorithm-confusion-improper-auth
	token, _, err := parser.ParseUnverified(tokenStr, claims)
	return token, err
}

func vulnIgnoredError(tokenStr string) *jwt.Token {
	// ruleid: jwt-algorithm-confusion-improper-auth
	token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret-key"), nil
	})
	return token
}

func vulnIgnoredErrorClaims(tokenStr string) *jwt.Token {
	claims := jwt.MapClaims{}
	// ruleid: jwt-algorithm-confusion-improper-auth
	token, _ := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret-key"), nil
	})
	return token
}

func vulnUnhandledError(tokenStr string) *jwt.Token {
	// ruleid: jwt-algorithm-confusion-improper-auth
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret-key"), nil
	})
	fmt.Println("Proceeding without checking err:", err)
	return token
}

func vulnUnhandledErrorClaims(tokenStr string) *jwt.Token {
	claims := jwt.MapClaims{}
	// ruleid: jwt-algorithm-confusion-improper-auth
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret-key"), nil
	})
	fmt.Println("Proceeding without checking err:", err)
	return token
}

func vulnAlgorithmConfusionParse(tokenStr string) (*jwt.Token, error) {
	// ruleid: jwt-algorithm-confusion-improper-auth
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("public-key-or-secret"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func vulnAlgorithmConfusionClaims(tokenStr string) (*jwt.Token, error) {
	claims := jwt.MapClaims{}
	// ruleid: jwt-algorithm-confusion-improper-auth
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("public-key-or-secret"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func vulnAllowNoneSignature(tokenStr string) (*jwt.Token, error) {
	// ruleid: jwt-algorithm-confusion-improper-auth
	key := jwt.UnsafeAllowNoneSignatureType
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
}

func okParseValidated(tokenStr string) (*jwt.Token, error) {
	// ok: jwt-algorithm-confusion-improper-auth
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret-key"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func okParseWithClaimsValidated(tokenStr string) (*jwt.Token, error) {
	claims := jwt.MapClaims{}
	// ok: jwt-algorithm-confusion-improper-auth
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret-key"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func okParseCheckedNil(tokenStr string) (*jwt.Token, error) {
	// ok: jwt-algorithm-confusion-improper-auth
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret-key"), nil
	})
	if err == nil {
		return token, nil
	}
	return nil, err
}

func okParseInlineIf(tokenStr string) (*jwt.Token, error) {
	// ok: jwt-algorithm-confusion-improper-auth
	if token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret-key"), nil
	}); err != nil {
		return nil, err
	} else {
		return token, nil
	}
}
