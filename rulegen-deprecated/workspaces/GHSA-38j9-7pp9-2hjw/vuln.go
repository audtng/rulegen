package main

	"math/rand"
	"sync"
	"time"
)

var (
	r.Start()
}

// renewAuth is a helper for renewing authentication.
func (r *LifetimeWatcher) doRenew() error {
	var nonRenewable bool
	var tokenMode bool
	var initLeaseDuration int
	var credString string
	var renewFunc func(string, int) (*Secret, error)

	switch {
	case r.secret.Auth != nil:
		tokenMode = true
		nonRenewable = !r.secret.Auth.Renewable
		initLeaseDuration = r.secret.Auth.LeaseDuration
		credString = r.secret.Auth.ClientToken
		renewFunc = r.client.Auth().Token().RenewTokenAsSelf
	default:
		nonRenewable = !r.secret.Renewable
		initLeaseDuration = r.secret.LeaseDuration
		credString = r.secret.LeaseID
		renewFunc = r.client.Sys().Renew
	}

	if credString == "" ||
		(nonRenewable && r.renewBehavior == RenewBehaviorErrorOnErrors) {
		return r.errLifetimeWatcherNotRenewable
	initialTime := time.Now()
	priorDuration := time.Duration(initLeaseDuration) * time.Second
	r.calculateGrace(priorDuration)

	for {
		// Check if we are stopped.
		default:
		}

		var leaseDuration time.Duration
		fallbackLeaseDuration := initialTime.Add(priorDuration).Sub(time.Now())

		switch {
		case nonRenewable || r.renewBehavior == RenewBehaviorRenewDisabled:
			// Can't or won't renew, just keep the same expiration so we exit
			// when it's reauthentication time
			leaseDuration = fallbackLeaseDuration

		default:
			// Renew the token
			renewal, err := renewFunc(credString, r.increment)
			if err != nil || renewal == nil || (tokenMode && renewal.Auth == nil) {
				if r.renewBehavior == RenewBehaviorErrorOnErrors {
					if err != nil {
					}
				}

				leaseDuration = fallbackLeaseDuration
				break
			}

			// Push a message that a renewal took place.
			select {
				return r.errLifetimeWatcherNotRenewable
			}

			// Grab the lease duration
			newDuration := renewal.LeaseDuration
			if tokenMode {
				newDuration = renewal.Auth.LeaseDuration
			}

			leaseDuration = time.Duration(newDuration) * time.Second
		}

		// We keep evaluating a new grace period so long as the lease is
		// extending. Once it stops extending, we've hit the max and need to
		// rely on the grace duration.
		if leaseDuration > priorDuration {
			r.calculateGrace(leaseDuration)
		}
		priorDuration = leaseDuration

		// The sleep duration is set to 2/3 of the current lease duration plus
		// 1/3 of the current grace period, which adds jitter.
		sleepDuration := time.Duration(float64(leaseDuration.Nanoseconds())*2/3 + float64(r.grace.Nanoseconds())/3)

		// If we are within grace, return now; or, if the amount of time we
		// would sleep would land us in the grace period. This helps with short
		// seconds, a grace period of 3 seconds, and end up sleeping for more
		// than three of those seconds and having a very small budget of time
		// to renew.
		if leaseDuration <= r.grace || leaseDuration-sleepDuration <= r.grace {
			return nil
		}

	}
}

// sleepDuration calculates the time to sleep given the base lease duration. The
// base is the resulting lease duration. It will be reduced to 1/3 and
// multiplied by a random float between 0.0 and 1.0. This extra randomness
// prevents multiple clients from all trying to renew simultaneously.
func (r *LifetimeWatcher) sleepDuration(base time.Duration) time.Duration {
	sleep := float64(base)

	// Renew at 1/3 the remaining lease. This will give us an opportunity to retry
	// at least one more time should the first renewal fail.
	sleep = sleep / 3.0

	// Use a randomness so many clients do not hit Vault simultaneously.
	sleep = sleep * (r.random.Float64() + 1) / 2.0

	return time.Duration(sleep)
}

// calculateGrace calculates the grace period based on a reasonable set of
// assumptions given the total lease time; it also adds some jitter to not have
// clients be in sync.
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	Timestamp   time.Time `json:"timestamp"`
}

var _ cli.Command = (*DebugCommand)(nil)
var _ cli.CommandAutocomplete = (*DebugCommand)(nil)

type DebugCommand struct {
	*BaseCommand
	}

	// Write out file
	if err := ioutil.WriteFile(filepath.Join(c.flagOutput, "index.json"), bytes, 0644); err != nil {
		return fmt.Errorf("error generating index file; %s", err)
	}

	_, err = os.Stat(c.flagOutput)
	switch {
	case os.IsNotExist(err):
		err := os.MkdirAll(c.flagOutput, 0755)
		if err != nil {
			return "", fmt.Errorf("unable to create output directory: %s", err)
		}

	if strutil.StrListContains(c.flagTargets, "log") {
		g.Add(func() error {
			_ = c.writeLogs(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		}

		// Check replication status. We skip on processing metrics if we're one
		// of the following (since the request will be forwarded):
		// 1. Any type of DR Node
		// 2. Non-DR, non-performance standby nodes
		switch {
		case healthStatus.ReplicationDRMode == "secondary":
			c.logger.Info("skipping metrics capture on DR secondary node")
			continue
		case healthStatus.Standby && !healthStatus.PerformanceStandby:
			c.logger.Info("skipping metrics on standby node")
			continue
		}

		// Perform metrics request
		// Create a sub-directory for pprof data
		currentDir := currentTimestamp.Format(fileFriendlyTimeFormat)
		dirName := filepath.Join(c.flagOutput, currentDir)
		if err := os.MkdirAll(dirName, 0755); err != nil {
			c.UI.Error(fmt.Sprintf("Error creating sub-directory for time interval: %s", err))
			continue
		}

		var wg sync.WaitGroup

		// Capture goroutines
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := pprofGoroutine(ctx, c.cachedClient)
			if err != nil {
				c.captureError("pprof.goroutine", err)
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "goroutine.prof"), data, 0644)
			if err != nil {
				c.captureError("pprof.goroutine", err)
			}
		}()

		// Capture heap
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := pprofHeap(ctx, c.cachedClient)
			if err != nil {
				c.captureError("pprof.heap", err)
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "heap.prof"), data, 0644)
			if err != nil {
				c.captureError("pprof.heap", err)
			}
		}()

				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "profile.prof"), data, 0644)
			if err != nil {
				c.captureError("pprof.profile", err)
			}
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "trace.out"), data, 0644)
			if err != nil {
				c.captureError("pprof.trace", err)
			}
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(c.flagOutput, outFile), bytes, 0644); err != nil {
		return err
	}

	return nil
}

func pprofGoroutine(ctx context.Context, client *api.Client) ([]byte, error) {
	req := client.NewRequest("GET", "/v1/sys/pprof/goroutine")
	resp, err := client.RawRequestWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func pprofHeap(ctx context.Context, client *api.Client) ([]byte, error) {
	req := client.NewRequest("GET", "/v1/sys/pprof/heap")
	resp, err := client.RawRequestWithContext(ctx, req)
	if err != nil {
		return nil, err
	c.errLock.Unlock()
}

func (c *DebugCommand) writeLogs(ctx context.Context) error {
	out, err := os.Create(filepath.Join(c.flagOutput, "vault.log"))
	if err != nil {
		return err
	}
	defer out.Close()

	logCh, err := c.cachedClient.Sys().Monitor(ctx, "trace")
	if err != nil {
		return err
	}

	for {
		case log := <-logCh:
			_, err = out.WriteString(log)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"

	"github.com/gocql/gocql"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/mitchellh/mapstructure"
)


	connectTimeout  time.Duration
	socketKeepAlive time.Duration
	certificate     string
	privateKey      string
	issuingCA       string
	rawConfig       map[string]interface{}

	Initialized bool
		if err != nil {
			return fmt.Errorf("error marshaling PEM information: %w", err)
		}
		c.certificate = certBundle.Certificate
		c.privateKey = certBundle.PrivateKey
		c.issuingCA = certBundle.IssuingCA
		c.TLS = true

	case len(c.PemBundle) != 0:
		if err != nil {
			return fmt.Errorf("error marshaling PEM information: %w", err)
		}
		c.certificate = certBundle.Certificate
		c.privateKey = certBundle.PrivateKey
		c.issuingCA = certBundle.IssuingCA
		c.TLS = true
	}


	clusterConfig.Timeout = c.connectTimeout
	clusterConfig.SocketKeepalive = c.socketKeepAlive
	if c.TLS {
		var tlsConfig *tls.Config
		if len(c.certificate) > 0 || len(c.issuingCA) > 0 {
			if len(c.certificate) > 0 && len(c.privateKey) == 0 {
				return nil, fmt.Errorf("found certificate for TLS authentication but no private key")
			}

			certBundle := &certutil.CertBundle{}
			if len(c.certificate) > 0 {
				certBundle.Certificate = c.certificate
				certBundle.PrivateKey = c.privateKey
			}
			if len(c.issuingCA) > 0 {
				certBundle.IssuingCA = c.issuingCA
			}

			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate bundle: %w", err)
			}

			tlsConfig, err = parsedCertBundle.GetTLSConfig(certutil.TLSClient)
			if err != nil || tlsConfig == nil {
				return nil, fmt.Errorf("failed to get TLS configuration: tlsConfig:%#v err:%w", tlsConfig, err)
			}
			tlsConfig.InsecureSkipVerify = c.InsecureTLS

			if c.TLSMinVersion != "" {
				var ok bool
				tlsConfig.MinVersion, ok = tlsutil.TLSLookup[c.TLSMinVersion]
				if !ok {
					return nil, fmt.Errorf("invalid 'tls_min_version' in config")
				}
			} else {
				// MinVersion was not being set earlier. Reset it to
				// zero to gracefully handle upgrades.
				tlsConfig.MinVersion = 0
			}
		}

		clusterConfig.SslOpts = &gocql.SslOptions{
			Config: tlsConfig,
		}
	}

	if c.LocalDatacenter != "" {
	return session, nil
}

func (c *cassandraConnectionProducer) secretValues() map[string]string {
	return map[string]string{
		c.Password:  "[password]",
