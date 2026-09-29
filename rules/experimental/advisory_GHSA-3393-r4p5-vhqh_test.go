package rules

import (
	"net/http"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
)

func testVulnIgnoredTotpValidate(passcode, secret string) {
	// ruleid: otp-unverified-authentication
	_ = totp.Validate(passcode, secret)
}

func testVulnIgnoredHotpValidate(passcode, secret string, count uint64) {
	// ruleid: otp-unverified-authentication
	_ = hotp.Validate(passcode, count, secret)
}

func testVulnIgnoredTotpCustomValidate(passcode, secret string) {
	opts := totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}
	// ruleid: otp-unverified-authentication
	_, _ = totp.ValidateCustom(passcode, secret, time.Now(), opts)
}

func testVulnIgnoredHotpCustomValidate(passcode, secret string, count uint64) {
	opts := hotp.ValidateOpts{
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}
	// ruleid: otp-unverified-authentication
	_, _ = hotp.ValidateCustom(passcode, count, secret, opts)
}

func testVulnIgnoredBoolCustomValidate(passcode, secret string) error {
	opts := totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}
	// ruleid: otp-unverified-authentication
	_, err := totp.ValidateCustom(passcode, secret, time.Now(), opts)
	return err
}

func testVulnEmptySecret(passcode string) bool {
	// ruleid: otp-unverified-authentication
	return totp.Validate(passcode, "")
}

func testVulnEmptyPasscode(secret string) bool {
	// ruleid: otp-unverified-authentication
	return totp.Validate("", secret)
}

func testVulnMissingReturnHttpError(w http.ResponseWriter, passcode, secret string) {
	ok := totp.Validate(passcode, secret)
	// ruleid: otp-unverified-authentication
	if !ok {
		http.Error(w, "Invalid OTP token", http.StatusUnauthorized)
	}
	w.Write([]byte("Access Granted"))
}

func testVulnMissingReturnWriteHeader(w http.ResponseWriter, passcode, secret string) {
	// ruleid: otp-unverified-authentication
	if !totp.Validate(passcode, secret) {
		w.WriteHeader(http.StatusForbidden)
	}
	w.Write([]byte("Access Granted"))
}

func testSafeTotpValidate(w http.ResponseWriter, passcode, secret string) {
	ok := totp.Validate(passcode, secret)
	// ok: otp-unverified-authentication
	if !ok {
		http.Error(w, "Invalid OTP token", http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Access Granted"))
}

func testSafeTotpInlineValidate(w http.ResponseWriter, passcode, secret string) {
	// ok: otp-unverified-authentication
	if !totp.Validate(passcode, secret) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Access Granted"))
}

func testSafeTotpCustomValidate(passcode, secret string) (bool, error) {
	opts := totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}
	// ok: otp-unverified-authentication
	valid, err := totp.ValidateCustom(passcode, secret, time.Now(), opts)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, nil
	}
	return true, nil
}

func testSafeHotpValidate(passcode, secret string, count uint64) bool {
	// ok: otp-unverified-authentication
	return hotp.Validate(passcode, count, secret)
}
