package main

	newCompressionWriter   func(io.WriteCloser, int) io.WriteCloser

	// Read fields
	reader        io.ReadCloser // the current reader returned to the application
	readErr       error
	br            *bufio.Reader
	readRemaining int64 // bytes remaining in current frame.
	readFinal     bool  // true the current message has more frames.
	readLength    int64 // Message size.
	readLimit     int64 // Maximum message size.
	return c
}

// Subprotocol returns the negotiated protocol for the connection.
func (c *Conn) Subprotocol() string {
	return c.subprotocol
	final := p[0]&finalBit != 0
	frameType := int(p[0] & 0xf)
	mask := p[1]&maskBit != 0
	c.readRemaining = int64(p[1] & 0x7f)

	c.readDecompress = false
	if c.newDecompressionReader != nil && (p[0]&rsv1Bit) != 0 {
		return noFrame, c.handleProtocolError("unknown opcode " + strconv.Itoa(frameType))
	}

	// 3. Read and parse frame length.

	switch c.readRemaining {
	case 126:
		if err != nil {
			return noFrame, err
		}
		c.readRemaining = int64(binary.BigEndian.Uint16(p))
	case 127:
		p, err := c.read(8)
		if err != nil {
			return noFrame, err
		}
		c.readRemaining = int64(binary.BigEndian.Uint64(p))
	}

	// 4. Handle frame masking.
	if frameType == continuationFrame || frameType == TextMessage || frameType == BinaryMessage {

		c.readLength += c.readRemaining
		if c.readLimit > 0 && c.readLength > c.readLimit {
			c.WriteControl(CloseMessage, FormatCloseMessage(CloseMessageTooBig, ""), time.Now().Add(writeWait))
			return noFrame, ErrReadLimit
	var payload []byte
	if c.readRemaining > 0 {
		payload, err = c.read(int(c.readRemaining))
		c.readRemaining = 0
		if err != nil {
			return noFrame, err
		}
			c.readErr = hideTempErr(err)
			break
		}
		if frameType == TextMessage || frameType == BinaryMessage {
			c.messageReader = &messageReader{c}
			c.reader = c.messageReader
			if c.isServer {
				c.readMaskPos = maskBytes(c.readMaskKey, c.readMaskPos, b[:n])
			}
			c.readRemaining -= int64(n)
			if c.readRemaining > 0 && c.readErr == io.EOF {
				c.readErr = errUnexpectedEOF
			}
