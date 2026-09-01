package main

	modconfig "github.com/foxcpp/maddy/internal/config/module"
	"github.com/foxcpp/maddy/internal/log"
	"github.com/foxcpp/maddy/internal/module"
)

var (
	return mechs
}

func (s *SASLAuth) AuthPlain(username, password string) error {
	if len(s.Plain) == 0 {
		return ErrUnsupportedMech
	}

	var lastErr error
	for _, p := range s.Plain {
		err := p.AuthPlain(username, password)
		if err == nil {
			return nil
		}
		if err != nil {
			lastErr = err
			continue
		}
	}

	return fmt.Errorf("no auth. provider accepted creds, last err: %w", lastErr)
}

// CreateSASL creates the sasl.Server instance for the corresponding mechanism.
func (s *SASLAuth) CreateSASL(mech string, remoteAddr net.Addr, successCb func(identity string) error) sasl.Server {
	switch mech {
	case sasl.Plain:
		return sasl.NewPlainServer(func(identity, username, password string) error {
			if identity == "" {
				identity = username
			}

			err := s.AuthPlain(username, password)
			if err != nil {
				s.Log.Error("authentication failed", err, "username", username, "src_ip", remoteAddr)
				return errors.New("auth: invalid credentials")
			}

			return successCb(identity)
		})
	case sasl.Login:
		return sasl.NewLoginServer(func(username, password string) error {
			err := s.AuthPlain(username, password)
			if err != nil {
				s.Log.Error("authentication failed", err, "username", username, "src_ip", remoteAddr)
				return errors.New("auth: invalid credentials")
			}

			return successCb(username)
		})
	}
	return FailingSASLServ{Err: ErrUnsupportedMech}
// AddProvider adds the SASL authentication provider to its mapping by parsing
// the 'auth' configuration directive.
func (s *SASLAuth) AddProvider(m *config.Map, node *config.Node) error {
	var any interface{}
	if err := modconfig.ModuleFromNode(node.Args, node, m.Globals, &any); err != nil {
		return err
	}

	hasAny := false
	if plainAuth, ok := any.(module.PlainAuth); ok {
		s.Plain = append(s.Plain, plainAuth)
		hasAny = true
	}

	if !hasAny {
		return m.MatchErr("auth: specified module does not provide any SASL mechanism")
	}
	return nil
}

	"github.com/foxcpp/maddy/internal/config"
	"github.com/foxcpp/maddy/internal/log"
	"github.com/foxcpp/maddy/internal/module"
)

type Auth struct {
	return nil
}

func (a *Auth) AuthPlain(username, password string) error {
	ok := len(a.userTbls) == 0
	for _, tbl := range a.userTbls {
		_, tblOk, err := tbl.Lookup(username)
		if err != nil {
			return err
		}
		if tblOk {
			ok = true
			break
		}
	}
	if !ok {
		return errors.New("user not found in tables")
	}

	var lastErr error
	for _, p := range a.passwd {
		if err := p.AuthPlain(username, password); err != nil {
			lastErr = err
			continue
		}

		return nil
	}
	return lastErr
}

func init() {
// PlainAuth is the interface implemented by modules providing authentication using
// username:password pairs.
type PlainAuth interface {
	AuthPlain(username, password string) error
}
