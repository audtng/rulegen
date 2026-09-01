package main

			log.Errorf("Error during basic auth for caldav: %v", err)
			return false, nil
		}
	}
	if u != nil && err == nil {
		c.Set("userBasicAuth", u)
