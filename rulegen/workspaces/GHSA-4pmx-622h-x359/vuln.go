package main

	return errors.Wrapf(ErrNoPermissions, "user `%s` does not have permission to make playbook `%s` public", userID, playbook.ID)
}

func (p *PermissionsService) RunCreate(userID string, playbook Playbook) error {
	if p.hasPermissionsToPlaybook(userID, playbook, model.PermissionRunCreate) {
		return nil
	}

	return errors.Wrapf(ErrNoPermissions, "user `%s` does not have permission to run playbook `%s`", userID, playbook.ID)
}

func (p *PermissionsService) RunManageProperties(userID, runID string) error {
			return nil, errors.New("playbook is archived, cannot create a new run using an archived playbook")
		}

		if err = h.permissions.RunCreate(userID, *playbook); err != nil {
			return nil, err
		}


	filteredPlaybooks := make([]Playbook, 0, len(playbooks))
	for _, playbook := range playbooks {
		if err := s.permissions.RunCreate(requesterID, playbook); err == nil {
			filteredPlaybooks = append(filteredPlaybooks, playbook)
		}
	}
