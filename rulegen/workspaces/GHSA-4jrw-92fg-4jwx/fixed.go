package main

	var lastEapID uint8
	if updateAuthenticationInfo.ResynchronizationInfo != nil {
		logger.UeAuthLog.Warningln("Auts: ", updateAuthenticationInfo.ResynchronizationInfo.Auts)
		if !ausf_context.CheckIfSuciSupiPairExists(supiOrSuci) {
			logger.UeAuthLog.Warningln("Resync failed: SUCI mapping not found for ", supiOrSuci)
			c.JSON(http.StatusNotFound, models.ProblemDetails{Status: 404, Cause: "USER_NOT_FOUND"})
			return
		}
		ausfCurrentSupi := ausf_context.GetSupiFromSuciSupiMap(supiOrSuci)
		logger.UeAuthLog.Warningln(ausfCurrentSupi)
		if !ausf_context.CheckIfAusfUeContextExists(ausfCurrentSupi) {
			logger.UeAuthLog.Warningln("Resync failed: AusfUeContext not found for SUPI: ", ausfCurrentSupi)
			c.JSON(http.StatusNotFound, models.ProblemDetails{Status: 404, Cause: "USER_NOT_FOUND"})
			return
		}
		ausfCurrentContext := ausf_context.GetAusfUeContext(ausfCurrentSupi)
		logger.UeAuthLog.Warningln(ausfCurrentContext.Rand)
		if updateAuthenticationInfo.ResynchronizationInfo.Rand == "" {
			updateAuthenticationInfo.ResynchronizationInfo.Rand = ausfCurrentContext.Rand
		}
		logger.UeAuthLog.Warningln("Rand: ", updateAuthenticationInfo.ResynchronizationInfo.Rand)
		authInfoReq.ResynchronizationInfo = updateAuthenticationInfo.ResynchronizationInfo
		lastEapID = ausfCurrentContext.EapID
	}

	udmUrl := p.Consumer().GetUdmUrl(self.NrfUri)
	var lastEapID uint8
	if updateAuthenticationInfo.ResynchronizationInfo != nil {
		logger.UeAuthLog.Warningln("Auts: ", updateAuthenticationInfo.ResynchronizationInfo.Auts)
		if ausf_context.CheckIfSuciSupiPairExists(supiOrSuci) {
			ausfCurrentSupi := ausf_context.GetSupiFromSuciSupiMap(supiOrSuci)
			logger.UeAuthLog.Warningln(ausfCurrentSupi)
			if ausf_context.CheckIfAusfUeContextExists(ausfCurrentSupi) {
				ausfCurrentContext := ausf_context.GetAusfUeContext(ausfCurrentSupi)
				logger.UeAuthLog.Warningln(ausfCurrentContext.Rand)
				if updateAuthenticationInfo.ResynchronizationInfo.Rand == "" {
					updateAuthenticationInfo.ResynchronizationInfo.Rand = ausfCurrentContext.Rand
				}
				logger.UeAuthLog.Warningln("Rand: ", updateAuthenticationInfo.ResynchronizationInfo.Rand)
				authInfoReq.ResynchronizationInfo = updateAuthenticationInfo.ResynchronizationInfo
				lastEapID = ausfCurrentContext.EapID
			} else {
				logger.UeAuthLog.Warningln("Resync failed: AusfUeContext not found for SUPI: ", ausfCurrentSupi)
				c.JSON(http.StatusNotFound, models.ProblemDetails{Status: 404, Cause: "USER_NOT_FOUND"})
				return
			}
		} else {
			logger.UeAuthLog.Warningln("Resync failed: SUCI mapping not found for ", supiOrSuci)
			c.JSON(http.StatusNotFound, models.ProblemDetails{Status: 404, Cause: "USER_NOT_FOUND"})
			return
		}
	}

	udmUrl := p.Consumer().GetUdmUrl(self.NrfUri)
internal/sbi/processor/ue_authentication.go | 34 ++++++++++-----------
1 file changed, 16 insertions(+), 18 deletions(-)
