package rules

type Server interface{}

var sasl struct {
	NewPlainServer func(func(identity, username, password string) error) Server
}

func testSafe2() {
	sasl.NewPlainServer(func(identity, username, password string) error {
		if identity == "" {
			identity = username
		}
		if identity != "" && identity != username {
			return nil
		}
		return nil
	})
}
