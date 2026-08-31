package main

}

type SecurityConfig struct {
	Enabled          bool          `yaml:"enabled"`
	AllowedCommands  []string      `yaml:"allowed_commands"`      // Deprecated: use AllowedExecutables
	BlockedCommands  []string      `yaml:"blocked_commands"`      // Deprecated: use validation instead
	BlockedPatterns  []string      `yaml:"blocked_patterns"`      // Deprecated: use validation instead
	AllowedExecutables []string    `yaml:"allowed_executables"`   // Secure: list of allowed executable paths
	MaxExecutionTime time.Duration `yaml:"max_execution_time"`
	WorkingDirectory string        `yaml:"working_directory"`
	RunAsUser        string        `yaml:"run_as_user"`
	MaxOutputSize    int           `yaml:"max_output_size"`
	AuditLog         bool          `yaml:"audit_log"`
	UseShellExecution bool         `yaml:"use_shell_execution"`   // Legacy mode - enables shell execution (DANGEROUS)
}

type ServerConfig struct {
	Output string
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Security: SecurityConfig{
			Enabled: false,
		},
		Server: ServerConfig{
			Name:    getEnv("MCP_SHELL_SERVER_NAME", "mcp-shell 🐚"),
			Version: version,

	var yamlConfig struct {
		Security struct {
			Enabled          bool     `yaml:"enabled"`
			AllowedCommands  []string `yaml:"allowed_commands"`
			BlockedCommands  []string `yaml:"blocked_commands"`
			BlockedPatterns  []string `yaml:"blocked_patterns"`
			AllowedExecutables []string `yaml:"allowed_executables"`
			MaxExecutionTime string   `yaml:"max_execution_time"`
			WorkingDirectory string   `yaml:"working_directory"`
			RunAsUser        string   `yaml:"run_as_user"`
			MaxOutputSize    int      `yaml:"max_output_size"`
			AuditLog         bool     `yaml:"audit_log"`
			UseShellExecution bool    `yaml:"use_shell_execution"`
		} `yaml:"security"`
	}

}

func newSecurityValidator(cfg SecurityConfig, logger zerolog.Logger) *SecurityValidator {
	return &SecurityValidator{
		config: cfg,
		logger: logger.With().Str("component", "security").Logger(),
	}
}

func (v *SecurityValidator) validateCommand(command string) error {
// containsShellMetacharacters checks if a string contains shell metacharacters
// that could be used for command injection
func containsShellMetacharacters(s string) bool {
	metachars := "|&;<>(){}[]$`\\"
	for _, char := range s {
		if strings.ContainsRune(metachars, char) {
			return true
// containsDangerousShellConstructs checks for potentially dangerous shell constructs
func containsDangerousShellConstructs(s string) bool {
	dangerous := []string{
		"$(", "`", "${", "&&", "||", ";", "|", ">", "<", ">>", "<<", "&",
	}
	for _, construct := range dangerous {
		if strings.Contains(s, construct) {
	}

	configFile := os.Getenv("MCP_SHELL_SEC_CONFIG_FILE")
	if configFile != "" {
		log.Info().Str("config_file", configFile).Msg("Loading security config")
	} else {
		log.Info().Msg("No security config file specified, security disabled")
	}

	log.Info().
