package main

	applyRoutes(smfCallbackGroup, smfCallbackRoutes)

	upiGroup := router.Group(factory.UpiUriPrefix)
	upiRoutes := s.getUPIRoutes()
	applyRoutes(upiGroup, upiRoutes)

