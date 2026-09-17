package main

		ackErr = errorsmod.Wrapf(ibcerrors.ErrInvalidType, "cannot unmarshal ICS-20 transfer packet data")
		logger.Error(fmt.Sprintf("%s sequence %d", ackErr.Error(), packet.Sequence))
		ack = channeltypes.NewErrorAcknowledgement(ackErr)
	}

	// only attempt the application logic if the packet data
