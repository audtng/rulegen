package main

)

type Config struct {
	LogLevel               string          `long:"loglevel" description:"Logging level for all subsystems" choice:"trace" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal"`
	KeyringBackend         string          `long:"keyring-type" description:"Type of keyring to use"`
	RPCListener            string          `long:"rpclistener" description:"the listener for RPC connections, e.g., 127.0.0.1:1234"`
	HMACKey                string          `long:"hmackey" description:"The HMAC key for authentication with FPD. If not provided, will use HMAC_KEY environment variable."`
	Metrics                *metrics.Config `group:"metrics" namespace:"metrics"`
	DisableUnsafeEndpoints *bool           `long:"disable-unsafe-endpoints" description:"Disable unsafe RPC endpoints (e.g., UnsafeSignEOTS) that bypass slashing protection. Defaults to true (disabled) if not set."`
	DatabaseConfig         *DBConfig       `group:"dbconfig" namespace:"dbconfig"`
}

// LoadConfig initializes and parses the config using a config file and command
}

func DefaultConfigWithHomePath(homePath string) *Config {
	disableUnsafe := true
	cfg := &Config{
		LogLevel:               defaultLogLevel,
		KeyringBackend:         defaultKeyringBackend,
		DatabaseConfig:         DefaultDBConfigWithHomePath(homePath),
		RPCListener:            defaultRPCListener,
		Metrics:                metrics.DefaultEotsConfig(),
		DisableUnsafeEndpoints: &disableUnsafe,
	}
	if err := cfg.Validate(); err != nil {
		panic(err)

	return cfg
}

// IsUnsafeEndpointsDisabled returns true if unsafe endpoints should be disabled.
// Defaults to true (safe) if not explicitly set.
func (cfg *Config) IsUnsafeEndpointsDisabled() bool {
	if cfg.DisableUnsafeEndpoints == nil {
		return true // Safe default: disabled
	}

	return *cfg.DisableUnsafeEndpoints
}
		return nil, fmt.Errorf("failed to get EOTS private key: %w", err)
	}

	// Verify the retrieved private key corresponds to the requested public key
	derivedPubKey := schnorr.SerializePubKey(privKey.PubKey())
	if !bytes.Equal(derivedPubKey, eotsPk) {
		return nil, fmt.Errorf("public key mismatch: requested key does not match stored key")
	}

	// Update metrics
	lm.metrics.IncrementEotsFpTotalEotsSignCounter(hex.EncodeToString(eotsPk))
	lm.metrics.SetEotsFpLastEotsSignHeight(hex.EncodeToString(eotsPk), float64(height))
		return nil, err
	}

	privKey, err := lm.eotsPrivKeyFromKeyName(keyName)
	if err != nil {
		return nil, err
	}

	// Verify the retrieved private key corresponds to the requested public key
	derivedPubKey := schnorr.SerializePubKey(privKey.PubKey())
	if !bytes.Equal(derivedPubKey, fpPk) {
		return nil, fmt.Errorf("public key mismatch: requested key does not match stored key")
	}

	return privKey, nil
}

func (lm *LocalEOTSManager) eotsPrivKeyFromKeyName(keyName string) (*btcec.PrivateKey, error) {
import (
	"context"

	"github.com/babylonlabs-io/finality-provider/eotsmanager/config"
	"github.com/btcsuite/btcd/btcec/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/babylonlabs-io/finality-provider/eotsmanager"
	"github.com/babylonlabs-io/finality-provider/eotsmanager/proto"
type rpcServer struct {
	proto.UnimplementedEOTSManagerServer

	em  *eotsmanager.LocalEOTSManager
	cfg *config.Config
}

// newRPCServer creates a new RPC sever from the set of input dependencies.
func newRPCServer(
	em *eotsmanager.LocalEOTSManager,
	cfg *config.Config,
) *rpcServer {
	return &rpcServer{
		em:  em,
		cfg: cfg,
	}
}

// UnsafeSignEOTS only used for testing purposes. Doesn't offer slashing protection!
func (r *rpcServer) UnsafeSignEOTS(_ context.Context, req *proto.SignEOTSRequest) (
	*proto.SignEOTSResponse, error) {
	if r.cfg.IsUnsafeEndpointsDisabled() {
		return nil, status.Error(codes.PermissionDenied, //nolint:wrapcheck
			"UnsafeSignEOTS endpoint is disabled in configuration for security reasons")
	}
	sig, err := r.em.UnsafeSignEOTS(req.Uid, req.ChainId, req.Msg, req.Height)
	if err != nil {
		return nil, err
