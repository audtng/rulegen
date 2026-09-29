package rules

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func testVulnIgnoredBlank(hash, password []byte) {
	// ruleid: bcrypt-unverified-authentication
	_ = bcrypt.CompareHashAndPassword(hash, password)
}

func testVulnIgnoredAssignment(hash, password []byte) {
	// ruleid: bcrypt-unverified-authentication
	err := bcrypt.CompareHashAndPassword(hash, password)
	fmt.Println("Proceeding with login despite error status", err)
}

func testVulnReassignIgnored(hash, password []byte) {
	var err error
	// ruleid: bcrypt-unverified-authentication
	err = bcrypt.CompareHashAndPassword(hash, password)
	fmt.Println("Proceeding with login despite error status", err)
}

func testVulnEmptyPassword(hash []byte) {
	// ruleid: bcrypt-unverified-authentication
	if err := bcrypt.CompareHashAndPassword(hash, []byte("")); err != nil {
		return
	}
}

func testVulnEmptyHash(password []byte) {
	// ruleid: bcrypt-unverified-authentication
	if err := bcrypt.CompareHashAndPassword([]byte(""), password); err != nil {
		return
	}
}

func testVulnMissingReturnHttpError(w http.ResponseWriter, hash, password []byte) {
	// ruleid: bcrypt-unverified-authentication
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	fmt.Fprintf(w, "Welcome authenticated user")
}

func testVulnMissingReturnWriteHeader(w http.ResponseWriter, hash, password []byte) {
	// ruleid: bcrypt-unverified-authentication
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	fmt.Fprintf(w, "Welcome authenticated user")
}

func testSafeCheckInline(hash, password []byte) error {
	// ok: bcrypt-unverified-authentication
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return err
	}
	return nil
}

func testSafeAssignedCheck(hash, password []byte) error {
	// ok: bcrypt-unverified-authentication
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return err
	}
	return nil
}

func testSafeHttpError(w http.ResponseWriter, hash, password []byte) {
	// ok: bcrypt-unverified-authentication
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Welcome authenticated user")
}

func testSafeWriteHeader(w http.ResponseWriter, hash, password []byte) {
	// ok: bcrypt-unverified-authentication
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Welcome authenticated user")
}
