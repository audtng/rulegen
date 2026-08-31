package main

}

type SecurityConfig struct {
	Enabled            bool          `yaml:"enabled"`
	AllowedCommands    []string      `yaml:"allowed_commands"`    // Deprecated: use AllowedExecutables
	BlockedCommands    []string      `yaml:"blocked_commands"`    // Deprecated: use validation instead
	BlockedPatterns    []string      `yaml:"blocked_patterns"`    // Deprecated: use validation instead
	AllowedExecutables []string      `yaml:"allowed_executables"` // Secure: list of allowed executable paths
	MaxExecutionTime   time.Duration `yaml:"max_execution_time"`
	WorkingDirectory   string        `yaml:"working_directory"`
	RunAsUser          string        `yaml:"run_as_user"`
	MaxOutputSize      int           `yaml:"max_output_size"`
	AuditLog           bool          `yaml:"audit_log"`
	UseShellExecution  bool          `yaml:"use_shell_execution"` // Legacy mode - enables shell execution (DANGEROUS)
}

type ServerConfig struct {
	Output string
}

// newDefaultSecurityConfig returns the built-in secure defaults applied when no
// MCP_SHELL_SEC_CONFIG_FILE is provided. Secure mode is the operating default:
// the server boots restricted to a narrow allowlist of utilities that cannot
// themselves spawn arbitrary processes. Shell/language interpreters (bash, sh,
// python, perl, ruby, git) are intentionally excluded - allowing one is
// equivalent to disabling the allowlist entirely.
func newDefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		Enabled:           true,
		UseShellExecution: false,
		AllowedExecutables: []string{
			"ls", "pwd", "whoami", "date", "echo",
			"cat", "grep", "find", "wc", "head", "tail", "sort", "uniq",
		},
		MaxExecutionTime: 30 * time.Second,
		MaxOutputSize:    1048576,
		WorkingDirectory: "/tmp",
		AuditLog:         true,
	}
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load()

	security := newDefaultSecurityConfig()
	// Unrestricted mode requires affirmative opt-in, never silence.
	if getBoolEnv("MCP_SHELL_ALLOW_UNSAFE", false) {
		security.Enabled = false
	}

	config := &Config{
		Security: security,
		Server: ServerConfig{
			Name:    getEnv("MCP_SHELL_SERVER_NAME", "mcp-shell 🐚"),
			Version: version,

	var yamlConfig struct {
		Security struct {
			Enabled            bool     `yaml:"enabled"`
			AllowedCommands    []string `yaml:"allowed_commands"`
			BlockedCommands    []string `yaml:"blocked_commands"`
			BlockedPatterns    []string `yaml:"blocked_patterns"`
			AllowedExecutables []string `yaml:"allowed_executables"`
			MaxExecutionTime   string   `yaml:"max_execution_time"`
			WorkingDirectory   string   `yaml:"working_directory"`
			RunAsUser          string   `yaml:"run_as_user"`
			MaxOutputSize      int      `yaml:"max_output_size"`
			AuditLog           bool     `yaml:"audit_log"`
			UseShellExecution  bool     `yaml:"use_shell_execution"`
		} `yaml:"security"`
	}

}

func newSecurityValidator(cfg SecurityConfig, logger zerolog.Logger) *SecurityValidator {
	v := &SecurityValidator{
		config: cfg,
		logger: logger.With().Str("component", "security").Logger(),
	}
	v.warnOnInterpreters()
	return v
}

// warnOnInterpreters flags allowlisted executables that can execute arbitrary
// commands regardless of metacharacter checks (shell/language interpreters, and
// git via `-c alias.x=!cmd`). Allowing one defeats secure mode; the warning
// surfaces the misconfiguration at startup instead of silently trusting it.
func (v *SecurityValidator) warnOnInterpreters() {
	if !v.config.Enabled || v.config.UseShellExecution {
		return
	}
	for _, exe := range v.config.AllowedExecutables {
		if isInterpreterExecutable(filepath.Base(exe)) {
			v.logger.Warn().
				Str("executable", exe).
				Msg("allowed executable can run arbitrary commands (interpreter or alias-capable) - this defeats secure mode regardless of metacharacter filtering")
		}
	}
}

// isInterpreterExecutable reports whether base names an executable that can
// itself run arbitrary commands, making executable-allowlisting ineffective.
func isInterpreterExecutable(base string) bool {
	switch base {
	case "bash", "sh", "zsh", "fish", "dash", "ksh", "csh", "tcsh",
		"python", "python2", "python3", "perl", "ruby", "node", "php", "git":
		return true
	}
	return false
}

func (v *SecurityValidator) validateCommand(command string) error {
// containsShellMetacharacters checks if a string contains shell metacharacters
// that could be used for command injection
func containsShellMetacharacters(s string) bool {
	// '!' is included because git interprets `-c alias.x=!cmd` as a shell alias,
	// turning an allowlisted git into arbitrary command execution.
	metachars := "|&;<>(){}[]$`\\!"
	for _, char := range s {
		if strings.ContainsRune(metachars, char) {
			return true
// containsDangerousShellConstructs checks for potentially dangerous shell constructs
func containsDangerousShellConstructs(s string) bool {
	dangerous := []string{
		"$(", "`", "${", "&&", "||", ";", "|", ">", "<", ">>", "<<", "&", "=!",
	}
	for _, construct := range dangerous {
		if strings.Contains(s, construct) {
	}

	configFile := os.Getenv("MCP_SHELL_SEC_CONFIG_FILE")
	switch {
	case configFile != "":
		log.Info().Str("config_file", configFile).Msg("Loading security config from file")
	case cfg.Security.Enabled:
		log.Info().Msg("No security config file specified, using built-in secure defaults")
	default:
		log.Warn().Msg("SECURITY DISABLED via MCP_SHELL_ALLOW_UNSAFE - all commands run unrestricted")
	}

	log.Info().
