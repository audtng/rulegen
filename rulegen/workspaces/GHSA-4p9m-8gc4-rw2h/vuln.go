package main

					handling := fmsg.handling
					useRevisedError := h.fsm.isTreatAsWithdraw

					ok, err := bgp.ValidateUpdateMsg(body, rfMap, h.fsm.isEBGP, h.fsm.isConfed, h.allowLoopback)
					if !ok {
						handling = h.handlingError(m, err, useRevisedError)
						fmsg.handling = handling
					}
					if handling == bgp.ERROR_HANDLING_SESSION_RESET {
						h.fsm.logger.Warn("Session will be reset due to malformed BGP update message",
