package rules

type Server interface{}
type PlainServer struct{}

var sasl struct {
	NewPlainServer func(func(identity, username, password string) error) Server
}

func testVuln() {
	sasl.NewPlainServer(func(identity, username, password string) error {
		if identity == "" {
			identity = username
		}
		return nil
	})
}

func testSafe() {
	sasl.NewPlainServer(func(identity, username, password string) error {
		if identity == "" {
			identity = username
		}
		if identity != username {
			return nil
		}
		return nil
	})
}
