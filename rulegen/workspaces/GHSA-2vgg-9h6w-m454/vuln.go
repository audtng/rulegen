package main

	// Maximum length of username, too keep the cache's memory signature low
	maxUsernameLength = 32
	// The default maximum session cache size
	defaultMaxCacheSize = 1000
	// The default number of maximum login failures before delay kicks in
	defaultMaxLoginFailures = 5
	// The default time in seconds for the failure window
	return expiredCount
}

// Updates the failure count for a given username. If failed is true, increases the counter. Otherwise, sets counter back to 0.
func (mgr *SessionManager) updateFailureCount(username string, failed bool) {

	// prevent overbloating the cache with fake entries, as this could lead to
	// memory exhaustion and ultimately in a DoS. We remove a single entry to
	// replace it with the new one.
	//
	// Chances are that we remove the one that is under active attack, but this
	// chance is low (1:cache_size)
	if failed && len(failures) >= getMaximumCacheSize() {
		log.Warnf("Session cache size exceeds %d entries, removing random entry", getMaximumCacheSize())
		idx := rand.Intn(len(failures) - 1)
		var rmUser string
		i := 0
		for key := range failures {
			if i == idx {
				rmUser = key
				delete(failures, key)
				break
			}
			i++
		}
		log.Infof("Deleted entry for user %s from cache", rmUser)
	}

	attempt, ok := failures[username]
		assert.ErrorIs(t, err, common.TokenVerificationErr)
	})
}
