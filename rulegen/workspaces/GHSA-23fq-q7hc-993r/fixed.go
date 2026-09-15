package main

	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	systemd "github.com/coreos/go-systemd/daemon"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/gatedwriter"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/internalshared/reloadutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/mlock"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/posener/complete"
	"go.uber.org/atomic"
	"golang.org/x/net/http/httpproxy"
	"google.golang.org/grpc/grpclog"
)

var (
	_ cli.Command             = (*ServerCommand)(nil)
	_ cli.CommandAutocomplete = (*ServerCommand)(nil)
)

var memProfilerEnabled = false

	return 1
}

const (
	storageMigrationLock = "core/migration"

	// Even though there are more types than the ones below, the following consts
	// are declared internally for value comparison and reusability.
	storageTypeRaft   = "raft"
	storageTypeConsul = "consul"
)

type ServerCommand struct {
	*BaseCommand

	logOutput   io.Writer
	gatedWriter *gatedwriter.Writer
	logger      log.InterceptLogger

	cleanupGuard sync.Once

	reloadFuncsLock   *sync.RWMutex
	reloadFuncs       *map[string][]reloadutil.ReloadFunc
	startedCh         chan (struct{}) // for tests
	reloadedCh        chan (struct{}) // for tests
	licenseReloadedCh chan (error)    // for tests

	allLoggers []log.Logger

	// new stuff
	flagConfigs            []string
	flagTestServerConfig   bool
	flagDevConsul          bool
	flagExitOnCoreShutdown bool
	flagDiagnose           string
}

func (c *ServerCommand) Synopsis() string {
			"Using a recovery operation token, \"sys/raw\" API can be used to manipulate the storage.",
	})

	// Disabled by default until functional
	if os.Getenv(OperatorDiagnoseEnableEnv) != "" {
		f.StringVar(&StringVar{
			Name:    "diagnose",
			Target:  &c.flagDiagnose,
			Default: notSetValue,
			Usage:   "Run diagnostics before starting Vault. Specify a filename to direct output to that file.",
		})
	} else {
		// Ensure diagnose is *not* run when feature flag is off.
		c.flagDiagnose = notSetValue
	}

	f = set.NewFlagSet("Dev Options")

	f.BoolVar(&BoolVar{
	return c.Flags().Completions()
}

func (c *ServerCommand) flushLog() {
	c.logger.(hclog.OutputResettable).ResetOutputWithFlush(&hclog.LoggerOptions{
		Output: c.logOutput,
	}, c.gatedWriter)
}

func (c *ServerCommand) parseConfig() (*server.Config, []configutil.ConfigError, error) {
	var configErrors []configutil.ConfigError
	// Load the configuration
	var config *server.Config
	for _, path := range c.flagConfigs {
		current, err := server.LoadConfig(path)
		if err != nil {
			return nil, nil, fmt.Errorf("error loading configuration from %s: %w", path, err)
		}

		configErrors = append(configErrors, current.Validate(path)...)

		if config == nil {
			config = current
		} else {
			config = config.Merge(current)
		}
	}
	return config, configErrors, nil
}

func (c *ServerCommand) runRecoveryMode() int {
	config, _, err := c.parseConfig()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
		return 1
	}

	c.logger = log.NewInterceptLogger(&log.LoggerOptions{
		Output: c.gatedWriter,
		Level:  level,
		// Note that if logFormat is either unspecified or standard, then
		JSONFormat: logFormat == logging.JSONFormat,
	})

	// Ensure logging is flushed if initialization fails
	defer c.flushLog()

	logLevelStr, err := c.adjustLogLevel(config, logLevelWasNotSet)
	if err != nil {
		c.UI.Error(err.Error())
		vault.DefaultMaxRequestDuration = config.DefaultMaxRequestDuration
	}

	logProxyEnvironmentVariables(c.logger)

	// Initialize the storage backend
	factory, exists := c.PhysicalBackends[config.Storage.Type]
		c.UI.Error(fmt.Sprintf("Unknown storage type %s", config.Storage.Type))
		return 1
	}
	if config.Storage.Type == storageTypeRaft || (config.HAStorage != nil && config.HAStorage.Type == storageTypeRaft) {
		if envCA := os.Getenv("VAULT_CLUSTER_ADDR"); envCA != "" {
			config.ClusterAddr = envCA
		}

	var barrierSeal vault.Seal
	var sealConfigError error
	var wrapper wrapping.Wrapper

	if len(config.Seals) == 0 {
		config.Seals = append(config.Seals, &configutil.KMS{Type: wrapping.Shamir})
	}

	if len(config.Seals) > 1 {
		sealType = configSeal.Type
	}

	infoKeys = append(infoKeys, "Seal Type")
	info["Seal Type"] = sealType

	var seal vault.Seal
	defaultSeal := vault.NewDefaultSeal(&vaultseal.Access{
		Wrapper: aeadwrapper.NewShamirWrapper(&wrapping.WrapperOptions{
			Logger: c.logger.Named("shamir"),
		}),
	})
	sealLogger := c.logger.ResetNamed(fmt.Sprintf("seal.%s", sealType))
	wrapper, sealConfigError = configutil.ConfigureWrapper(configSeal, &infoKeys, &info, sealLogger)
	if sealConfigError != nil {
		if !errwrap.ContainsType(sealConfigError, new(logical.KeyNotFoundError)) {
			c.UI.Error(fmt.Sprintf(
			return 1
		}
	}
	if wrapper == nil {
		seal = defaultSeal
	} else {
		seal = vault.NewAutoSeal(&vaultseal.Access{
			Wrapper: wrapper,
		})
	}
	barrierSeal = seal

	// Ensure that the seal finalizer is called, even if using verify-only
	}

	// Initialize the listeners
	lns := make([]listenerutil.Listener, 0, len(config.Listeners))
	for _, lnConfig := range config.Listeners {
		ln, _, _, err := server.NewListener(lnConfig, c.gatedWriter, c.UI)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing listener of type %s: %s", lnConfig.Type, err))
			return 1
		}

		lns = append(lns, listenerutil.Listener{
			Listener: ln,
			Config:   lnConfig,
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

	infoKeys = append(infoKeys, "go version")
	info["go version"] = runtime.Version()

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
			ListenerConfig:        ln.Config,
			DisablePrintableCheck: config.DisablePrintableCheck,
			RecoveryMode:          c.flagRecovery,
			RecoveryToken:         atomic.NewString(""),
	}

	if sealConfigError != nil {
		init, err := core.InitializedLocally(context.Background())
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error checking if core is initialized: %v", err))
			return 1
		c.UI.Output("==> Vault server started! Log data will stream in below:\n")
	}

	c.flushLog()

	for {
		select {
	return 0
}

func logProxyEnvironmentVariables(logger hclog.Logger) {
	proxyCfg := httpproxy.FromEnvironment()
	cfgMap := map[string]string{
		"http_proxy":  proxyCfg.HTTPProxy,
		"https_proxy": proxyCfg.HTTPSProxy,
		"no_proxy":    proxyCfg.NoProxy,
	}
	for k, v := range cfgMap {
		u, err := url.Parse(v)
		if err != nil {
			// Env vars may contain URLs or host:port values.  We only care
			// about the former.
			continue
		}
		if _, ok := u.User.Password(); ok {
			u.User = url.UserPassword("redacted-username", "redacted-password")
		} else if user := u.User.Username(); user != "" {
			u.User = url.User("redacted-username")
		}
		cfgMap[k] = u.String()
	}
	logger.Info("proxy environment", "http_proxy", cfgMap["http_proxy"],
		"https_proxy", cfgMap["https_proxy"], "no_proxy", cfgMap["no_proxy"])
}

func (c *ServerCommand) adjustLogLevel(config *server.Config, logLevelWasNotSet bool) (string, error) {
	var logLevelString string
	if config.LogLevel != "" && logLevelWasNotSet {
	return level, logLevelString, logLevelWasNotSet, logFormat, nil
}

type quiescenceSink struct {
	t *time.Timer
}

func (q quiescenceSink) Accept(name string, level log.Level, msg string, args ...interface{}) {
	q.t.Reset(100 * time.Millisecond)
}

func (c *ServerCommand) setupStorage(config *server.Config) (physical.Backend, error) {
	// Ensure that a backend is provided
	if config.Storage == nil {
		return nil, fmt.Errorf("A storage backend must be specified")
	}

	// Initialize the backend
	factory, exists := c.PhysicalBackends[config.Storage.Type]
	if !exists {
		return nil, fmt.Errorf("Unknown storage type %s", config.Storage.Type)
	}

	// Do any custom configuration needed per backend
	switch config.Storage.Type {
	case storageTypeConsul:
		if config.ServiceRegistration == nil {
			// If Consul is configured for storage and service registration is unconfigured,
			// use Consul for service registration without requiring additional configuration.
			// This maintains backward-compatibility.
			config.ServiceRegistration = &server.ServiceRegistration{
				Type:   "consul",
				Config: config.Storage.Config,
			}
		}
	case storageTypeRaft:
		if envCA := os.Getenv("VAULT_CLUSTER_ADDR"); envCA != "" {
			config.ClusterAddr = envCA
		}
		if len(config.ClusterAddr) == 0 {
			return nil, fmt.Errorf("Cluster address must be set when using raft storage")
		}
	}

	namedStorageLogger := c.logger.Named("storage." + config.Storage.Type)
	c.allLoggers = append(c.allLoggers, namedStorageLogger)
	backend, err := factory(config.Storage.Config, namedStorageLogger)
	if err != nil {
		return nil, fmt.Errorf("Error initializing storage of type %s: %w", config.Storage.Type, err)
	}

	return backend, nil
}

func beginServiceRegistration(c *ServerCommand, config *server.Config) (sr.ServiceRegistration, error) {
	sdFactory, ok := c.ServiceRegistrations[config.ServiceRegistration.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown service_registration type %s", config.ServiceRegistration.Type)
	}

	namedSDLogger := c.logger.Named("service_registration." + config.ServiceRegistration.Type)
	c.allLoggers = append(c.allLoggers, namedSDLogger)

	// Since we haven't even begun starting Vault's core yet,
	// we know that Vault is in its pre-running state.
	state := sr.State{
		VaultVersion:         version.GetVersion().VersionNumber(),
		IsInitialized:        false,
		IsSealed:             true,
		IsActive:             false,
		IsPerformanceStandby: false,
	}
	var err error
	configSR, err := sdFactory(config.ServiceRegistration.Config, namedSDLogger, state)
	if err != nil {
		return nil, fmt.Errorf("Error initializing service_registration of type %s: %s", config.ServiceRegistration.Type, err)
	}

	return configSR, nil
}

// InitListeners returns a response code, error message, Listeners, and a TCP Address list.
func (c *ServerCommand) InitListeners(config *server.Config, disableClustering bool, infoKeys *[]string, info *map[string]string) (int, []listenerutil.Listener, []*net.TCPAddr, error) {
	clusterAddrs := []*net.TCPAddr{}

	// Initialize the listeners
	lns := make([]listenerutil.Listener, 0, len(config.Listeners))

	c.reloadFuncsLock.Lock()

	defer c.reloadFuncsLock.Unlock()

	var errMsg error
	for i, lnConfig := range config.Listeners {
		ln, props, reloadFunc, err := server.NewListener(lnConfig, c.gatedWriter, c.UI)
		if err != nil {
			errMsg = fmt.Errorf("Error initializing listener of type %s: %s", lnConfig.Type, err)
			return 1, nil, nil, errMsg
		}

		if reloadFunc != nil {
			relSlice := (*c.reloadFuncs)["listener|"+lnConfig.Type]
			relSlice = append(relSlice, reloadFunc)
			(*c.reloadFuncs)["listener|"+lnConfig.Type] = relSlice
		}

		if !disableClustering && lnConfig.Type == "tcp" {
			addr := lnConfig.ClusterAddress
			if addr != "" {
				tcpAddr, err := net.ResolveTCPAddr("tcp", lnConfig.ClusterAddress)
				if err != nil {
					errMsg = fmt.Errorf("Error resolving cluster_address: %s", err)
					return 1, nil, nil, errMsg
				}
				clusterAddrs = append(clusterAddrs, tcpAddr)
			} else {
				tcpAddr, ok := ln.Addr().(*net.TCPAddr)
				if !ok {
					errMsg = fmt.Errorf("Failed to parse tcp listener")
					return 1, nil, nil, errMsg
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

		if lnConfig.MaxRequestSize == 0 {
			lnConfig.MaxRequestSize = vaulthttp.DefaultMaxRequestSize
		}
		props["max_request_size"] = fmt.Sprintf("%d", lnConfig.MaxRequestSize)

		if lnConfig.MaxRequestDuration == 0 {
			lnConfig.MaxRequestDuration = vault.DefaultMaxRequestDuration
		}
		props["max_request_duration"] = lnConfig.MaxRequestDuration.String()

		lns = append(lns, listenerutil.Listener{
			Listener: ln,
			Config:   lnConfig,
		})

		// Store the listener props for output later
		key := fmt.Sprintf("listener %d", i+1)
		propsList := make([]string, 0, len(props))
		for k, v := range props {
			propsList = append(propsList, fmt.Sprintf(
				"%s: %q", k, v))
		}
		sort.Strings(propsList)
		*infoKeys = append(*infoKeys, key)
		(*info)[key] = fmt.Sprintf(
			"%s (%s)", lnConfig.Type, strings.Join(propsList, ", "))

	}
	if !disableClustering {
		if c.logger.IsDebug() {
			c.logger.Debug("cluster listener addresses synthesized", "cluster_addresses", clusterAddrs)
		}
	}
	return 0, lns, clusterAddrs, nil
}

func (c *ServerCommand) Run(args []string) int {
	f := c.Flags()

		}
	}

	if c.flagDiagnose != notSetValue {
		if c.flagDev {
			c.UI.Error("Cannot run diagnose on Vault in dev mode.")
			return 1
		}
		// TODO: add a file output flag to Diagnose
		diagnose := &OperatorDiagnoseCommand{
			BaseCommand: c.BaseCommand,
			flagDebug:   false,
			flagSkips:   []string{},
			flagConfigs: c.flagConfigs,
		}
		diagnose.RunWithParsedFlags()
	}

	// Load the configuration
	var config *server.Config
	var err error
	if c.flagDev {
		var devStorageType string
		switch {
		default:
			devStorageType = "inmem"
		}
		config, err = server.DevConfig(devStorageType)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		if c.flagDevListenAddr != "" {
			config.Listeners[0].Address = c.flagDevListenAddr
		}
		config.Listeners[0].Telemetry.UnauthenticatedMetricsAccess = true
	}

	parsedConfig, _, err := c.parseConfig()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
		return 1
	}

	config.LogFormat = logFormat.String()

	if c.flagDevThreeNode || c.flagDevFourCluster {
		c.logger = log.NewInterceptLogger(&log.LoggerOptions{
			Mutex:  &sync.Mutex{},
			Output: c.gatedWriter,
			Level:  log.Trace,
		})
	} else {
		c.logger = log.NewInterceptLogger(&log.LoggerOptions{
			Output: c.gatedWriter,
			Level:  level,
			// Note that if logFormat is either unspecified or standard, then
		})
	}

	// Ensure logging is flushed if initialization fails
	defer c.flushLog()

	c.allLoggers = []log.Logger{c.logger}

	logLevelStr, err := c.adjustLogLevel(config, logLevelWasNotSet)
	if err != nil {

	// create GRPC logger
	namedGRPCLogFaker := c.logger.Named("grpclogfaker")
	c.allLoggers = append(c.allLoggers, namedGRPCLogFaker)
	grpclog.SetLogger(&grpclogFaker{
		logger: namedGRPCLogFaker,
		log:    os.Getenv("VAULT_GRPC_LOGGING") != "",
		c.startMemProfiler()
	}

	if config.DefaultMaxRequestDuration != 0 {
		vault.DefaultMaxRequestDuration = config.DefaultMaxRequestDuration
	}

	logProxyEnvironmentVariables(c.logger)

	if envMlock := os.Getenv("VAULT_DISABLE_MLOCK"); envMlock != "" {
		var err error
		config.DisableMlock, err = strconv.ParseBool(envMlock)
		if err != nil {
			c.UI.Output("Error parsing the environment variable VAULT_DISABLE_MLOCK")
			return 1
		}
	}

	if envLicensePath := os.Getenv(EnvVaultLicensePath); envLicensePath != "" {
		config.LicensePath = envLicensePath
	}
	if envLicense := os.Getenv(EnvVaultLicense); envLicense != "" {
		config.License = envLicense
	}

	// If mlockall(2) isn't supported, show a warning. We disable this in dev
	// because it is quite scary to see when first using Vault. We also disable
				"in a Docker container, provide the IPC_LOCK cap to the container."))
	}

	inmemMetrics, metricSink, prometheusEnabled, err := configutil.SetupTelemetry(&configutil.SetupTelemetryOpts{
		Config:      config.Telemetry,
		Ui:          c.UI,
		ServiceName: "vault",
		DisplayName: "Vault",
		UserAgent:   useragent.String(),
		ClusterName: config.ClusterName,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}
	metricsHelper := metricsutil.NewMetricsHelper(inmemMetrics, prometheusEnabled)

	// Initialize the storage backend
	backend, err := c.setupStorage(config)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Prevent server startup if migration is active
	// TODO: Use OpenTelemetry to integrate this into Diagnose
	if c.storageMigrationActive(backend) {
		return 1
	}

	// Initialize the Service Discovery, if there is one
	var configSR sr.ServiceRegistration
	if config.ServiceRegistration != nil {
		configSR, err = beginServiceRegistration(c, config)
		if err != nil {
			c.UI.Output(err.Error())
			return 1
		}
	}
	info := make(map[string]string)
	info["log level"] = logLevelString
	infoKeys = append(infoKeys, "log level")
	barrierSeal, barrierWrapper, unwrapSeal, seals, sealConfigError, err := setSeal(c, config, infoKeys, info)

	// Check error here
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if seals != nil {
		for _, seal := range seals {
			// Ensure that the seal finalizer is called, even if using verify-only
			defer func(seal *vault.Seal) {
				err = (*seal).Finalize(context.Background())
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error finalizing seals: %v", err))
				}
			}(&seal)
		}
	}

	}

	// prepare a secure random reader for core
	secureRandomReader, err := configutil.CreateSecureRandomReaderFunc(config.SharedConfig, barrierWrapper)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	coreConfig := createCoreConfig(c, config, backend, configSR, barrierSeal, unwrapSeal, metricsHelper, metricSink, secureRandomReader)
	if c.flagDevThreeNode {
		return c.enableThreeNodeDevCluster(&coreConfig, info, infoKeys, c.flagDevListenAddr, os.Getenv("VAULT_DEV_TEMP_DIR"))
	}

	if c.flagDevFourCluster {
		return enableFourClusterDev(c, &coreConfig, info, infoKeys, c.flagDevListenAddr, os.Getenv("VAULT_DEV_TEMP_DIR"))
	}

	// Initialize the separate HA storage backend, if it exists
	disableClustering, err := initHaBackend(c, config, &coreConfig, backend)
	if err != nil {
		c.UI.Output(err.Error())
		return 1
	}

	// Determine the redirect address from environment variables
	err = determineRedirectAddr(c, &coreConfig, config)
	if err != nil {
		c.UI.Output(err.Error())
	}

	// After the redirect bits are sorted out, if no cluster address was
	// explicitly given, derive one from the redirect addr
	err = findClusterAddress(c, &coreConfig, config, disableClustering)
	if err != nil {
		c.UI.Output(err.Error())
		return 1
	}

	// Override the UI enabling config by the environment variable
	if enableUI := os.Getenv("VAULT_UI"); enableUI != "" {
		var err error

	// If ServiceRegistration is configured, then the backend must support HA
	isBackendHA := coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled()
	if !c.flagDev && (coreConfig.GetServiceRegistration() != nil) && !isBackendHA {
		c.UI.Output("service_registration is configured, but storage does not support HA")
		return 1
	}

	// Apply any enterprise configuration onto the coreConfig.
	adjustCoreConfigForEnt(config, &coreConfig)

	// Initialize the core
	core, newCoreError := vault.NewCore(&coreConfig)
	if newCoreError != nil {
		if vault.IsFatalError(newCoreError) {
			c.UI.Error(fmt.Sprintf("Error initializing core: %s", newCoreError))
			return 1
		}
		c.UI.Warn(wrapAtLength(
			"WARNING! A non-fatal error occurred during initialization. Please " +
				"check the logs for more information."))
		c.UI.Warn("")

	}

	// Copy the reload funcs pointers back
		}
	}

	status, lns, clusterAddrs, errMsg := c.InitListeners(config, disableClustering, &infoKeys, &info)

	if status != 0 {
		c.UI.Output("Error parsing listener configuration.")
		c.UI.Error(errMsg.Error())
		return 1
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

	infoKeys = append(infoKeys, "go version")
	info["go version"] = runtime.Version()

	sort.Strings(infoKeys)
	c.UI.Output("==> Vault server configuration:\n")

	for _, k := range infoKeys {
		c.UI.Output(fmt.Sprintf(
			"%24s: %s",
			strings.Title(k),
			info[k]))
	}

	c.UI.Output("")

	// Tests might not want to start a vault server and just want to verify
		Core: core,
	}))

	// Attempt unsealing in a background goroutine. This is needed for when a
	// Vault cluster with multiple servers is configured with auto-unseal but is
	// uninitialized. Once one server initializes the storage backend, this
	// goroutine will pick up the unseal keys and unseal this instance.
	if !core.IsInSealMigrationMode() {
		go runUnseal(c, core, context.Background())
	}

	// When the underlying storage is raft, kick off retry join if it was specified
	// in the configuration
	// TODO: Should we also support retry_join for ha_storage?
	if config.Storage.Type == storageTypeRaft {
		if err := core.InitiateRetryJoin(context.Background()); err != nil {
			c.UI.Error(fmt.Sprintf("Failed to initiate raft retry join, %q", err.Error()))
			return 1
	}

	// Perform initialization of HTTP server after the verifyOnly check.

	// Instantiate the wait group
	c.WaitGroup = &sync.WaitGroup{}

	// If service discovery is available, run service discovery
	err = runListeners(c, &coreConfig, config, configSR)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// If we're in Dev mode, then initialize the core
	err = initDevCore(c, &coreConfig, config, core)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Initialize the HTTP servers
	err = startHttpServers(c, core, config, lns)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if c.flagTestServerConfig {
		return 0
	}

	if sealConfigError != nil {
		init, err := core.InitializedLocally(context.Background())
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error checking if core is initialized: %v", err))
			return 1
		}
		if init {
			c.UI.Error("Vault is initialized but no Seal key could be loaded")
			return 1
		}
	}

	// Output the header that the server has started
	if !c.flagCombineLogs {
	}

	// Release the log gate.
	c.flushLog()

	// Write out the PID to the file now that server has successfully started
	if err := c.storePidFile(config.PidFile); err != nil {
		return 1
	}

	// Notify systemd that the server is ready (if applicable)
	c.notifySystemd(systemd.SdNotifyReady)

	defer func() {
		if err := c.removePidFile(config.PidFile); err != nil {
			c.UI.Error(fmt.Sprintf("Error deleting the PID file: %s", err))
		case <-c.SighupCh:
			c.UI.Output("==> Vault reload triggered")

			// Notify systemd that the server is reloading config
			c.notifySystemd(systemd.SdNotifyReloading)

			// Check for new log level
			var config *server.Config
			var level log.Level
				c.UI.Error(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
			}

			// Reload license file
			if err = vault.LicenseReload(core); err != nil {
				c.UI.Error(err.Error())
			}

			select {
			case c.licenseReloadedCh <- err:
			default:
			}

		case <-c.SigUSR2Ch:
			logWriter := c.logger.StandardWriter(&hclog.StandardLoggerOptions{})
			pprof.Lookup("goroutine").WriteTo(logWriter, 2)
		}
	}
	// Notify systemd that the server is shutting down
	c.notifySystemd(systemd.SdNotifyStopping)

	// Stop the listeners so that we don't process further client requests.
	c.cleanupGuard.Do(listenerCloseFunc)

	// Finalize will wait until after Vault is sealed, which means the
	// request forwarding listeners will also be closed (and also
	// waited for).
	if err := core.Shutdown(); err != nil {
	return retCode
}

func (c *ServerCommand) notifySystemd(status string) {
	sent, err := systemd.SdNotify(false, status)
	if err != nil {
		c.logger.Error("error notifying systemd", "error", err)
	} else {
		if sent {
			c.logger.Debug("sent systemd notification", "notification", status)
		} else {
			c.logger.Debug("would have sent systemd notification (systemd not present)", "notification", status)
		}
	}
}

func (c *ServerCommand) enableDev(core *vault.Core, coreConfig *vault.CoreConfig) (*vault.InitResult, error) {
	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)


	isLeader, _, _, err := core.Leader()
	if err != nil && err != vault.ErrHANotEnabled {
		return nil, fmt.Errorf("failed to check active status: %w", err)
	}
	if err == nil {
		leaderCount := 5
			time.Sleep(1 * time.Second)
			isLeader, _, _, err = core.Leader()
			if err != nil {
				return nil, fmt.Errorf("failed to check active status: %w", err)
			}
			leaderCount--
		}
		}
		resp, err := core.HandleRequest(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to create root token with ID %q: %w", coreConfig.DevToken, err)
		}
		if resp == nil {
			return nil, fmt.Errorf("nil response when creating root token with ID %q", coreConfig.DevToken)
		req.Data = nil
		resp, err = core.HandleRequest(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to revoke initial root token: %w", err)
		}
	}

	}
	resp, err := core.HandleRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating default K/V store: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to create default K/V store: %w", resp.Error())
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

	infoKeys = append(infoKeys, "go version")
	info["go version"] = runtime.Version()

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

	if err := ioutil.WriteFile(filepath.Join(testCluster.TempDir, "root_token"), []byte(testCluster.RootToken), 0o755); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing token to tempfile: %s", err))
		return 1
	}
	}

	// Release the log gate.
	c.flushLog()

	// Wait for shutdown
	shutdownTriggered := false
			// Stop the listeners so that we don't process further client requests.
			c.cleanupGuard.Do(testCluster.Cleanup)

			// Finalize will wait until after Vault is sealed, which means the
			// request forwarding listeners will also be closed (and also
			// waited for).
			for _, core := range testCluster.Cores {
		}

		// Check if TLS is disabled
		if list.TLSDisable {
			scheme = "http"
		}

		// Check for address override
		addr := list.Address
		if addr == "" {
			addr = "127.0.0.1:8200"
		}

		// Check for localhost
	return url.String(), nil
}

func (c *ServerCommand) Reload(lock *sync.RWMutex, reloadFuncs *map[string][]reloadutil.ReloadFunc, configPath []string) error {
	lock.RLock()
	defer lock.RUnlock()
		case strings.HasPrefix(k, "listener|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(); err != nil {
						reloadErrors = multierror.Append(reloadErrors, fmt.Errorf("error encountered reloading listener: %w", err))
					}
				}
			}
		case strings.HasPrefix(k, "audit_file|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(); err != nil {
						reloadErrors = multierror.Append(reloadErrors, fmt.Errorf("error encountered reloading file audit device at path %q: %w", strings.TrimPrefix(k, "audit_file|"), err))
					}
				}
			}
	}

	// Open the PID file
	pidFile, err := os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("could not open pid file: %w", err)
	}
	defer pidFile.Close()

	pid := os.Getpid()
	_, err = pidFile.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return fmt.Errorf("could not write to pid file: %w", err)
	}
	return nil
}
			c.UI.Warn("\nWARNING! Unable to read storage migration status.")

			// unexpected state, so stop buffering log messages
			c.flushLog()
		}
		c.logger.Warn("storage migration check error", "error", err.Error())


func CheckStorageMigration(b physical.Backend) (*StorageMigrationStatus, error) {
	entry, err := b.Get(context.Background(), storageMigrationLock)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// setSeal return barrierSeal, barrierWrapper, unwrapSeal, and all the created seals from the configs so we can close them in Run
// The two errors are the sealConfigError and the regular error
func setSeal(c *ServerCommand, config *server.Config, infoKeys []string, info map[string]string) (vault.Seal, wrapping.Wrapper, vault.Seal, []vault.Seal, error, error) {
	var barrierSeal vault.Seal
	var unwrapSeal vault.Seal

	var sealConfigError error
	var wrapper wrapping.Wrapper
	var barrierWrapper wrapping.Wrapper
	if c.flagDevAutoSeal {
		barrierSeal = vault.NewAutoSeal(vaultseal.NewTestSeal(nil))
		return barrierSeal, nil, nil, nil, nil, nil
	}

	// Handle the case where no seal is provided
	switch len(config.Seals) {
	case 0:
		config.Seals = append(config.Seals, &configutil.KMS{Type: wrapping.Shamir})
	case 1:
		// If there's only one seal and it's disabled assume they want to
		// migrate to a shamir seal and simply didn't provide it
		if config.Seals[0].Disabled {
			config.Seals = append(config.Seals, &configutil.KMS{Type: wrapping.Shamir})
		}
	}
	var createdSeals []vault.Seal = make([]vault.Seal, len(config.Seals))
	for _, configSeal := range config.Seals {
		sealType := wrapping.Shamir
		if !configSeal.Disabled && os.Getenv("VAULT_SEAL_TYPE") != "" {
			sealType = os.Getenv("VAULT_SEAL_TYPE")
			configSeal.Type = sealType
		} else {
			sealType = configSeal.Type
		}

		var seal vault.Seal
		sealLogger := c.logger.ResetNamed(fmt.Sprintf("seal.%s", sealType))
		c.allLoggers = append(c.allLoggers, sealLogger)
		defaultSeal := vault.NewDefaultSeal(&vaultseal.Access{
			Wrapper: aeadwrapper.NewShamirWrapper(&wrapping.WrapperOptions{
				Logger: c.logger.Named("shamir"),
			}),
		})
		var sealInfoKeys []string
		sealInfoMap := map[string]string{}
		wrapper, sealConfigError = configutil.ConfigureWrapper(configSeal, &sealInfoKeys, &sealInfoMap, sealLogger)
		if sealConfigError != nil {
			if !errwrap.ContainsType(sealConfigError, new(logical.KeyNotFoundError)) {
				return barrierSeal, barrierWrapper, unwrapSeal, createdSeals, sealConfigError, fmt.Errorf(
					"Error parsing Seal configuration: %s", sealConfigError)
			}
		}
		if wrapper == nil {
			seal = defaultSeal
		} else {
			seal = vault.NewAutoSeal(&vaultseal.Access{
				Wrapper: wrapper,
			})
		}
		infoPrefix := ""
		if configSeal.Disabled {
			unwrapSeal = seal
			infoPrefix = "Old "
		} else {
			barrierSeal = seal
			barrierWrapper = wrapper
		}
		for _, k := range sealInfoKeys {
			infoKeys = append(infoKeys, infoPrefix+k)
			info[infoPrefix+k] = sealInfoMap[k]
		}
		createdSeals = append(createdSeals, seal)
	}
	return barrierSeal, barrierWrapper, unwrapSeal, createdSeals, sealConfigError, nil
}

func initHaBackend(c *ServerCommand, config *server.Config, coreConfig *vault.CoreConfig, backend physical.Backend) (bool, error) {
	// Initialize the separate HA storage backend, if it exists
	var ok bool
	if config.HAStorage != nil {
		if config.Storage.Type == storageTypeRaft && config.HAStorage.Type == storageTypeRaft {
			return false, fmt.Errorf("Raft cannot be set both as 'storage' and 'ha_storage'. Setting 'storage' to 'raft' will automatically set it up for HA operations as well")
		}

		if config.Storage.Type == storageTypeRaft {
			return false, fmt.Errorf("HA storage cannot be declared when Raft is the storage type")
		}

		factory, exists := c.PhysicalBackends[config.HAStorage.Type]
		if !exists {
			return false, fmt.Errorf("Unknown HA storage type %s", config.HAStorage.Type)
		}

		namedHALogger := c.logger.Named("ha." + config.HAStorage.Type)
		c.allLoggers = append(c.allLoggers, namedHALogger)
		habackend, err := factory(config.HAStorage.Config, namedHALogger)
		if err != nil {
			return false, fmt.Errorf("Error initializing HA storage of type %s: %s", config.HAStorage.Type, err)
		}

		if coreConfig.HAPhysical, ok = habackend.(physical.HABackend); !ok {
			return false, fmt.Errorf("Specified HA storage does not support HA")
		}

		if !coreConfig.HAPhysical.HAEnabled() {
			return false, fmt.Errorf("Specified HA storage has HA support disabled; please consult documentation")
		}

		coreConfig.RedirectAddr = config.HAStorage.RedirectAddr
		disableClustering := config.HAStorage.DisableClustering

		if config.HAStorage.Type == storageTypeRaft && disableClustering {
			return disableClustering, fmt.Errorf("Disable clustering cannot be set to true when Raft is the HA storage type")
		}

		if !disableClustering {
			coreConfig.ClusterAddr = config.HAStorage.ClusterAddr
		}
	} else {
		if coreConfig.HAPhysical, ok = backend.(physical.HABackend); ok {
			coreConfig.RedirectAddr = config.Storage.RedirectAddr
			disableClustering := config.Storage.DisableClustering

			if (config.Storage.Type == storageTypeRaft) && disableClustering {
				return disableClustering, fmt.Errorf("Disable clustering cannot be set to true when Raft is the storage type")
			}

			if !disableClustering {
				coreConfig.ClusterAddr = config.Storage.ClusterAddr
			}
		}
	}
	return config.DisableClustering, nil
}

func determineRedirectAddr(c *ServerCommand, coreConfig *vault.CoreConfig, config *server.Config) error {
	var retErr error
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

	var ok bool
	var detect physical.RedirectDetect
	if coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled() {
		detect, ok = coreConfig.HAPhysical.(physical.RedirectDetect)
	} else {
		detect, ok = coreConfig.Physical.(physical.RedirectDetect)
	}
	if ok && coreConfig.RedirectAddr == "" {
		redirect, err := c.detectRedirect(detect, config)
		// the following errors did not cause Run to return, so I'm not returning these
		// as errors.
		if err != nil {
			retErr = fmt.Errorf("Error detecting api address: %s", err)
		} else if redirect == "" {
			retErr = fmt.Errorf("Failed to detect api address")
		} else {
			coreConfig.RedirectAddr = redirect
		}
	}
	if coreConfig.RedirectAddr == "" && c.flagDev {
		coreConfig.RedirectAddr = fmt.Sprintf("http://%s", config.Listeners[0].Address)
	}
	return retErr
}

func findClusterAddress(c *ServerCommand, coreConfig *vault.CoreConfig, config *server.Config, disableClustering bool) error {
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
			addrToUse = fmt.Sprintf("http://%s", config.Listeners[0].Address)
		default:
			goto CLUSTER_SYNTHESIS_COMPLETE
		}
		u, err := url.ParseRequestURI(addrToUse)
		if err != nil {
			return fmt.Errorf("Error parsing synthesized cluster address %s: %v", addrToUse, err)
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			// This sucks, as it's a const in the function but not exported in the package
			if strings.Contains(err.Error(), "missing port in address") {
				host = u.Host
				port = "443"
			} else {
				return fmt.Errorf("Error parsing api address: %v", err)
			}
		}
		nPort, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("Error parsing synthesized address; failed to convert %q to a numeric: %v", port, err)
		}
		u.Host = net.JoinHostPort(host, strconv.Itoa(nPort+1))
		// Will always be TLS-secured
		u.Scheme = "https"
		coreConfig.ClusterAddr = u.String()
	}

CLUSTER_SYNTHESIS_COMPLETE:

	if coreConfig.RedirectAddr == coreConfig.ClusterAddr && len(coreConfig.RedirectAddr) != 0 {
		return fmt.Errorf("Address %q used for both API and cluster addresses", coreConfig.RedirectAddr)
	}

	if coreConfig.ClusterAddr != "" {
		// Force https as we'll always be TLS-secured
		u, err := url.ParseRequestURI(coreConfig.ClusterAddr)
		if err != nil {
			return fmt.Errorf("Error parsing cluster address %s: %v", coreConfig.ClusterAddr, err)
		}
		u.Scheme = "https"
		coreConfig.ClusterAddr = u.String()
	}
	return nil
}

func runUnseal(c *ServerCommand, core *vault.Core, ctx context.Context) {
	for {
		err := core.UnsealWithStoredKeys(ctx)
		if err == nil {
			return
		}

		if vault.IsFatalError(err) {
			c.logger.Error("error unsealing core", "error", err)
			return
		}
		c.logger.Warn("failed to unseal core", "error", err)

		select {
		case <-c.ShutdownCh:
			return
		case <-time.After(5 * time.Second):
		}
	}
}

func createCoreConfig(c *ServerCommand, config *server.Config, backend physical.Backend, configSR sr.ServiceRegistration, barrierSeal, unwrapSeal vault.Seal,
	metricsHelper *metricsutil.MetricsHelper, metricSink *metricsutil.ClusterMetricSink, secureRandomReader io.Reader) vault.CoreConfig {
	coreConfig := &vault.CoreConfig{
		RawConfig:                      config,
		Physical:                       backend,
		RedirectAddr:                   config.Storage.RedirectAddr,
		StorageType:                    config.Storage.Type,
		HAPhysical:                     nil,
		ServiceRegistration:            configSR,
		Seal:                           barrierSeal,
		UnwrapSeal:                     unwrapSeal,
		AuditBackends:                  c.AuditBackends,
		CredentialBackends:             c.CredentialBackends,
		LogicalBackends:                c.LogicalBackends,
		Logger:                         c.logger,
		DisableSentinelTrace:           config.DisableSentinelTrace,
		DisableCache:                   config.DisableCache,
		DisableMlock:                   config.DisableMlock,
		MaxLeaseTTL:                    config.MaxLeaseTTL,
		DefaultLeaseTTL:                config.DefaultLeaseTTL,
		ClusterName:                    config.ClusterName,
		CacheSize:                      config.CacheSize,
		PluginDirectory:                config.PluginDirectory,
		EnableUI:                       config.EnableUI,
		EnableRaw:                      config.EnableRawEndpoint,
		DisableSealWrap:                config.DisableSealWrap,
		DisablePerformanceStandby:      config.DisablePerformanceStandby,
		DisableIndexing:                config.DisableIndexing,
		AllLoggers:                     c.allLoggers,
		BuiltinRegistry:                builtinplugins.Registry,
		DisableKeyEncodingChecks:       config.DisablePrintableCheck,
		MetricsHelper:                  metricsHelper,
		MetricSink:                     metricSink,
		SecureRandomReader:             secureRandomReader,
		EnableResponseHeaderHostname:   config.EnableResponseHeaderHostname,
		EnableResponseHeaderRaftNodeID: config.EnableResponseHeaderRaftNodeID,
		License:                        config.License,
		LicensePath:                    config.LicensePath,
	}
	if c.flagDev {
		coreConfig.EnableRaw = true
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
	return *coreConfig
}

func runListeners(c *ServerCommand, coreConfig *vault.CoreConfig, config *server.Config, configSR sr.ServiceRegistration) error {
	if sd := coreConfig.GetServiceRegistration(); sd != nil {
		if err := configSR.Run(c.ShutdownCh, c.WaitGroup, coreConfig.RedirectAddr); err != nil {
			return fmt.Errorf("Error running service_registration of type %s: %s", config.ServiceRegistration.Type, err)
		}
	}
	return nil
}

func initDevCore(c *ServerCommand, coreConfig *vault.CoreConfig, config *server.Config, core *vault.Core) error {
	if c.flagDev && !c.flagDevSkipInit {

		init, err := c.enableDev(core, coreConfig)
		if err != nil {
			return fmt.Errorf("Error initializing Dev mode: %s", err)
		}

		var plugins, pluginsNotLoaded []string
		if c.flagDevPluginDir != "" && c.flagDevPluginInit {

			f, err := os.Open(c.flagDevPluginDir)
			if err != nil {
				return fmt.Errorf("Error reading plugin dir: %s", err)
			}

			list, err := f.Readdirnames(0)
			f.Close()
			if err != nil {
				return fmt.Errorf("Error listing plugins: %s", err)
			}

			for _, name := range list {
				path := filepath.Join(f.Name(), name)
				if err := c.addPlugin(path, init.RootToken, core); err != nil {
					if !errwrap.Contains(err, vault.ErrPluginBadType.Error()) {
						return fmt.Errorf("Error enabling plugin %s: %s", name, err)
					}
					pluginsNotLoaded = append(pluginsNotLoaded, name)
					continue
				}
				plugins = append(plugins, name)
			}

			sort.Strings(plugins)
		}

		var qw *quiescenceSink
		var qwo sync.Once
		qw = &quiescenceSink{
			t: time.AfterFunc(100*time.Millisecond, func() {
				qwo.Do(func() {
					c.logger.DeregisterSink(qw)

					// Print the big dev mode warning!
					c.UI.Warn(wrapAtLength(
						"WARNING! dev mode is enabled! In this mode, Vault runs entirely " +
							"in-memory and starts unsealed with a single unseal key. The root " +
							"token is already authenticated to the CLI, so you can immediately " +
							"begin using Vault."))
					c.UI.Warn("")
					c.UI.Warn("You may need to set the following environment variable:")
					c.UI.Warn("")

					endpointURL := "http://" + config.Listeners[0].Address
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
				})
			}),
		}
		c.logger.RegisterSink(qw)
	}
	return nil
}

// Initialize the HTTP servers
func startHttpServers(c *ServerCommand, core *vault.Core, config *server.Config, lns []listenerutil.Listener) error {
	for _, ln := range lns {
		if ln.Config == nil {
			return fmt.Errorf("Found nil listener config after parsing")
		}
		handler := vaulthttp.Handler(&vault.HandlerProperties{
			Core:                  core,
			ListenerConfig:        ln.Config,
			DisablePrintableCheck: config.DisablePrintableCheck,
			RecoveryMode:          c.flagRecovery,
		})

		if len(ln.Config.XForwardedForAuthorizedAddrs) > 0 {
			handler = vaulthttp.WrapForwardedForHandler(handler, ln.Config)
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
		if ln.Config.HTTPReadHeaderTimeout > 0 {
			server.ReadHeaderTimeout = ln.Config.HTTPReadHeaderTimeout
		}
		if ln.Config.HTTPReadTimeout > 0 {
			server.ReadTimeout = ln.Config.HTTPReadTimeout
		}
		if ln.Config.HTTPWriteTimeout > 0 {
			server.WriteTimeout = ln.Config.HTTPWriteTimeout
		}
		if ln.Config.HTTPIdleTimeout > 0 {
			server.IdleTimeout = ln.Config.HTTPIdleTimeout
		}

		// server config tests can exit now
		if c.flagTestServerConfig {
			continue
		}

		go server.Serve(ln.Listener)
	}
	return nil
}

func SetStorageMigration(b physical.Backend, active bool) error {
	if !active {
		return b.Delete(context.Background(), storageMigrationLock)
package docker

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-multierror"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/testing/stepwise"
	"golang.org/x/net/http2"
)

var _ stepwise.Environment = (*DockerCluster)(nil)

const dockerVersion = "1.40"

// DockerCluster is used to managing the lifecycle of the test Vault cluster
type DockerCluster struct {
	ID string
	// PluginName is the input from the test case
	PluginName string
	// ClusterName is a UUID name of the cluster.
	ClusterName string

	// MountOptions are a set of options for registering and mounting the plugin
	MountOptions stepwise.MountOptions

	RaftStorage  bool
	ClusterNodes []*dockerClusterNode

	// Certificate fields
	CACert        *x509.Certificate
	CACertBytes   []byte
	CACertPEM     []byte
	CACertPEMFile string
	CAKey         *ecdsa.PrivateKey
	CAKeyPEM      []byte
	RootCAs       *x509.CertPool

	// networkID tracks the network ID of the created docker network
	networkID string

	barrierKeys  [][]byte
	recoveryKeys [][]byte
	tmpDir       string

	clientAuthRequired bool
	// the mountpath of the plugin under test
	mountPath string
	// rootToken is the initial root token created when the Vault cluster is
	// created.
	rootToken string
}

// Teardown stops all the containers.
func (dc *DockerCluster) Teardown() error {
	var result *multierror.Error
	for _, node := range dc.ClusterNodes {
		if err := node.Cleanup(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// clean up networks
	if dc.networkID != "" {
		cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithVersion(dockerVersion))
		if err != nil {
			return multierror.Append(result, err)
		}
		if err := cli.NetworkRemove(context.Background(), dc.networkID); err != nil {
			return multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

// MountPath returns the path that the plugin under test is mounted at. If a
// MountPathPrefix was given, the mount path uses the prefix with a uuid
// appended. The default is the given PluginName with a uuid suffix.
func (dc *DockerCluster) MountPath() string {
	if dc.mountPath != "" {
		return dc.mountPath
	}

	uuidStr, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}

	prefix := dc.PluginName
	if dc.MountOptions.MountPathPrefix != "" {
		prefix = dc.MountOptions.MountPathPrefix
	}

	dc.mountPath = fmt.Sprintf("%s_%s", prefix, uuidStr)
	if dc.MountOptions.PluginType == stepwise.PluginTypeCredential {
		dc.mountPath = path.Join("auth", dc.mountPath)
	}

	return dc.mountPath
}

// RootToken returns the root token of the cluster, if set
func (dc *DockerCluster) RootToken() string {
	return dc.rootToken
}

// Name returns the name of this environment
func (dc *DockerCluster) Name() string {
	return "docker"
}

// Client returns a clone of the configured Vault API client.
func (dc *DockerCluster) Client() (*api.Client, error) {
	if len(dc.ClusterNodes) > 0 {
		if dc.ClusterNodes[0].Client != nil {
			c, err := dc.ClusterNodes[0].Client.Clone()
			if err != nil {
				return nil, err
			}
			c.SetToken(dc.ClusterNodes[0].Client.Token())
			return c, nil
		}
	}

	return nil, errors.New("no configured client found")
}

func (n *dockerClusterNode) Name() string {
	return n.Cluster.ClusterName + "-" + n.NodeID
}

func (dc *DockerCluster) Initialize(ctx context.Context) error {
	client, err := dc.ClusterNodes[0].NewAPIClient()
	if err != nil {
		return err
	}

	var resp *api.InitResponse
	for ctx.Err() == nil {
		resp, err = client.Sys().Init(&api.InitRequest{
			SecretShares:    3,
			SecretThreshold: 3,
		})
		if err == nil && resp != nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("nil response to init request")
	}

	for _, k := range resp.Keys {
		raw, err := hex.DecodeString(k)
		if err != nil {
			return err
		}
		dc.barrierKeys = append(dc.barrierKeys, raw)
	}

	for _, k := range resp.RecoveryKeys {
		raw, err := hex.DecodeString(k)
		if err != nil {
			return err
		}
		dc.recoveryKeys = append(dc.recoveryKeys, raw)
	}

	dc.rootToken = resp.RootToken

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Unseal
	for j, node := range dc.ClusterNodes {
		// copy the index value, so we're not reusing it in deeper scopes
		i := j
		client, err := node.NewAPIClient()
		if err != nil {
			return err
		}
		node.Client = client

		if i > 0 && dc.RaftStorage {
			leader := dc.ClusterNodes[0]
			resp, err := client.Sys().RaftJoin(&api.RaftJoinRequest{
				LeaderAPIAddr:    fmt.Sprintf("https://%s:%d", dc.ClusterNodes[0].Name(), leader.Address.Port),
				LeaderCACert:     string(dc.CACertPEM),
				LeaderClientCert: string(node.ServerCertPEM),
				LeaderClientKey:  string(node.ServerKeyPEM),
			})
			if err != nil {
				return err
			}
			if resp == nil || !resp.Joined {
				return fmt.Errorf("nil or negative response from raft join request: %v", resp)
			}
		}

		var unsealed bool
		for _, key := range dc.barrierKeys {
			resp, err := client.Sys().Unseal(hex.EncodeToString(key))
			if err != nil {
				return err
			}
			unsealed = !resp.Sealed
		}
		if i == 0 && !unsealed {
			return fmt.Errorf("could not unseal node %d", i)
		}
		client.SetToken(dc.rootToken)

		err = ensureHealthMatches(ctx, node.Client, func(health *api.HealthResponse) error {
			if health.Sealed {
				return fmt.Errorf("node %d is sealed: %#v", i, health)
			}
			if health.ClusterID == "" {
				return fmt.Errorf("node %d has no cluster ID", i)
			}

			dc.ID = health.ClusterID
			return nil
		})
		if err != nil {
			return err
		}

		if i == 0 {
			err = ensureLeaderMatches(ctx, node.Client, func(leader *api.LeaderResponse) error {
				if !leader.IsSelf {
					return fmt.Errorf("node %d leader=%v, expected=%v", i, leader.IsSelf, true)
				}

				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	for i, node := range dc.ClusterNodes {
		expectLeader := i == 0
		err = ensureLeaderMatches(ctx, node.Client, func(leader *api.LeaderResponse) error {
			if expectLeader != leader.IsSelf {
				return fmt.Errorf("node %d leader=%v, expected=%v", i, leader.IsSelf, expectLeader)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (dc *DockerCluster) setupCA(opts *DockerClusterOptions) error {
	var err error

	certIPs := []net.IP{
		net.IPv6loopback,
		net.ParseIP("127.0.0.1"),
	}

	var caKey *ecdsa.PrivateKey
	if opts != nil && opts.CAKey != nil {
		caKey = opts.CAKey
	} else {
		caKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
	}
	dc.CAKey = caKey

	var caBytes []byte
	if opts != nil && len(opts.CACert) > 0 {
		caBytes = opts.CACert
	} else {
		serialNumber := mathrand.New(mathrand.NewSource(time.Now().UnixNano())).Int63()
		CACertTemplate := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			DNSNames:              []string{"localhost"},
			IPAddresses:           certIPs,
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
			SerialNumber:          big.NewInt(serialNumber),
			NotBefore:             time.Now().Add(-30 * time.Second),
			NotAfter:              time.Now().Add(262980 * time.Hour),
			BasicConstraintsValid: true,
			IsCA:                  true,
		}
		caBytes, err = x509.CreateCertificate(rand.Reader, CACertTemplate, CACertTemplate, caKey.Public(), caKey)
		if err != nil {
			return err
		}
	}
	CACert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return err
	}
	dc.CACert = CACert
	dc.CACertBytes = caBytes

	dc.RootCAs = x509.NewCertPool()
	dc.RootCAs.AddCert(CACert)

	CACertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	dc.CACertPEM = pem.EncodeToMemory(CACertPEMBlock)

	dc.CACertPEMFile = filepath.Join(dc.tmpDir, "ca", "ca.pem")
	err = ioutil.WriteFile(dc.CACertPEMFile, dc.CACertPEM, 0o755)
	if err != nil {
		return err
	}

	marshaledCAKey, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		return err
	}
	CAKeyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledCAKey,
	}
	dc.CAKeyPEM = pem.EncodeToMemory(CAKeyPEMBlock)

	return nil
}

// Don't call this until n.Address.IP is populated
func (n *dockerClusterNode) setupCert() error {
	var err error

	n.ServerKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	serialNumber := mathrand.New(mathrand.NewSource(time.Now().UnixNano())).Int63()
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: n.Name(),
		},
		DNSNames:    []string{"localhost", n.Name()},
		IPAddresses: []net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(serialNumber),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}
	n.ServerCertBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, n.Cluster.CACert, n.ServerKey.Public(), n.Cluster.CAKey)
	if err != nil {
		return err
	}
	n.ServerCert, err = x509.ParseCertificate(n.ServerCertBytes)
	if err != nil {
		return err
	}
	n.ServerCertPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: n.ServerCertBytes,
	})

	marshaledKey, err := x509.MarshalECPrivateKey(n.ServerKey)
	if err != nil {
		return err
	}
	n.ServerKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledKey,
	})

	n.ServerCertPEMFile = filepath.Join(n.WorkDir, "cert.pem")
	err = ioutil.WriteFile(n.ServerCertPEMFile, n.ServerCertPEM, 0o755)
	if err != nil {
		return err
	}

	n.ServerKeyPEMFile = filepath.Join(n.WorkDir, "key.pem")
	err = ioutil.WriteFile(n.ServerKeyPEMFile, n.ServerKeyPEM, 0o755)
	if err != nil {
		return err
	}

	tlsCert, err := tls.X509KeyPair(n.ServerCertPEM, n.ServerKeyPEM)
	if err != nil {
		return err
	}

	certGetter := stepwise.NewCertificateGetter(n.ServerCertPEMFile, n.ServerKeyPEMFile, "")
	if err := certGetter.Reload(); err != nil {
		return err
	}
	tlsConfig := &tls.Config{
		Certificates:   []tls.Certificate{tlsCert},
		RootCAs:        n.Cluster.RootCAs,
		ClientCAs:      n.Cluster.RootCAs,
		ClientAuth:     tls.RequestClientCert,
		NextProtos:     []string{"h2", "http/1.1"},
		GetCertificate: certGetter.GetCertificate,
	}

	if n.Cluster.clientAuthRequired {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	n.TLSConfig = tlsConfig

	return nil
}

// NewEnvironment creats a new Stepwise Environment for executing tests
func NewEnvironment(name string, options *stepwise.MountOptions) *DockerCluster {
	if options == nil {
		return nil
	}

	clusterUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}

	return &DockerCluster{
		PluginName:   options.PluginName,
		ClusterName:  fmt.Sprintf("test-%s-%s", name, clusterUUID),
		RaftStorage:  true,
		MountOptions: *options,
	}
}

// DockerClusterNode represents a single instance of Vault in a cluster
type dockerClusterNode struct {
	NodeID            string
	Address           *net.TCPAddr
	HostPort          string
	Client            *api.Client
	ServerCert        *x509.Certificate
	ServerCertBytes   []byte
	ServerCertPEM     []byte
	ServerCertPEMFile string
	ServerKey         *ecdsa.PrivateKey
	ServerKeyPEM      []byte
	ServerKeyPEMFile  string
	TLSConfig         *tls.Config
	WorkDir           string
	Cluster           *DockerCluster
	container         *types.ContainerJSON
	dockerAPI         *docker.Client
}

// NewAPIClient creates and configures a Vault API client to communicate with
// the running Vault Cluster for this DockerClusterNode
func (n *dockerClusterNode) NewAPIClient() (*api.Client, error) {
	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = n.TLSConfig.Clone()
	if err := http2.ConfigureTransport(transport); err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			// This can of course be overridden per-test by using its own client
			return fmt.Errorf("redirects not allowed in these tests")
		},
	}
	config := api.DefaultConfig()
	if config.Error != nil {
		return nil, config.Error
	}
	config.Address = fmt.Sprintf("https://127.0.0.1:%s", n.HostPort)
	config.HttpClient = client
	config.MaxRetries = 0
	apiClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	apiClient.SetToken(n.Cluster.RootToken())
	return apiClient, nil
}

// Cleanup kills the container of the node
func (n *dockerClusterNode) Cleanup() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return n.dockerAPI.ContainerKill(ctx, n.container.ID, "KILL")
}

func (n *dockerClusterNode) start(cli *docker.Client, caDir, netName string, netCIDR *dockerClusterNode, pluginBinPath string) error {
	n.dockerAPI = cli

	err := n.setupCert()
	if err != nil {
		return err
	}

	vaultCfg := map[string]interface{}{
		"listener": map[string]interface{}{
			"tcp": map[string]interface{}{
				"address":       fmt.Sprintf("%s:%d", "0.0.0.0", 8200),
				"tls_cert_file": "/vault/config/cert.pem",
				"tls_key_file":  "/vault/config/key.pem",
				"telemetry": map[string]interface{}{
					"unauthenticated_metrics_access": true,
				},
			},
		},
		"telemetry": map[string]interface{}{
			"disable_hostname": true,
		},
		"storage": map[string]interface{}{
			"raft": map[string]interface{}{
				"path":    "/vault/file",
				"node_id": n.NodeID,
			},
		},
		"cluster_name":         netName,
		"log_level":            "TRACE",
		"raw_storage_endpoint": true,
		"plugin_directory":     "/vault/config",
		// disable_mlock is required for working in the Docker environment with
		// custom plugins
		"disable_mlock": true,
		// These are being provided by docker-entrypoint now, since we don't know
		// the address before the container starts.
		//"api_addr": fmt.Sprintf("https://%s:%d", n.Address.IP, n.Address.Port),
		//"cluster_addr": fmt.Sprintf("https://%s:%d", n.Address.IP, n.Address.Port+1),
	}
	cfgJSON, err := json.Marshal(vaultCfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(n.WorkDir, "local.json"), cfgJSON, 0o644)
	if err != nil {
		return err
	}
	// setup plugin bin copy if needed
	copyFromTo := map[string]string{
		n.WorkDir: "/vault/config",
		caDir:     "/usr/local/share/ca-certificates/",
	}
	if pluginBinPath != "" {
		base := path.Base(pluginBinPath)
		copyFromTo[pluginBinPath] = filepath.Join("/vault/config", base)
	}

	r := &Runner{
		dockerAPI: cli,
		ContainerConfig: &container.Config{
			Image: "vault",
			Entrypoint: []string{"/bin/sh", "-c", "update-ca-certificates && " +
				"exec /usr/local/bin/docker-entrypoint.sh vault server -log-level=trace -dev-plugin-dir=./vault/config -config /vault/config/local.json"},
			Env: []string{
				"VAULT_CLUSTER_INTERFACE=eth0",
				"VAULT_API_ADDR=https://127.0.0.1:8200",
				fmt.Sprintf("VAULT_REDIRECT_ADDR=https://%s:8200", n.Name()),
			},
			Labels:       nil,
			ExposedPorts: nat.PortSet{"8200/tcp": {}, "8201/tcp": {}},
		},
		ContainerName: n.Name(),
		NetName:       netName,
		CopyFromTo:    copyFromTo,
	}

	n.container, err = r.Start(context.Background())
	if err != nil {
		return err
	}

	n.Address = &net.TCPAddr{
		IP:   net.ParseIP(n.container.NetworkSettings.IPAddress),
		Port: 8200,
	}
	ports := n.container.NetworkSettings.NetworkSettingsBase.Ports[nat.Port("8200/tcp")]
	if len(ports) == 0 {
		n.Cleanup()
		return fmt.Errorf("could not find port binding for 8200/tcp")
	}
	n.HostPort = ports[0].HostPort

	return nil
}

// DockerClusterOptions has options for setting up the docker cluster
type DockerClusterOptions struct {
	KeepStandbysSealed bool
	RequireClientAuth  bool
	SkipInit           bool
	CACert             []byte
	NumCores           int
	tmpDir             string
	PluginTestBin      string
	// SetupFunc is called after the cluster is started.
	SetupFunc func(t testing.T, c *DockerCluster)
	CAKey     *ecdsa.PrivateKey
}

//
// helper methods/functions
//

// ensureHealthMatches checks health
func ensureHealthMatches(ctx context.Context, client *api.Client, ready func(response *api.HealthResponse) error) error {
	var health *api.HealthResponse
	var err error
	for ctx.Err() == nil {
		health, err = client.Sys().Health()
		switch {
		case err != nil:
		case health == nil:
			err = fmt.Errorf("nil response to health check")
		default:
			err = ready(health)
			if err == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("error checking health: %v", err)
}

func ensureLeaderMatches(ctx context.Context, client *api.Client, ready func(response *api.LeaderResponse) error) error {
	var leader *api.LeaderResponse
	var err error
	for ctx.Err() == nil {
		leader, err = client.Sys().Leader()
		switch {
		case err != nil:
		case leader == nil:
			err = fmt.Errorf("nil response to leader check")
		default:
			err = ready(leader)
			if err == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("error checking leader: %v", err)
}

// end test helper methods

// TODO: allow number of cores/servers to be configurable
var DefaultNumCores = 1

// creates a managed docker container running Vault
func (cluster *DockerCluster) setupDockerCluster(opts *DockerClusterOptions) error {
	if opts != nil && opts.tmpDir != "" {
		if _, err := os.Stat(opts.tmpDir); os.IsNotExist(err) {
			if err := os.MkdirAll(opts.tmpDir, 0o700); err != nil {
				return err
			}
		}
		cluster.tmpDir = opts.tmpDir
	} else {
		tempDir, err := ioutil.TempDir("", "vault-test-cluster-")
		if err != nil {
			return err
		}
		cluster.tmpDir = tempDir
	}
	caDir := filepath.Join(cluster.tmpDir, "ca")
	if err := os.MkdirAll(caDir, 0o755); err != nil {
		return err
	}

	var numCores int
	if opts == nil || opts.NumCores == 0 {
		numCores = DefaultNumCores
	} else {
		numCores = opts.NumCores
	}

	if opts != nil && opts.RequireClientAuth {
		cluster.clientAuthRequired = true
	}

	for i := 0; i < numCores; i++ {
		nodeID := fmt.Sprintf("vault-%d", i)
		node := &dockerClusterNode{
			NodeID:  nodeID,
			Cluster: cluster,
			WorkDir: filepath.Join(cluster.tmpDir, nodeID),
		}
		cluster.ClusterNodes = append(cluster.ClusterNodes, node)
		if err := os.MkdirAll(node.WorkDir, 0o700); err != nil {
			return err
		}
	}

	err := cluster.setupCA(opts)
	if err != nil {
		return err
	}

	cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithVersion(dockerVersion))
	if err != nil {
		return err
	}

	netUUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	netName := fmt.Sprintf("%s-%s", "vault-test", netUUID)
	netID, err := setupNetwork(cli, netName)
	if err != nil {
		return err
	}
	cluster.networkID = netID

	for _, node := range cluster.ClusterNodes {
		pluginBinPath := ""
		if opts != nil {
			pluginBinPath = opts.PluginTestBin
		}

		err := node.start(cli, caDir, netName, node, pluginBinPath)
		if err != nil {
			return err
		}
	}

	if opts == nil || !opts.SkipInit {
		if err := cluster.Initialize(context.Background()); err != nil {
			return err
		}
	}

	return nil
}

// Docker networking functions
// setupNetwork establishes networking for the Docker container
func setupNetwork(cli *docker.Client, netName string) (string, error) {
	id, err := createNetwork(cli, netName)
	if err != nil {
		return "", fmt.Errorf("couldn't create network %s: %w", netName, err)
	}
	return id, nil
}

func createNetwork(cli *docker.Client, netName string) (string, error) {
	resp, err := cli.NetworkCreate(context.Background(), netName, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		Options:        map[string]string{},
		IPAM: &network.IPAM{
			Driver:  "default",
			Options: map[string]string{},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// Setup creates any temp directories needed and compiles the binary for copying to Docker
func (dc *DockerCluster) Setup() error {
	registryName := dc.MountOptions.RegistryName
	pluginName := dc.MountOptions.PluginName

	// get the working directory of the plugin being tested.
	srcDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// tmpDir gets cleaned up when the cluster is cleaned up
	tmpDir, err := ioutil.TempDir("", "bin")
	if err != nil {
		return err
	}

	binName, binPath, sha256value, err := stepwise.CompilePlugin(registryName, pluginName, srcDir, tmpDir)
	if err != nil {
		return err
	}

	dOpts := &DockerClusterOptions{PluginTestBin: binPath}
	if err := dc.setupDockerCluster(dOpts); err != nil {
		return err
	}

	cores := dc.ClusterNodes
	client := cores[0].Client

	// use client to mount plugin
	err = client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    registryName,
		Type:    consts.PluginType(dc.MountOptions.PluginType),
		Command: binName,
		SHA256:  sha256value,
	})
	if err != nil {
		return err
	}

	switch dc.MountOptions.PluginType {
	case stepwise.PluginTypeCredential:
		// the mount path includes "auth/" for credential type plugins. For enabling
		// auth mounts via the /sys endpoint, we need to remove that prefix
		authPath := strings.TrimPrefix(dc.MountPath(), "auth/")
		err = client.Sys().EnableAuthWithOptions(authPath, &api.EnableAuthOptions{
			Type: registryName,
		})
	case stepwise.PluginTypeDatabase:
	case stepwise.PluginTypeSecrets:
		err = client.Sys().Mount(dc.MountPath(), &api.MountInput{
			Type: registryName,
		})
	default:
		return fmt.Errorf("unknown plugin type: %s", dc.MountOptions.PluginType.String())
	}
	return err
}

	"github.com/armon/go-metrics"
	"github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-raftchunking"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/raft"
	autopilot "github.com/hashicorp/raft-autopilot"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	snapshot "github.com/hashicorp/raft-snapshot"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/vault/seal"
	bolt "go.etcd.io/bbolt"
)

// EnvVaultRaftNodeID is used to fetch the Raft node ID from the environment.
const EnvVaultRaftPath = "VAULT_RAFT_PATH"

// Verify RaftBackend satisfies the correct interfaces
var (
	_ physical.Backend       = (*RaftBackend)(nil)
	_ physical.Transactional = (*RaftBackend)(nil)
	_ physical.HABackend     = (*RaftBackend)(nil)
	_ physical.Lock          = (*RaftLock)(nil)
)

var (
	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512

	raftState     = "raft/"
	peersFileName = "peers.json"

	restoreOpDelayDuration = 5 * time.Second

	defaultMaxEntrySize = uint64(2 * raftchunking.ChunkSize)
)

// RaftBackend implements the backend interfaces and uses the raft protocol to
	// raft is the instance of raft we will operate on.
	raft *raft.Raft

	// raftInitCh is used to block during HA lock acquisition if raft
	// has not been initialized yet, which can occur if raft is being
	// used for HA-only.
	raftInitCh chan struct{}

	// raftNotifyCh is used to receive updates about leadership changes
	// regarding this node.
	raftNotifyCh chan bool

	// permitPool is used to limit the number of concurrent storage calls.
	permitPool *physical.PermitPool

	// maxEntrySize imposes a size limit (in bytes) on a raft entry (put or transaction).
	// It is suggested to use a value of 2x the Raft chunking size for optimal
	// performance.
	maxEntrySize uint64

	// autopilot is the instance of raft-autopilot library implementation of the
	// autopilot features. This will be instantiated in both leader and followers.
	// However, only active node will have a "running" autopilot.
	autopilot *autopilot.Autopilot

	// autopilotConfig represents the configuration required to instantiate autopilot.
	autopilotConfig *AutopilotConfig

	// followerStates represents the information about all the peers of the raft
	// leader. This is used to track some state of the peers and as well as used
	// to see if the peers are "alive" using the heartbeat received from them.
	followerStates *FollowerStates

	// followerHeartbeatTicker is used to compute dead servers using follower
	// state heartbeats.
	followerHeartbeatTicker *time.Ticker

	// disableAutopilot if set will not put autopilot implementation to use. The
	// fallback will be to interact with the raft instance directly. This can only
	// be set during startup via the environment variable
	// VAULT_RAFT_AUTOPILOT_DISABLE during startup and can't be updated once the
	// node is up and running.
	disableAutopilot bool

	autopilotReconcileInterval time.Duration
}

// LeaderJoinInfo contains information required by a node to join itself as a
// follower to an existing raft cluster
type LeaderJoinInfo struct {
	// AutoJoin defines any cloud auto-join metadata. If supplied, Vault will
	// attempt to automatically discover peers in addition to what can be provided
	// via 'leader_api_addr'.
	AutoJoin string `json:"auto_join"`

	// AutoJoinScheme defines the optional URI protocol scheme for addresses
	// discovered via auto-join.
	AutoJoinScheme string `json:"auto_join_scheme"`

	// AutoJoinPort defines the optional port used for addressed discovered via
	// auto-join.
	AutoJoinPort uint `json:"auto_join_port"`

	// LeaderAPIAddr is the address of the leader node to connect to
	LeaderAPIAddr string `json:"leader_api_addr"`

	// LeaderCACert is the CA cert of the leader node
	LeaderCACert string `json:"leader_ca_cert"`

	// LeaderClientCert is the client certificate for the follower node to
	// establish client authentication during TLS
	LeaderClientCert string `json:"leader_client_cert"`

	// LeaderClientKey is the client key for the follower node to establish
	// client authentication during TLS.
	LeaderClientKey string `json:"leader_client_key"`

	// LeaderCACertFile is the path on disk to the the CA cert file of the
	// leader node. This should only be provided via Vault's configuration file.
	LeaderCACertFile string `json:"leader_ca_cert_file"`

	// LeaderClientCertFile is the path on disk to the client certificate file
	// for the follower node to establish client authentication during TLS. This
	// should only be provided via Vault's configuration file.
	LeaderClientCertFile string `json:"leader_client_cert_file"`

	// LeaderClientKeyFile is the path on disk to the client key file for the
	// follower node to establish client authentication during TLS. This should
	// only be provided via Vault's configuration file.
	LeaderClientKeyFile string `json:"leader_client_key_file"`

	// LeaderTLSServerName is the optional ServerName to expect in the leader's
	// certificate, instead of the host/IP we're actually connecting to.
	LeaderTLSServerName string `json:"leader_tls_servername"`

	// Retry indicates if the join process should automatically be retried
	Retry bool `json:"-"`

	var leaderInfos []*LeaderJoinInfo
	err := jsonutil.DecodeJSON([]byte(config), &leaderInfos)
	if err != nil {
		return nil, fmt.Errorf("failed to decode retry_join config: %w", err)
	}

	if len(leaderInfos) == 0 {
		return nil, errors.New("invalid retry_join config")
	}

	for i, info := range leaderInfos {
		if len(info.AutoJoin) != 0 && len(info.LeaderAPIAddr) != 0 {
			return nil, errors.New("cannot provide both a leader_api_addr and auto_join")
		}

		if info.AutoJoinScheme != "" && (info.AutoJoinScheme != "http" && info.AutoJoinScheme != "https") {
			return nil, fmt.Errorf("invalid scheme '%s'; must either be http or https", info.AutoJoinScheme)
		}

		info.Retry = true
		info.TLSConfig, err = parseTLSInfo(info)
		if err != nil {
			return nil, fmt.Errorf("failed to create tls config to communicate with leader node (retry_join index: %d): %w", i, err)
		}
	}

	return leaderInfos, nil
}

// parseTLSInfo is a helper for parses the TLS information, preferring file
// paths over raw certificate content.
func parseTLSInfo(leaderInfo *LeaderJoinInfo) (*tls.Config, error) {
	var tlsConfig *tls.Config
	var err error
	if len(leaderInfo.LeaderCACertFile) != 0 || len(leaderInfo.LeaderClientCertFile) != 0 || len(leaderInfo.LeaderClientKeyFile) != 0 {
		tlsConfig, err = tlsutil.LoadClientTLSConfig(leaderInfo.LeaderCACertFile, leaderInfo.LeaderClientCertFile, leaderInfo.LeaderClientKeyFile)
		if err != nil {
			return nil, err
		}
	} else if len(leaderInfo.LeaderCACert) != 0 || len(leaderInfo.LeaderClientCert) != 0 || len(leaderInfo.LeaderClientKey) != 0 {
		tlsConfig, err = tlsutil.ClientTLSConfig([]byte(leaderInfo.LeaderCACert), []byte(leaderInfo.LeaderClientCert), []byte(leaderInfo.LeaderClientKey))
		if err != nil {
			return nil, err
		}
	}
	if tlsConfig != nil {
		tlsConfig.ServerName = leaderInfo.LeaderTLSServerName
	}

	return tlsConfig, nil
}

// EnsurePath is used to make sure a path exists
func EnsurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}
	return os.MkdirAll(path, 0o755)
}

// NewRaftBackend constructs a RaftBackend using the given directory
func NewRaftBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	path := os.Getenv(EnvVaultRaftPath)
	if path == "" {
		pathFromConfig, ok := conf["path"]
		path = pathFromConfig
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

			if err := ioutil.WriteFile(filepath.Join(path, "node-id"), []byte(id), 0o600); err != nil {
				return nil, err
			}

			localID = id
		}
	}

	// Create the FSM.
	fsm, err := NewFSM(path, localID, logger.Named("fsm"))
	if err != nil {
		return nil, fmt.Errorf("failed to create fsm: %v", err)
	}

	if delayRaw, ok := conf["apply_delay"]; ok {
		delay, err := time.ParseDuration(delayRaw)
		if err != nil {
			return nil, fmt.Errorf("apply_delay does not parse as a duration: %w", err)
		}
		fsm.applyCallback = func() {
			time.Sleep(delay)
		}
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
		freelistType, noFreelistSync := freelistOptions()
		raftOptions := raftboltdb.Options{
			Path: filepath.Join(path, "raft.db"),
			BoltOptions: &bolt.Options{
				FreelistType:   freelistType,
				NoFreelistSync: noFreelistSync,
			},
		}
		store, err := raftboltdb.New(raftOptions)
		if err != nil {
			return nil, err
		}
		log = cacheStore

		// Create the snapshot store.
		snapshots, err := NewBoltSnapshotStore(path, logger.Named("snapshot"), fsm)
		if err != nil {
			return nil, err
		}
		snap = snapshots
	}

	if delayRaw, ok := conf["snapshot_delay"]; ok {
		delay, err := time.ParseDuration(delayRaw)
		if err != nil {
			return nil, fmt.Errorf("snapshot_delay does not parse as a duration: %w", err)
		}
		snap = newSnapshotStoreDelay(snap, delay, logger)
	}

	maxEntrySize := defaultMaxEntrySize
	if maxEntrySizeCfg := conf["max_entry_size"]; len(maxEntrySizeCfg) != 0 {
		i, err := strconv.Atoi(maxEntrySizeCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse 'max_entry_size': %w", err)
		}

		maxEntrySize = uint64(i)
	}

	var reconcileInterval time.Duration
	if interval := conf["autopilot_reconcile_interval"]; interval != "" {
		interval, err := time.ParseDuration(interval)
		if err != nil {
			return nil, fmt.Errorf("autopilot_reconcile_interval does not parse as a duration: %w", err)
		}
		reconcileInterval = interval
	}

	return &RaftBackend{
		logger:                     logger,
		fsm:                        fsm,
		raftInitCh:                 make(chan struct{}),
		conf:                       conf,
		logStore:                   log,
		stableStore:                stable,
		snapStore:                  snap,
		dataDir:                    path,
		localID:                    localID,
		permitPool:                 physical.NewPermitPool(physical.DefaultParallelOperations),
		maxEntrySize:               maxEntrySize,
		followerHeartbeatTicker:    time.NewTicker(time.Second),
		autopilotReconcileInterval: reconcileInterval,
	}, nil
}

type snapshotStoreDelay struct {
	logger  log.Logger
	wrapped raft.SnapshotStore
	delay   time.Duration
}

func (s snapshotStoreDelay) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
	s.logger.Trace("delaying before creating snapshot", "delay", s.delay)
	time.Sleep(s.delay)
	return s.wrapped.Create(version, index, term, configuration, configurationIndex, trans)
}

func (s snapshotStoreDelay) List() ([]*raft.SnapshotMeta, error) {
	return s.wrapped.List()
}

func (s snapshotStoreDelay) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	return s.wrapped.Open(id)
}

var _ raft.SnapshotStore = &snapshotStoreDelay{}

func newSnapshotStoreDelay(snap raft.SnapshotStore, delay time.Duration, logger log.Logger) *snapshotStoreDelay {
	return &snapshotStoreDelay{
		logger:  logger,
		wrapped: snap,
		delay:   delay,
	}
}

// Close is used to gracefully close all file resources.  N.B. This method
// should only be called if you are sure the RaftBackend will never be used
// again.
func (b *RaftBackend) Close() error {
	b.l.Lock()
	defer b.l.Unlock()

	if err := b.fsm.db.Close(); err != nil {
		return err
	}

	if err := b.stableStore.(*raftboltdb.BoltStore).Close(); err != nil {
		return err
	}

	return nil
}

func (b *RaftBackend) CollectMetrics(sink *metricsutil.ClusterMetricSink) {
	b.l.RLock()
	logstoreStats := b.stableStore.(*raftboltdb.BoltStore).Stats()
	fsmStats := b.fsm.db.Stats()
	b.l.RUnlock()
	b.collectMetricsWithStats(logstoreStats, sink, "logstore")
	b.collectMetricsWithStats(fsmStats, sink, "fsm")
}

func (b *RaftBackend) collectMetricsWithStats(stats bolt.Stats, sink *metricsutil.ClusterMetricSink, database string) {
	txstats := stats.TxStats
	labels := []metricsutil.Label{{"database", database}}
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "freelist", "free_pages"}, float32(stats.FreePageN), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "freelist", "pending_pages"}, float32(stats.PendingPageN), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "freelist", "allocated_bytes"}, float32(stats.FreeAlloc), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "freelist", "used_bytes"}, float32(stats.FreelistInuse), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "transaction", "started_read_transactions"}, float32(stats.TxN), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "transaction", "currently_open_read_transactions"}, float32(stats.OpenTxN), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "page", "count"}, float32(txstats.PageCount), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "page", "bytes_allocated"}, float32(txstats.PageAlloc), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "cursor", "count"}, float32(txstats.CursorCount), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "node", "count"}, float32(txstats.NodeCount), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "node", "dereferences"}, float32(txstats.NodeDeref), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "rebalance", "count"}, float32(txstats.Rebalance), labels)
	sink.AddSampleWithLabels([]string{"raft_storage", "bolt", "rebalance", "time"}, float32(txstats.RebalanceTime), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "split", "count"}, float32(txstats.Split), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "spill", "count"}, float32(txstats.Spill), labels)
	sink.AddSampleWithLabels([]string{"raft_storage", "bolt", "spill", "time"}, float32(txstats.SpillTime), labels)
	sink.SetGaugeWithLabels([]string{"raft_storage", "bolt", "write", "count"}, float32(txstats.Write), labels)
	sink.AddSampleWithLabels([]string{"raft_storage", "bolt", "write", "time"}, float32(txstats.WriteTime), labels)
}

// RaftServer has information about a server in the Raft configuration
type RaftServer struct {
	// NodeID is the name of the server

// Peer defines the ID and Address for a given member of the raft cluster.
type Peer struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Suffrage int    `json:"suffrage"`
}

// NodeID returns the identifier of the node
}

// Bootstrap prepares the given peers to be part of the raft cluster
func (b *RaftBackend) Bootstrap(peers []Peer) error {
	b.l.Lock()
	defer b.l.Unlock()


	for i, p := range peers {
		raftConfig.Servers[i] = raft.Server{
			ID:       raft.ServerID(p.ID),
			Address:  raft.ServerAddress(p.Address),
			Suffrage: raft.ServerSuffrage(p.Suffrage),
		}
	}

		}
		config.TrailingLogs = uint64(trailingLogs)
	}
	snapshotIntervalRaw, ok := b.conf["snapshot_interval"]
	if ok {
		var err error
		snapshotInterval, err := time.ParseDuration(snapshotIntervalRaw)
		if err != nil {
			return err
		}
		config.SnapshotInterval = snapshotInterval
	}

	config.NoSnapshotRestoreOnStart = true
	config.MaxAppendEntries = 64

	// Setting BatchApplyCh allows the raft library to enqueue up to
	// MaxAppendEntries into each raft apply rather than relying on the
	// scheduler.
	config.BatchApplyCh = true

	return nil
}

	})
}

func (b *RaftBackend) HasState() (bool, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	return raft.HasExistingState(b.logStore, b.stableStore, b.snapStore)
}

// SetupCluster starts the raft cluster and enables the networking needed for
// the raft nodes to communicate.
func (b *RaftBackend) SetupCluster(ctx context.Context, opts SetupOpts) error {
		return err
	}

	listenerIsNil := func(cl cluster.ClusterHook) bool {
		switch {
		case opts.ClusterListener == nil:
			return true
		default:
			// Concrete type checks
			switch cl.(type) {
			case *cluster.Listener:
				return cl.(*cluster.Listener) == nil
			}
		}
		return false
	}

	switch {
	case opts.TLSKeyring == nil && listenerIsNil(opts.ClusterListener):
		// If we don't have a provided network we use an in-memory one.
		// This allows us to bootstrap a node without bringing up a cluster
		// network. This will be true during bootstrap, tests and dev modes.
		_, b.raftTransport = raft.NewInmemTransportWithTimeout(raft.ServerAddress(b.localID), time.Second)
	case opts.TLSKeyring == nil:
		return errors.New("no keyring provided")
	case listenerIsNil(opts.ClusterListener):
		return errors.New("no cluster listener provided")
	default:
		// Set the local address and localID in the streaming layer and the raft config.
			MaxPool:               3,
			Timeout:               10 * time.Second,
			ServerAddressProvider: b.serverAddressProvider,
			Logger:                b.logger.Named("raft-net"),
		}
		transport := raft.NewNetworkTransportWithConfig(transConfig)

	raftConfig.LocalID = raft.ServerID(b.localID)

	// Set up a channel for reliable leader notifications.
	raftNotifyCh := make(chan bool, 10)
	raftConfig.NotifyCh = raftNotifyCh

	// If we have a bootstrapConfig set we should bootstrap now.
		if err := raft.BootstrapCluster(raftConfig, b.logStore, b.stableStore, b.snapStore, b.raftTransport, *bootstrapConfig); err != nil {
			return err
		}
	}

	// Setup the Raft store.
	b.fsm.SetNoopRestore(true)


		recoveryConfig, err := raft.ReadConfigJSON(peersFile)
		if err != nil {
			return fmt.Errorf("raft recovery failed to parse peers.json: %w", err)
		}

		// Non-voting servers are only allowed in enterprise. If Suffrage is disabled,
		// error out to indicate that it isn't allowed.
		for idx := range recoveryConfig.Servers {
			if !nonVotersAllowed && recoveryConfig.Servers[idx].Suffrage == raft.Nonvoter {
				return fmt.Errorf("raft recovery failed to parse configuration for node %q: setting `non_voter` is only supported in enterprise", recoveryConfig.Servers[idx].ID)
			}
		}

		b.logger.Info("raft recovery found new config", "config", recoveryConfig)

		err = raft.RecoverCluster(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport, recoveryConfig)
		if err != nil {
			return fmt.Errorf("raft recovery failed: %w", err)
		}

		err = os.Remove(peersFile)
		if err != nil {
			return fmt.Errorf("raft recovery failed to delete peers.json; please delete manually: %w", err)
		}
		b.logger.Info("raft recovery deleted peers.json")
	}
	if opts.RecoveryModeConfig != nil {
		err = raft.RecoverCluster(raftConfig, b.fsm, b.logStore, b.stableStore, b.snapStore, b.raftTransport, *opts.RecoveryModeConfig)
		if err != nil {
			return fmt.Errorf("recovering raft cluster failed: %w", err)
		}
	}

	b.logger.Info("creating Raft", "config", fmt.Sprintf("%#v", raftConfig))
	raftObj, err := raft.NewRaft(raftConfig, b.fsm.chunker, b.logStore, b.stableStore, b.snapStore, b.raftTransport)
	b.fsm.SetNoopRestore(false)
	if err != nil {
		return err
	}

	// If we are expecting to start as leader wait until we win the election.
	// This should happen quickly since there is only one node in the cluster.
	// StartAsLeader is only set during init, recovery mode, storage migration,
	// and tests.
	if opts.StartAsLeader {
		for {
			if raftObj.State() == raft.Leader {
				break
			}

			select {
			case <-ctx.Done():
				future := raftObj.Shutdown()
				if future.Error() != nil {
					return fmt.Errorf("shutdown while waiting for leadership: %w", future.Error())
				}

				return errors.New("shutdown while waiting for leadership")
			case <-time.After(10 * time.Millisecond):
			}
		}
	}

	b.raft = raftObj
	b.raftNotifyCh = raftNotifyCh

	if err := b.fsm.upgradeLocalNodeConfig(); err != nil {
		b.logger.Error("failed to upgrade local node configuration")
		return err
	}

	if b.streamLayer != nil {
		// Add Handler to the cluster.
		opts.ClusterListener.AddHandler(consts.RaftStorageALPN, b.streamLayer)
		opts.ClusterListener.AddClient(consts.RaftStorageALPN, b.streamLayer)
	}

	// Close the init channel to signal setup has been completed
	close(b.raftInitCh)

	b.logger.Trace("finished setting up raft cluster")
	return nil
}

	}

	b.l.Lock()

	// Perform shutdown only if the raft object is non-nil. The object could be nil
	// if the node is unsealed but has not joined the peer set.
	var future raft.Future
	if b.raft != nil {
		future = b.raft.Shutdown()
	}

	b.raft = nil

	// If we're tearing down, then we need to recreate the raftInitCh
	b.raftInitCh = make(chan struct{})
	b.l.Unlock()

	if future != nil {
		return future.Error()
	}

	return nil
}

// CommittedIndex returns the latest index committed to stable storage
func (b *RaftBackend) CommittedIndex() uint64 {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return 0
	}

	return b.raft.LastIndex()
}

// AppliedIndex returns the latest index applied to the FSM
	b.l.RLock()
	defer b.l.RUnlock()

	if b.fsm == nil {
		return 0
	}

	// We use the latest index that the FSM has seen here, which may be behind
	// raft.AppliedIndex() due to the async nature of the raft library.
	indexState, _ := b.fsm.LatestState()
	return indexState.Index
}

// Term returns the raft term of this node.
func (b *RaftBackend) Term() uint64 {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.fsm == nil {
		return 0
	}

	// We use the latest index that the FSM has seen here, which may be behind
	// raft.AppliedIndex() due to the async nature of the raft library.
	indexState, _ := b.fsm.LatestState()
	return indexState.Term
}

// RemovePeer removes the given peer ID from the raft cluster. If the node is
	b.l.RLock()
	defer b.l.RUnlock()

	if b.disableAutopilot {
		if b.raft == nil {
			return errors.New("raft storage is not initialized")
		}
		b.logger.Trace("removing server from raft", "id", peerID)
		future := b.raft.RemoveServer(raft.ServerID(peerID), 0, 0)
		return future.Error()
	}

	if b.autopilot == nil {
		return errors.New("raft storage autopilot is not initialized")
	}

	b.logger.Trace("removing server from raft via autopilot", "id", peerID)
	return b.autopilot.RemoveServer(raft.ServerID(peerID))
}

// GetConfigurationOffline is used to read the stale, last known raft
// configuration to this node. It accesses the last state written into the
// FSM. When a server is online use GetConfiguration instead.
func (b *RaftBackend) GetConfigurationOffline() (*RaftConfigurationResponse, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft != nil {
		return nil, errors.New("raft storage is initialized, used GetConfiguration instead")
	}

	if b.fsm == nil {
		return nil, nil
	}

	state, configuration := b.fsm.LatestState()
	config := &RaftConfigurationResponse{
		Index: state.Index,
	}

	if configuration == nil || configuration.Servers == nil {
		return config, nil
	}

	for _, server := range configuration.Servers {
		entry := &RaftServer{
			NodeID:  server.Id,
			Address: server.Address,
			// Since we are offline no node is the leader.
			Leader: false,
			Voter:  raft.ServerSuffrage(server.Suffrage) == raft.Voter,
		}
		config.Servers = append(config.Servers, entry)
	}

	return config, nil
}

func (b *RaftBackend) GetConfiguration(ctx context.Context) (*RaftConfigurationResponse, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.disableAutopilot {
		if b.raft == nil {
			return errors.New("raft storage is not initialized")
		}
		b.logger.Trace("adding server to raft", "id", peerID)
		future := b.raft.AddVoter(raft.ServerID(peerID), raft.ServerAddress(clusterAddr), 0, 0)
		return future.Error()
	}

	if b.autopilot == nil {
		return errors.New("raft storage autopilot is not initialized")
	}

	b.logger.Trace("adding server to raft via autopilot", "id", peerID)
	return b.autopilot.AddServer(&autopilot.Server{
		ID:          raft.ServerID(peerID),
		Name:        peerID,
		Address:     raft.ServerAddress(clusterAddr),
		RaftVersion: raft.ProtocolVersionMax,
		NodeType:    autopilot.NodeVoter,
	})
}

// Peers returns all the servers present in the raft cluster
	defer b.l.RUnlock()

	if b.raft == nil {
		return nil, errors.New("raft storage is not initialized")
	}

	future := b.raft.GetConfiguration()
	ret := make([]Peer, len(future.Configuration().Servers))
	for i, s := range future.Configuration().Servers {
		ret[i] = Peer{
			ID:       string(s.ID),
			Address:  string(s.Address),
			Suffrage: int(s.Suffrage),
		}
	}

	return ret, nil
}

// SnapshotHTTP is a wrapper for Snapshot that sends the snapshot as an HTTP
// response.
func (b *RaftBackend) SnapshotHTTP(out *logical.HTTPResponseWriter, access *seal.Access) error {
	out.Header().Add("Content-Disposition", "attachment")
	out.Header().Add("Content-Type", "application/gzip")

	return b.Snapshot(out, access)
}

// Snapshot takes a raft snapshot, packages it into a archive file and writes it
// to the provided writer. Seal access is used to encrypt the SHASUM file so we
// can validate the snapshot was taken using the same master keys or not.
func (b *RaftBackend) Snapshot(out io.Writer, access *seal.Access) error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.raft == nil {
		return errors.New("raft storage is sealed")
	}

	// If we have access to the seal create a sealer object
		}
	}

	return snapshot.Write(b.logger.Named("snapshot"), b.raft, s, out)
}

// WriteSnapshotToTemp reads a snapshot archive off the provided reader,

	var metadata raft.SnapshotMeta
	if b.raft == nil {
		return nil, nil, metadata, errors.New("raft storage is sealed")
	}

	// If we have access to the seal create a sealer object
	// snapshot applied to a quorum of nodes.
	command := &LogData{
		Operations: []*LogOperation{
			{
				OpType: restoreCallbackOp,
			},
		},
	}

	err := b.applyLog(ctx, command)

	// Do a best-effort attempt to let the standbys apply the restoreCallbackOp
	// before we continue.
	defer metrics.MeasureSince([]string{"raft-storage", "delete"}, time.Now())
	command := &LogData{
		Operations: []*LogOperation{
			{
				OpType: deleteOp,
				Key:    path,
			},
	b.permitPool.Acquire()
	defer b.permitPool.Release()

	entry, err := b.fsm.Get(ctx, path)
	if entry != nil {
		valueLen := len(entry.Value)
		if uint64(valueLen) > b.maxEntrySize {
			b.logger.Warn("retrieved entry value is too large, has raft's max_entry_size been reduced?",
				"size", valueLen, "max_entry_size", b.maxEntrySize)
		}
	}

	return entry, err
}

// Put inserts an entry in the log for the put operation. It will return an
// error if the resulting entry encoding exceeds the configured max_entry_size
// or if the call to applyLog fails.
func (b *RaftBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"raft-storage", "put"}, time.Now())
	command := &LogData{
		Operations: []*LogOperation{
			{
				OpType: putOp,
				Key:    entry.Key,
				Value:  entry.Value,
// persisted to the local FSM. Caller should hold the backend's read lock.
func (b *RaftBackend) applyLog(ctx context.Context, command *LogData) error {
	if b.raft == nil {
		return errors.New("raft storage is not initialized")
	}

	commandBytes, err := proto.Marshal(command)
		return err
	}

	cmdSize := len(commandBytes)
	if uint64(cmdSize) > b.maxEntrySize {
		return fmt.Errorf("%s; got %d bytes, max: %d bytes", physical.ErrValueTooLarge, cmdSize, b.maxEntrySize)
	}

	defer metrics.AddSample([]string{"raft-storage", "entry_size"}, float32(cmdSize))

	var chunked bool
	var applyFuture raft.ApplyFuture
	switch {
	return nil
}

// HAEnabled is the implementation of the HABackend interface
func (b *RaftBackend) HAEnabled() bool { return true }

// HAEnabled is the implementation of the HABackend interface
func (b *RaftBackend) LockWith(key, value string) (physical.Lock, error) {
	return &RaftLock{
		key:   key,
	}, nil
}

// SetDesiredSuffrage sets a field in the fsm indicating the suffrage intent for
// this node.
func (b *RaftBackend) SetDesiredSuffrage(nonVoter bool) error {
	b.l.Lock()
	defer b.l.Unlock()

	var desiredSuffrage string
	switch nonVoter {
	case true:
		desiredSuffrage = "non-voter"
	default:
		desiredSuffrage = "voter"
	}

	err := b.fsm.recordSuffrage(desiredSuffrage)
	if err != nil {
		return err
	}

	return nil
}

func (b *RaftBackend) DesiredSuffrage() string {
	return b.fsm.DesiredSuffrage()
}

// RaftLock implements the physical Lock interface and enables HA for this
// backend. The Lock uses the raftNotifyCh for receiving leadership edge
// triggers. Vault's active duty matches raft's leadership.
func (l *RaftLock) monitorLeadership(stopCh <-chan struct{}, leaderNotifyCh <-chan bool) <-chan struct{} {
	leaderLost := make(chan struct{})
	go func() {
		for {
			select {
			case isLeader := <-leaderNotifyCh:
				// leaderNotifyCh may deliver a true value initially if this
				// server is already the leader prior to RaftLock.Lock call
				// (the true message was already queued). The next message is
				// always going to be false. The for loop should loop at most
				// twice.
				if !isLeader {
					close(leaderLost)
					return
				}
			case <-stopCh:
				return
			}
		}
	}()
	return leaderLost
// Lock blocks until we become leader or are shutdown. It returns a channel that
// is closed when we detect a loss of leadership.
func (l *RaftLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	// If not initialized, block until it is
	if !l.b.Initialized() {
		select {
		case <-l.b.raftInitCh:
		case <-stopCh:
			return nil, nil
		}
	}

	l.b.l.RLock()

	// Ensure that we still have a raft instance after grabbing the read lock
	if l.b.raft == nil {
		l.b.l.RUnlock()
		return nil, errors.New("attempted to grab a lock on a nil raft backend")
	}

	// Cache the notifyCh locally
	leaderNotifyCh := l.b.raftNotifyCh

	// Check to see if we are already leader.
	if l.b.raft.State() == raft.Leader {
		err := l.b.applyLog(context.Background(), &LogData{
			Operations: []*LogOperation{
				{
					OpType: putOp,
					Key:    l.key,
					Value:  l.value,
				l.b.l.RLock()
				err := l.b.applyLog(context.Background(), &LogData{
					Operations: []*LogOperation{
						{
							OpType: putOp,
							Key:    l.key,
							Value:  l.value,
			return nil, nil
		}
	}
}

// Unlock gives up leadership.
func (l *RaftLock) Unlock() error {
	if l.b.raft == nil {
		return nil
	}

	return l.b.raft.LeadershipTransfer().Error()
}


	return s.access.Decrypt(ctx, &eblob, nil)
}

// freelistOptions returns the freelist type and nofreelistsync values to use
// when opening boltdb files, based on our preferred defaults, and the possible
// presence of overriding environment variables.
func freelistOptions() (bolt.FreelistType, bool) {
	freelistType := bolt.FreelistMapType
	noFreelistSync := true

	if os.Getenv("VAULT_RAFT_FREELIST_TYPE") == "array" {
		freelistType = bolt.FreelistArrayType
	}

	if os.Getenv("VAULT_RAFT_FREELIST_SYNC") != "" {
		noFreelistSync = false
	}

	return freelistType, noFreelistSync
}
