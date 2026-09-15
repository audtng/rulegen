package main

// UserPassword holds a user password. Used to update it.
type UserPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password" valid:"bcrypt_password" minLength:"8" maxLength:"72"`
}

// UserChangePassword is the handler to change a users password
		return echo.NewHTTPError(http.StatusBadRequest, "No password provided.").Wrap(err)
	}

	// Validate the new password
	if err := c.Validate(newPW); err != nil {
		return err
	}

	if newPW.OldPassword == "" {
		return user.ErrEmptyOldPassword{}
	}
		return echo.NewHTTPError(http.StatusBadRequest, "No password provided.").Wrap(err)
	}

	// Validate the password
	if err := c.Validate(pwReset); err != nil {
		return err
	}

	s := db.NewSession()
	defer s.Close()

	// The previously issued reset token.
	Token string `json:"token"`
	// The new password for this user.
	NewPassword string `json:"new_password" valid:"bcrypt_password" minLength:"8" maxLength:"72"`
}

// ResetPassword resets a users password. It returns the ID of the user whose
