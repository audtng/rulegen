package rules

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func CheckV5() {
	var claims jwt.MapClaims
	parser := jwt.NewParser()
	_, _, _ = parser.ParseUnverified("token", claims)
	_, _ = jwt.Parse("token", func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("bad alg")
		}
		return []byte("secret"), nil
	})
	_ = jwt.UnsafeAllowNoneSignatureType
}
