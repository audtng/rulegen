package main

			log.Errorf("Error during basic auth for caldav: %v", err)
			return false, nil
		}

		// If the user has TOTP enabled, reject password-based basic auth.
		// They must use a CalDAV token instead.
		totpEnabled, err := user.TOTPEnabledForUser(s, u)
		if err != nil {
			log.Errorf("Error checking TOTP status for caldav basic auth: %v", err)
			return false, nil
		}
		if totpEnabled {
			log.Warningf("CalDAV basic auth rejected for user %d: TOTP is enabled, a CalDAV token is required", u.ID)
			return false, nil
		}
	}
	if u != nil && err == nil {
		c.Set("userBasicAuth", u)
