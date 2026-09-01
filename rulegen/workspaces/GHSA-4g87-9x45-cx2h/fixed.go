package main

		return
	}

	if err := json.NewEncoder(w).Encode(c.App.SanitizeTeams(*c.AppContext.Session(), teams)); err != nil {
		c.Logger.Warn("Error while writing response from getGroupMessageMembersCommonTeams", mlog.Err(err))
	}
}
