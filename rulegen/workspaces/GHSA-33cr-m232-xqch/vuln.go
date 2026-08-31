package main

		ackErr = errorsmod.Wrapf(ibcerrors.ErrInvalidType, "cannot unmarshal ICS-20 transfer packet data")
		logger.Error(fmt.Sprintf("%s sequence %d", ackErr.Error(), packet.Sequence))
		ack = channeltypes.NewErrorAcknowledgement(ackErr)
	} else {
		bz := data.GetBytes()
		if !bytes.Equal(bz, packet.GetData()) {
			ackErr = errorsmod.Wrapf(ibcerrors.ErrInvalidType, "packet data did not marshal to expected bytes: %X ≠ %X", bz, packet.GetData())
			ack = channeltypes.NewErrorAcknowledgement(ackErr)
			im.keeper.Logger(ctx).Error(fmt.Sprintf("%s sequence %d", ackErr.Error(), packet.Sequence))
		}
	}

	// only attempt the application logic if the packet data
