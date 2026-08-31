package main

	// Maximum length of username, too keep the cache's memory signature low
	maxUsernameLength = 32
	// The default maximum session cache size
	defaultMaxCacheSize = 10000
	// The default number of maximum login failures before delay kicks in
	defaultMaxLoginFailures = 5
	// The default time in seconds for the failure window
	return expiredCount
}

// Protect admin user from login attempt reset caused by attempts to overflow cache in a brute force attack. Instead remove random non-admin to make room in cache. 
func pickRandomNonAdminLoginFailure(failures map[string]LoginAttempts, username string) *string {
	idx := rand.Intn(len(failures) - 1)
	i := 0
	for key := range failures {
		if i == idx {
			if key == common.ArgoCDAdminUsername || key == username {
				return pickRandomNonAdminLoginFailure(failures, username)
			}
			return &key
		}
		i++
	}
	return nil
}

// Updates the failure count for a given username. If failed is true, increases the counter. Otherwise, sets counter back to 0.
func (mgr *SessionManager) updateFailureCount(username string, failed bool) {

	// prevent overbloating the cache with fake entries, as this could lead to
	// memory exhaustion and ultimately in a DoS. We remove a single entry to
	// replace it with the new one.
	if failed && len(failures) >= getMaximumCacheSize() {
		log.Warnf("Session cache size exceeds %d entries, removing random entry", getMaximumCacheSize())
		rmUser := pickRandomNonAdminLoginFailure(failures, username)
		if rmUser != nil {
			delete(failures, *rmUser)
			log.Infof("Deleted entry for user %s from cache", *rmUser)
		}
	}

	attempt, ok := failures[username]
		assert.ErrorIs(t, err, common.TokenVerificationErr)
	})
}

func Test_PickFailureAttemptWhenOverflowed(t *testing.T) {
	t.Run("Not pick admin user from the queue", func(t *testing.T) {
		failures := map[string]LoginAttempts{
			"admin": {
				FailCount: 1,
			},
			"test2": {
				FailCount: 1,
			},
		}

		// inside pickRandomNonAdminLoginFailure, it uses random, so we need to test it multiple times
		for i := 0; i < 1000; i++ {
			user := pickRandomNonAdminLoginFailure(failures, "test")
			assert.Equal(t, "test2", *user)
		}
	})

	t.Run("Not pick admin user and current user from the queue", func(t *testing.T) {
		failures := map[string]LoginAttempts{
			"test": {
				FailCount: 1,
			},
			"admin": {
				FailCount: 1,
			},
			"test2": {
				FailCount: 1,
			},
		}

		// inside pickRandomNonAdminLoginFailure, it uses random, so we need to test it multiple times
		for i := 0; i < 1000; i++ {
			user := pickRandomNonAdminLoginFailure(failures, "test")
			assert.Equal(t, "test2", *user)
		}
	})
}
