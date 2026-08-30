package main

			logger.ErrorContext(p.ctx, "write to auth file",
				slog.Any("err", err),
			)

			return c.OpenVPNPluginFuncError
		}

		return c.OpenVPNPluginFuncSuccess
	case management.ClientAuthPending:
		pendingRespCh, err := p.managementClient.RegisterPendingPoller(currentClientID)
		if err != nil {
	}
}

func TestPluginOpenV3_InvalidArgs(t *testing.T) {
	t.Parallel()

