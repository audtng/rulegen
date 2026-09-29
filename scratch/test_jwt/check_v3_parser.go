package rules

import (
	"github.com/golang-jwt/jwt/v5"
)

func CheckV5Parser() {
    var claims jwt.MapClaims
    parser := jwt.NewParser()
    token, parts, err := parser.ParseUnverified("token", claims)
    _ = token
    _ = parts
    _ = err
}
