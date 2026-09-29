package rules

import (
	"errors"
)

type Server interface{}

type Session struct {
	User string
}

var sasl struct {
	NewPlainServer func(func(identity, username, password string) error) Server
}

func checkCredentials(username, password string) error {
	if username == "admin" && password == "secret" {
		return nil
	}
	return errors.New("invalid credentials")
}

// Case 1: Normalizing identity = username, but failing to verify identity == username (CVE-2023-27582)
func testVulnIdentityNormalization(s *Session) {
	// ruleid: sasl-plain-authorization-bypass
	sasl.NewPlainServer(func(identity, username, password string) error {
		if identity == "" {
			identity = username
		}
		if err := checkCredentials(username, password); err != nil {
			return err
		}
		s.User = identity
		return nil
	})
}

// Safe Case 1: Identity normalized and verified against username
func testSafeIdentityNormalization(s *Session) {
	// ok: sasl-plain-authorization-bypass
	sasl.NewPlainServer(func(identity, username, password string) error {
		if identity == "" {
			identity = username
		}
		if identity != username {
			return errors.New("authorization identity does not match authenticated user")
		}
		if err := checkCredentials(username, password); err != nil {
			return err
		}
		s.User = identity
		return nil
	})
}

// Case 2: Directly assigning identity to session without any verification
func testVulnDirectIdentityAssignment(s *Session) {
	// ruleid: sasl-plain-authorization-bypass
	sasl.NewPlainServer(func(identity, username, password string) error {
		if err := checkCredentials(username, password); err != nil {
			return err
		}
		s.User = identity
		return nil
	})
}

// Safe Case 2: Checking authorization identity before assignment
func testSafeDirectIdentityAssignment(s *Session) {
	// ok: sasl-plain-authorization-bypass
	sasl.NewPlainServer(func(identity, username, password string) error {
		if identity != "" && identity != username {
			return errors.New("unauthorized identity")
		}
		if err := checkCredentials(username, password); err != nil {
			return err
		}
		s.User = identity
		return nil
	})
}

// Case 3: Calling authentication function but ignoring the returned error with blank identifier
func testVulnIgnoredAuthError(s *Session) {
	// ruleid: sasl-plain-authorization-bypass
	sasl.NewPlainServer(func(identity, username, password string) error {
		_ = checkCredentials(username, password)
		s.User = username
		return nil
	})
}

// Safe Case 3: Error is properly checked
func testSafeAuthErrorChecked(s *Session) {
	// ok: sasl-plain-authorization-bypass
	sasl.NewPlainServer(func(identity, username, password string) error {
		if err := checkCredentials(username, password); err != nil {
			return err
		}
		s.User = username
		return nil
	})
}

// Case 4: Unconditional return nil without any validation checks
func testVulnUnconditionalSuccess() {
	// ruleid: sasl-plain-authorization-bypass
	sasl.NewPlainServer(func(identity, username, password string) error {
		return nil
	})
}

// Safe Case 4: Validating credentials with if condition
func testSafeCredentialValidation() {
	// ok: sasl-plain-authorization-bypass
	sasl.NewPlainServer(func(identity, username, password string) error {
		if username != "admin" || password != "secret" {
			return errors.New("unauthorized")
		}
		return nil
	})
}

// Safe Case 5: Completely ignoring identity and using authenticated username
func testSafeIgnoreIdentity(s *Session) {
	// ok: sasl-plain-authorization-bypass
	sasl.NewPlainServer(func(identity, username, password string) error {
		if err := checkCredentials(username, password); err != nil {
			return err
		}
		s.User = username
		return nil
	})
}
