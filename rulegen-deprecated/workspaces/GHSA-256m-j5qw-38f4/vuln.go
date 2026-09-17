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
