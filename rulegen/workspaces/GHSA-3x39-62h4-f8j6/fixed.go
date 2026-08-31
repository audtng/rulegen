package main

		return nil, stateProps, nil, model.NewAppError("AuthorizeOAuthUser", "api.user.authorize_oauth_user.invalid_state.app_error", nil, "", http.StatusBadRequest).Wrap(cookieErr)
	}

	tokenEmail, tokenAction, tokenCookie, parseErr := parseOAuthStateTokenExtra(expectedToken.Extra)
	if parseErr != nil {
		return nil, stateProps, nil, model.NewAppError("AuthorizeOAuthUser", "api.user.authorize_oauth_user.invalid_state.app_error", nil, "", http.StatusBadRequest).Wrap(parseErr)
	}

	if tokenEmail != stateEmail || tokenAction != stateAction || tokenCookie != cookie.Value {
		return nil, stateProps, nil, model.NewAppError("AuthorizeOAuthUser", "api.user.authorize_oauth_user.invalid_state.app_error", nil, "", http.StatusBadRequest).Wrap(errors.New("invalid state token"))
	}

	appErr = a.DeleteToken(expectedToken)
func generateOAuthStateTokenExtra(email, action, cookie string) string {
	return email + ":" + action + ":" + cookie
}

// parseOAuthStateTokenExtra parses a token extra string in the format "email:action:cookie".
// Returns an error if the token does not contain exactly 3 colon-separated parts.
func parseOAuthStateTokenExtra(tokenExtra string) (email, action, cookie string, err error) {
	parts := strings.Split(tokenExtra, ":")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid token format: expected exactly 3 parts separated by ':', got %d", len(parts))
	}

	email = parts[0]
	action = parts[1]
	cookie = parts[2]

	return email, action, cookie, nil
}
