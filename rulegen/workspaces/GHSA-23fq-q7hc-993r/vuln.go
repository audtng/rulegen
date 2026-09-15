package main

	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/circonus"
	"github.com/armon/go-metrics/datadog"
	"github.com/armon/go-metrics/prometheus"
	stackdriver "github.com/google/go-metrics-stackdriver"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command/server"
	serverseal "github.com/hashicorp/vault/command/server/seal"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/gatedwriter"
	"github.com/hashicorp/vault/internalshared/reloadutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/mlock"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/posener/complete"
	"go.uber.org/atomic"
	"golang.org/x/net/http/httpproxy"
	"google.golang.org/api/option"
	"google.golang.org/grpc/grpclog"
)

var _ cli.Command = (*ServerCommand)(nil)
var _ cli.CommandAutocomplete = (*ServerCommand)(nil)

var memProfilerEnabled = false

	return 1
}

const storageMigrationLock = "core/migration"

type ServerCommand struct {
	*BaseCommand

	logOutput   io.Writer
	gatedWriter *gatedwriter.Writer
	logger      log.Logger

	cleanupGuard sync.Once

	reloadFuncsLock *sync.RWMutex
	reloadFuncs     *map[string][]reloadutil.ReloadFunc
	startedCh       chan (struct{}) // for tests
	reloadedCh      chan (struct{}) // for tests

	// new stuff
	flagConfigs            []string
	flagTestServerConfig   bool
	flagDevConsul          bool
	flagExitOnCoreShutdown bool
}

type ServerListener struct {
	net.Listener
	config                       map[string]interface{}
	maxRequestSize               int64
	maxRequestDuration           time.Duration
	unauthenticatedMetricsAccess bool
}

func (c *ServerCommand) Synopsis() string {
			"Using a recovery operation token, \"sys/raw\" API can be used to manipulate the storage.",
	})

	f = set.NewFlagSet("Dev Options")

	f.BoolVar(&BoolVar{
	return c.Flags().Completions()
}

func (c *ServerCommand) parseConfig() (*server.Config, error) {
	// Load the configuration
	var config *server.Config
	for _, path := range c.flagConfigs {
		current, err := server.LoadConfig(path)
		if err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("error loading configuration from %s: {{err}}", path), err)
		}

		if config == nil {
			config = current
		} else {
			config = config.Merge(current)
		}
	}
	return config, nil
}

func (c *ServerCommand) runRecoveryMode() int {
	config, err := c.parseConfig()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
		return 1
	}

	c.logger = log.New(&log.LoggerOptions{
		Output: c.gatedWriter,
		Level:  level,
		// Note that if logFormat is either unspecified or standard, then
		JSONFormat: logFormat == logging.JSONFormat,
	})

	logLevelStr, err := c.adjustLogLevel(config, logLevelWasNotSet)
	if err != nil {
		c.UI.Error(err.Error())
		vault.DefaultMaxRequestDuration = config.DefaultMaxRequestDuration
	}

	proxyCfg := httpproxy.FromEnvironment()
	c.logger.Info("proxy environment", "http_proxy", proxyCfg.HTTPProxy,
		"https_proxy", proxyCfg.HTTPSProxy, "no_proxy", proxyCfg.NoProxy)

	// Initialize the storage backend
	factory, exists := c.PhysicalBackends[config.Storage.Type]
		c.UI.Error(fmt.Sprintf("Unknown storage type %s", config.Storage.Type))
		return 1
	}
	if config.Storage.Type == "raft" {
		if envCA := os.Getenv("VAULT_CLUSTER_ADDR"); envCA != "" {
			config.ClusterAddr = envCA
		}

	var barrierSeal vault.Seal
	var sealConfigError error

	if len(config.Seals) == 0 {
		config.Seals = append(config.Seals, &server.Seal{Type: wrapping.Shamir})
	}

	if len(config.Seals) > 1 {
		sealType = configSeal.Type
	}

	var seal vault.Seal
	sealLogger := c.logger.Named(sealType)
	seal, sealConfigError = serverseal.ConfigureSeal(configSeal, &infoKeys, &info, sealLogger, vault.NewDefaultSeal(&vaultseal.Access{
		Wrapper: aeadwrapper.NewWrapper(&wrapping.WrapperOptions{
			Logger: c.logger.Named("shamir"),
		}),
	}))
	if sealConfigError != nil {
		if !errwrap.ContainsType(sealConfigError, new(logical.KeyNotFoundError)) {
			c.UI.Error(fmt.Sprintf(
			return 1
		}
	}
	if seal == nil {
		c.UI.Error(fmt.Sprintf(
			"After configuring seal nil returned, seal type was %s", sealType))
		return 1
	}

	barrierSeal = seal

	// Ensure that the seal finalizer is called, even if using verify-only
	}

	// Initialize the listeners
	lns := make([]ServerListener, 0, len(config.Listeners))
	for _, lnConfig := range config.Listeners {
		ln, _, _, err := server.NewListener(lnConfig.Type, lnConfig.Config, c.gatedWriter, c.UI)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing listener of type %s: %s", lnConfig.Type, err))
			return 1
		}

		lns = append(lns, ServerListener{
			Listener: ln,
			config:   lnConfig.Config,
		})
	}

	infoKeys = append(infoKeys, "version")
	verInfo := version.GetVersion()
	info["version"] = verInfo.FullVersionNumber(false)
	if verInfo.Revision != "" {
		info["version sha"] = strings.Trim(verInfo.Revision, "'")
		infoKeys = append(infoKeys, "version sha")
	infoKeys = append(infoKeys, "recovery mode")
	info["recovery mode"] = "true"

	// Server configuration output
	padding := 24
	sort.Strings(infoKeys)
	c.UI.Output("==> Vault server configuration:\n")
	for _, k := range infoKeys {
		c.UI.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Title(k),
			info[k]))
	}
	c.UI.Output("")

	for _, ln := range lns {
		handler := vaulthttp.Handler(&vault.HandlerProperties{
			Core:                  core,
			MaxRequestSize:        ln.maxRequestSize,
			MaxRequestDuration:    ln.maxRequestDuration,
			DisablePrintableCheck: config.DisablePrintableCheck,
			RecoveryMode:          c.flagRecovery,
			RecoveryToken:         atomic.NewString(""),
	}

	if sealConfigError != nil {
		init, err := core.Initialized(context.Background())
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error checking if core is initialized: %v", err))
			return 1
		c.UI.Output("==> Vault server started! Log data will stream in below:\n")
	}

	c.logger.(hclog.OutputResettable).ResetOutputWithFlush(&hclog.LoggerOptions{
		Output: c.logOutput,
	}, c.gatedWriter)

	for {
		select {
	return 0
}

func (c *ServerCommand) adjustLogLevel(config *server.Config, logLevelWasNotSet bool) (string, error) {
	var logLevelString string
	if config.LogLevel != "" && logLevelWasNotSet {
	return level, logLevelString, logLevelWasNotSet, logFormat, nil
}

func (c *ServerCommand) Run(args []string) int {
	f := c.Flags()

		}
	}

	// Load the configuration
	var config *server.Config
	if c.flagDev {
		var devStorageType string
		switch {
		default:
			devStorageType = "inmem"
		}
		config = server.DevConfig(devStorageType)
		if c.flagDevListenAddr != "" {
			config.Listeners[0].Config["address"] = c.flagDevListenAddr
		}
	}

	parsedConfig, err := c.parseConfig()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
		return 1
	}

	if c.flagDevThreeNode || c.flagDevFourCluster {
		c.logger = log.New(&log.LoggerOptions{
			Mutex:  &sync.Mutex{},
			Output: c.gatedWriter,
			Level:  log.Trace,
		})
	} else {
		c.logger = log.New(&log.LoggerOptions{
			Output: c.gatedWriter,
			Level:  level,
			// Note that if logFormat is either unspecified or standard, then
		})
	}

	allLoggers := []log.Logger{c.logger}

	logLevelStr, err := c.adjustLogLevel(config, logLevelWasNotSet)
	if err != nil {

	// create GRPC logger
	namedGRPCLogFaker := c.logger.Named("grpclogfaker")
	allLoggers = append(allLoggers, namedGRPCLogFaker)
	grpclog.SetLogger(&grpclogFaker{
		logger: namedGRPCLogFaker,
		log:    os.Getenv("VAULT_GRPC_LOGGING") != "",
		c.startMemProfiler()
	}

	// Ensure that a backend is provided
	if config.Storage == nil {
		c.UI.Output("A storage backend must be specified")
		return 1
	}

	if config.DefaultMaxRequestDuration != 0 {
		vault.DefaultMaxRequestDuration = config.DefaultMaxRequestDuration
	}

	// log proxy settings
	proxyCfg := httpproxy.FromEnvironment()
	c.logger.Info("proxy environment", "http_proxy", proxyCfg.HTTPProxy,
		"https_proxy", proxyCfg.HTTPSProxy, "no_proxy", proxyCfg.NoProxy)

	// If mlockall(2) isn't supported, show a warning. We disable this in dev
	// because it is quite scary to see when first using Vault. We also disable
				"in a Docker container, provide the IPC_LOCK cap to the container."))
	}

	metricsHelper, err := c.setupTelemetry(config)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}

	// Initialize the backend
	factory, exists := c.PhysicalBackends[config.Storage.Type]
	if !exists {
		c.UI.Error(fmt.Sprintf("Unknown storage type %s", config.Storage.Type))
		return 1
	}

	// Do any custom configuration needed per backend
	switch config.Storage.Type {
	case "consul":
		if config.ServiceRegistration == nil {
			// If Consul is configured for storage and service registration is unconfigured,
			// use Consul for service registration without requiring additional configuration.
			// This maintains backward-compatibility.
			config.ServiceRegistration = &server.ServiceRegistration{
				Type:   "consul",
				Config: config.Storage.Config,
			}
		}
	case "raft":
		if envCA := os.Getenv("VAULT_CLUSTER_ADDR"); envCA != "" {
			config.ClusterAddr = envCA
		}
		if len(config.ClusterAddr) == 0 {
			c.UI.Error("Cluster address must be set when using raft storage")
			return 1
		}
	}

	namedStorageLogger := c.logger.Named("storage." + config.Storage.Type)
	allLoggers = append(allLoggers, namedStorageLogger)
	backend, err := factory(config.Storage.Config, namedStorageLogger)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing storage of type %s: %s", config.Storage.Type, err))
		return 1
	}

	// Prevent server startup if migration is active
	if c.storageMigrationActive(backend) {
		return 1
	}

	// Instantiate the wait group
	c.WaitGroup = &sync.WaitGroup{}

	// Initialize the Service Discovery, if there is one
	var configSR sr.ServiceRegistration
	if config.ServiceRegistration != nil {
		sdFactory, ok := c.ServiceRegistrations[config.ServiceRegistration.Type]
		if !ok {
			c.UI.Error(fmt.Sprintf("Unknown service_registration type %s", config.ServiceRegistration.Type))
			return 1
		}

		namedSDLogger := c.logger.Named("service_registration." + config.ServiceRegistration.Type)
		allLoggers = append(allLoggers, namedSDLogger)

		// Since we haven't even begun starting Vault's core yet,
		// we know that Vault is in its pre-running state.
		state := sr.State{
			VaultVersion:         version.GetVersion().VersionNumber(),
			IsInitialized:        false,
			IsSealed:             true,
			IsActive:             false,
			IsPerformanceStandby: false,
		}
		configSR, err = sdFactory(config.ServiceRegistration.Config, namedSDLogger, state, config.Storage.RedirectAddr)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing service_registration of type %s: %s", config.ServiceRegistration.Type, err))
			return 1
		}
		if err := configSR.Run(c.ShutdownCh, c.WaitGroup); err != nil {
			c.UI.Error(fmt.Sprintf("Error running service_registration of type %s: %s", config.ServiceRegistration.Type, err))
			return 1
		}
	}
	info := make(map[string]string)
	info["log level"] = logLevelString
	infoKeys = append(infoKeys, "log level")

	var barrierSeal vault.Seal
	var unwrapSeal vault.Seal

	var sealConfigError error
	if c.flagDevAutoSeal {
		barrierSeal = vault.NewAutoSeal(vaultseal.NewTestSeal(nil))
	} else {
		// Handle the case where no seal is provided
		switch len(config.Seals) {
		case 0:
			config.Seals = append(config.Seals, &server.Seal{Type: wrapping.Shamir})
		case 1:
			// If there's only one seal and it's disabled assume they want to
			// migrate to a shamir seal and simply didn't provide it
			if config.Seals[0].Disabled {
				config.Seals = append(config.Seals, &server.Seal{Type: wrapping.Shamir})
			}
		}
		for _, configSeal := range config.Seals {
			sealType := wrapping.Shamir
			if !configSeal.Disabled && os.Getenv("VAULT_SEAL_TYPE") != "" {
				sealType = os.Getenv("VAULT_SEAL_TYPE")
				configSeal.Type = sealType
			} else {
				sealType = configSeal.Type
			}

			var seal vault.Seal
			sealLogger := c.logger.Named(sealType)
			allLoggers = append(allLoggers, sealLogger)
			seal, sealConfigError = serverseal.ConfigureSeal(configSeal, &infoKeys, &info, sealLogger, vault.NewDefaultSeal(&vaultseal.Access{
				Wrapper: aeadwrapper.NewWrapper(&wrapping.WrapperOptions{
					Logger: c.logger.Named("shamir"),
				}),
			}))
			if sealConfigError != nil {
				if !errwrap.ContainsType(sealConfigError, new(logical.KeyNotFoundError)) {
					c.UI.Error(fmt.Sprintf(
						"Error parsing Seal configuration: %s", sealConfigError))
					return 1
				}
			}
			if seal == nil {
				c.UI.Error(fmt.Sprintf(
					"After configuring seal nil returned, seal type was %s", sealType))
				return 1
			}

			if configSeal.Disabled {
				unwrapSeal = seal
			} else {
				barrierSeal = seal
			}

			// Ensure that the seal finalizer is called, even if using verify-only
			defer func() {
				err = seal.Finalize(context.Background())
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error finalizing seals: %v", err))
				}
			}()

		}
	}

	}

	// prepare a secure random reader for core
	secureRandomReader, err := createSecureRandomReaderFunc(config, &barrierSeal)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	coreConfig := &vault.CoreConfig{
		RawConfig:                 config,
		Physical:                  backend,
		RedirectAddr:              config.Storage.RedirectAddr,
		StorageType:               config.Storage.Type,
		HAPhysical:                nil,
		ServiceRegistration:       configSR,
		Seal:                      barrierSeal,
		AuditBackends:             c.AuditBackends,
		CredentialBackends:        c.CredentialBackends,
		LogicalBackends:           c.LogicalBackends,
		Logger:                    c.logger,
		DisableCache:              config.DisableCache,
		DisableMlock:              config.DisableMlock,
		MaxLeaseTTL:               config.MaxLeaseTTL,
		DefaultLeaseTTL:           config.DefaultLeaseTTL,
		ClusterName:               config.ClusterName,
		CacheSize:                 config.CacheSize,
		PluginDirectory:           config.PluginDirectory,
		EnableUI:                  config.EnableUI,
		EnableRaw:                 config.EnableRawEndpoint,
		DisableSealWrap:           config.DisableSealWrap,
		DisablePerformanceStandby: config.DisablePerformanceStandby,
		DisableIndexing:           config.DisableIndexing,
		AllLoggers:                allLoggers,
		BuiltinRegistry:           builtinplugins.Registry,
		DisableKeyEncodingChecks:  config.DisablePrintableCheck,
		MetricsHelper:             metricsHelper,
		SecureRandomReader:        secureRandomReader,
	}
	if c.flagDev {
		coreConfig.DevToken = c.flagDevRootTokenID
		if c.flagDevLeasedKV {
			coreConfig.LogicalBackends["kv"] = vault.LeasedPassthroughBackendFactory
		}
		if c.flagDevPluginDir != "" {
			coreConfig.PluginDirectory = c.flagDevPluginDir
		}
		if c.flagDevLatency > 0 {
			injectLatency := time.Duration(c.flagDevLatency) * time.Millisecond
			if _, txnOK := backend.(physical.Transactional); txnOK {
				coreConfig.Physical = physical.NewTransactionalLatencyInjector(backend, injectLatency, c.flagDevLatencyJitter, c.logger)
			} else {
				coreConfig.Physical = physical.NewLatencyInjector(backend, injectLatency, c.flagDevLatencyJitter, c.logger)
			}
		}
	}

	if c.flagDevThreeNode {
		return c.enableThreeNodeDevCluster(coreConfig, info, infoKeys, c.flagDevListenAddr, os.Getenv("VAULT_DEV_TEMP_DIR"))
	}

	if c.flagDevFourCluster {
		return enableFourClusterDev(c, coreConfig, info, infoKeys, c.flagDevListenAddr, os.Getenv("VAULT_DEV_TEMP_DIR"))
	}

	var disableClustering bool

	// Initialize the separate HA storage backend, if it exists
	var ok bool
	if config.HAStorage != nil {
		// TODO: Remove when Raft can server as the ha_storage backend.
		// See https://github.com/hashicorp/vault/issues/8206
		if config.HAStorage.Type == "raft" {
			c.UI.Error("Raft cannot be used as seperate HA storage at this time")
			return 1
		}
		factory, exists := c.PhysicalBackends[config.HAStorage.Type]
		if !exists {
			c.UI.Error(fmt.Sprintf("Unknown HA storage type %s", config.HAStorage.Type))
			return 1

		}
		habackend, err := factory(config.HAStorage.Config, c.logger)
		if err != nil {
			c.UI.Error(fmt.Sprintf(
				"Error initializing HA storage of type %s: %s", config.HAStorage.Type, err))
			return 1

		}

		if coreConfig.HAPhysical, ok = habackend.(physical.HABackend); !ok {
			c.UI.Error("Specified HA storage does not support HA")
			return 1
		}

		if !coreConfig.HAPhysical.HAEnabled() {
			c.UI.Error("Specified HA storage has HA support disabled; please consult documentation")
			return 1
		}

		coreConfig.RedirectAddr = config.HAStorage.RedirectAddr
		disableClustering = config.HAStorage.DisableClustering
		if !disableClustering {
			coreConfig.ClusterAddr = config.HAStorage.ClusterAddr
		}
	} else {
		if coreConfig.HAPhysical, ok = backend.(physical.HABackend); ok {
			coreConfig.RedirectAddr = config.Storage.RedirectAddr
			disableClustering = config.Storage.DisableClustering
			if !disableClustering {
				coreConfig.ClusterAddr = config.Storage.ClusterAddr
			}
		}
	}

	if envRA := os.Getenv("VAULT_API_ADDR"); envRA != "" {
		coreConfig.RedirectAddr = envRA
	} else if envRA := os.Getenv("VAULT_REDIRECT_ADDR"); envRA != "" {
		coreConfig.RedirectAddr = envRA
	} else if envAA := os.Getenv("VAULT_ADVERTISE_ADDR"); envAA != "" {
		coreConfig.RedirectAddr = envAA
	}

	// Attempt to detect the redirect address, if possible
	if coreConfig.RedirectAddr == "" {
		c.logger.Warn("no `api_addr` value specified in config or in VAULT_API_ADDR; falling back to detection if possible, but this value should be manually set")
	}
	var detect physical.RedirectDetect
	if coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled() {
		detect, ok = coreConfig.HAPhysical.(physical.RedirectDetect)
	} else {
		detect, ok = coreConfig.Physical.(physical.RedirectDetect)
	}
	if ok && coreConfig.RedirectAddr == "" {
		redirect, err := c.detectRedirect(detect, config)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error detecting api address: %s", err))
		} else if redirect == "" {
			c.UI.Error("Failed to detect api address")
		} else {
			coreConfig.RedirectAddr = redirect
		}
	}
	if coreConfig.RedirectAddr == "" && c.flagDev {
		coreConfig.RedirectAddr = fmt.Sprintf("http://%s", config.Listeners[0].Config["address"])
	}

	// After the redirect bits are sorted out, if no cluster address was
	// explicitly given, derive one from the redirect addr
	if disableClustering {
		coreConfig.ClusterAddr = ""
	} else if envCA := os.Getenv("VAULT_CLUSTER_ADDR"); envCA != "" {
		coreConfig.ClusterAddr = envCA
	} else {
		var addrToUse string
		switch {
		case coreConfig.ClusterAddr == "" && coreConfig.RedirectAddr != "":
			addrToUse = coreConfig.RedirectAddr
		case c.flagDev:
			addrToUse = fmt.Sprintf("http://%s", config.Listeners[0].Config["address"])
		default:
			goto CLUSTER_SYNTHESIS_COMPLETE
		}
		u, err := url.ParseRequestURI(addrToUse)
		if err != nil {
			c.UI.Error(fmt.Sprintf(
				"Error parsing synthesized cluster address %s: %v", addrToUse, err))
			return 1
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			// This sucks, as it's a const in the function but not exported in the package
			if strings.Contains(err.Error(), "missing port in address") {
				host = u.Host
				port = "443"
			} else {
				c.UI.Error(fmt.Sprintf("Error parsing api address: %v", err))
				return 1
			}
		}
		nPort, err := strconv.Atoi(port)
		if err != nil {
			c.UI.Error(fmt.Sprintf(
				"Error parsing synthesized address; failed to convert %q to a numeric: %v", port, err))
			return 1
		}
		u.Host = net.JoinHostPort(host, strconv.Itoa(nPort+1))
		// Will always be TLS-secured
		u.Scheme = "https"
		coreConfig.ClusterAddr = u.String()
	}

CLUSTER_SYNTHESIS_COMPLETE:

	if coreConfig.RedirectAddr == coreConfig.ClusterAddr && len(coreConfig.RedirectAddr) != 0 {
		c.UI.Error(fmt.Sprintf(
			"Address %q used for both API and cluster addresses", coreConfig.RedirectAddr))
		return 1
	}

	if coreConfig.ClusterAddr != "" {
		// Force https as we'll always be TLS-secured
		u, err := url.ParseRequestURI(coreConfig.ClusterAddr)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing cluster address %s: %v", coreConfig.ClusterAddr, err))
			return 11
		}
		u.Scheme = "https"
		coreConfig.ClusterAddr = u.String()
	}

	// Override the UI enabling config by the environment variable
	if enableUI := os.Getenv("VAULT_UI"); enableUI != "" {
		var err error

	// If ServiceRegistration is configured, then the backend must support HA
	isBackendHA := coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled()
	if !c.flagDev && (coreConfig.ServiceRegistration != nil) && !isBackendHA {
		c.UI.Output("service_registration is configured, but storage does not support HA")
		return 1
	}

	// Apply any enterprise configuration onto the coreConfig.
	adjustCoreConfigForEnt(config, coreConfig)

	// Initialize the core
	core, newCoreError := vault.NewCore(coreConfig)
	if newCoreError != nil {
		if vault.IsFatalError(newCoreError) {
			c.UI.Error(fmt.Sprintf("Error initializing core: %s", newCoreError))
			return 1
		}
	}

	// Copy the reload funcs pointers back
		}
	}

	clusterAddrs := []*net.TCPAddr{}

	// Initialize the listeners
	lns := make([]ServerListener, 0, len(config.Listeners))
	c.reloadFuncsLock.Lock()
	for i, lnConfig := range config.Listeners {
		ln, props, reloadFunc, err := server.NewListener(lnConfig.Type, lnConfig.Config, c.gatedWriter, c.UI)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing listener of type %s: %s", lnConfig.Type, err))
			return 1
		}

		if reloadFunc != nil {
			relSlice := (*c.reloadFuncs)["listener|"+lnConfig.Type]
			relSlice = append(relSlice, reloadFunc)
			(*c.reloadFuncs)["listener|"+lnConfig.Type] = relSlice
		}

		if !disableClustering && lnConfig.Type == "tcp" {
			var addrRaw interface{}
			var addr string
			var ok bool
			if addrRaw, ok = lnConfig.Config["cluster_address"]; ok {
				addr = addrRaw.(string)
				tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error resolving cluster_address: %s", err))
					return 1
				}
				clusterAddrs = append(clusterAddrs, tcpAddr)
			} else {
				tcpAddr, ok := ln.Addr().(*net.TCPAddr)
				if !ok {
					c.UI.Error("Failed to parse tcp listener")
					return 1
				}
				clusterAddr := &net.TCPAddr{
					IP:   tcpAddr.IP,
					Port: tcpAddr.Port + 1,
				}
				clusterAddrs = append(clusterAddrs, clusterAddr)
				addr = clusterAddr.String()
			}
			props["cluster address"] = addr
		}

		var maxRequestSize int64 = vaulthttp.DefaultMaxRequestSize
		if valRaw, ok := lnConfig.Config["max_request_size"]; ok {
			val, err := parseutil.ParseInt(valRaw)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Could not parse max_request_size value %v", valRaw))
				return 1
			}

			if val >= 0 {
				maxRequestSize = val
			}
		}
		props["max_request_size"] = fmt.Sprintf("%d", maxRequestSize)

		maxRequestDuration := vault.DefaultMaxRequestDuration
		if valRaw, ok := lnConfig.Config["max_request_duration"]; ok {
			val, err := parseutil.ParseDurationSecond(valRaw)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Could not parse max_request_duration value %v", valRaw))
				return 1
			}

			if val >= 0 {
				maxRequestDuration = val
			}
		}
		props["max_request_duration"] = fmt.Sprintf("%s", maxRequestDuration.String())

		var unauthenticatedMetricsAccess bool
		if telemetryRaw, ok := lnConfig.Config["telemetry"]; ok {
			telemetry, ok := telemetryRaw.([]map[string]interface{})
			if !ok {
				c.UI.Error(fmt.Sprintf("Could not parse telemetry sink value %v", telemetryRaw))
				return 1
			}

			for _, item := range telemetry {
				if valRaw, ok := item["unauthenticated_metrics_access"]; ok {
					unauthenticatedMetricsAccess, err = parseutil.ParseBool(valRaw)
					if err != nil {
						c.UI.Error(fmt.Sprintf("Could not parse unauthenticated_metrics_access value %v", valRaw))
						return 1
					}
				}
			}
		}

		lns = append(lns, ServerListener{
			Listener:                     ln,
			config:                       lnConfig.Config,
			maxRequestSize:               maxRequestSize,
			maxRequestDuration:           maxRequestDuration,
			unauthenticatedMetricsAccess: unauthenticatedMetricsAccess,
		})

		// Store the listener props for output later
		key := fmt.Sprintf("listener %d", i+1)
		propsList := make([]string, 0, len(props))
		for k, v := range props {
			propsList = append(propsList, fmt.Sprintf(
				"%s: %q", k, v))
		}
		sort.Strings(propsList)
		infoKeys = append(infoKeys, key)
		info[key] = fmt.Sprintf(
			"%s (%s)", lnConfig.Type, strings.Join(propsList, ", "))

	}
	c.reloadFuncsLock.Unlock()
	if !disableClustering {
		if c.logger.IsDebug() {
			c.logger.Debug("cluster listener addresses synthesized", "cluster_addresses", clusterAddrs)
		}
	}

	// Make sure we close all listeners from this point on
		info["version sha"] = strings.Trim(verInfo.Revision, "'")
		infoKeys = append(infoKeys, "version sha")
	}
	infoKeys = append(infoKeys, "cgo")
	info["cgo"] = "disabled"
	if version.CgoEnabled {
	infoKeys = append(infoKeys, "recovery mode")
	info["recovery mode"] = "false"

	// Server configuration output
	padding := 24
	sort.Strings(infoKeys)
	c.UI.Output("==> Vault server configuration:\n")
	for _, k := range infoKeys {
		c.UI.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}
	c.UI.Output("")

	// Tests might not want to start a vault server and just want to verify
		Core: core,
	}))

	// Before unsealing with stored keys, setup seal migration if needed
	if err := adjustCoreForSealMigration(c.logger, core, barrierSeal, unwrapSeal); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Attempt unsealing in a background goroutine. This is needed for when a
	// Vault cluster with multiple servers is configured with auto-unseal but is
	// uninitialized. Once one server initializes the storage backend, this
	// goroutine will pick up the unseal keys and unseal this instance.
	if !core.IsInSealMigration() {
		go func() {
			for {
				err := core.UnsealWithStoredKeys(context.Background())
				if err == nil {
					return
				}

				if vault.IsFatalError(err) {
					c.logger.Error("error unsealing core", "error", err)
					return
				} else {
					c.logger.Warn("failed to unseal core", "error", err)
				}

				select {
				case <-c.ShutdownCh:
					return
				case <-time.After(5 * time.Second):
				}
			}
		}()
	}

	// When the underlying storage is raft, kick off retry join if it was specified
	// in the configuration
	if config.Storage.Type == "raft" {
		if err := core.InitiateRetryJoin(context.Background()); err != nil {
			c.UI.Error(fmt.Sprintf("Failed to initiate raft retry join, %q", err.Error()))
			return 1
	}

	// Perform initialization of HTTP server after the verifyOnly check.
	// If we're in Dev mode, then initialize the core
	if c.flagDev && !c.flagDevSkipInit {
		init, err := c.enableDev(core, coreConfig)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing Dev mode: %s", err))
			return 1
		}

		var plugins, pluginsNotLoaded []string
		if c.flagDevPluginDir != "" && c.flagDevPluginInit {

			f, err := os.Open(c.flagDevPluginDir)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error reading plugin dir: %s", err))
				return 1
			}

			list, err := f.Readdirnames(0)
			f.Close()
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error listing plugins: %s", err))
				return 1
			}

			for _, name := range list {
				path := filepath.Join(f.Name(), name)
				if err := c.addPlugin(path, init.RootToken, core); err != nil {
					if !errwrap.Contains(err, vault.ErrPluginBadType.Error()) {
						c.UI.Error(fmt.Sprintf("Error enabling plugin %s: %s", name, err))
						return 1
					}
					pluginsNotLoaded = append(pluginsNotLoaded, name)
					continue
				}
				plugins = append(plugins, name)
			}

			sort.Strings(plugins)
		}

		// Print the big dev mode warning!
		c.UI.Warn(wrapAtLength(
			"WARNING! dev mode is enabled! In this mode, Vault runs entirely " +
				"in-memory and starts unsealed with a single unseal key. The root " +
				"token is already authenticated to the CLI, so you can immediately " +
				"begin using Vault."))
		c.UI.Warn("")
		c.UI.Warn("You may need to set the following environment variable:")
		c.UI.Warn("")

		endpointURL := "http://" + config.Listeners[0].Config["address"].(string)
		if runtime.GOOS == "windows" {
			c.UI.Warn("PowerShell:")
			c.UI.Warn(fmt.Sprintf("    $env:VAULT_ADDR=\"%s\"", endpointURL))
			c.UI.Warn("cmd.exe:")
			c.UI.Warn(fmt.Sprintf("    set VAULT_ADDR=%s", endpointURL))
		} else {
			c.UI.Warn(fmt.Sprintf("    $ export VAULT_ADDR='%s'", endpointURL))
		}

		// Unseal key is not returned if stored shares is supported
		if len(init.SecretShares) > 0 {
			c.UI.Warn("")
			c.UI.Warn(wrapAtLength(
				"The unseal key and root token are displayed below in case you want " +
					"to seal/unseal the Vault or re-authenticate."))
			c.UI.Warn("")
			c.UI.Warn(fmt.Sprintf("Unseal Key: %s", base64.StdEncoding.EncodeToString(init.SecretShares[0])))
		}

		if len(init.RecoveryShares) > 0 {
			c.UI.Warn("")
			c.UI.Warn(wrapAtLength(
				"The recovery key and root token are displayed below in case you want " +
					"to seal/unseal the Vault or re-authenticate."))
			c.UI.Warn("")
			c.UI.Warn(fmt.Sprintf("Recovery Key: %s", base64.StdEncoding.EncodeToString(init.RecoveryShares[0])))
		}

		c.UI.Warn(fmt.Sprintf("Root Token: %s", init.RootToken))

		if len(plugins) > 0 {
			c.UI.Warn("")
			c.UI.Warn(wrapAtLength(
				"The following dev plugins are registered in the catalog:"))
			for _, p := range plugins {
				c.UI.Warn(fmt.Sprintf("    - %s", p))
			}
		}

		if len(pluginsNotLoaded) > 0 {
			c.UI.Warn("")
			c.UI.Warn(wrapAtLength(
				"The following dev plugins FAILED to be registered in the catalog due to unknown type:"))
			for _, p := range pluginsNotLoaded {
				c.UI.Warn(fmt.Sprintf("    - %s", p))
			}
		}

		c.UI.Warn("")
		c.UI.Warn(wrapAtLength(
			"Development mode should NOT be used in production installations!"))
		c.UI.Warn("")
	}

	// Initialize the HTTP servers
	for _, ln := range lns {
		handler := vaulthttp.Handler(&vault.HandlerProperties{
			Core:                         core,
			MaxRequestSize:               ln.maxRequestSize,
			MaxRequestDuration:           ln.maxRequestDuration,
			DisablePrintableCheck:        config.DisablePrintableCheck,
			UnauthenticatedMetricsAccess: ln.unauthenticatedMetricsAccess,
			RecoveryMode:                 c.flagRecovery,
		})

		// We perform validation on the config earlier, we can just cast here
		if _, ok := ln.config["x_forwarded_for_authorized_addrs"]; ok {
			hopSkips := ln.config["x_forwarded_for_hop_skips"].(int)
			authzdAddrs := ln.config["x_forwarded_for_authorized_addrs"].([]*sockaddr.SockAddrMarshaler)
			rejectNotPresent := ln.config["x_forwarded_for_reject_not_present"].(bool)
			rejectNonAuthz := ln.config["x_forwarded_for_reject_not_authorized"].(bool)
			if len(authzdAddrs) > 0 {
				handler = vaulthttp.WrapForwardedForHandler(handler, authzdAddrs, rejectNotPresent, rejectNonAuthz, hopSkips)
			}
		}

		// server defaults
		server := &http.Server{
			Handler:           handler,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			IdleTimeout:       5 * time.Minute,
			ErrorLog:          c.logger.StandardLogger(nil),
		}

		// override server defaults with config values for read/write/idle timeouts if configured
		if readHeaderTimeoutInterface, ok := ln.config["http_read_header_timeout"]; ok {
			readHeaderTimeout, err := parseutil.ParseDurationSecond(readHeaderTimeoutInterface)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Could not parse a time value for http_read_header_timeout %v", readHeaderTimeout))
				return 1
			}
			server.ReadHeaderTimeout = readHeaderTimeout
		}

		if readTimeoutInterface, ok := ln.config["http_read_timeout"]; ok {
			readTimeout, err := parseutil.ParseDurationSecond(readTimeoutInterface)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Could not parse a time value for http_read_timeout %v", readTimeout))
				return 1
			}
			server.ReadTimeout = readTimeout
		}

		if writeTimeoutInterface, ok := ln.config["http_write_timeout"]; ok {
			writeTimeout, err := parseutil.ParseDurationSecond(writeTimeoutInterface)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Could not parse a time value for http_write_timeout %v", writeTimeout))
				return 1
			}
			server.WriteTimeout = writeTimeout
		}

		if idleTimeoutInterface, ok := ln.config["http_idle_timeout"]; ok {
			idleTimeout, err := parseutil.ParseDurationSecond(idleTimeoutInterface)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Could not parse a time value for http_idle_timeout %v", idleTimeout))
				return 1
			}
			server.IdleTimeout = idleTimeout
		}

		// server config tests can exit now
		if c.flagTestServerConfig {
			continue
		}

		go server.Serve(ln.Listener)
	}

	if c.flagTestServerConfig {
		return 0
	}

	if sealConfigError != nil {
		init, err := core.Initialized(context.Background())
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error checking if core is initialized: %v", err))
			return 1
		}
		if init {
			c.UI.Error("Vault is initialized but no Seal key could be loaded")
			return 1
		}
	}

	if newCoreError != nil {
		c.UI.Warn(wrapAtLength(
			"WARNING! A non-fatal error occurred during initialization. Please " +
				"check the logs for more information."))
		c.UI.Warn("")
	}

	// Output the header that the server has started
	if !c.flagCombineLogs {
	}

	// Release the log gate.
	c.logger.(hclog.OutputResettable).ResetOutputWithFlush(&hclog.LoggerOptions{
		Output: c.logOutput,
	}, c.gatedWriter)

	// Write out the PID to the file now that server has successfully started
	if err := c.storePidFile(config.PidFile); err != nil {
		return 1
	}

	defer func() {
		if err := c.removePidFile(config.PidFile); err != nil {
			c.UI.Error(fmt.Sprintf("Error deleting the PID file: %s", err))
		case <-c.SighupCh:
			c.UI.Output("==> Vault reload triggered")

			// Check for new log level
			var config *server.Config
			var level log.Level
				c.UI.Error(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
			}

		case <-c.SigUSR2Ch:
			buf := make([]byte, 32*1024*1024)
			n := runtime.Stack(buf[:], true)
			c.logger.Info("goroutine trace", "stack", string(buf[:n]))
		}
	}

	// Stop the listeners so that we don't process further client requests.
	c.cleanupGuard.Do(listenerCloseFunc)

	// Shutdown will wait until after Vault is sealed, which means the
	// request forwarding listeners will also be closed (and also
	// waited for).
	if err := core.Shutdown(); err != nil {
	return retCode
}

func (c *ServerCommand) enableDev(core *vault.Core, coreConfig *vault.CoreConfig) (*vault.InitResult, error) {
	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)


	isLeader, _, _, err := core.Leader()
	if err != nil && err != vault.ErrHANotEnabled {
		return nil, errwrap.Wrapf("failed to check active status: {{err}}", err)
	}
	if err == nil {
		leaderCount := 5
			time.Sleep(1 * time.Second)
			isLeader, _, _, err = core.Leader()
			if err != nil {
				return nil, errwrap.Wrapf("failed to check active status: {{err}}", err)
			}
			leaderCount--
		}
		}
		resp, err := core.HandleRequest(ctx, req)
		if err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("failed to create root token with ID %q: {{err}}", coreConfig.DevToken), err)
		}
		if resp == nil {
			return nil, fmt.Errorf("nil response when creating root token with ID %q", coreConfig.DevToken)
		req.Data = nil
		resp, err = core.HandleRequest(ctx, req)
		if err != nil {
			return nil, errwrap.Wrapf("failed to revoke initial root token: {{err}}", err)
		}
	}

	}
	resp, err := core.HandleRequest(ctx, req)
	if err != nil {
		return nil, errwrap.Wrapf("error creating default K/V store: {{err}}", err)
	}
	if resp.IsError() {
		return nil, errwrap.Wrapf("failed to create default K/V store: {{err}}", resp.Error())
	}

	return init, nil
		info["version sha"] = strings.Trim(verInfo.Revision, "'")
		infoKeys = append(infoKeys, "version sha")
	}
	infoKeys = append(infoKeys, "cgo")
	info["cgo"] = "disabled"
	if version.CgoEnabled {
		info["cgo"] = "enabled"
	}

	// Server configuration output
	padding := 24
	sort.Strings(infoKeys)
	c.UI.Output("==> Vault server configuration:\n")
	for _, k := range infoKeys {
		c.UI.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Title(k),
			info[k]))
	}
	c.UI.Output("")

	for _, core := range testCluster.Cores {
		return 1
	}

	if err := ioutil.WriteFile(filepath.Join(testCluster.TempDir, "root_token"), []byte(testCluster.RootToken), 0755); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing token to tempfile: %s", err))
		return 1
	}
	}

	// Release the log gate.
	c.logger.(hclog.OutputResettable).ResetOutputWithFlush(&hclog.LoggerOptions{
		Output: c.logOutput,
	}, c.gatedWriter)

	// Wait for shutdown
	shutdownTriggered := false
			// Stop the listeners so that we don't process further client requests.
			c.cleanupGuard.Do(testCluster.Cleanup)

			// Shutdown will wait until after Vault is sealed, which means the
			// request forwarding listeners will also be closed (and also
			// waited for).
			for _, core := range testCluster.Cores {
		}

		// Check if TLS is disabled
		if val, ok := list.Config["tls_disable"]; ok {
			disable, err := parseutil.ParseBool(val)
			if err != nil {
				return "", errwrap.Wrapf("tls_disable: {{err}}", err)
			}

			if disable {
				scheme = "http"
			}
		}

		// Check for address override
		var addr string
		addrRaw, ok := list.Config["address"]
		if !ok {
			addr = "127.0.0.1:8200"
		} else {
			addr = addrRaw.(string)
		}

		// Check for localhost
	return url.String(), nil
}

// setupTelemetry is used to setup the telemetry sub-systems and returns the in-memory sink to be used in http configuration
func (c *ServerCommand) setupTelemetry(config *server.Config) (*metricsutil.MetricsHelper, error) {
	/* Setup telemetry
	Aggregate on 10 second intervals for 1 minute. Expose the
	metrics over stderr when there is a SIGUSR1 received.
	*/
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	metrics.DefaultInmemSignal(inm)

	var telConfig *server.Telemetry
	if config.Telemetry != nil {
		telConfig = config.Telemetry
	} else {
		telConfig = &server.Telemetry{}
	}

	serviceName := "vault"
	if telConfig.MetricsPrefix != "" {
		serviceName = telConfig.MetricsPrefix
	}

	metricsConf := metrics.DefaultConfig(serviceName)
	metricsConf.EnableHostname = !telConfig.DisableHostname
	metricsConf.EnableHostnameLabel = telConfig.EnableHostnameLabel

	// Configure the statsite sink
	var fanout metrics.FanoutSink
	var prometheusEnabled bool

	// Configure the Prometheus sink
	if telConfig.PrometheusRetentionTime != 0 {
		prometheusEnabled = true
		prometheusOpts := prometheus.PrometheusOpts{
			Expiration: telConfig.PrometheusRetentionTime,
		}

		sink, err := prometheus.NewPrometheusSinkFrom(prometheusOpts)
		if err != nil {
			return nil, err
		}
		fanout = append(fanout, sink)
	}

	metricHelper := metricsutil.NewMetricsHelper(inm, prometheusEnabled)

	if telConfig.StatsiteAddr != "" {
		sink, err := metrics.NewStatsiteSink(telConfig.StatsiteAddr)
		if err != nil {
			return nil, err
		}
		fanout = append(fanout, sink)
	}

	// Configure the statsd sink
	if telConfig.StatsdAddr != "" {
		sink, err := metrics.NewStatsdSink(telConfig.StatsdAddr)
		if err != nil {
			return nil, err
		}
		fanout = append(fanout, sink)
	}

	// Configure the Circonus sink
	if telConfig.CirconusAPIToken != "" || telConfig.CirconusCheckSubmissionURL != "" {
		cfg := &circonus.Config{}
		cfg.Interval = telConfig.CirconusSubmissionInterval
		cfg.CheckManager.API.TokenKey = telConfig.CirconusAPIToken
		cfg.CheckManager.API.TokenApp = telConfig.CirconusAPIApp
		cfg.CheckManager.API.URL = telConfig.CirconusAPIURL
		cfg.CheckManager.Check.SubmissionURL = telConfig.CirconusCheckSubmissionURL
		cfg.CheckManager.Check.ID = telConfig.CirconusCheckID
		cfg.CheckManager.Check.ForceMetricActivation = telConfig.CirconusCheckForceMetricActivation
		cfg.CheckManager.Check.InstanceID = telConfig.CirconusCheckInstanceID
		cfg.CheckManager.Check.SearchTag = telConfig.CirconusCheckSearchTag
		cfg.CheckManager.Check.DisplayName = telConfig.CirconusCheckDisplayName
		cfg.CheckManager.Check.Tags = telConfig.CirconusCheckTags
		cfg.CheckManager.Broker.ID = telConfig.CirconusBrokerID
		cfg.CheckManager.Broker.SelectTag = telConfig.CirconusBrokerSelectTag

		if cfg.CheckManager.API.TokenApp == "" {
			cfg.CheckManager.API.TokenApp = "vault"
		}

		if cfg.CheckManager.Check.DisplayName == "" {
			cfg.CheckManager.Check.DisplayName = "Vault"
		}

		if cfg.CheckManager.Check.SearchTag == "" {
			cfg.CheckManager.Check.SearchTag = "service:vault"
		}

		sink, err := circonus.NewCirconusSink(cfg)
		if err != nil {
			return nil, err
		}
		sink.Start()
		fanout = append(fanout, sink)
	}

	if telConfig.DogStatsDAddr != "" {
		var tags []string

		if telConfig.DogStatsDTags != nil {
			tags = telConfig.DogStatsDTags
		}

		sink, err := datadog.NewDogStatsdSink(telConfig.DogStatsDAddr, metricsConf.HostName)
		if err != nil {
			return nil, errwrap.Wrapf("failed to start DogStatsD sink: {{err}}", err)
		}
		sink.SetTags(tags)
		fanout = append(fanout, sink)
	}

	// Configure the stackdriver sink
	if telConfig.StackdriverProjectID != "" {
		client, err := monitoring.NewMetricClient(context.Background(), option.WithUserAgent(useragent.String()))
		if err != nil {
			return nil, fmt.Errorf("Failed to create stackdriver client: %v", err)
		}
		sink := stackdriver.NewSink(client, &stackdriver.Config{
			ProjectID: telConfig.StackdriverProjectID,
			Location:  telConfig.StackdriverLocation,
			Namespace: telConfig.StackdriverNamespace,
		})
		fanout = append(fanout, sink)
	}

	// Initialize the global sink
	if len(fanout) > 1 {
		// Hostname enabled will create poor quality metrics name for prometheus
		if !telConfig.DisableHostname {
			c.UI.Warn("telemetry.disable_hostname has been set to false. Recommended setting is true for Prometheus to avoid poorly named metrics.")
		}
	} else {
		metricsConf.EnableHostname = false
	}
	fanout = append(fanout, inm)
	_, err := metrics.NewGlobal(metricsConf, fanout)

	if err != nil {
		return nil, err
	}

	return metricHelper, nil
}

func (c *ServerCommand) Reload(lock *sync.RWMutex, reloadFuncs *map[string][]reloadutil.ReloadFunc, configPath []string) error {
	lock.RLock()
	defer lock.RUnlock()
		case strings.HasPrefix(k, "listener|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(nil); err != nil {
						reloadErrors = multierror.Append(reloadErrors, errwrap.Wrapf("error encountered reloading listener: {{err}}", err))
					}
				}
			}
		case strings.HasPrefix(k, "audit_file|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(nil); err != nil {
						reloadErrors = multierror.Append(reloadErrors, errwrap.Wrapf(fmt.Sprintf("error encountered reloading file audit device at path %q: {{err}}", strings.TrimPrefix(k, "audit_file|")), err))
					}
				}
			}
	}

	// Open the PID file
	pidFile, err := os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errwrap.Wrapf("could not open pid file: {{err}}", err)
	}
	defer pidFile.Close()

	pid := os.Getpid()
	_, err = pidFile.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return errwrap.Wrapf("could not write to pid file: {{err}}", err)
	}
	return nil
}
			c.UI.Warn("\nWARNING! Unable to read storage migration status.")

			// unexpected state, so stop buffering log messages
			c.logger.(hclog.OutputResettable).ResetOutputWithFlush(&hclog.LoggerOptions{
				Output: c.logOutput,
			}, c.gatedWriter)
		}
		c.logger.Warn("storage migration check error", "error", err.Error())


func CheckStorageMigration(b physical.Backend) (*StorageMigrationStatus, error) {
	entry, err := b.Get(context.Background(), storageMigrationLock)

	if err != nil {
		return nil, err
	}
	return &status, nil
}

func SetStorageMigration(b physical.Backend, active bool) error {
	if !active {
		return b.Delete(context.Background(), storageMigrationLock)

	"github.com/armon/go-metrics"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-raftchunking"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/raft"
	snapshot "github.com/hashicorp/raft-snapshot"
	raftboltdb "github.com/hashicorp/vault/physical/raft/logstore"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/vault/seal"
)

// EnvVaultRaftNodeID is used to fetch the Raft node ID from the environment.
const EnvVaultRaftPath = "VAULT_RAFT_PATH"

// Verify RaftBackend satisfies the correct interfaces
var _ physical.Backend = (*RaftBackend)(nil)
var _ physical.Transactional = (*RaftBackend)(nil)

var (
	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512

	raftState         = "raft/"
	peersFileName     = "peers.json"
	snapshotsRetained = 2

	restoreOpDelayDuration = 5 * time.Second
)

// RaftBackend implements the backend interfaces and uses the raft protocol to
	// raft is the instance of raft we will operate on.
	raft *raft.Raft

	// raftNotifyCh is used to receive updates about leadership changes
	// regarding this node.
	raftNotifyCh chan bool

	// permitPool is used to limit the number of concurrent storage calls.
	permitPool *physical.PermitPool
}

// LeaderJoinInfo contains information required by a node to join itself as a
// follower to an existing raft cluster
type LeaderJoinInfo struct {
	// LeaderAPIAddr is the address of the leader node to connect to
	LeaderAPIAddr string `json:"leader_api_addr"`

	// LeaderCACert is the CA cert of the leader node
	LeaderCACert string `json:"leader_ca_cert"`

	// LeaderClientCert is the client certificate for the follower node to establish
	// client authentication during TLS
	LeaderClientCert string `json:"leader_client_cert"`

	// LeaderClientKey is the client key for the follower node to establish client
	// authentication during TLS
	LeaderClientKey string `json:"leader_client_key"`

	// Retry indicates if the join process should automatically be retried
	Retry bool `json:"-"`

	var leaderInfos []*LeaderJoinInfo
	err := jsonutil.DecodeJSON([]byte(config), &leaderInfos)
	if err != nil {
		return nil, errwrap.Wrapf("failed to decode retry_join config: {{err}}", err)
	}

	if len(leaderInfos) == 0 {
		return nil, errors.New("invalid retry_join config")
	}

	for _, info := range leaderInfos {
		info.Retry = true
		var tlsConfig *tls.Config
		var err error
		if len(info.LeaderCACert) != 0 || len(info.LeaderClientCert) != 0 || len(info.LeaderClientKey) != 0 {
			tlsConfig, err = tlsutil.ClientTLSConfig([]byte(info.LeaderCACert), []byte(info.LeaderClientCert), []byte(info.LeaderClientKey))
			if err != nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("failed to create tls config to communicate with leader node %q: {{err}}", info.LeaderAPIAddr), err)
			}
		}
		info.TLSConfig = tlsConfig
	}

	return leaderInfos, nil
}

// EnsurePath is used to make sure a path exists
func EnsurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}
	return os.MkdirAll(path, 0755)
}

// NewRaftBackend constructs a RaftBackend using the given directory
func NewRaftBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Create the FSM.
	fsm, err := NewFSM(conf, logger.Named("fsm"))
	if err != nil {
		return nil, fmt.Errorf("failed to create fsm: %v", err)
	}

	path := os.Getenv(EnvVaultRaftPath)
	if path == "" {
		pathFromConfig, ok := conf["path"]
		path = pathFromConfig
	}

	// Build an all in-memory setup for dev mode, otherwise prepare a full
	// disk-based setup.
	var log raft.LogStore
	var stable raft.StableStore
	var snap raft.SnapshotStore
	var devMode bool
	if devMode {
		store := raft.NewInmemStore()
		}

		// Create the backend raft store for logs and stable storage.
		store, err := raftboltdb.NewBoltStore(filepath.Join(path, "raft.db"))
		if err != nil {
			return nil, err
		}
		log = cacheStore

		// Create the snapshot store.
		snapshots, err := NewBoltSnapshotStore(path, snapshotsRetained, logger.Named("snapshot"), fsm)
		if err != nil {
			return nil, err
		}
		snap = snapshots
	}

	var localID string
	{
		// Determine the local node ID from the environment.
		if raftNodeID := os.Getenv(EnvVaultRaftNodeID); raftNodeID != "" {
			localID = raftNodeID
		}

		// If not set in the environment check the configuration file.
		if len(localID) == 0 {
			localID = conf["node_id"]
		}

		// If not set in the config check the "node-id" file.
		if len(localID) == 0 {
			localIDRaw, err := ioutil.ReadFile(filepath.Join(path, "node-id"))
			switch {
			case err == nil:
				if len(localIDRaw) > 0 {
					localID = string(localIDRaw)
				}
			case os.IsNotExist(err):
			default:
				return nil, err
			}
		}

		// If all of the above fails generate a UUID and persist it to the
		// "node-id" file.
		if len(localID) == 0 {
			id, err := uuid.GenerateUUID()
			if err != nil {
				return nil, err
			}

			if err := ioutil.WriteFile(filepath.Join(path, "node-id"), []byte(id), 0600); err != nil {
				return nil, err
			}

			localID = id
		}
	}

	return &RaftBackend{
		logger:      logger,
		fsm:         fsm,
		conf:        conf,
		logStore:    log,
		stableStore: stable,
		snapStore:   snap,
		dataDir:     path,
		localID:     localID,
		permitPool:  physical.NewPermitPool(physical.DefaultParallelOperations),
	}, nil
}

// RaftServer has information about a server in the Raft configuration
type RaftServer struct {
	// NodeID is the name of the server

// Peer defines the ID and Address for a given member of the raft cluster.
type Peer struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

// NodeID returns the identifier of the node
}

// Bootstrap prepares the given peers to be part of the raft cluster
func (b *RaftBackend) Bootstrap(ctx context.Context, peers []Peer) error {
	b.l.Lock()
	defer b.l.Unlock()


	for i, p := range peers {
		raftConfig.Servers[i] = raft.Server{
			ID:      raft.ServerID(p.ID),
			Address: raft.ServerAddress(p.Address),
		}
	}

		}
		config.TrailingLogs = uint64(trailingLogs)
	}

	config.NoSnapshotRestoreOnStart = true
	config.MaxAppendEntries = 64

	return nil
}

	})
}

// SetupCluster starts the raft cluster and enables the networking needed for
// the raft nodes to communicate.
func (b *RaftBackend) SetupCluster(ctx context.Context, opts SetupOpts) error {
		return err
	}

	switch {
	case opts.TLSKeyring == nil && opts.ClusterListener == nil:
		// If we don't have a provided network we use an in-memory one.
		// This allows us to bootstrap a node without bringing up a cluster
		// network. This will be true during bootstrap, tests and dev modes.
		_, b.raftTransport = raft.NewInmemTransportWithTimeout(raft.ServerAddress(b.localID), time.Second)
	case opts.TLSKeyring == nil:
		return errors.New("no keyring provided")
	case opts.ClusterListener == nil:
		return errors.New("no cluster listener provided")
	default:
		// Set the local address and localID in the streaming layer and the raft config.
			MaxPool:               3,
			Timeout:               10 * time.Second,
			ServerAddressProvider: b.serverAddressProvider,
		}
		transport := raft.NewNetworkTransportWithConfig(transConfig)

	raftConfig.LocalID = raft.ServerID(b.localID)

	// Set up a channel for reliable leader notifications.
	raftNotifyCh := make(chan bool, 1)
	raftConfig.NotifyCh = raftNotifyCh

	// If we have a bootstrapConfig set we should bootstrap now.
		if err := raft.BootstrapCluster(raftConfig, b.logStore, b.stableStore, b.snapStore, b.raftTransport, *bootstrapConfig); err != nil {
			return err
		}
		// If we are the only node we should start as the leader.
		if len(bootstrapConfig.Servers) == 1 {
			opts.StartAsLeader = true
		}
	}

	raftConfig.StartAsLeader = opts.StartAsLeader
	// Setup the Raft store.
	b.fsm.SetNoopRestore(true)


		recoveryConfig, err := raft.ReadConfigJSON(peersFile)
		if err != nil {
			return errwrap.Wrapf("raft recovery failed to parse peers.json: {{err}}", err)
		}

		b.logger.Info("raft recovery: found new config", "config", recoveryConfig)
		err = raft.RecoverCluster(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport, recoveryConfig)
		if err != nil {
			return errwrap.Wrapf("raft recovery failed: {{err}}", err)
		}

		err = os.Remove(peersFile)
		if err != nil {
			return errwrap.Wrapf("raft recovery failed to delete peers.json; please delete manually: {{err}}", err)
		}
		b.logger.Info("raft recovery deleted peers.json")
	}
	if opts.RecoveryModeConfig != nil {
		err = raft.RecoverCluster(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport, *opts.RecoveryModeConfig)
		if err != nil {
			return errwrap.Wrapf("recovering raft cluster failed: {{err}}", err)
		}
	}

	raftObj, err := raft.NewRaft(raftConfig, b.fsm.chunker, b.logStore, b.stableStore, b.snapStore, b.raftTransport)
	b.fsm.SetNoopRestore(false)
	if err != nil {
		return err
	}
	b.raft = raftObj
	b.raftNotifyCh = raftNotifyCh

	if b.streamLayer != nil {
		// Add Handler to the cluster.
		opts.ClusterListener.AddHandler(consts.RaftStorageALPN, b.streamLayer)
		opts.ClusterListener.AddClient(consts.RaftStorageALPN, b.streamLayer)
	}

	return nil
}

	}

	b.l.Lock()
	future := b.raft.Shutdown()
	b.raft = nil
	b.l.Unlock()

	return future.Error()
}

// AppliedIndex returns the latest index applied to the FSM
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return 0
	}

	return b.raft.AppliedIndex()
}

// RemovePeer removes the given peer ID from the raft cluster. If the node is
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	future := b.raft.RemoveServer(raft.ServerID(peerID), 0, 0)

	return future.Error()
}

func (b *RaftBackend) GetConfiguration(ctx context.Context) (*RaftConfigurationResponse, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	b.logger.Debug("adding raft peer", "node_id", peerID, "cluster_addr", clusterAddr)

	future := b.raft.AddVoter(raft.ServerID(peerID), raft.ServerAddress(clusterAddr), 0, 0)
	return future.Error()
}

// Peers returns all the servers present in the raft cluster
	defer b.l.RUnlock()

	if b.raft == nil {
		return nil, errors.New("raft storage backend is not initialized")
	}

	future := b.raft.GetConfiguration()
	ret := make([]Peer, len(future.Configuration().Servers))
	for i, s := range future.Configuration().Servers {
		ret[i] = Peer{
			ID:      string(s.ID),
			Address: string(s.Address),
		}
	}

	return ret, nil
}

// Snapshot takes a raft snapshot, packages it into a archive file and writes it
// to the provided writer. Seal access is used to encrypt the SHASUM file so we
// can validate the snapshot was taken using the same master keys or not.
func (b *RaftBackend) Snapshot(out *logical.HTTPResponseWriter, access *seal.Access) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage backend is sealed")
	}

	// If we have access to the seal create a sealer object
		}
	}

	snap, err := snapshot.NewWithSealer(b.logger.Named("snapshot"), b.raft, s)
	if err != nil {
		return err
	}
	defer snap.Close()

	size, err := snap.Size()
	if err != nil {
		return err
	}

	out.Header().Add("Content-Disposition", "attachment")
	out.Header().Add("Content-Length", fmt.Sprintf("%d", size))
	out.Header().Add("Content-Type", "application/gzip")
	_, err = io.Copy(out, snap)
	if err != nil {
		return err
	}

	return nil
}

// WriteSnapshotToTemp reads a snapshot archive off the provided reader,

	var metadata raft.SnapshotMeta
	if b.raft == nil {
		return nil, nil, metadata, errors.New("raft storage backend is sealed")
	}

	// If we have access to the seal create a sealer object
	// snapshot applied to a quorum of nodes.
	command := &LogData{
		Operations: []*LogOperation{
			&LogOperation{
				OpType: restoreCallbackOp,
			},
		},
	}

	b.l.RLock()
	err := b.applyLog(ctx, command)
	b.l.RUnlock()

	// Do a best-effort attempt to let the standbys apply the restoreCallbackOp
	// before we continue.
	defer metrics.MeasureSince([]string{"raft-storage", "delete"}, time.Now())
	command := &LogData{
		Operations: []*LogOperation{
			&LogOperation{
				OpType: deleteOp,
				Key:    path,
			},
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	return b.fsm.Get(ctx, path)
}

// Put inserts an entry in the log for the put operation
func (b *RaftBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"raft-storage", "put"}, time.Now())
	command := &LogData{
		Operations: []*LogOperation{
			&LogOperation{
				OpType: putOp,
				Key:    entry.Key,
				Value:  entry.Value,
// persisted to the local FSM. Caller should hold the backend's read lock.
func (b *RaftBackend) applyLog(ctx context.Context, command *LogData) error {
	if b.raft == nil {
		return errors.New("raft storage backend is not initialized")
	}

	commandBytes, err := proto.Marshal(command)
		return err
	}

	var chunked bool
	var applyFuture raft.ApplyFuture
	switch {
	return nil
}

// HAEnabled is the implemention of the HABackend interface
func (b *RaftBackend) HAEnabled() bool { return true }

// HAEnabled is the implemention of the HABackend interface
func (b *RaftBackend) LockWith(key, value string) (physical.Lock, error) {
	return &RaftLock{
		key:   key,
	}, nil
}

// RaftLock implements the physical Lock interface and enables HA for this
// backend. The Lock uses the raftNotifyCh for receiving leadership edge
// triggers. Vault's active duty matches raft's leadership.
func (l *RaftLock) monitorLeadership(stopCh <-chan struct{}, leaderNotifyCh <-chan bool) <-chan struct{} {
	leaderLost := make(chan struct{})
	go func() {
		select {
		case isLeader := <-leaderNotifyCh:
			if !isLeader {
				close(leaderLost)
			}
		case <-stopCh:
		}
	}()
	return leaderLost
// Lock blocks until we become leader or are shutdown. It returns a channel that
// is closed when we detect a loss of leadership.
func (l *RaftLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	l.b.l.RLock()

	// Cache the notifyCh locally
	leaderNotifyCh := l.b.raftNotifyCh

	// TODO: Remove when Raft can server as the ha_storage backend. The internal
	// raft pointer should not be nil here, but the nil check is a guard against
	// https://github.com/hashicorp/vault/issues/8206
	if l.b.raft == nil {
		return nil, errors.New("attempted to grab a lock on a nil raft backend")
	}

	// Check to see if we are already leader.
	if l.b.raft.State() == raft.Leader {
		err := l.b.applyLog(context.Background(), &LogData{
			Operations: []*LogOperation{
				&LogOperation{
					OpType: putOp,
					Key:    l.key,
					Value:  l.value,
				l.b.l.RLock()
				err := l.b.applyLog(context.Background(), &LogData{
					Operations: []*LogOperation{
						&LogOperation{
							OpType: putOp,
							Key:    l.key,
							Value:  l.value,
			return nil, nil
		}
	}

	return nil, nil
}

// Unlock gives up leadership.
func (l *RaftLock) Unlock() error {
	return l.b.raft.LeadershipTransfer().Error()
}


	return s.access.Decrypt(ctx, &eblob, nil)
}
