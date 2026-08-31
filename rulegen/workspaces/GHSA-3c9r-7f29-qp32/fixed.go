package main

	enableUsername := *a.Config().EmailSettings.EnableSignInWithUsername
	enableEmail := *a.Config().EmailSettings.EnableSignInWithEmail

	if enableEmail || enableUsername {
		// If we are given a userID then fail if we can't find a user with that ID
		if id != "" {
			user, err := a.GetUser(id)
			if err != nil {
				if err.Id != MissingAccountError {
					err.StatusCode = http.StatusInternalServerError
					return nil, err
				}
				err.StatusCode = http.StatusBadRequest
				return nil, err
			}
			return user, nil
		}

		// Try to get the user by username/email
		if user, err := a.Srv().Store().User().GetForLogin(loginId, enableUsername, enableEmail); err == nil {
			return user, nil
		}
	}

	// Try to get the user with LDAP if enabled
