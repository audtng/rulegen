package main

	// Apply memory configuration FIRST (before heavy allocations)
	// First, try to load from config file, then fall back to environment variables
	var cfg *config.Config
	explicitConfigPath, _ := cmd.Flags().GetString("config") // persistent
	configPath := strings.TrimSpace(explicitConfigPath)
	if configPath == "" {
			fmt.Printf("⚠️  Warning: failed to load config from %s: %v\n", configPath, err)
			cfg = config.LoadFromEnv()
		} else {
			fmt.Printf("📄 Loaded config from: %s\n", configPath)
		}
	}

	// YAML config file is the source of truth for embedding settings
	// Always use config file values if they are set (non-zero/non-empty)
	if cfg.Memory.EmbeddingDimensions > 0 {
	// Create and start HTTP server
	serverConfig := server.DefaultConfig()
	serverConfig.Port = httpPort
	serverConfig.Address = address
	// MCP server configuration
	serverConfig.MCPEnabled = mcpEnabled
	// Pass embedding settings to server (from loaded config)

	// Create and start Bolt server for Neo4j driver compatibility
	boltConfig := bolt.DefaultConfig()
	boltConfig.Port = boltPort
	boltConfig.LogQueries = logQueries
	boltConfig.ServerAnnouncement = cfg.Server.BoltServerAnnouncement
	fmt.Println("✅ NornicDB is ready!")
	fmt.Println()
	// Determine the display address for user-friendly output
	displayAddr := address
	if address == "0.0.0.0" {
		displayAddr = "localhost" // 0.0.0.0 is all interfaces, show localhost for convenience
	}
	fmt.Println("Endpoints:")
	return nil
}

func startStdioLogCompactor(maxKB int, interval time.Duration) func() {
	if maxKB <= 0 {
		return func() {}
//	config = bolt.DefaultConfig()
//	config.Port = 7688 // Use different port
type Config struct {
	Port            int
	MaxConnections  int
	ReadBufferSize  int
//	server := bolt.New(config, executor)
func DefaultConfig() *Config {
	return &Config{
		Port:            7687,
		MaxConnections:  100,
		ReadBufferSize:  8192,
//
// The server will print its listening address when started successfully.
func (s *Server) ListenAndServe() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener

	fmt.Printf("Bolt server listening on bolt://localhost:%d\n", s.config.Port)

	return s.serve()
}
