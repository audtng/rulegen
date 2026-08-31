package main

		return false, err
	}

	if _, err := SendPasswordResetEmail(email, token, user.Locale, siteURL); err != nil {
		return false, model.NewLocAppError("SendPasswordReset", "api.user.send_password_reset.send.app_error", nil, "err="+err.Message)
	}

