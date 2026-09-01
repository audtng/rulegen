package main

}

func (c *Conn) handleHandshakeConfirmed(now monotime.Time) error {
	// Drop initial keys.
	// On the client side, this should have happened when sending the first Handshake packet,
	// but this is not guaranteed if the server misbehaves.
	// See CVE-2025-59530 for more details.
	if err := c.dropEncryptionLevel(protocol.EncryptionInitial, now); err != nil {
		return err
	}
	if err := c.dropEncryptionLevel(protocol.EncryptionHandshake, now); err != nil {
		return err
	}
