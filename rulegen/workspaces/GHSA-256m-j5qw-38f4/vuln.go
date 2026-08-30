package main

	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	// start here
	username := params["username"]
	user, err := logic.GetUser(username)
	if err != nil {
		logger.Log(0, username,
		logic.ReturnErrorResponse(w, r, logic.FormatError(err, "badrequest"))
		return
	}
	userchange.Networks = nil
	user, err = logic.UpdateUser(&userchange, user)
	if err != nil {
import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	return "", err
}

// VerifyUserToken func will used to Verify the JWT Token while using APIS
func VerifyUserToken(tokenString string) (username string, networks []string, isadmin bool, err error) {
	claims := &models.UserClaims{}
controllers/user.go | 2 +-
logic/jwts.go       | 2 +-
2 files changed, 2 insertions(+), 2 deletions(-)
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	// start here
	jwtUser, _, isadmin, err := logic.VerifyJWS(r.Header.Get("Authorization"))
	if err != nil {
		logger.Log(0, "verifyJWT error", err.Error())
		logic.ReturnErrorResponse(w, r, logic.FormatError(err, "internal"))
}

// VerifyJWT verifies Auth Header
func VerifyJWS(bearerToken string) (username string, networks []string, isadmin bool, err error) {
	token := ""
	tokenSplit := strings.Split(bearerToken, " ")
	if len(tokenSplit) > 1 {
