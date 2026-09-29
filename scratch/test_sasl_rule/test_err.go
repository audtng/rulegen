package rules

type Server interface{}
var sasl struct {
    NewPlainServer func(func(string, string, string) error) Server
}
func auth(u, p string) error { return nil }

func testVuln() {
    sasl.NewPlainServer(func(identity, username, password string) error {
        err := auth(username, password)
        return nil
    })
}

func testSafe() {
    sasl.NewPlainServer(func(identity, username, password string) error {
        if err := auth(username, password); err != nil {
            return err
        }
        return nil
    })
}

func testSafe2() {
    sasl.NewPlainServer(func(identity, username, password string) error {
        err := auth(username, password)
        if err != nil {
            return err
        }
        return nil
    })
}
