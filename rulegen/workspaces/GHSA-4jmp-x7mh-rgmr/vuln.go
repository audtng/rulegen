package main

)

type Config struct {
	LogLevel       string          `long:"loglevel" description:"Logging level for all subsystems" choice:"trace" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal"`
	KeyringBackend string          `long:"keyring-type" description:"Type of keyring to use"`
	RPCListener    string          `long:"rpclistener" description:"the listener for RPC connections, e.g., 127.0.0.1:1234"`
	HMACKey        string          `long:"hmackey" description:"The HMAC key for authentication with FPD. If not provided, will use HMAC_KEY environment variable."`
	Metrics        *metrics.Config `group:"metrics" namespace:"metrics"`

	DatabaseConfig *DBConfig `group:"dbconfig" namespace:"dbconfig"`
}

// LoadConfig initializes and parses the config using a config file and command
}

func DefaultConfigWithHomePath(homePath string) *Config {
	cfg := &Config{
		LogLevel:       defaultLogLevel,
		KeyringBackend: defaultKeyringBackend,
		DatabaseConfig: DefaultDBConfigWithHomePath(homePath),
		RPCListener:    defaultRPCListener,
		Metrics:        metrics.DefaultEotsConfig(),
	}
	if err := cfg.Validate(); err != nil {
		panic(err)

	return cfg
}
		return nil, fmt.Errorf("failed to get EOTS private key: %w", err)
	}

	// Update metrics
	lm.metrics.IncrementEotsFpTotalEotsSignCounter(hex.EncodeToString(eotsPk))
	lm.metrics.SetEotsFpLastEotsSignHeight(hex.EncodeToString(eotsPk), float64(height))
		return nil, err
	}

	return lm.eotsPrivKeyFromKeyName(keyName)
}

func (lm *LocalEOTSManager) eotsPrivKeyFromKeyName(keyName string) (*btcec.PrivateKey, error) {
import (
	"context"

	"github.com/btcsuite/btcd/btcec/v2"
	"google.golang.org/grpc"

	"github.com/babylonlabs-io/finality-provider/eotsmanager"
	"github.com/babylonlabs-io/finality-provider/eotsmanager/proto"
type rpcServer struct {
	proto.UnimplementedEOTSManagerServer

	em *eotsmanager.LocalEOTSManager
}

// newRPCServer creates a new RPC sever from the set of input dependencies.
func newRPCServer(
	em *eotsmanager.LocalEOTSManager,
) *rpcServer {
	return &rpcServer{
		em: em,
	}
}

// UnsafeSignEOTS only used for testing purposes. Doesn't offer slashing protection!
func (r *rpcServer) UnsafeSignEOTS(_ context.Context, req *proto.SignEOTSRequest) (
	*proto.SignEOTSResponse, error) {
	sig, err := r.em.UnsafeSignEOTS(req.Uid, req.ChainId, req.Msg, req.Height)
	if err != nil {
		return nil, err
