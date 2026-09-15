package main

	// Apply memory configuration FIRST (before heavy allocations)
	// First, try to load from config file, then fall back to environment variables
	var cfg *config.Config
	loadedConfigFile := false
	explicitConfigPath, _ := cmd.Flags().GetString("config") // persistent
	configPath := strings.TrimSpace(explicitConfigPath)
	if configPath == "" {
			fmt.Printf("⚠️  Warning: failed to load config from %s: %v\n", configPath, err)
			cfg = config.LoadFromEnv()
		} else {
			loadedConfigFile = true
			fmt.Printf("📄 Loaded config from: %s\n", configPath)
		}
	}

	resolvedAddress := resolveBindAddress(cmd, cfg, address, loadedConfigFile)
	cfg.Server.HTTPAddress = resolvedAddress
	cfg.Server.BoltAddress = resolvedAddress

	// YAML config file is the source of truth for embedding settings
	// Always use config file values if they are set (non-zero/non-empty)
	if cfg.Memory.EmbeddingDimensions > 0 {
	// Create and start HTTP server
	serverConfig := server.DefaultConfig()
	serverConfig.Port = httpPort
	serverConfig.Address = resolvedAddress
	// MCP server configuration
	serverConfig.MCPEnabled = mcpEnabled
	// Pass embedding settings to server (from loaded config)

	// Create and start Bolt server for Neo4j driver compatibility
	boltConfig := bolt.DefaultConfig()
	boltConfig.Host = resolvedAddress
	boltConfig.Port = boltPort
	boltConfig.LogQueries = logQueries
	boltConfig.ServerAnnouncement = cfg.Server.BoltServerAnnouncement
	fmt.Println("✅ NornicDB is ready!")
	fmt.Println()
	// Determine the display address for user-friendly output
	displayAddr := resolvedAddress
	if resolvedAddress == "0.0.0.0" || resolvedAddress == "::" {
		displayAddr = "localhost" // 0.0.0.0 is all interfaces, show localhost for convenience
	}
	fmt.Println("Endpoints:")
	return nil
}

func resolveBindAddress(cmd *cobra.Command, cfg *config.Config, cliAddress string, loadedConfigFile bool) string {
	resolvedAddress := strings.TrimSpace(cliAddress)
	if cmd != nil && !cmd.Flags().Changed("address") && cfg != nil {
		if loadedConfigFile && cfg.Server.HTTPAddress != "" {
			resolvedAddress = cfg.Server.HTTPAddress
		} else if loadedConfigFile && cfg.Server.BoltAddress != "" {
			resolvedAddress = cfg.Server.BoltAddress
		} else if hasExplicitProtocolBindEnv() {
			if cfg.Server.HTTPAddress != "" {
				resolvedAddress = cfg.Server.HTTPAddress
			} else if cfg.Server.BoltAddress != "" {
				resolvedAddress = cfg.Server.BoltAddress
			}
		}
	}
	if strings.TrimSpace(resolvedAddress) == "" {
		return "127.0.0.1"
	}
	return strings.TrimSpace(resolvedAddress)
}

func hasExplicitProtocolBindEnv() bool {
	for _, envName := range []string{
		"NORNICDB_BOLT_ADDRESS",
		"NORNICDB_HTTP_ADDRESS",
		"NEO4J_dbms_connector_bolt_listen__address",
		"NEO4J_dbms_connector_http_listen__address",
	} {
		if strings.TrimSpace(os.Getenv(envName)) != "" {
			return true
		}
	}
	return false
}

func startStdioLogCompactor(maxKB int, interval time.Duration) func() {
	if maxKB <= 0 {
		return func() {}
//	config = bolt.DefaultConfig()
//	config.Port = 7688 // Use different port
type Config struct {
	Host            string
	Port            int
	MaxConnections  int
	ReadBufferSize  int
//	server := bolt.New(config, executor)
func DefaultConfig() *Config {
	return &Config{
		Host:            "127.0.0.1",
		Port:            7687,
		MaxConnections:  100,
		ReadBufferSize:  8192,
//
// The server will print its listening address when started successfully.
func (s *Server) ListenAndServe() error {
	host := strings.TrimSpace(s.config.Host)
	if host == "" {
		host = "127.0.0.1"
	}
	addr := net.JoinHostPort(host, strconv.Itoa(s.config.Port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener

	announceHost := host
	if host == "0.0.0.0" || host == "::" || host == "" {
		announceHost = "localhost"
	}
	actualPort := s.config.Port
	if tcpAddr, ok := listener.Addr().(*net.TCPAddr); ok && tcpAddr.Port > 0 {
		actualPort = tcpAddr.Port
	}
	fmt.Printf("Bolt server listening on bolt://%s:%d\n", announceHost, actualPort)

	return s.serve()
}
