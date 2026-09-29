package rules

import (
	"errors"
	"fmt"
)

type AuthManager struct{}

func (a *AuthManager) Authenticate(user, pass string) error {
	if user == "admin" && pass == "secret" {
		return nil
	}
	return errors.New("auth failed")
}

func (a *AuthManager) CheckPassword(user, pass string) (bool, error) {
	if pass == "secret" {
		return true, nil
	}
	return false, errors.New("invalid password")
}

func (a *AuthManager) VerifyCredentials(token string) (string, error) {
	if token == "valid" {
		return "user", nil
	}
	return "", errors.New("invalid credentials")
}

func (a *AuthManager) AuthWithoutVerification(user string) bool {
	return true
}

func (a *AuthManager) ParseUnverified(token string, claims map[string]string) error {
	return nil
}

func testIgnoredErrorBlank(a *AuthManager) {
	// ruleid: ignored-authentication-error
	_ = a.Authenticate("admin", "pass")
}

func testIgnoredErrorAssign(a *AuthManager) {
	// ruleid: ignored-authentication-error
	_, _ = a.CheckPassword("admin", "pass")
}

func testIgnoredErrorToken(a *AuthManager) {
	// ruleid: ignored-authentication-error
	user, _ := a.VerifyCredentials("bad-token")
	fmt.Println("User:", user)
}

func testUncheckedError(a *AuthManager) {
	// ruleid: ignored-authentication-error
	err := a.Authenticate("admin", "pass")
	fmt.Println("Proceeding without checking err:", err)
}

func testUnverifiedAuth(a *AuthManager) {
	// ruleid: ignored-authentication-error
	ok := a.AuthWithoutVerification("guest")
	fmt.Println("Auth without verification:", ok)
}

func testParseUnverified(a *AuthManager) {
	claims := make(map[string]string)
	// ruleid: ignored-authentication-error
	err := a.ParseUnverified("raw.jwt.token", claims)
	_ = err
}

func testSafeInlineCheck(a *AuthManager) error {
	// ok: ignored-authentication-error
	if err := a.Authenticate("admin", "pass"); err != nil {
		return err
	}
	return nil
}

func testSafeSeparateCheck(a *AuthManager) error {
	// ok: ignored-authentication-error
	err := a.Authenticate("admin", "pass")
	if err != nil {
		return err
	}
	return nil
}

func testSafeCheckNil(a *AuthManager) {
	// ok: ignored-authentication-error
	err := a.Authenticate("admin", "pass")
	if err == nil {
		fmt.Println("authenticated successfully")
	}
}

func testSafeReturn(a *AuthManager) error {
	// ok: ignored-authentication-error
	err := a.Authenticate("admin", "pass")
	return err
}
