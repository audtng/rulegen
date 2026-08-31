package main

	}

	group.POST("/oidc/authorize", authMiddleware.WithAdminNotRequired().Add(), oc.authorizeHandler)
	group.POST("/oidc/authorization-required", authMiddleware.WithAdminNotRequired().Add(), oc.authorizationConfirmationRequiredHandler)
	group.GET("/oidc/par-request-info", authMiddleware.WithAdminNotRequired().Add(), oc.parRequestInfoHandler)

			c.JSON(http.StatusOK, gin.H{
				"error":            err.Error(),
				"requiresRedirect": true,
			})
			return
		}
	c.JSON(http.StatusOK, response)
}

// isOidcPromptError checks if an error is a prompt-related OIDC error that should trigger a redirect
func isOidcPromptError(err error) bool {
	var loginReq *common.OidcLoginRequiredError
	Issuer      string `json:"issuer"`
}

type AuthorizationRequiredDto struct {
	ClientID   string `json:"clientID" binding:"required"`
	Scope      string `json:"scope"`
			return "", "", err
		}
		if !hasAlreadyAuthorized {
			return "", "", &common.OidcConsentRequiredError{}
		}
	}

		return "", err
	}

	if inputCallbackURL == "" || len(client.CallbackURLs) == 0 {
		return "", &common.OidcMissingCallbackURLError{}
	}

