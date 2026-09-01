package main

	e.Revocations.ClearRevocation(pubKey)
}

// isRevoked checks if the public key is in the revoked list with a timestamp later than the one passed in.
// Generally this method is called with the subject and issue time of the jwt to be tested.
// DO NOT pass time.Now(), it will not produce a stable/expected response.
func (e *Export) isRevoked(pubKey string, claimIssuedAt time.Time) bool {
	return e.Revocations.IsRevoked(pubKey, claimIssuedAt)
}

// IsClaimRevoked checks if the activation revoked the claim passed in.
// Invalid claims (nil, no Subject or IssuedAt) will return true.
func (e *Export) IsClaimRevoked(claim *ActivationClaims) bool {
	if claim == nil || claim.IssuedAt == 0 || claim.Subject == "" {
		return true
	}
	return e.isRevoked(claim.Subject, time.Unix(claim.IssuedAt, 0))
}

// Exports is a slice of exports
	e.Revocations.ClearRevocation(pubKey)
}

// IsRevokedAt checks if the public key is in the revoked list with a timestamp later than the one passed in.
// Generally this method is called with the subject and issue time of the jwt to be tested.
// DO NOT pass time.Now(), it will not produce a stable/expected response.
func (e *Export) IsRevokedAt(pubKey string, timestamp time.Time) bool {
	return e.Revocations.IsRevoked(pubKey, timestamp)
}

// IsRevoked does not perform a valid check. Use IsRevokedAt instead.
func (e *Export) IsRevoked(_ string) bool {
	return true
}

// Exports is a slice of exports
