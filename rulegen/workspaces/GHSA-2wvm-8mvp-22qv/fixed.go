package main

	}

	group.POST("/oidc/authorize", authMiddleware.WithAdminNotRequired().Add(), oc.authorizeHandler)
	group.POST("/oidc/authorize/callback-url", oc.authorizeCallbackURLHandler)
	group.POST("/oidc/authorization-required", authMiddleware.WithAdminNotRequired().Add(), oc.authorizationConfirmationRequiredHandler)
	group.GET("/oidc/par-request-info", authMiddleware.WithAdminNotRequired().Add(), oc.parRequestInfoHandler)

			c.JSON(http.StatusOK, gin.H{
				"error":            err.Error(),
				"requiresRedirect": true,
				"callbackURL":      callbackURL,
			})
			return
		}
	c.JSON(http.StatusOK, response)
}

// authorizeCallbackURLHandler godoc
// @Summary Resolve a validated callback URL for an OIDC authorization request
// @Description Resolves the redirect URI against the client's configured callback URLs without consuming PAR records.
// @Tags OIDC
// @Accept json
// @Produce json
// @Param request body dto.AuthorizeOidcClientCallbackRequestDto true "Authorization callback parameters"
// @Success 200 {object} dto.AuthorizeOidcClientCallbackResponseDto "Resolved callback URL"
// @Router /api/oidc/authorize/callback-url [post]
func (oc *OidcController) authorizeCallbackURLHandler(c *gin.Context) {
	var input dto.AuthorizeOidcClientCallbackRequestDto
	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(err)
		return
	}

	callbackURL, err := oc.oidcService.ResolveAuthorizeCallbackURL(
		c.Request.Context(),
		input.ClientID,
		input.CallbackURL,
		input.RequestURI,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.AuthorizeOidcClientCallbackResponseDto{CallbackURL: callbackURL})
}

// isOidcPromptError checks if an error is a prompt-related OIDC error that should trigger a redirect
func isOidcPromptError(err error) bool {
	var loginReq *common.OidcLoginRequiredError
	Issuer      string `json:"issuer"`
}

type AuthorizeOidcClientCallbackRequestDto struct {
	ClientID    string `json:"clientID" binding:"required"`
	CallbackURL string `json:"callbackURL"`
	RequestURI  string `json:"requestURI"`
}

type AuthorizeOidcClientCallbackResponseDto struct {
	CallbackURL string `json:"callbackURL"`
}

type AuthorizationRequiredDto struct {
	ClientID   string `json:"clientID" binding:"required"`
	Scope      string `json:"scope"`
			return "", "", err
		}
		if !hasAlreadyAuthorized {
			return "", callbackURL, &common.OidcConsentRequiredError{}
		}
	}

		return "", err
	}

	return resolveConfiguredCallbackURL(&client, inputCallbackURL)
}

// ResolveAuthorizeCallbackURL resolves the callback URL for a browser authorization
// request without authorizing the user or consuming PAR state.
func (s *OidcService) ResolveAuthorizeCallbackURL(ctx context.Context, clientID, inputCallbackURL, requestURI string) (string, error) {
	client, err := s.GetClient(ctx, clientID)
	if err != nil {
		return "", err
	}

	if client.RequiresPushedAuthorizationRequests && requestURI == "" {
		return "", &common.OidcPARRequiredError{}
	}

	if requestURI != "" {
		par, err := s.GetPushedAuthorizationRequest(ctx, clientID, requestURI)
		if err != nil {
			return "", err
		}
		inputCallbackURL = par.Parameters.RedirectURI
	}

	return resolveConfiguredCallbackURL(&client, inputCallbackURL)
}

func resolveConfiguredCallbackURL(client *model.OidcClient, inputCallbackURL string) (string, error) {
	if inputCallbackURL == "" {
		if len(client.CallbackURLs) > 0 {
			return client.CallbackURLs[0], nil
		}
		return "", &common.OidcMissingCallbackURLError{}
	}

	if len(client.CallbackURLs) == 0 {
		return "", &common.OidcMissingCallbackURLError{}
	}

