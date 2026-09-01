package main

		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), c.Params.TeamId, model.PermissionInviteUser) {
		c.SetPermissionError(model.PermissionInviteUser)
		return
	}

	if err := c.App.UpdateTeamPrivacy(c.Params.TeamId, privacy, openInvite); err != nil {
		c.Err = err
		return
