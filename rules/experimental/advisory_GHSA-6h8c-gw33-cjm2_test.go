package rules

import (
	"fmt"
	"net/http"
)

// Upgrader mocks gorilla/websocket.Upgrader interface for standalone compilation
type Upgrader struct {
	CheckOrigin func(r *http.Request) bool
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (interface{}, error) {
	return nil, nil
}

// 1. Insecure: CheckOrigin in struct literal returns true unconditionally
func testVulnCheckOriginStruct() {
	var upgrader = Upgrader{
		// ruleid: websocket-improper-authentication
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	_ = upgrader
}

// 2. Insecure: CheckOrigin assignment returns true unconditionally
func testVulnCheckOriginAssign() {
	var upgrader Upgrader
	// ruleid: websocket-improper-authentication
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	_ = upgrader
}

// 3. Insecure: Ignored Upgrade error using blank identifier
func testVulnUpgradeIgnoredError(w http.ResponseWriter, r *http.Request) {
	var upgrader Upgrader
	// ruleid: websocket-improper-authentication
	conn, _ := upgrader.Upgrade(w, r, nil)
	_ = conn
}

// 4. Insecure: Unchecked Upgrade error proceeding to connection usage
func testVulnUpgradeUncheckedError(w http.ResponseWriter, r *http.Request) {
	var upgrader Upgrader
	// ruleid: websocket-improper-authentication
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("Proceeding without checking err:", err)
	_ = conn
}

// 5. Safe: CheckOrigin validates request origin against allowed host
func testSafeCheckOriginValidating(w http.ResponseWriter, r *http.Request) {
	var upgrader = Upgrader{
		// ok: websocket-improper-authentication
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "https://trusted.example.com"
		},
	}
	// ok: websocket-improper-authentication
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	_ = conn
}

// 6. Safe: CheckOrigin assignment validates origin against expected host
func testSafeCheckOriginAssign() {
	var upgrader Upgrader
	// ok: websocket-improper-authentication
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return r.Host == "localhost:8090"
	}
	_ = upgrader
}

// 7. Safe: Upgrade error is handled with error response and early return
func testSafeUpgradeProperCheck(w http.ResponseWriter, r *http.Request) {
	var upgrader Upgrader
	// ok: websocket-improper-authentication
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Upgrade failed", http.StatusBadRequest)
		return
	}
	_ = conn
}

// 8. Safe: Upgrade called inside if initializer with error check
func testSafeUpgradeIfInit(w http.ResponseWriter, r *http.Request) {
	var upgrader Upgrader
	// ok: websocket-improper-authentication
	if conn, err := upgrader.Upgrade(w, r, nil); err != nil {
		http.Error(w, "Upgrade failed", http.StatusUnauthorized)
		return
	} else {
		_ = conn
	}
}
