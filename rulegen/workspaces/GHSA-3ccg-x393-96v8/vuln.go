package main

// UserPassword holds a user password. Used to update it.
type UserPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UserChangePassword is the handler to change a users password
		return echo.NewHTTPError(http.StatusBadRequest, "No password provided.").Wrap(err)
	}

	if newPW.OldPassword == "" {
		return user.ErrEmptyOldPassword{}
	}
		return echo.NewHTTPError(http.StatusBadRequest, "No password provided.").Wrap(err)
	}

	s := db.NewSession()
	defer s.Close()

	// The previously issued reset token.
	Token string `json:"token"`
	// The new password for this user.
	NewPassword string `json:"new_password"`
}

// ResetPassword resets a users password. It returns the ID of the user whose
