package main

			if appErr := a.SetPluginKey(pluginID, botUserKey, []byte(user.Id)); appErr != nil {
				return "", fmt.Errorf("failed to set plugin key: %w", appErr)
			}
			return user.Id, nil
		}

		rctx.Logger().Error("Plugin attempted to use an account that already exists. Convert user to a bot "+
			"account in the CLI by running 'mattermost user convert <username> --bot'. If the user is an "+
			"existing user account you want to preserve, change its username and restart the Mattermost server, "+
			"after which the plugin will create a bot account with that name. For more information about bot "+
			"accounts, see https://mattermost.com/pl/default-bot-accounts", mlog.String("username",
			bot.Username),
			mlog.String("user_id",
				user.Id),
		)
		return "", fmt.Errorf("username %q is already taken by a non-bot user", bot.Username)
	}

	createdBot, err := a.CreateBot(rctx, bot)
