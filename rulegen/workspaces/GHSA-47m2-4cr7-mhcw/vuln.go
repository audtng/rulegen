package main

}

func (c *Conn) handleHandshakeConfirmed(now monotime.Time) error {
	if err := c.dropEncryptionLevel(protocol.EncryptionHandshake, now); err != nil {
		return err
	}
