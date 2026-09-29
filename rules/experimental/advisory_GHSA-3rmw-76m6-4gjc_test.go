package rules

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func vulnParserParseUnverified(tokenStr string) (*jwt.Token, error) {
	// ruleid: go-jwt-improper-authentication
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	return token, err
}

func vulnIgnoredError(tokenStr string) *jwt.Token {
	// ruleid: go-jwt-improper-authentication
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	return token
}

func vulnIgnoredErrorClaims(tokenStr string) *jwt.Token {
	// ruleid: go-jwt-improper-authentication
	token, _ := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	return token
}

func vulnUnhandledError(tokenStr string) *jwt.Token {
	// ruleid: go-jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	fmt.Println("Proceeding without checking error:", err)
	return token
}

func vulnUnhandledErrorClaims(tokenStr string) *jwt.Token {
	// ruleid: go-jwt-improper-authentication
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	fmt.Println("Proceeding without checking error:", err)
	return token
}

func vulnInsecureKeyfunc(tokenStr string) (*jwt.Token, error) {
	// ruleid: go-jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func okParseChecked(tokenStr string) (*jwt.Token, error) {
	// ok: go-jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func okParseWithClaimsChecked(tokenStr string) (*jwt.Token, error) {
	// ok: go-jwt-improper-authentication
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func okCheckedNil(tokenStr string) (*jwt.Token, error) {
	// ok: go-jwt-improper-authentication
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err == nil {
		return token, nil
	}
	return nil, err
}
