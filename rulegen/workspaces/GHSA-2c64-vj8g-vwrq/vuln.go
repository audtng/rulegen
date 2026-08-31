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

	account.Exports.Add(e)

	pubKey := "bar"
	now := time.Now()

	// test that clear is safe before we add any
	e.ClearRevocation(pubKey)

	if e.IsRevokedAt(pubKey, now) {
		t.Errorf("no revocation was added so is revoked should be false")
	}

	e.RevokeAt(pubKey, now.Add(time.Second*100))

	if !e.IsRevokedAt(pubKey, now) {
		t.Errorf("revocation should hold when timestamp is in the future")
	}

	if e.IsRevokedAt(pubKey, now.Add(time.Second*150)) {
		t.Errorf("revocation should time out")
	}

	e.RevokeAt(pubKey, now.Add(time.Second*50)) // shouldn't change the revocation, you can't move it in

	if !e.IsRevokedAt(pubKey, now.Add(time.Second*60)) {
		t.Errorf("revocation should hold, 100 > 50")
	}

	encoded, _ := account.Encode(akp)
	decoded, _ := DecodeAccountClaims(encoded)

	if !decoded.Exports[0].IsRevokedAt(pubKey, now.Add(time.Second*60)) {
		t.Errorf("revocation should last across encoding")
	}

	e.ClearRevocation(pubKey)

	if e.IsRevokedAt(pubKey, now) {
		t.Errorf("revocations should be cleared")
	}

	e.RevokeAt(pubKey, now.Add(time.Second*1000))

	if !e.IsRevoked(pubKey) {
		t.Errorf("revocation be true we revoked in the future")
	}
}
