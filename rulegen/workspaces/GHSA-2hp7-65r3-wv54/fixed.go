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
package main

import (
	"testing"

	"github.com/orneryd/nornicdb/pkg/config"
	"github.com/spf13/cobra"
)

func TestResolveBindAddress(t *testing.T) {
	t.Run("uses_cli_address_when_flag_changed", func(t *testing.T) {
		cfg := config.LoadDefaults()
		cfg.Server.BoltAddress = "0.0.0.0"
		cfg.Server.HTTPAddress = "0.0.0.0"

		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("address", "127.0.0.1", "")
		if err := cmd.Flags().Set("address", "127.0.0.1"); err != nil {
			t.Fatalf("set address flag: %v", err)
		}

		resolved := resolveBindAddress(cmd, cfg, "127.0.0.1", false)
		if resolved != "127.0.0.1" {
			t.Fatalf("expected CLI address to win, got %q", resolved)
		}
	})

	t.Run("uses_loaded_server_address_when_config_file_sets_host", func(t *testing.T) {
		cfg := config.LoadDefaults()
		cfg.Server.BoltAddress = "127.0.0.2"
		cfg.Server.HTTPAddress = "127.0.0.2"

		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("address", "127.0.0.1", "")

		resolved := resolveBindAddress(cmd, cfg, "127.0.0.1", true)
		if resolved != "127.0.0.2" {
			t.Fatalf("expected loaded config address, got %q", resolved)
		}
	})

	t.Run("keeps_secure_default_when_no_explicit_config_exists", func(t *testing.T) {
		cfg := config.LoadDefaults()

		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("address", "127.0.0.1", "")

		resolved := resolveBindAddress(cmd, cfg, "127.0.0.1", false)
		if resolved != "127.0.0.1" {
			t.Fatalf("expected loopback CLI default, got %q", resolved)
		}
	})

	t.Run("falls_back_to_protocol_address_when_env_explicitly_sets_it", func(t *testing.T) {
		cfg := config.LoadDefaults()
		cfg.Server.HTTPAddress = ""
		cfg.Server.BoltAddress = "127.0.0.2"
		t.Setenv("NORNICDB_BOLT_ADDRESS", "127.0.0.2")

		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("address", "127.0.0.1", "")

		resolved := resolveBindAddress(cmd, cfg, "127.0.0.1", false)
		if resolved != "127.0.0.2" {
			t.Fatalf("expected Bolt address fallback, got %q", resolved)
		}
	})

	t.Run("defaults_to_loopback_when_empty", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("address", "", "")

		resolved := resolveBindAddress(cmd, nil, "", false)
		if resolved != "127.0.0.1" {
			t.Fatalf("expected loopback default, got %q", resolved)
		}
	})
}
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
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %q", config.Host)
	}
	if config.Port != 7687 {
		t.Errorf("expected port 7687, got %d", config.Port)
	}
		}
	})

	t.Run("binds_configured_host", func(t *testing.T) {
		config := &Config{Host: "127.0.0.1", Port: 0, MaxConnections: 10}
		server := New(config, &mockExecutor{})

		done := make(chan error, 1)
		go func() {
			done <- server.ListenAndServe()
		}()

		time.Sleep(50 * time.Millisecond)

		tcpAddr, ok := server.listener.Addr().(*net.TCPAddr)
		if !ok {
			t.Fatalf("expected TCP listener address, got %T", server.listener.Addr())
		}
		if !tcpAddr.IP.IsLoopback() {
			t.Fatalf("expected loopback bind address, got %v", tcpAddr.IP)
		}

		if err := server.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}

		select {
		case <-done:
		case <-time.After(500 * time.Millisecond):
			t.Fatal("server did not shut down")
		}
	})

	t.Run("listen_error", func(t *testing.T) {
		// Try to listen on an invalid port
		config := &Config{Port: -1}
