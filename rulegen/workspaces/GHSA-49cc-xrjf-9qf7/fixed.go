package main


type eventActionPage struct {
	basePage
	Action          dataprovider.BaseEventAction
	ActionTypes     []dataprovider.EnumMapping
	FsActions       []dataprovider.EnumMapping
	HTTPMethods     []string
	EnabledCommands []string
	RedactedSecret  string
	Error           *util.I18nError
	Mode            genericPageMode
}

type eventRulePage struct {
	}

	data := eventActionPage{
		basePage:        s.getBasePageData(title, currentURL, w, r),
		Action:          action,
		ActionTypes:     dataprovider.EventActionTypes,
		FsActions:       dataprovider.FsActionTypes,
		HTTPMethods:     dataprovider.SupportedHTTPActionMethods,
		EnabledCommands: dataprovider.EnabledActionCommands,
		RedactedSecret:  redactedSecret,
		Error:           getI18nError(err),
		Mode:            mode,
	}
	renderAdminTemplate(w, templateEventAction, data)
}
	if err := c.initializeProxyProtocol(); err != nil {
		return err
	}
	if err := c.EventManager.validate(); err != nil {
		return err
	}
	vfs.SetTempPath(c.TempPath)
	dataprovider.SetTempPath(c.TempPath)
	vfs.SetAllowSelfConnections(c.AllowSelfConnections)
	vfs.SetResumeMaxSize(c.ResumeMaxSize)
	vfs.SetUploadMode(c.UploadMode)
	dataprovider.SetAllowSelfConnections(c.AllowSelfConnections)
	dataprovider.EnabledActionCommands = c.EventManager.EnabledCommands
	transfersChecker = getTransfersChecker(isShared)
	return nil
}
	DLSize        int64  `json:"-"`
}

// EventManagerConfig defines the configuration for the EventManager
type EventManagerConfig struct {
	// EnabledCommands defines the system commands that can be executed via EventManager,
	// an empty list means that any command is allowed to be executed.
	// Commands must be set as an absolute path
	EnabledCommands []string `json:"enabled_commands" mapstructure:"enabled_commands"`
}

func (c *EventManagerConfig) validate() error {
	for _, c := range c.EnabledCommands {
		if !filepath.IsAbs(c) {
			return fmt.Errorf("invalid command %q: it must be an absolute path", c)
		}
	}
	return nil
}

// MetadataConfig defines how to handle metadata for cloud storage backends
type MetadataConfig struct {
	// If not zero the metadata will be read before downloads and will be
	// server's local time, otherwise UTC will be used.
	TZ string `json:"tz" mapstructure:"tz"`
	// Metadata configuration
	Metadata MetadataConfig `json:"metadata" mapstructure:"metadata"`
	// EventManager configuration
	EventManager          EventManagerConfig `json:"event_manager" mapstructure:"event_manager"`
	idleTimeoutAsDuration time.Duration
	idleLoginTimeout      time.Duration
	defender              Defender
		ActionTypeBackup, ActionTypeUserQuotaReset, ActionTypeFolderQuotaReset, ActionTypeTransferQuotaReset,
		ActionTypeDataRetentionCheck, ActionTypePasswordExpirationCheck, ActionTypeUserExpirationCheck,
		ActionTypeUserInactivityCheck, ActionTypeIDPAccountCheck, ActionTypeRotateLogs}
	// EnabledActionCommands defines the system commands that can be executed via EventManager,
	// an empty list means that any command is allowed to be executed.
	EnabledActionCommands []string
)

func isActionTypeValid(action int) bool {
	return client
}

// IsActionCommandAllowed returns true if the specified command is allowed
func IsActionCommandAllowed(cmd string) bool {
	if len(EnabledActionCommands) == 0 {
		return true
	}
	return slices.Contains(EnabledActionCommands, cmd)
}

// EventActionCommandConfig defines the configuration for a command event target
type EventActionCommandConfig struct {
	Cmd     string     `json:"cmd,omitempty"`
	if c.Cmd == "" {
		return util.NewI18nError(util.NewValidationError("command is required"), util.I18nErrorCommandRequired)
	}
	if !IsActionCommandAllowed(c.Cmd) {
		return util.NewValidationError(fmt.Sprintf("command %q is not allowed", c.Cmd))
	}
	if !filepath.IsAbs(c.Cmd) {
		return util.NewI18nError(
			util.NewValidationError("invalid command, it must be an absolute path"),
