package main


// SlackImporter is a service that allows to import slack dumps into mattermost
type SlackImporter struct {
	store         store.Store
	actions       Actions
	config        *model.Config
	isAdminImport bool
}

// New creates a new SlackImporter service instance. It receive a store, a set of actions and the current config.
	}
}

// NewWithAdminFlag creates a new SlackImporter service instance with information about whether this is an admin import.
// This allows for enhanced security controls based on the importing user's role.
func NewWithAdminFlag(store store.Store, actions Actions, config *model.Config, isAdminImport bool) *SlackImporter {
	return &SlackImporter{
		store:         store,
		actions:       actions,
		config:        config,
		isAdminImport: isAdminImport,
	}
}

func (si *SlackImporter) SlackImport(rctx request.CTX, fileData multipart.File, fileSize int64, teamID string) (*model.AppError, *bytes.Buffer) {
	// Create log file
	log := bytes.NewBufferString(i18n.T("api.slackimport.slack_import.log"))
		return nil
	}

	// Only system admins can automatically verify emails during import
	if si.isAdminImport {
		if _, err := si.store.User().VerifyEmail(ruser.Id, ruser.Email); err != nil {
			rctx.Logger().Warn("Failed to set email verified for admin import.", mlog.Err(err))
		}
	} else {
		// Non-admin users: emails remain unverified
		rctx.Logger().Debug("Email verification skipped for non-admin import.",
			mlog.String("user_email", ruser.Email))
	}

	if _, err := si.actions.JoinUserToTeam(team, user, ""); err != nil {
		},
	}

	// Determine if this is an Admin import:
	// mattermost cmd imports (no session) are treated as admin imports since only server admins can run them
	// Web imports (include mmctl calls) check the actual user's role
	isAdminImport := false

	if c.Session() == nil {
		// no session means it's being run directly on the server and only
		// server admins can run CLI commands, so treat as admin import
		isAdminImport = true
		c.Logger().Info("Slack import initiated via CLI, treating as admin import")
	} else if c.Session().UserId != "" {
		// Web API + mmctl import - check if the user is a system admin
		if user, err := a.GetUser(c.Session().UserId); err == nil {
			isAdminImport = user.IsSystemAdmin()
		}
	}

	importer := slackimport.NewWithAdminFlag(a.Srv().Store(), actions, a.Config(), isAdminImport)
	return importer.SlackImport(c, fileData, fileSize, teamID)
}

