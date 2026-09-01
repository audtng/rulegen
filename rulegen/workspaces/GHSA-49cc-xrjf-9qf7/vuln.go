package main


type eventActionPage struct {
	basePage
	Action         dataprovider.BaseEventAction
	ActionTypes    []dataprovider.EnumMapping
	FsActions      []dataprovider.EnumMapping
	HTTPMethods    []string
	RedactedSecret string
	Error          *util.I18nError
	Mode           genericPageMode
}

type eventRulePage struct {
	}

	data := eventActionPage{
		basePage:       s.getBasePageData(title, currentURL, w, r),
		Action:         action,
		ActionTypes:    dataprovider.EventActionTypes,
		FsActions:      dataprovider.FsActionTypes,
		HTTPMethods:    dataprovider.SupportedHTTPActionMethods,
		RedactedSecret: redactedSecret,
		Error:          getI18nError(err),
		Mode:           mode,
	}
	renderAdminTemplate(w, templateEventAction, data)
}
	if err := c.initializeProxyProtocol(); err != nil {
		return err
	}
	vfs.SetTempPath(c.TempPath)
	dataprovider.SetTempPath(c.TempPath)
	vfs.SetAllowSelfConnections(c.AllowSelfConnections)
	vfs.SetResumeMaxSize(c.ResumeMaxSize)
	vfs.SetUploadMode(c.UploadMode)
	dataprovider.SetAllowSelfConnections(c.AllowSelfConnections)
	transfersChecker = getTransfersChecker(isShared)
	return nil
}
	DLSize        int64  `json:"-"`
}

// MetadataConfig defines how to handle metadata for cloud storage backends
type MetadataConfig struct {
	// If not zero the metadata will be read before downloads and will be
	// server's local time, otherwise UTC will be used.
	TZ string `json:"tz" mapstructure:"tz"`
	// Metadata configuration
	Metadata              MetadataConfig `json:"metadata" mapstructure:"metadata"`
	idleTimeoutAsDuration time.Duration
	idleLoginTimeout      time.Duration
	defender              Defender
		ActionTypeBackup, ActionTypeUserQuotaReset, ActionTypeFolderQuotaReset, ActionTypeTransferQuotaReset,
		ActionTypeDataRetentionCheck, ActionTypePasswordExpirationCheck, ActionTypeUserExpirationCheck,
		ActionTypeUserInactivityCheck, ActionTypeIDPAccountCheck, ActionTypeRotateLogs}
)

func isActionTypeValid(action int) bool {
	return client
}

// EventActionCommandConfig defines the configuration for a command event target
type EventActionCommandConfig struct {
	Cmd     string     `json:"cmd,omitempty"`
	if c.Cmd == "" {
		return util.NewI18nError(util.NewValidationError("command is required"), util.I18nErrorCommandRequired)
	}
	if !filepath.IsAbs(c.Cmd) {
		return util.NewI18nError(
			util.NewValidationError("invalid command, it must be an absolute path"),
