package rules

import (
	"fmt"
	"net/http"
)

type SimpleBindRequest struct {
	Username string
	Password string
}

type LDAPConn struct{}

func (c *LDAPConn) Bind(username, password string) error {
	return nil
}

func (c *LDAPConn) UnauthenticatedBind(username string) error {
	return nil
}

func (c *LDAPConn) SimpleBind(req *SimpleBindRequest) error {
	return nil
}

func NewSimpleBindRequest(username, password string) *SimpleBindRequest {
	return &SimpleBindRequest{Username: username, Password: password}
}

func testVulnEmptyPassword(conn *LDAPConn, user string) {
	// ruleid: ldap-improper-authentication
	_ = conn.Bind(user, "")
}

func testVulnEmptyUsername(conn *LDAPConn, pass string) {
	// ruleid: ldap-improper-authentication
	_ = conn.Bind("", pass)
}

func testVulnUnauthenticatedBind(conn *LDAPConn, user string) {
	// ruleid: ldap-improper-authentication
	_ = conn.UnauthenticatedBind(user)
}

func testVulnSimpleBindEmptyPass() {
	// ruleid: ldap-improper-authentication
	_ = &SimpleBindRequest{
		Username: "admin",
		Password: "",
	}
}

func testVulnNewSimpleBindEmptyPass(user string) {
	// ruleid: ldap-improper-authentication
	_ = NewSimpleBindRequest(user, "")
}

func testVulnIgnoredError(conn *LDAPConn, user, pass string) {
	// ruleid: ldap-improper-authentication
	err := conn.Bind(user, pass)
	fmt.Println("Proceeding without checking err:", err)
}

func testVulnMissingReturnHttpError(w http.ResponseWriter, conn *LDAPConn, user, pass string) {
	// ruleid: ldap-improper-authentication
	if err := conn.Bind(user, pass); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	fmt.Println("Proceeding after failed auth")
}

func testSafeBind(conn *LDAPConn, user, pass string) error {
	// ok: ldap-improper-authentication
	if err := conn.Bind(user, pass); err != nil {
		return err
	}
	return nil
}

func testSafeAssignedCheck(conn *LDAPConn, user, pass string) error {
	// ok: ldap-improper-authentication
	err := conn.Bind(user, pass)
	if err != nil {
		return err
	}
	return nil
}

func testSafeHttpHandler(w http.ResponseWriter, conn *LDAPConn, user, pass string) {
	// ok: ldap-improper-authentication
	if err := conn.Bind(user, pass); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println("Authenticated successfully")
}
