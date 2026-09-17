package main

	"math/rand"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v3"
)

var (
	r.Start()
}

type renewFunc func(string, int) (*Secret, error)

// doRenew is a helper for renewing authentication.
func (r *LifetimeWatcher) doRenew() error {
	defaultInitialRetryInterval := 10 * time.Second
	switch {
	case r.secret.Auth != nil:
		return r.doRenewWithOptions(true, !r.secret.Auth.Renewable,
			r.secret.Auth.LeaseDuration, r.secret.Auth.ClientToken,
			r.client.Auth().Token().RenewTokenAsSelf, defaultInitialRetryInterval)
	default:
		return r.doRenewWithOptions(false, !r.secret.Renewable,
			r.secret.LeaseDuration, r.secret.LeaseID,
			r.client.Sys().Renew, defaultInitialRetryInterval)
	}
}

func (r *LifetimeWatcher) doRenewWithOptions(tokenMode bool, nonRenewable bool, initLeaseDuration int, credString string,
	renew renewFunc, initialRetryInterval time.Duration) error {
	if credString == "" ||
		(nonRenewable && r.renewBehavior == RenewBehaviorErrorOnErrors) {
		return r.errLifetimeWatcherNotRenewable
	initialTime := time.Now()
	priorDuration := time.Duration(initLeaseDuration) * time.Second
	r.calculateGrace(priorDuration)
	var errorBackoff backoff.BackOff

	for {
		// Check if we are stopped.
		default:
		}

		var remainingLeaseDuration time.Duration
		fallbackLeaseDuration := initialTime.Add(priorDuration).Sub(time.Now())
		var renewal *Secret
		var err error

		switch {
		case nonRenewable || r.renewBehavior == RenewBehaviorRenewDisabled:
			// Can't or won't renew, just keep the same expiration so we exit
			// when it's reauthentication time
			remainingLeaseDuration = fallbackLeaseDuration

		default:
			// Renew the token
			renewal, err = renew(credString, r.increment)
			if err != nil || renewal == nil || (tokenMode && renewal.Auth == nil) {
				if r.renewBehavior == RenewBehaviorErrorOnErrors {
					if err != nil {
					}
				}

				// Calculate remaining duration until initial token lease expires
				remainingLeaseDuration = initialTime.Add(time.Duration(initLeaseDuration) * time.Second).Sub(time.Now())
				if errorBackoff == nil {
					errorBackoff = &backoff.ExponentialBackOff{
						MaxElapsedTime:      remainingLeaseDuration,
						RandomizationFactor: backoff.DefaultRandomizationFactor,
						InitialInterval:     initialRetryInterval,
						MaxInterval:         5 * time.Minute,
						Multiplier:          2,
						Clock:               backoff.SystemClock,
					}
					errorBackoff.Reset()
				}
				break
			}
			errorBackoff = nil

			// Push a message that a renewal took place.
			select {
				return r.errLifetimeWatcherNotRenewable
			}

			// Reset initial time
			initialTime = time.Now()

			// Grab the lease duration
			initLeaseDuration = renewal.LeaseDuration
			if tokenMode {
				initLeaseDuration = renewal.Auth.LeaseDuration
			}

			remainingLeaseDuration = time.Duration(initLeaseDuration) * time.Second
		}

		var sleepDuration time.Duration

		if errorBackoff != nil {
			sleepDuration = errorBackoff.NextBackOff()
			if sleepDuration == backoff.Stop {
				return err
			}
		} else {
			// We keep evaluating a new grace period so long as the lease is
			// extending. Once it stops extending, we've hit the max and need to
			// rely on the grace duration.
			if remainingLeaseDuration > priorDuration {
				r.calculateGrace(remainingLeaseDuration)
			}
			priorDuration = remainingLeaseDuration

			// The sleep duration is set to 2/3 of the current lease duration plus
			// 1/3 of the current grace period, which adds jitter.
			sleepDuration = time.Duration(float64(remainingLeaseDuration.Nanoseconds())*2/3 + float64(r.grace.Nanoseconds())/3)
		}

		// If we are within grace, return now; or, if the amount of time we
		// would sleep would land us in the grace period. This helps with short
		// seconds, a grace period of 3 seconds, and end up sleeping for more
		// than three of those seconds and having a very small budget of time
		// to renew.
		if remainingLeaseDuration <= r.grace || remainingLeaseDuration-sleepDuration <= r.grace {
			return nil
		}

	}
}

// calculateGrace calculates the grace period based on a reasonable set of
// assumptions given the total lease time; it also adds some jitter to not have
// clients be in sync.
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	Timestamp   time.Time `json:"timestamp"`
}

var (
	_ cli.Command             = (*DebugCommand)(nil)
	_ cli.CommandAutocomplete = (*DebugCommand)(nil)
)

type DebugCommand struct {
	*BaseCommand
	}

	// Write out file
	if err := ioutil.WriteFile(filepath.Join(c.flagOutput, "index.json"), bytes, 0o644); err != nil {
		return fmt.Errorf("error generating index file; %s", err)
	}

	_, err = os.Stat(c.flagOutput)
	switch {
	case os.IsNotExist(err):
		err := os.MkdirAll(c.flagOutput, 0o755)
		if err != nil {
			return "", fmt.Errorf("unable to create output directory: %s", err)
		}

	if strutil.StrListContains(c.flagTargets, "log") {
		g.Add(func() error {
			c.writeLogs(ctx)
			// If writeLogs returned earlier due to an error, wait for context
			// to terminate so we don't abort everything.
			<-ctx.Done()
			return nil
		}, func(error) {
			cancelFunc()
		}

		// Check replication status. We skip on processing metrics if we're one
		// a DR node, though non-perf standbys will fail if they aren't using
		// unauthenticated_metrics_access.
		switch {
		case healthStatus.ReplicationDRMode == "secondary":
			c.logger.Info("skipping metrics capture on DR secondary node")
			continue
		}

		// Perform metrics request
		// Create a sub-directory for pprof data
		currentDir := currentTimestamp.Format(fileFriendlyTimeFormat)
		dirName := filepath.Join(c.flagOutput, currentDir)
		if err := os.MkdirAll(dirName, 0o755); err != nil {
			c.UI.Error(fmt.Sprintf("Error creating sub-directory for time interval: %s", err))
			continue
		}

		var wg sync.WaitGroup

		for _, target := range []string{"threadcreate", "allocs", "block", "mutex", "goroutine", "heap"} {
			wg.Add(1)
			go func(target string) {
				defer wg.Done()
				data, err := pprofTarget(ctx, c.cachedClient, target, nil)
				if err != nil {
					c.captureError("pprof."+target, err)
					return
				}

				err = ioutil.WriteFile(filepath.Join(dirName, target+".prof"), data, 0o644)
				if err != nil {
					c.captureError("pprof."+target, err)
				}
			}(target)
		}

		// As a convenience, we'll also fetch the goroutine target using debug=2, which yields a text
		// version of the stack traces that don't require using `go tool pprof` to view.
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := pprofTarget(ctx, c.cachedClient, "goroutine", url.Values{"debug": []string{"2"}})
			if err != nil {
				c.captureError("pprof.goroutines-text", err)
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "goroutines.txt"), data, 0o644)
			if err != nil {
				c.captureError("pprof.goroutines-text", err)
			}
		}()

				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "profile.prof"), data, 0o644)
			if err != nil {
				c.captureError("pprof.profile", err)
			}
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "trace.out"), data, 0o644)
			if err != nil {
				c.captureError("pprof.trace", err)
			}
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(c.flagOutput, outFile), bytes, 0o644); err != nil {
		return err
	}

	return nil
}

func pprofTarget(ctx context.Context, client *api.Client, target string, params url.Values) ([]byte, error) {
	req := client.NewRequest("GET", "/v1/sys/pprof/"+target)
	if params != nil {
		req.Params = params
	}
	resp, err := client.RawRequestWithContext(ctx, req)
	if err != nil {
		return nil, err
	c.errLock.Unlock()
}

func (c *DebugCommand) writeLogs(ctx context.Context) {
	out, err := os.Create(filepath.Join(c.flagOutput, "vault.log"))
	if err != nil {
		c.captureError("log", err)
		return
	}
	defer out.Close()

	logCh, err := c.cachedClient.Sys().Monitor(ctx, "trace")
	if err != nil {
		c.captureError("log", err)
		return
	}

	for {
		case log := <-logCh:
			_, err = out.WriteString(log)
			if err != nil {
				c.captureError("log", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
	"sync"
	"time"

	"github.com/gocql/gocql"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
	"github.com/mitchellh/mapstructure"
)


	connectTimeout  time.Duration
	socketKeepAlive time.Duration
	certBundle      *certutil.CertBundle
	rawConfig       map[string]interface{}

	Initialized bool
		if err != nil {
			return fmt.Errorf("error marshaling PEM information: %w", err)
		}
		c.certBundle = certBundle
		c.TLS = true

	case len(c.PemBundle) != 0:
		if err != nil {
			return fmt.Errorf("error marshaling PEM information: %w", err)
		}
		c.certBundle = certBundle
		c.TLS = true
	}

	if c.InsecureTLS {
		c.TLS = true
	}


	clusterConfig.Timeout = c.connectTimeout
	clusterConfig.SocketKeepalive = c.socketKeepAlive

	if c.TLS {
		sslOpts, err := getSslOpts(c.certBundle, c.TLSMinVersion, c.InsecureTLS)
		if err != nil {
			return nil, err
		}
		clusterConfig.SslOpts = sslOpts
	}

	if c.LocalDatacenter != "" {
	return session, nil
}

func getSslOpts(certBundle *certutil.CertBundle, minTLSVersion string, insecureSkipVerify bool) (*gocql.SslOptions, error) {
	tlsConfig := &tls.Config{}
	if certBundle != nil {
		if certBundle.Certificate == "" && certBundle.PrivateKey != "" {
			return nil, fmt.Errorf("found private key for TLS authentication but no certificate")
		}
		if certBundle.Certificate != "" && certBundle.PrivateKey == "" {
			return nil, fmt.Errorf("found certificate for TLS authentication but no private key")
		}

		parsedCertBundle, err := certBundle.ToParsedCertBundle()
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate bundle: %w", err)
		}

		tlsConfig, err = parsedCertBundle.GetTLSConfig(certutil.TLSClient)
		if err != nil {
			return nil, fmt.Errorf("failed to get TLS configuration: tlsConfig:%#v err:%w", tlsConfig, err)
		}
	}

	tlsConfig.InsecureSkipVerify = insecureSkipVerify

	if minTLSVersion != "" {
		var ok bool
		tlsConfig.MinVersion, ok = tlsutil.TLSLookup[minTLSVersion]
		if !ok {
			return nil, fmt.Errorf("invalid 'tls_min_version' in config")
		}
	} else {
		// MinVersion was not being set earlier. Reset it to
		// zero to gracefully handle upgrades.
		tlsConfig.MinVersion = 0
	}

	opts := &gocql.SslOptions{
		Config:                 tlsConfig,
		EnableHostVerification: !insecureSkipVerify,
	}
	return opts, nil
}

func (c *cassandraConnectionProducer) secretValues() map[string]string {
	return map[string]string{
		c.Password:  "[password]",
