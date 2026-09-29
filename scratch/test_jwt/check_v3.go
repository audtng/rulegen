package rules

import (
	"github.com/golang-jwt/jwt"
)

func CheckV3() {
    _ = jwt.ParseUnverified
}
