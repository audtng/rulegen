package rules

import (
	"net/http"
)

type Account struct {
	Username      string
	RecoveryCodes []string
	MFAEnabled    bool
}

func (a *Account) ValidateRecoveryCode(code string) bool {
	return code == "valid-code"
}

func (a *Account) VerifyRecoveryCode(code string) bool {
	return code == "valid-code"
}

type AuthService struct{}

func (s *AuthService) UseRecoveryCode(a *Account, code string) (bool, error) {
	return true, nil
}

func (s *AuthService) GenerateRecoveryCodes(a *Account) []string {
	return []string{"code1", "code2"}
}

func generateRecoveryCodes(a *Account) []string {
	return []string{"code1", "code2"}
}

// Edge Case 1: Ignored boolean / result from recovery code validation
func testVulnIgnoredRecoveryCodeValidation(user *Account, code string) {
	// ruleid: mfa-recovery-code-auth-bypass
	_ = user.ValidateRecoveryCode(code)
}

func testSafeRecoveryCodeValidation(user *Account, code string) bool {
	// ok: mfa-recovery-code-auth-bypass
	if !user.ValidateRecoveryCode(code) {
		return false
	}
	return true
}

// Edge Case 2: Missing return after recovery code validation failure
func testVulnMissingReturnRecoveryFailure(w http.ResponseWriter, r *http.Request, user *Account, code string) {
	// ruleid: mfa-recovery-code-auth-bypass
	if !user.VerifyRecoveryCode(code) {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.Write([]byte("sensitive data accessed"))
}

func testSafeReturnRecoveryFailure(w http.ResponseWriter, r *http.Request, user *Account, code string) {
	// ok: mfa-recovery-code-auth-bypass
	if !user.VerifyRecoveryCode(code) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte("sensitive data accessed"))
}

// Edge Case 3: Generating/assigning recovery codes without MFA enabled check
func testVulnGenerateRecoveryCodesWithoutMFA(user *Account) {
	// ruleid: mfa-recovery-code-auth-bypass
	user.RecoveryCodes = generateRecoveryCodes(user)
}

func testSafeGenerateRecoveryCodesWithMFA(user *Account) {
	if !user.MFAEnabled {
		return
	}
	// ok: mfa-recovery-code-auth-bypass
	user.RecoveryCodes = generateRecoveryCodes(user)
}

// Edge Case 4: Ignored error when consuming recovery code
func testVulnIgnoredErrorConsumeRecovery(auth *AuthService, user *Account, code string) {
	// ruleid: mfa-recovery-code-auth-bypass
	_, _ = auth.UseRecoveryCode(user, code)
}

func testSafeCheckedErrorConsumeRecovery(auth *AuthService, user *Account, code string) bool {
	// ok: mfa-recovery-code-auth-bypass
	ok, err := auth.UseRecoveryCode(user, code)
	if err != nil || !ok {
		return false
	}
	return true
}

// Edge Case 5: Missing return with http.Error after recovery validation error
func testVulnMissingReturnHttpError(w http.ResponseWriter, r *http.Request, user *Account, code string) {
	// ruleid: mfa-recovery-code-auth-bypass
	if !user.ValidateRecoveryCode(code) {
		http.Error(w, "Invalid recovery code", http.StatusForbidden)
	}
	w.Write([]byte("unauthorized access"))
}

func testSafeReturnHttpError(w http.ResponseWriter, r *http.Request, user *Account, code string) {
	// ok: mfa-recovery-code-auth-bypass
	if !user.ValidateRecoveryCode(code) {
		http.Error(w, "Invalid recovery code", http.StatusForbidden)
		return
	}
	w.Write([]byte("authorized access"))
}
