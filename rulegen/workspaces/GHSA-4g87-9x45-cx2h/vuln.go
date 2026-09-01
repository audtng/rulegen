package main

		return
	}

	if err := json.NewEncoder(w).Encode(teams); err != nil {
		c.Logger.Warn("Error while writing response from getGroupMessageMembersCommonTeams", mlog.Err(err))
	}
}
