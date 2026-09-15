package main

	e.Revocations.ClearRevocation(pubKey)
}

// IsRevokedAt checks if the public key is in the revoked list with a timestamp later than
// the one passed in. Generally this method is called with time.Now() but other time's can
// be used for testing.
func (e *Export) IsRevokedAt(pubKey string, timestamp time.Time) bool {
	return e.Revocations.IsRevoked(pubKey, timestamp)
}

// IsRevoked checks if the public key is in the revoked list with time.Now()
func (e *Export) IsRevoked(pubKey string) bool {
	return e.Revocations.IsRevoked(pubKey, time.Now())
}

// Exports is a slice of exports
	e.Revocations.ClearRevocation(pubKey)
}

// IsRevokedAt checks if the public key is in the revoked list with a timestamp later than
// the one passed in. Generally this method is called with time.Now() but other time's can
// be used for testing.
func (e *Export) IsRevokedAt(pubKey string, timestamp time.Time) bool {
	return e.Revocations.IsRevoked(pubKey, timestamp)
}

// IsRevoked checks if the public key is in the revoked list with time.Now()
func (e *Export) IsRevoked(pubKey string) bool {
	return e.Revocations.IsRevoked(pubKey, time.Now())
}

// Exports is a slice of exports
