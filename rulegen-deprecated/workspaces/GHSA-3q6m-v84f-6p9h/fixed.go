package main

}

func (s *connection) handleHandshakeComplete() error {
	defer s.handshakeCtxCancel()
	// Once the handshake completes, we have derived 1-RTT keys.
	// There's no point in queueing undecryptable packets for later decryption anymore.
	s.undecryptablePackets = nil

	s.connIDManager.SetHandshakeComplete()
	s.connIDGenerator.SetHandshakeComplete()

	// The server applies transport parameters right away, but the client side has to wait for handshake completion.
	// During a 0-RTT connection, the client is only allowed to use the new transport parameters for 1-RTT packets.
	if s.perspective == protocol.PerspectiveClient {
		s.applyTransportParameters()
		return nil
	}

	// All these only apply to the server side.
	if err := s.handleHandshakeConfirmed(); err != nil {
		return err
	}
	if log != nil {
		frames = make([]logging.Frame, 0, 4)
	}
	handshakeWasComplete := s.handshakeComplete
	var handleErr error
	for len(data) > 0 {
		l, frame, err := s.frameParser.ParseNext(data, encLevel, s.version)
			return false, handleErr
		}
	}

	// Handle completion of the handshake after processing all the frames.
	// This ensures that we correctly handle the following case on the server side:
	// We receive a Handshake packet that contains the CRYPTO frame that allows us to complete the handshake,
	// and an ACK serialized after that CRYPTO frame. In this case, we still want to process the ACK frame.
	if !handshakeWasComplete && s.handshakeComplete {
		if err := s.handleHandshakeComplete(); err != nil {
			return false, err
		}
	}

	return
}

		case handshake.EventNoEvent:
			return nil
		case handshake.EventHandshakeComplete:
			// Don't call handleHandshakeComplete yet.
			// It's advantageous to process ACK frames that might be serialized after the CRYPTO frame first.
			s.handshakeComplete = true
		case handshake.EventReceivedTransportParameters:
			err = s.handleTransportParameters(ev.TransportParameters)
		case handshake.EventRestoredTransportParameters:
	if !acked1RTTPacket {
		return nil
	}
	// On the client side: If the packet acknowledged a 1-RTT packet, this confirms the handshake.
	// This is only possible if the ACK was sent in a 1-RTT packet.
	// This is an optimization over simply waiting for a HANDSHAKE_DONE frame, see section 4.1.2 of RFC 9001.
	if s.perspective == protocol.PerspectiveClient && !s.handshakeConfirmed {
		if err := s.handleHandshakeConfirmed(); err != nil {
			return err
}

func (s *connection) handleTransportParameters(params *wire.TransportParameters) error {
	if s.tracer != nil {
		s.tracer.ReceivedTransportParameters(params)
	}
	if err := s.checkTransportParameters(params); err != nil {
		return &qerr.TransportError{
			ErrorCode:    qerr.TransportParameterError,
	if s.logger.Debug() {
		s.logger.Debugf("Processed Transport Parameters: %s", params)
	}

	// check the initial_source_connection_id
	if params.InitialSourceConnectionID != s.handshakeDestConnID {
