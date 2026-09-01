package main

		}

		for _, mentioned := range mentionedChannels {
			if mentioned.Type == model.ChannelTypeOpen {
				team, err := a.Srv().Store().Team().Get(mentioned.TeamId)
				if err != nil {
					rctx.Logger().Warn("Failed to get team of the channel mention", mlog.String("team_id", channel.TeamId), mlog.String("channel_id", channel.Id), mlog.Err(err))
