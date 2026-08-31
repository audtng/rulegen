package main


// SlackImporter is a service that allows to import slack dumps into mattermost
type SlackImporter struct {
	store   store.Store
	actions Actions
	config  *model.Config
}

// New creates a new SlackImporter service instance. It receive a store, a set of actions and the current config.
	}
}

func (si *SlackImporter) SlackImport(rctx request.CTX, fileData multipart.File, fileSize int64, teamID string) (*model.AppError, *bytes.Buffer) {
	// Create log file
	log := bytes.NewBufferString(i18n.T("api.slackimport.slack_import.log"))
		return nil
	}

	if _, err := si.store.User().VerifyEmail(ruser.Id, ruser.Email); err != nil {
		rctx.Logger().Warn("Failed to set email verified.", mlog.Err(err))
	}

	if _, err := si.actions.JoinUserToTeam(team, user, ""); err != nil {
		},
	}

	importer := slackimport.New(a.Srv().Store(), actions, a.Config())
	return importer.SlackImport(c, fileData, fileSize, teamID)
}

