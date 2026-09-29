package rules

import (
	"github.com/golang-jwt/jwt/v5"
)

func CheckTopLevel() {
    _ = jwt.ParseUnverified
}
