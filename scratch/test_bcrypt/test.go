package rules

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func testVulnIgnoredBlank(hash, password []byte) {
	// Discarding error via blank identifier
	_ = bcrypt.CompareHashAndPassword(hash, password)
}

func testVulnIgnoredAssignment(hash, password []byte) {
	// Assigning error but never checking it
	err := bcrypt.CompareHashAndPassword(hash, password)
	fmt.Println("Authenticated", err)
}

func testVulnReassignIgnored(hash, password []byte) {
	var err error
	err = bcrypt.CompareHashAndPassword(hash, password)
	fmt.Println("Authenticated", err)
}

func testVulnEmptyPassword(hash []byte) {
	_ = bcrypt.CompareHashAndPassword(hash, []byte(""))
}

func testVulnEmptyHash(password []byte) {
	_ = bcrypt.CompareHashAndPassword([]byte(""), password)
}

func testSafeCheck(hash, password []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return err
	}
	return nil
}

func testSafeAssignedCheck(hash, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return err
	}
	return nil
}
