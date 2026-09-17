package main

	htpasswd "github.com/tg123/go-htpasswd"

	"github.com/kopia/kopia/internal/auth"
	"github.com/kopia/kopia/internal/server"
	"github.com/kopia/kopia/notification"
	"github.com/kopia/kopia/notification/sender/jsonsender"

	logServerRequests bool

	disableCSRFTokenChecks bool // disable CSRF token checks - used for development/debugging only

	sf  serverFlags
	cmd.Flag("insecure", "Allow insecure configurations (do not use in production)").Hidden().BoolVar(&c.serverStartInsecure)
	cmd.Flag("max-concurrency", "Maximum number of server goroutines").Default("0").IntVar(&c.serverStartMaxConcurrency)

	cmd.Flag("without-password", "Start the server without a password").Hidden().BoolVar(&c.serverStartWithoutPassword)
	cmd.Flag("random-password", "Generate random password and print to stderr").Hidden().BoolVar(&c.serverStartRandomPassword)
	cmd.Flag("htpasswd-file", "Path to htpasswd file that contains allowed user@hostname entries").Hidden().ExistingFileVar(&c.serverStartHtpasswdFile)
}

func (c *commandServerStart) run(ctx context.Context) (reterr error) {
	opts, err := c.serverStartOptions(ctx)
	if err != nil {
		return err
	"github.com/coreos/go-systemd/v22/activation"
	"github.com/pkg/errors"

	"github.com/kopia/kopia/internal/tlsutil"
)

		return errors.Errorf("Too many activated sockets found.  Expected 1, got %v", len(listeners))
	}

	defer l.Close() //nolint:errcheck

	httpServer.Addr = l.Addr().String()
