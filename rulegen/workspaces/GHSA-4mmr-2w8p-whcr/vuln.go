package main

		return
	}

	if err := c.App.UpdateTeamPrivacy(c.Params.TeamId, privacy, openInvite); err != nil {
		c.Err = err
		return
