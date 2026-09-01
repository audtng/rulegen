package main

	return errors.Wrapf(ErrNoPermissions, "user `%s` does not have permission to make playbook `%s` public", userID, playbook.ID)
}

func (p *PermissionsService) RunCreate(userID string, playbook Playbook, targetTeamID string) error {
	if !p.hasPermissionsToPlaybook(userID, playbook, model.PermissionRunCreate) {
		return errors.Wrapf(ErrNoPermissions, "user `%s` does not have permission to run playbook `%s`", userID, playbook.ID)
	}

	if targetTeamID != "" && targetTeamID != playbook.TeamID {
		if !p.pluginAPI.User.HasPermissionToTeam(userID, targetTeamID, model.PermissionRunCreate) {
			return errors.Wrapf(ErrNoPermissions, "user `%s` does not have permission to create a run in team `%s`", userID, targetTeamID)
		}
	}

	return nil
}

func (p *PermissionsService) RunManageProperties(userID, runID string) error {
			return nil, errors.New("playbook is archived, cannot create a new run using an archived playbook")
		}

		if err = h.permissions.RunCreate(userID, *playbook, playbookRun.TeamID); err != nil {
			return nil, err
		}


	filteredPlaybooks := make([]Playbook, 0, len(playbooks))
	for _, playbook := range playbooks {
		if err := s.permissions.RunCreate(requesterID, playbook, ""); err == nil {
			filteredPlaybooks = append(filteredPlaybooks, playbook)
		}
	}
