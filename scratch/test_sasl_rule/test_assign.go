package rules

type Session struct {
    User string
}

type Server interface{}
var sasl struct {
    NewPlainServer func(func(string, string, string) error) Server
}

func testField(s *Session) {
    sasl.NewPlainServer(func(identity, username, password string) error {
        s.User = identity
        return nil
    })
}
