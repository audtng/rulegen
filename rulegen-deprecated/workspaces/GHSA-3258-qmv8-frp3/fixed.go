package main

	applyRoutes(smfCallbackGroup, smfCallbackRoutes)

	upiGroup := router.Group(factory.UpiUriPrefix)
	upiAuthCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NSMF_OAM)
	upiGroup.Use(func(c *gin.Context) {
		upiAuthCheck.Check(c, smf_context.GetSelf())
	})
	upiRoutes := s.getUPIRoutes()
	applyRoutes(upiGroup, upiRoutes)

