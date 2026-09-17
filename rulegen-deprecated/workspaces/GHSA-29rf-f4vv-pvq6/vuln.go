package main

		scopeString := sessionSplit[3]
		scopes := parseScopes(scopeString)
		var user *schemas.User
		oauthCode := ctx.Request.FormValue("code")
		if oauthCode == "" {
			log.Debug().Err(err).Msg("Invalid oauth code")
		}
		switch provider {
		case constants.AuthRecipeMethodGoogle:
			user, err = h.processGoogleUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodGithub:
			user, err = h.processGithubUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodFacebook:
			user, err = h.processFacebookUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodLinkedIn:
			user, err = h.processLinkedInUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodApple:
			var appleUser *AppleUserInfo
			appleUser, err = parseAppleUserField(ctx.Request.FormValue("user"))
				ctx.JSON(400, gin.H{"error": "invalid apple user info"})
				return
			}
			user, err = h.processAppleUserInfo(ctx, oauthCode, appleUser)
		case constants.AuthRecipeMethodDiscord:
			user, err = h.processDiscordUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodTwitter:
			// Twitter/X uses PKCE: retrieve the verifier stored at login keyed by state.
			verifier, verr := h.MemoryStoreProvider.GetAndRemoveState(pkceVerifierKeyPrefix + state)
				ctx.JSON(400, gin.H{"error": "invalid oauth state"})
				return
			}
			user, err = h.processTwitterUserInfo(ctx, oauthCode, verifier)
		case constants.AuthRecipeMethodMicrosoft:
			user, err = h.processMicrosoftUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodTwitch:
			user, err = h.processTwitchUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodRoblox:
			user, err = h.processRobloxUserInfo(ctx, oauthCode)
		default:
			log.Debug().Err(err).Msg("Invalid oauth provider")
			err = fmt.Errorf(`invalid oauth provider`)
		log := log.With().Str("email", refs.StringValue(user.Email)).Logger()
		isSignUp := false

		if err != nil {
			isSignupEnabled := h.Config.EnableSignup
			if !isSignupEnabled {
			}

			user.Roles = strings.Join(inputRoles, ",")
			now := time.Now().Unix()
			user.EmailVerifiedAt = &now
			user, err = h.StorageProvider.AddUser(ctx, user)
			if err != nil {
				log.Debug().Err(err).Msg("Failed to add user")
			// was never verified, do not link the OAuth identity to it.
			// Instead, delete the unverified account and treat as a new signup
			// for the OAuth user who actually controls the email address.
			if existingUser.EmailVerifiedAt == nil {
				log.Info().Msg("Removing unverified pre-existing account before OAuth signup")
				if err := h.StorageProvider.DeleteUser(ctx, existingUser); err != nil {
					log.Debug().Err(err).Msg("Failed to delete unverified user")
					ctx.JSON(500, gin.H{"error": "failed to process OAuth login"})
	}
}

// oidcClaims is the allow-list of OpenID Connect standard claims Authorizer
// maps onto a user. ID tokens are decoded into this and never straight into
// schemas.User, for two reasons:
//     signup_methods, is_active, created_at ...) merely by sharing its json
//     tag.
type oidcClaims struct {
	Email       string `json:"email"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	return user
}

func (h *httpProvider) processGoogleUserInfo(ctx *gin.Context, code string) (*schemas.User, error) {
	log := h.Log.With().Str("func", "processGoogleUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodGoogle)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, fmt.Errorf("error getting oauth config: %s", err.Error())
	}
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, fmt.Errorf("invalid google exchange code: %s", err.Error())
	}

	issuer := "https://accounts.google.com"
	}
	oidcProvider, err := getOIDCProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to create oidc provider: %s", err.Error())
	}
	verifier := oidcProvider.Verifier(&oidc.Config{ClientID: h.GoogleClientID})
	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Debug().Err(err).Msg("Failed to extract ID Token from OAuth2 token")
		return nil, fmt.Errorf("unable to extract id_token")
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to verify ID Token")
		return nil, fmt.Errorf("unable to verify id_token: %s", err.Error())
	}
	claims := &oidcClaims{}
	if err := idToken.Claims(claims); err != nil {
		log.Debug().Err(err).Msg("Failed to parse ID Token claims")
		return nil, fmt.Errorf("unable to extract claims")
	}

	return claims.toUser(), nil
}

// setGithubHeaders applies the headers GitHub's REST API docs ask every
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}

func (h *httpProvider) processGithubUserInfo(ctx *gin.Context, code string) (*schemas.User, error) {
	log := h.Log.With().Str("func", "processGithubUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodGithub)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, fmt.Errorf("error getting oauth config: %s", err.Error())
	}

	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, fmt.Errorf("invalid github exchange code: %s", err.Error())
	}
	userInfoURL := constants.GithubUserInfoURL
	emailsURL := constants.GithubUserEmails
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create github user info request")
		return nil, fmt.Errorf("error creating github user info request: %s", err.Error())
	}
	setGithubHeaders(req, oauth2Token.AccessToken)

	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to request github user info")
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read github user info response body")
		return nil, fmt.Errorf("failed to read github response body: %s", err.Error())
	}
	if response.StatusCode >= 400 {
		log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request github user info")
		return nil, fmt.Errorf("failed to request github user info: %s", string(body))
	}

	// Only the three fields below are used. A typed struct (rather than a
	}
	if err := json.Unmarshal(body, &userRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal github user info")
		return nil, fmt.Errorf("failed to parse github user info: %s", err.Error())
	}

	name := strings.Split(userRawData.Name, " ")
		req, err := http.NewRequest(http.MethodGet, emailsURL, nil)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to create github emails request")
			return nil, fmt.Errorf("error creating github user info request: %s", err.Error())
		}
		setGithubHeaders(req, oauth2Token.AccessToken)

		response, err := client.Do(req)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to request github user email")
			return nil, err
		}

		defer func() { _ = response.Body.Close() }()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to read github user email response body")
			return nil, fmt.Errorf("failed to read github response body: %s", err.Error())
		}
		if response.StatusCode >= 400 {
			log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request github user email")
			return nil, fmt.Errorf("failed to request github user info: %s", string(body))
		}

		emailData := []GithubUserEmails{}
		err = json.Unmarshal(body, &emailData)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to parse github user email")
			return nil, fmt.Errorf("failed to parse github user email: %s", err.Error())
		}

		// GET /user/emails lists every address on the account, verified or
		}
		if email == "" {
			log.Debug().Msg("No verified email on github account")
			return nil, fmt.Errorf("failed to get a verified email address from github")
		}
	}

		Email:      &email,
	}

	return user, nil
}

func (h *httpProvider) processFacebookUserInfo(ctx *gin.Context, code string) (*schemas.User, error) {
	log := h.Log.With().Str("func", "processFacebookUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodFacebook)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, fmt.Errorf("error getting oauth config: %s", err.Error())
	}
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Invalid facebook exchange code")
		return nil, fmt.Errorf("invalid facebook exchange code: %s", err.Error())
	}
	userInfoURL := constants.FacebookUserInfoURL
	if mockBase := h.TestOAuthBaseURL(constants.AuthRecipeMethodFacebook); mockBase != "" {
	req, err := http.NewRequest("GET", userInfoURL+oauth2Token.AccessToken, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Error creating facebook user info request")
		return nil, fmt.Errorf("error creating facebook user info request: %s", err.Error())
	}

	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to process facebook user")
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read facebook response")
		return nil, fmt.Errorf("failed to read facebook response body: %s", err.Error())
	}
	if response.StatusCode >= 400 {
		log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request facebook user info")
		return nil, fmt.Errorf("failed to request facebook user info: %s", string(body))
	}
	// Typed decode, not fmt.Sprintf over a map: Graph API omits `email`
	// entirely when "no valid email address is available" (user/reference/user),
	}
	if err := json.Unmarshal(body, &userRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal facebook user info")
		return nil, fmt.Errorf("failed to parse facebook user info: %s", err.Error())
	}

	email := userRawData.Email
	if email == "" {
		log.Debug().Msg("Facebook user info has no email")
		return nil, fmt.Errorf("failed to get email from facebook user info: the account has no available email address")
	}

	picture := userRawData.Picture.Data.URL
		Email:      &email,
	}

	return user, nil
}

func (h *httpProvider) processLinkedInUserInfo(ctx *gin.Context, code string) (*schemas.User, error) {
	log := h.Log.With().Str("func", "processLinkedInUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodLinkedIn)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, fmt.Errorf("error getting oauth config: %s", err.Error())
	}

	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, fmt.Errorf("invalid linkedin exchange code: %s", err.Error())
	}

	userInfoURL := constants.LinkedInUserInfoURL
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create linkedin user info request")
		return nil, fmt.Errorf("error creating linkedin user info request: %s", err.Error())
	}
	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", oauth2Token.AccessToken)},
	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to request linkedin user info")
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read linkedin user info response body")
		return nil, fmt.Errorf("failed to read linkedin response body: %s", err.Error())
	}

	if response.StatusCode >= 400 {
		log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request linkedin user info")
		return nil, fmt.Errorf("failed to request linkedin user info: %s", string(body))
	}

	// OIDC userinfo shape (sub/name/given_name/family_name/picture/locale/
	// email/email_verified) - one call, no separate /v2/emailAddress hop.
	var userRawData struct {
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture    string `json:"picture"`
		Email      string `json:"email"`
	}
	if err := json.Unmarshal(body, &userRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal linkedin user info")
		return nil, fmt.Errorf("failed to parse linkedin user info: %s", err.Error())
	}

	// `email` is documented as optional - it is only present when the member
	// than a synthetic-email fallback.
	if userRawData.Email == "" {
		log.Debug().Msg("LinkedIn user info has no email")
		return nil, fmt.Errorf("failed to extract email from linkedin response")
	}

	user := &schemas.User{
		Email:      &userRawData.Email,
	}

	return user, nil
}

func (h *httpProvider) processAppleUserInfo(ctx *gin.Context, code string, appleUser *AppleUserInfo) (*schemas.User, error) {
	log := h.Log.With().Str("func", "processAppleUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodApple)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, fmt.Errorf("error getting oauth config: %s", err.Error())
	}

	var user = &schemas.User{}
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return user, fmt.Errorf("invalid apple exchange code: %s", err.Error())
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Debug().Err(err).Msg("Failed to extract ID Token from OAuth2 token")
		return user, fmt.Errorf("unable to extract id_token")
	}

	// Verify the Apple ID token signature, issuer, and audience using OIDC discovery
	oidcProvider, err := getOIDCProvider(ctx, issuer)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create Apple OIDC provider")
		return user, fmt.Errorf("failed to create oidc provider: %s", err.Error())
	}
	verifier := oidcProvider.Verifier(&oidc.Config{ClientID: h.AppleClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to verify Apple ID Token")
		return user, fmt.Errorf("unable to verify id_token: %s", err.Error())
	}

	claims := make(map[string]interface{})
	if err := idToken.Claims(&claims); err != nil {
		log.Debug().Err(err).Msg("Failed to parse Apple ID Token claims")
		return user, fmt.Errorf("failed to parse claims: %s", err.Error())
	}

	if val, ok := claims["email"]; !ok || val == nil {
		log.Debug().Msg("Failed to extract email from claims.")
		return user, fmt.Errorf("unable to extract email, please check the scopes enabled for your app. It needs `email`, `name` scopes")
	} else {
		email, _ := val.(string)
		user.Email = &email
	}

	user.GivenName = &appleUser.Name.FirstName
	user.FamilyName = &appleUser.Name.LastName

	return user, nil
}

// processDiscordUserInfo exchanges the Discord OAuth code for the user's
// creating a duplicate account - the same fallback discipline
// processTwitterUserInfo uses above for X, which never returns a real email
// at all.
func (h *httpProvider) processDiscordUserInfo(ctx *gin.Context, code string) (*schemas.User, error) {
	log := h.Log.With().Str("func", "processDiscordUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodDiscord)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, fmt.Errorf("error getting oauth config: %s", err.Error())
	}
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, fmt.Errorf("invalid discord exchange code: %s", err.Error())
	}

	userInfoURL := constants.DiscordUserInfoURL
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create Discord user info request")
		return nil, fmt.Errorf("error creating Discord user info request: %s", err.Error())
	}
	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", oauth2Token.AccessToken)},
	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to request Discord user info")
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read Discord user info response body")
		return nil, fmt.Errorf("failed to read Discord response body: %s", err.Error())
	}

	if response.StatusCode >= 400 {
		log.Debug().Err(err).Msg("Failed to request Discord user info")
		return nil, fmt.Errorf("failed to request Discord user info: %s", string(body))
	}

	// Unmarshal the response body into a map. GET /users/@me returns a flat
	userRawData := make(map[string]interface{})
	if err := json.Unmarshal(body, &userRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal Discord response")
		return nil, fmt.Errorf("failed to unmarshal Discord response: %s", err.Error())
	}

	// Extract the username
	firstName, ok := userRawData["username"].(string)
	if !ok {
		log.Debug().Err(err).Msg("Username is not in expected format or missing in user data")
		return nil, fmt.Errorf("username is not in expected format or missing in user data")
	}
	discordID, ok := userRawData["id"].(string)
	if !ok || discordID == "" {
		log.Debug().Msg("Discord user info missing id")
		return nil, fmt.Errorf("discord response missing id field")
	}
	// `avatar` is nullable (?string in Discord's user object) for accounts on
	// the default avatar - building the CDN URL from an empty hash yields a
	}

	email := resolveDiscordEmail(discordID, userRawData)

	user := &schemas.User{
		GivenName: &firstName,
		Email:     &email,
	}

	return user, nil
}

// resolveDiscordEmail prefers the real email Discord returns; falls back to
// returning Twitter user instead of creating a duplicate account on every
// login. Operators who opt into X's `users.email` scope + app permission get
// a real confirmed_email instead (see TwitterUserInfoURL's doc comment).
func (h *httpProvider) processTwitterUserInfo(ctx *gin.Context, code, verifier string) (*schemas.User, error) {
	log := h.Log.With().Str("func", "processTwitterUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodTwitter)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, fmt.Errorf("error getting oauth config: %s", err.Error())
	}

	oauth2Token, err := cfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, fmt.Errorf("invalid twitter exchange code: %s", err.Error())
	}

	userInfoURL := constants.TwitterUserInfoURL
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create Twitter user info request")
		return nil, fmt.Errorf("error creating Twitter user info request: %s", err.Error())
	}
	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", oauth2Token.AccessToken)},
	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to request Twitter user info")
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read Twitter user info response body")
		return nil, fmt.Errorf("failed to read Twitter response body: %s", err.Error())
	}

	if response.StatusCode >= 400 {
		log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request Twitter user info")
		return nil, fmt.Errorf("failed to request Twitter user info: %s", string(body))
	}

	responseRawData := make(map[string]interface{})
	if err := json.Unmarshal(body, &responseRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal twitter user info")
		return nil, fmt.Errorf("failed to parse twitter user info: %s", err.Error())
	}

	userRawData, ok := responseRawData["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("twitter response missing data field")
	}

	// Twitter API does not return E-Mail adresses by default. For that case special privileges have
	twitterID, ok := userRawData["id"].(string)
	if !ok || twitterID == "" {
		log.Debug().Msg("Twitter user info missing id")
		return nil, fmt.Errorf("twitter response missing id field")
	}

	// Currently Twitter API only provides the full name of a user. To fill givenName and familyName
	profilePicture, _ := userRawData["profile_image_url"].(string)

	email := resolveTwitterEmail(twitterID, userRawData)

	user := &schemas.User{
		Email:      &email,
		Nickname:   &nickname,
	}

	return user, nil
}

// twitterSyntheticEmail derives a stable, non-routable synthetic email from
}

// process microsoft user information
func (h *httpProvider) processMicrosoftUserInfo(ctx *gin.Context, code string) (*schemas.User, error) {
	log := h.Log.With().Str("func", "processMicrosoftUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodMicrosoft)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, fmt.Errorf("error getting oauth config: %s", err.Error())
	}
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, fmt.Errorf("invalid microsoft exchange code: %s", err.Error())
	}
	issuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", h.MicrosoftTenantID)
	if mockBase := h.TestOAuthBaseURL(constants.AuthRecipeMethodMicrosoft); mockBase != "" {
	}
	oidcProvider, err := getOIDCProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to create oidc provider: %s", err.Error())
	}
	// we need to skip issuer check because for common tenant it will return internal issuer which does not match
	verifier := oidcProvider.Verifier(&oidc.Config{
		ClientID:        h.MicrosoftClientID,
		SkipIssuerCheck: true,
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Debug().Err(err).Msg("Failed to extract ID Token from OAuth2 token")
		return nil, fmt.Errorf("unable to extract id_token")
	}
	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to verify ID Token")
		return nil, fmt.Errorf("unable to verify id_token: %s", err.Error())
	}
	claims := &oidcClaims{}
	if err := idToken.Claims(claims); err != nil {
		log.Debug().Err(err).Msg("Failed to parse ID Token claims")
		return nil, fmt.Errorf("unable to extract claims")
	}

	return claims.toUser(), nil
}

// process twitch user information
func (h *httpProvider) processTwitchUserInfo(ctx *gin.Context, code string) (*schemas.User, error) {
	log := h.Log.With().Str("func", "processTwitchUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodTwitch)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, fmt.Errorf("error getting oauth config: %s", err.Error())
	}

	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, fmt.Errorf("invalid twitch exchange code: %s", err.Error())
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Debug().Err(err).Msg("Failed to extract ID Token from OAuth2 token")
		return nil, fmt.Errorf("unable to extract id_token")
	}
	issuer := "https://id.twitch.tv/oauth2"
	if mockBase := h.TestOAuthBaseURL(constants.AuthRecipeMethodTwitch); mockBase != "" {
	oidcProvider, err := getOIDCProvider(ctx, issuer)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create OIDC provider")
		return nil, fmt.Errorf("failed to create oidc provider: %s", err.Error())
	}
	verifier := oidcProvider.Verifier(&oidc.Config{
		ClientID:        h.TwitchClientID,
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to verify ID Token")
		return nil, fmt.Errorf("unable to verify id_token: %s", err.Error())
	}

	claims := &oidcClaims{}
	if err := idToken.Claims(claims); err != nil {
		log.Debug().Err(err).Msg("Failed to parse ID Token claims")
		return nil, fmt.Errorf("unable to extract claims")
	}

	return claims.toUser(), nil
}

// process roblox user information
func (h *httpProvider) processRobloxUserInfo(ctx *gin.Context, code string) (*schemas.User, error) {
	log := h.Log.With().Str("func", "processRobloxUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodRoblox)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, fmt.Errorf("error getting oauth config: %s", err.Error())
	}
	// Roblox is a confidential client (client_secret set); PKCE is optional and
	// no code_challenge is sent at login, so no code_verifier is replayed here.
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, fmt.Errorf("invalid roblox exchange code: %s", err.Error())
	}

	userInfoURL := constants.RobloxUserInfoURL
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create roblox user info request")
		return nil, fmt.Errorf("error creating roblox user info request: %s", err.Error())
	}
	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", oauth2Token.AccessToken)},
	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to request roblox user info")
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read roblox user info response body")
		return nil, fmt.Errorf("failed to read roblox response body: %s", err.Error())
	}

	if response.StatusCode >= 400 {
		log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request roblox user info")
		return nil, fmt.Errorf("failed to request roblox user info: %s", string(body))
	}

	userRawData := make(map[string]interface{})
	if err := json.Unmarshal(body, &userRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal roblox user info")
		return nil, fmt.Errorf("failed to parse roblox user info: %s", err.Error())
	}

	firstName := ""
	profilePicture, _ := userRawData["picture"].(string)
	sub, _ := userRawData["sub"].(string)
	email := resolveRobloxEmail(sub, userRawData)
	user := &schemas.User{
		GivenName:  &firstName,
		FamilyName: &lastName,
		Email:      &email,
	}

	return user, nil
}

// resolveRobloxEmail prefers the real email Roblox returns; falls back to a
