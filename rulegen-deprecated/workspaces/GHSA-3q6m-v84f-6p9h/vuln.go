package main

}

func (s *connection) handleHandshakeComplete() error {
	s.handshakeComplete = true
	defer s.handshakeCtxCancel()
	// Once the handshake completes, we have derived 1-RTT keys.
	// There's no point in queueing undecryptable packets for later decryption any more.
	s.undecryptablePackets = nil

	s.connIDManager.SetHandshakeComplete()
	s.connIDGenerator.SetHandshakeComplete()

	if s.perspective == protocol.PerspectiveClient {
		s.applyTransportParameters()
		return nil
	}

	if err := s.handleHandshakeConfirmed(); err != nil {
		return err
	}
	if log != nil {
		frames = make([]logging.Frame, 0, 4)
	}
	var handleErr error
	for len(data) > 0 {
		l, frame, err := s.frameParser.ParseNext(data, encLevel, s.version)
			return false, handleErr
		}
	}
	return
}

		case handshake.EventNoEvent:
			return nil
		case handshake.EventHandshakeComplete:
			err = s.handleHandshakeComplete()
		case handshake.EventReceivedTransportParameters:
			err = s.handleTransportParameters(ev.TransportParameters)
		case handshake.EventRestoredTransportParameters:
	if !acked1RTTPacket {
		return nil
	}
	if s.perspective == protocol.PerspectiveClient && !s.handshakeConfirmed {
		if err := s.handleHandshakeConfirmed(); err != nil {
			return err
}

func (s *connection) handleTransportParameters(params *wire.TransportParameters) error {
	if err := s.checkTransportParameters(params); err != nil {
		return &qerr.TransportError{
			ErrorCode:    qerr.TransportParameterError,
	if s.logger.Debug() {
		s.logger.Debugf("Processed Transport Parameters: %s", params)
	}
	if s.tracer != nil {
		s.tracer.ReceivedTransportParameters(params)
	}

	// check the initial_source_connection_id
	if params.InitialSourceConnectionID != s.handshakeDestConnID {
