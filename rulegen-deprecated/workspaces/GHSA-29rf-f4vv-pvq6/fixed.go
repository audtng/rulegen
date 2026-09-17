package main

		scopeString := sessionSplit[3]
		scopes := parseScopes(scopeString)
		var user *schemas.User
		// providerEmailVerified is the provider's own assertion that the
		// principal controls the email it returned. It gates every path that
		// could attach this login to a pre-existing local account — see the
		// linking branch below.
		var providerEmailVerified bool
		oauthCode := ctx.Request.FormValue("code")
		if oauthCode == "" {
			log.Debug().Err(err).Msg("Invalid oauth code")
		}
		switch provider {
		case constants.AuthRecipeMethodGoogle:
			user, providerEmailVerified, err = h.processGoogleUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodGithub:
			user, providerEmailVerified, err = h.processGithubUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodFacebook:
			user, providerEmailVerified, err = h.processFacebookUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodLinkedIn:
			user, providerEmailVerified, err = h.processLinkedInUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodApple:
			var appleUser *AppleUserInfo
			appleUser, err = parseAppleUserField(ctx.Request.FormValue("user"))
				ctx.JSON(400, gin.H{"error": "invalid apple user info"})
				return
			}
			user, providerEmailVerified, err = h.processAppleUserInfo(ctx, oauthCode, appleUser)
		case constants.AuthRecipeMethodDiscord:
			user, providerEmailVerified, err = h.processDiscordUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodTwitter:
			// Twitter/X uses PKCE: retrieve the verifier stored at login keyed by state.
			verifier, verr := h.MemoryStoreProvider.GetAndRemoveState(pkceVerifierKeyPrefix + state)
				ctx.JSON(400, gin.H{"error": "invalid oauth state"})
				return
			}
			user, providerEmailVerified, err = h.processTwitterUserInfo(ctx, oauthCode, verifier)
		case constants.AuthRecipeMethodMicrosoft:
			user, providerEmailVerified, err = h.processMicrosoftUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodTwitch:
			user, providerEmailVerified, err = h.processTwitchUserInfo(ctx, oauthCode)
		case constants.AuthRecipeMethodRoblox:
			user, providerEmailVerified, err = h.processRobloxUserInfo(ctx, oauthCode)
		default:
			log.Debug().Err(err).Msg("Invalid oauth provider")
			err = fmt.Errorf(`invalid oauth provider`)
		log := log.With().Str("email", refs.StringValue(user.Email)).Logger()
		isSignUp := false

		// An email the identity provider has not attested is attacker-controlled
		// input, and every branch below keys the local account off it — the
		// lookup above decides signup vs. login, and the login branch merges
		// this federated identity into whatever account already holds the
		// address. That is the nOAuth account-takeover class: register a free
		// Entra tenant, set a user's mutable `email` attribute to the victim's
		// address, sign in, and land in the victim's session. The pre-hijack
		// guard further down does not help — it only removes *unverified* local
		// accounts, and verified accounts are exactly what gets stolen.
		//
		// OAuth itself proves nothing about email; that is precisely why OIDC
		// carries a separate `email_verified` claim (Core §5.1), and why Auth0
		// documents checking it before linking accounts.
		if !providerEmailVerified && !h.allowUnverifiedProviderEmail(provider, existingUser, err == nil) {
			log.Debug().Str("provider", provider).Msg("Provider did not attest the email address; refusing to resolve a local account")
			metrics.RecordAuthEvent(metrics.EventOAuthCallback, metrics.StatusFailure)
			metrics.RecordSecurityEvent("oauth_email_unverified", provider)
			h.AuditProvider.LogEvent(audit.Event{
				Action:       constants.AuditOAuthCallbackFailedEvent,
				ActorType:    constants.AuditActorTypeUser,
				ResourceType: constants.AuditResourceTypeSession,
				Metadata:     provider,
				IPAddress:    utils.GetIP(ctx.Request),
				UserAgent:    utils.GetUserAgent(ctx.Request),
			})
			ctx.JSON(400, gin.H{
				"error":             "email_not_verified",
				"error_description": "The identity provider did not confirm that you own this email address.",
			})
			return
		}

		if err != nil {
			isSignupEnabled := h.Config.EnableSignup
			if !isSignupEnabled {
			}

			user.Roles = strings.Join(inputRoles, ",")
			// Only record the address as verified when the provider actually
			// attested it. This used to be unconditional, which meant a
			// compatibility-mode signup from an unattested address wrote
			// email_verified=true into our own database — a claim we cannot
			// back, and one that downstream consumers trust (SAML IdP issuance
			// refuses to assert an unverified email as the Subject NameID).
			if providerEmailVerified {
				now := time.Now().Unix()
				user.EmailVerifiedAt = &now
			}
			user, err = h.StorageProvider.AddUser(ctx, user)
			if err != nil {
				log.Debug().Err(err).Msg("Failed to add user")
			// was never verified, do not link the OAuth identity to it.
			// Instead, delete the unverified account and treat as a new signup
			// for the OAuth user who actually controls the email address.
			//
			// Scoped to accounts some OTHER credential created. An unverified
			// account this same provider already owns is not a squatter — it is
			// this same principal's own account, created on a previous pass
			// through the signup branch above (which, correctly, no longer marks
			// an unattested address verified). Deleting it would recreate the
			// account on every single login, silently dropping its id, roles and
			// org memberships each time.
			if existingUser.EmailVerifiedAt == nil && !signupMethodsContain(existingUser.SignupMethods, provider) {
				// Deleting is only safe for an account that is actually a
				// squatter — created to intercept this address and never used.
				// The cascade is clean (#749) but total: an account carrying
				// real state would lose its org memberships, enrolled
				// authenticators and federated identities outright, and its FGA
				// grants would be orphaned, since the tuple purge lives in the
				// service layer and this is a direct StorageProvider call. All
				// of that on the say-so of an unauthenticated callback.
				// Refusing is recoverable; deleting is not.
				if hasState, what := h.accountHasState(ctx, existingUser); hasState {
					log.Warn().
						Str("reason", what).
						Str("existing_user_id", existingUser.ID).
						Msg("Refusing OAuth login: an unverified account with this email holds state and must not be replaced")
					metrics.RecordAuthEvent(metrics.EventOAuthCallback, metrics.StatusFailure)
					metrics.RecordSecurityEvent("oauth_email_collision_stateful_account", provider)
					h.AuditProvider.LogEvent(audit.Event{
						Action:       constants.AuditOAuthCallbackFailedEvent,
						ActorID:      existingUser.ID,
						ActorType:    constants.AuditActorTypeUser,
						ActorEmail:   refs.StringValue(existingUser.Email),
						ResourceType: constants.AuditResourceTypeSession,
						Metadata:     provider,
						IPAddress:    utils.GetIP(ctx.Request),
						UserAgent:    utils.GetUserAgent(ctx.Request),
					})
					ctx.JSON(400, gin.H{
						"error":             "email_already_registered",
						"error_description": "An unverified account already exists for this email address. Verify it first — request a new verification email for this address, or sign in with the method that created the account.",
					})
					return
				}
				log.Info().Str("existing_user_id", existingUser.ID).Msg("Removing unverified pre-existing account before OAuth signup")
				// Audited: this destroys an account row, which is
				// security-material even when the account was empty.
				h.AuditProvider.LogEvent(audit.Event{
					Action:       constants.AuditOAuthUnverifiedAccountReplacedEvent,
					ActorID:      existingUser.ID,
					ActorType:    constants.AuditActorTypeUser,
					ActorEmail:   refs.StringValue(existingUser.Email),
					ResourceType: constants.AuditResourceTypeUser,
					ResourceID:   existingUser.ID,
					Metadata:     provider,
					IPAddress:    utils.GetIP(ctx.Request),
					UserAgent:    utils.GetUserAgent(ctx.Request),
				})
				if err := h.StorageProvider.DeleteUser(ctx, existingUser); err != nil {
					log.Debug().Err(err).Msg("Failed to delete unverified user")
					ctx.JSON(500, gin.H{"error": "failed to process OAuth login"})
	}
}

// allowUnverifiedProviderEmail decides whether a federated login whose provider
// did NOT attest the email address may still resolve a local account.
//
// Default (--oauth-allow-unverified-provider-email=false): never. The address is
// attacker-controlled and it is what selects the account.
//
// Compatibility mode (=true) exists so a deployment upgrading from 2.3.x is not
// locked out the moment it restarts, but it is deliberately NOT a plain "turn
// the check off" switch — that would restore the CVE verbatim. Even in this
// mode, an unattested address may only:
//
//   - create a brand-new account (it selects nobody, so it harms nobody), or
//   - return to an account THIS SAME PROVIDER already owns — a returning user.
//
// It may never merge into an account some other credential owns. That single
// restriction removes the entire cross-credential takeover: an Entra tenant
// cannot reach a password account, a Google account, or any other provider's
// account, which is every practical form of the attack.
//
// The residual risk it does not cover, and the reason this mode is documented
// as temporary: two principals of the SAME unattested provider (two Entra
// tenants both asserting one address) can still collide. Pinning
// --microsoft-tenant-id or setting --microsoft-allowed-tenants closes that, and
// is the actual fix.
func (h *httpProvider) allowUnverifiedProviderEmail(provider string, existingUser *schemas.User, found bool) bool {
	if !h.Config.OAuthAllowUnverifiedProviderEmail {
		return false
	}
	if !found {
		// First-time signup: no account is being selected away from anyone.
		return true
	}
	if existingUser == nil {
		// "Found" with no row is an inconsistent storage result. Fail closed
		// rather than guess which case it was.
		return false
	}
	// Returning user of this same provider, or an attempt to cross into an
	// account another credential owns.
	return signupMethodsContain(existingUser.SignupMethods, provider)
}

// signupMethodsContain reports whether a stored comma-separated signup-methods
// list contains an exact method. Deliberately not strings.Contains: the
// provider names include the near-miss pair twitch/twitter, and a substring
// match on a security decision would let one provider inherit the other's
// accounts.
func signupMethodsContain(signupMethods, provider string) bool {
	for _, m := range strings.Split(signupMethods, ",") {
		if strings.TrimSpace(m) == provider {
			return true
		}
	}
	return false
}

// flexBool decodes a JSON boolean that some IdPs send quoted. Apple documents
// `email_verified` as "a string or Boolean value", and LinkedIn's userinfo has
// shipped both shapes; decoding either into a plain bool fails the whole claim
// set, which would silently turn a verified email into an unverified one.
type flexBool bool

// UnmarshalJSON accepts true/false, "true"/"false", or anything else (which
// decodes to false — unrecognised is never "verified").
func (b *flexBool) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*b = flexBool(claimTruthy(raw))
	return nil
}

// claimTruthy reads a boolean claim out of an untyped JSON value, tolerating
// the quoted-string form some IdPs emit. Used by the providers whose payloads
// are decoded into a map rather than oidcClaims (Apple, Discord, Roblox).
func claimTruthy(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return strings.EqualFold(t, "true")
	}
	return false
}

// oidcClaims is the allow-list of OpenID Connect standard claims Authorizer
// maps onto a user. ID tokens are decoded into this and never straight into
// schemas.User, for two reasons:
//     signup_methods, is_active, created_at ...) merely by sharing its json
//     tag.
type oidcClaims struct {
	// Subject is the provider-asserted stable identifier for the principal.
	// Unlike `email` it is not user-mutable, so it is the only claim safe to
	// treat as an identity key.
	Subject string `json:"sub"`
	// Issuer and TenantID back the Microsoft tenant checks; see
	// processMicrosoftUserInfo. Ignored for every other provider.
	Issuer   string `json:"iss"`
	TenantID string `json:"tid"`
	// EmailVerified is the provider's assertion that the principal actually
	// controls `Email`. Absent decodes to false — an IdP that does not say
	// "verified" has not verified anything, and this claim is what stops a
	// federated login from linking to somebody else's account.
	EmailVerified flexBool `json:"email_verified"`
	// XmsEdov ("email domain owner verified") is Microsoft Entra's equivalent.
	// Entra v2 ID tokens carry no `email_verified` claim at all, and their
	// `email` is a mutable, unverified profile attribute — the distinction
	// that makes the nOAuth attack work.
	XmsEdov flexBool `json:"xms_edov"`

	Email       string `json:"email"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	return user
}

func (h *httpProvider) processGoogleUserInfo(ctx *gin.Context, code string) (*schemas.User, bool, error) {
	log := h.Log.With().Str("func", "processGoogleUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodGoogle)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, false, fmt.Errorf("error getting oauth config: %s", err.Error())
	}
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, false, fmt.Errorf("invalid google exchange code: %s", err.Error())
	}

	issuer := "https://accounts.google.com"
	}
	oidcProvider, err := getOIDCProvider(ctx, issuer)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create oidc provider: %s", err.Error())
	}
	verifier := oidcProvider.Verifier(&oidc.Config{ClientID: h.GoogleClientID})
	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Debug().Err(err).Msg("Failed to extract ID Token from OAuth2 token")
		return nil, false, fmt.Errorf("unable to extract id_token")
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to verify ID Token")
		return nil, false, fmt.Errorf("unable to verify id_token: %s", err.Error())
	}
	claims := &oidcClaims{}
	if err := idToken.Claims(claims); err != nil {
		log.Debug().Err(err).Msg("Failed to parse ID Token claims")
		return nil, false, fmt.Errorf("unable to extract claims")
	}

	// Google asserts control of the address via `email_verified`.
	return claims.toUser(), bool(claims.EmailVerified), nil
}

// setGithubHeaders applies the headers GitHub's REST API docs ask every
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}

func (h *httpProvider) processGithubUserInfo(ctx *gin.Context, code string) (*schemas.User, bool, error) {
	log := h.Log.With().Str("func", "processGithubUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodGithub)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, false, fmt.Errorf("error getting oauth config: %s", err.Error())
	}

	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, false, fmt.Errorf("invalid github exchange code: %s", err.Error())
	}
	userInfoURL := constants.GithubUserInfoURL
	emailsURL := constants.GithubUserEmails
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create github user info request")
		return nil, false, fmt.Errorf("error creating github user info request: %s", err.Error())
	}
	setGithubHeaders(req, oauth2Token.AccessToken)

	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to request github user info")
		return nil, false, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read github user info response body")
		return nil, false, fmt.Errorf("failed to read github response body: %s", err.Error())
	}
	if response.StatusCode >= 400 {
		log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request github user info")
		return nil, false, fmt.Errorf("failed to request github user info: %s", string(body))
	}

	// Only the three fields below are used. A typed struct (rather than a
	}
	if err := json.Unmarshal(body, &userRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal github user info")
		return nil, false, fmt.Errorf("failed to parse github user info: %s", err.Error())
	}

	name := strings.Split(userRawData.Name, " ")
		req, err := http.NewRequest(http.MethodGet, emailsURL, nil)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to create github emails request")
			return nil, false, fmt.Errorf("error creating github user info request: %s", err.Error())
		}
		setGithubHeaders(req, oauth2Token.AccessToken)

		response, err := client.Do(req)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to request github user email")
			return nil, false, err
		}

		defer func() { _ = response.Body.Close() }()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to read github user email response body")
			return nil, false, fmt.Errorf("failed to read github response body: %s", err.Error())
		}
		if response.StatusCode >= 400 {
			log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request github user email")
			return nil, false, fmt.Errorf("failed to request github user info: %s", string(body))
		}

		emailData := []GithubUserEmails{}
		err = json.Unmarshal(body, &emailData)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to parse github user email")
			return nil, false, fmt.Errorf("failed to parse github user email: %s", err.Error())
		}

		// GET /user/emails lists every address on the account, verified or
		}
		if email == "" {
			log.Debug().Msg("No verified email on github account")
			return nil, false, fmt.Errorf("failed to get a verified email address from github")
		}
	}

		Email:      &email,
	}

	return user, true, nil
}

func (h *httpProvider) processFacebookUserInfo(ctx *gin.Context, code string) (*schemas.User, bool, error) {
	log := h.Log.With().Str("func", "processFacebookUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodFacebook)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, false, fmt.Errorf("error getting oauth config: %s", err.Error())
	}
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Invalid facebook exchange code")
		return nil, false, fmt.Errorf("invalid facebook exchange code: %s", err.Error())
	}
	userInfoURL := constants.FacebookUserInfoURL
	if mockBase := h.TestOAuthBaseURL(constants.AuthRecipeMethodFacebook); mockBase != "" {
	req, err := http.NewRequest("GET", userInfoURL+oauth2Token.AccessToken, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Error creating facebook user info request")
		return nil, false, fmt.Errorf("error creating facebook user info request: %s", err.Error())
	}

	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to process facebook user")
		return nil, false, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read facebook response")
		return nil, false, fmt.Errorf("failed to read facebook response body: %s", err.Error())
	}
	if response.StatusCode >= 400 {
		log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request facebook user info")
		return nil, false, fmt.Errorf("failed to request facebook user info: %s", string(body))
	}
	// Typed decode, not fmt.Sprintf over a map: Graph API omits `email`
	// entirely when "no valid email address is available" (user/reference/user),
	}
	if err := json.Unmarshal(body, &userRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal facebook user info")
		return nil, false, fmt.Errorf("failed to parse facebook user info: %s", err.Error())
	}

	email := userRawData.Email
	if email == "" {
		log.Debug().Msg("Facebook user info has no email")
		return nil, false, fmt.Errorf("failed to get email from facebook user info: the account has no available email address")
	}

	picture := userRawData.Picture.Data.URL
		Email:      &email,
	}

	return user, true, nil
}

func (h *httpProvider) processLinkedInUserInfo(ctx *gin.Context, code string) (*schemas.User, bool, error) {
	log := h.Log.With().Str("func", "processLinkedInUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodLinkedIn)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, false, fmt.Errorf("error getting oauth config: %s", err.Error())
	}

	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, false, fmt.Errorf("invalid linkedin exchange code: %s", err.Error())
	}

	userInfoURL := constants.LinkedInUserInfoURL
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create linkedin user info request")
		return nil, false, fmt.Errorf("error creating linkedin user info request: %s", err.Error())
	}
	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", oauth2Token.AccessToken)},
	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to request linkedin user info")
		return nil, false, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read linkedin user info response body")
		return nil, false, fmt.Errorf("failed to read linkedin response body: %s", err.Error())
	}

	if response.StatusCode >= 400 {
		log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request linkedin user info")
		return nil, false, fmt.Errorf("failed to request linkedin user info: %s", string(body))
	}

	// OIDC userinfo shape (sub/name/given_name/family_name/picture/locale/
	// email/email_verified) - one call, no separate /v2/emailAddress hop.
	var userRawData struct {
		GivenName     string   `json:"given_name"`
		FamilyName    string   `json:"family_name"`
		Picture       string   `json:"picture"`
		Email         string   `json:"email"`
		EmailVerified flexBool `json:"email_verified"`
	}
	if err := json.Unmarshal(body, &userRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal linkedin user info")
		return nil, false, fmt.Errorf("failed to parse linkedin user info: %s", err.Error())
	}

	// `email` is documented as optional - it is only present when the member
	// than a synthetic-email fallback.
	if userRawData.Email == "" {
		log.Debug().Msg("LinkedIn user info has no email")
		return nil, false, fmt.Errorf("failed to extract email from linkedin response")
	}

	user := &schemas.User{
		Email:      &userRawData.Email,
	}

	return user, bool(userRawData.EmailVerified), nil
}

func (h *httpProvider) processAppleUserInfo(ctx *gin.Context, code string, appleUser *AppleUserInfo) (*schemas.User, bool, error) {
	log := h.Log.With().Str("func", "processAppleUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodApple)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, false, fmt.Errorf("error getting oauth config: %s", err.Error())
	}

	var user = &schemas.User{}
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return user, false, fmt.Errorf("invalid apple exchange code: %s", err.Error())
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Debug().Err(err).Msg("Failed to extract ID Token from OAuth2 token")
		return user, false, fmt.Errorf("unable to extract id_token")
	}

	// Verify the Apple ID token signature, issuer, and audience using OIDC discovery
	oidcProvider, err := getOIDCProvider(ctx, issuer)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create Apple OIDC provider")
		return user, false, fmt.Errorf("failed to create oidc provider: %s", err.Error())
	}
	verifier := oidcProvider.Verifier(&oidc.Config{ClientID: h.AppleClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to verify Apple ID Token")
		return user, false, fmt.Errorf("unable to verify id_token: %s", err.Error())
	}

	claims := make(map[string]interface{})
	if err := idToken.Claims(&claims); err != nil {
		log.Debug().Err(err).Msg("Failed to parse Apple ID Token claims")
		return user, false, fmt.Errorf("failed to parse claims: %s", err.Error())
	}

	if val, ok := claims["email"]; !ok || val == nil {
		log.Debug().Msg("Failed to extract email from claims.")
		return user, false, fmt.Errorf("unable to extract email, please check the scopes enabled for your app. It needs `email`, `name` scopes")
	} else {
		email, _ := val.(string)
		user.Email = &email
	}

	// Apple documents `email_verified` as "a string or Boolean value", so it
	// arrives as either true or "true" — claimTruthy accepts both. Absent means
	// unverified.
	emailVerified := claimTruthy(claims["email_verified"])

	user.GivenName = &appleUser.Name.FirstName
	user.FamilyName = &appleUser.Name.LastName

	return user, emailVerified, nil
}

// processDiscordUserInfo exchanges the Discord OAuth code for the user's
// creating a duplicate account - the same fallback discipline
// processTwitterUserInfo uses above for X, which never returns a real email
// at all.
func (h *httpProvider) processDiscordUserInfo(ctx *gin.Context, code string) (*schemas.User, bool, error) {
	log := h.Log.With().Str("func", "processDiscordUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodDiscord)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, false, fmt.Errorf("error getting oauth config: %s", err.Error())
	}
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, false, fmt.Errorf("invalid discord exchange code: %s", err.Error())
	}

	userInfoURL := constants.DiscordUserInfoURL
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create Discord user info request")
		return nil, false, fmt.Errorf("error creating Discord user info request: %s", err.Error())
	}
	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", oauth2Token.AccessToken)},
	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to request Discord user info")
		return nil, false, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read Discord user info response body")
		return nil, false, fmt.Errorf("failed to read Discord response body: %s", err.Error())
	}

	if response.StatusCode >= 400 {
		log.Debug().Err(err).Msg("Failed to request Discord user info")
		return nil, false, fmt.Errorf("failed to request Discord user info: %s", string(body))
	}

	// Unmarshal the response body into a map. GET /users/@me returns a flat
	userRawData := make(map[string]interface{})
	if err := json.Unmarshal(body, &userRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal Discord response")
		return nil, false, fmt.Errorf("failed to unmarshal Discord response: %s", err.Error())
	}

	// Extract the username
	firstName, ok := userRawData["username"].(string)
	if !ok {
		log.Debug().Err(err).Msg("Username is not in expected format or missing in user data")
		return nil, false, fmt.Errorf("username is not in expected format or missing in user data")
	}
	discordID, ok := userRawData["id"].(string)
	if !ok || discordID == "" {
		log.Debug().Msg("Discord user info missing id")
		return nil, false, fmt.Errorf("discord response missing id field")
	}
	// `avatar` is nullable (?string in Discord's user object) for accounts on
	// the default avatar - building the CDN URL from an empty hash yields a
	}

	email := resolveDiscordEmail(discordID, userRawData)
	// GET /users/@me carries a `verified` flag for the account's email. The
	// synthetic fallback is trusted by construction: it lives on a reserved
	// non-routable domain keyed by Discord's permanent id, so it can never
	// collide with an address a real person could prove they own.
	emailVerified := claimTruthy(userRawData["verified"])
	if email == discordSyntheticEmail(discordID) {
		emailVerified = true
	}

	user := &schemas.User{
		GivenName: &firstName,
		Email:     &email,
	}

	return user, emailVerified, nil
}

// resolveDiscordEmail prefers the real email Discord returns; falls back to
// returning Twitter user instead of creating a duplicate account on every
// login. Operators who opt into X's `users.email` scope + app permission get
// a real confirmed_email instead (see TwitterUserInfoURL's doc comment).
func (h *httpProvider) processTwitterUserInfo(ctx *gin.Context, code, verifier string) (*schemas.User, bool, error) {
	log := h.Log.With().Str("func", "processTwitterUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodTwitter)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, false, fmt.Errorf("error getting oauth config: %s", err.Error())
	}

	oauth2Token, err := cfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, false, fmt.Errorf("invalid twitter exchange code: %s", err.Error())
	}

	userInfoURL := constants.TwitterUserInfoURL
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create Twitter user info request")
		return nil, false, fmt.Errorf("error creating Twitter user info request: %s", err.Error())
	}
	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", oauth2Token.AccessToken)},
	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to request Twitter user info")
		return nil, false, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read Twitter user info response body")
		return nil, false, fmt.Errorf("failed to read Twitter response body: %s", err.Error())
	}

	if response.StatusCode >= 400 {
		log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request Twitter user info")
		return nil, false, fmt.Errorf("failed to request Twitter user info: %s", string(body))
	}

	responseRawData := make(map[string]interface{})
	if err := json.Unmarshal(body, &responseRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal twitter user info")
		return nil, false, fmt.Errorf("failed to parse twitter user info: %s", err.Error())
	}

	userRawData, ok := responseRawData["data"].(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("twitter response missing data field")
	}

	// Twitter API does not return E-Mail adresses by default. For that case special privileges have
	twitterID, ok := userRawData["id"].(string)
	if !ok || twitterID == "" {
		log.Debug().Msg("Twitter user info missing id")
		return nil, false, fmt.Errorf("twitter response missing id field")
	}

	// Currently Twitter API only provides the full name of a user. To fill givenName and familyName
	profilePicture, _ := userRawData["profile_image_url"].(string)

	email := resolveTwitterEmail(twitterID, userRawData)
	// X only returns `confirmed_email` to apps granted the `users.email` scope
	// plus the app-dashboard permission, and the name says it: X has confirmed
	// it. The synthetic fallback is trusted by construction — a reserved
	// non-routable domain keyed by X's permanent numeric id, which no real
	// mailbox can occupy.
	emailVerified := true

	user := &schemas.User{
		Email:      &email,
		Nickname:   &nickname,
	}

	return user, emailVerified, nil
}

// twitterSyntheticEmail derives a stable, non-routable synthetic email from
}

// process microsoft user information
func (h *httpProvider) processMicrosoftUserInfo(ctx *gin.Context, code string) (*schemas.User, bool, error) {
	log := h.Log.With().Str("func", "processMicrosoftUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodMicrosoft)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, false, fmt.Errorf("error getting oauth config: %s", err.Error())
	}
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, false, fmt.Errorf("invalid microsoft exchange code: %s", err.Error())
	}
	issuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", h.MicrosoftTenantID)
	if mockBase := h.TestOAuthBaseURL(constants.AuthRecipeMethodMicrosoft); mockBase != "" {
	}
	oidcProvider, err := getOIDCProvider(ctx, issuer)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create oidc provider: %s", err.Error())
	}
	// The multi-tenant discovery documents ("common"/"organizations"/
	// "consumers") advertise a templated issuer containing {tenantid}, which
	// never literally equals the `iss` of a real token, so go-oidc's built-in
	// comparison cannot be used. Skipping it is not the same as not checking:
	// validateMicrosoftTenant below reconstructs the expected issuer from the
	// token's own `tid` and enforces it, plus the operator's tenant policy.
	verifier := oidcProvider.Verifier(&oidc.Config{
		ClientID:        h.MicrosoftClientID,
		SkipIssuerCheck: true,
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Debug().Err(err).Msg("Failed to extract ID Token from OAuth2 token")
		return nil, false, fmt.Errorf("unable to extract id_token")
	}
	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to verify ID Token")
		return nil, false, fmt.Errorf("unable to verify id_token: %s", err.Error())
	}
	claims := &oidcClaims{}
	if err := idToken.Claims(claims); err != nil {
		log.Debug().Err(err).Msg("Failed to parse ID Token claims")
		return nil, false, fmt.Errorf("unable to extract claims")
	}

	// The test double issues tokens from a stand-in issuer with no tenant
	// model at all; tenant policy is meaningless there.
	if mockBase := h.TestOAuthBaseURL(constants.AuthRecipeMethodMicrosoft); mockBase != "" {
		return claims.toUser(), bool(claims.EmailVerified), nil
	}

	tenantPinned, err := validateMicrosoftTenant(claims, h.MicrosoftTenantID, h.Config.MicrosoftAllowedTenants)
	if err != nil {
		log.Debug().Err(err).Str("tid", claims.TenantID).Msg("Microsoft tenant validation failed")
		return nil, false, err
	}

	// Entra v2 ID tokens have no `email_verified` claim, and `email` is a
	// mutable, unverified profile attribute any tenant admin can set to any
	// string — including somebody else's address. Two signals make it
	// trustworthy, and nothing else does:
	//
	//   - xms_edov ("email domain owner verified"), Microsoft's own attestation
	//     that the token's tenant owns the email's domain. It is an optional
	//     claim; operators enable it in the app registration.
	//   - the tenant being pinned or allowlisted, which means the address can
	//     only have come from a directory the operator already trusts.
	//
	// Without either, this is the nOAuth setup: an attacker registers a free
	// Entra tenant, sets a user's `email` to the victim's address, and signs in.
	return claims.toUser(), tenantPinned || bool(claims.XmsEdov), nil
}

// microsoftMultiTenantAliases are the Entra endpoint aliases that accept tokens
// from tenants the operator has never heard of. Any other configured value is a
// specific tenant (a GUID or a verified domain name) and is therefore pinned.
var microsoftMultiTenantAliases = map[string]bool{
	"common":        true,
	"organizations": true,
	"consumers":     true,
}

// validateMicrosoftTenant enforces the operator's tenant policy on a verified
// Entra ID token and reports whether the originating tenant is one the operator
// explicitly trusts.
//
// go-oidc has already checked the signature and that `aud` equals our client id
// — neither of which constrains WHICH tenant minted the token, because the
// multi-tenant endpoints sign with Microsoft's global keys. The tenant is the
// only thing that does, so it is checked here:
//
//   - `iss` must be the issuer the token's own `tid` implies, so a token cannot
//     claim one tenant in `iss` and another in `tid`;
//   - a pinned `--microsoft-tenant-id` must match `tid` exactly;
//   - a non-empty `--microsoft-allowed-tenants` must contain `tid`.
//
// Returns true when the tenant was pinned or allowlisted.
func validateMicrosoftTenant(claims *oidcClaims, configuredTenant string, allowedTenants []string) (bool, error) {
	tid := strings.TrimSpace(claims.TenantID)
	if tid == "" {
		return false, fmt.Errorf("microsoft id_token is missing the tid claim")
	}
	if expected := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tid); claims.Issuer != expected {
		return false, fmt.Errorf("microsoft id_token issuer does not match its tenant")
	}

	if len(allowedTenants) > 0 {
		if !utils.StringSliceContains(allowedTenants, tid) {
			return false, fmt.Errorf("microsoft tenant is not allowed")
		}
		return true, nil
	}

	configuredTenant = strings.TrimSpace(configuredTenant)
	if !microsoftMultiTenantAliases[strings.ToLower(configuredTenant)] {
		// A specific tenant was configured: only that directory may sign in.
		if !strings.EqualFold(configuredTenant, tid) {
			return false, fmt.Errorf("microsoft id_token was issued by an unexpected tenant")
		}
		return true, nil
	}

	// Multi-tenant with no allowlist. The login is still permitted — this is a
	// documented deployment mode — but the tenant is not trusted, so the caller
	// must not treat the email as proof of anything.
	return false, nil
}

// process twitch user information
func (h *httpProvider) processTwitchUserInfo(ctx *gin.Context, code string) (*schemas.User, bool, error) {
	log := h.Log.With().Str("func", "processTwitchUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodTwitch)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, false, fmt.Errorf("error getting oauth config: %s", err.Error())
	}

	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, false, fmt.Errorf("invalid twitch exchange code: %s", err.Error())
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Debug().Err(err).Msg("Failed to extract ID Token from OAuth2 token")
		return nil, false, fmt.Errorf("unable to extract id_token")
	}
	issuer := "https://id.twitch.tv/oauth2"
	if mockBase := h.TestOAuthBaseURL(constants.AuthRecipeMethodTwitch); mockBase != "" {
	oidcProvider, err := getOIDCProvider(ctx, issuer)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create OIDC provider")
		return nil, false, fmt.Errorf("failed to create oidc provider: %s", err.Error())
	}
	verifier := oidcProvider.Verifier(&oidc.Config{
		ClientID:        h.TwitchClientID,
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to verify ID Token")
		return nil, false, fmt.Errorf("unable to verify id_token: %s", err.Error())
	}

	claims := &oidcClaims{}
	if err := idToken.Claims(claims); err != nil {
		log.Debug().Err(err).Msg("Failed to parse ID Token claims")
		return nil, false, fmt.Errorf("unable to extract claims")
	}

	// Twitch is single-issuer (SkipIssuerCheck above is harmless — signature
	// and `aud` already pin the token to Twitch), and its ID token carries the
	// standard `email_verified`.
	return claims.toUser(), bool(claims.EmailVerified), nil
}

// process roblox user information
func (h *httpProvider) processRobloxUserInfo(ctx *gin.Context, code string) (*schemas.User, bool, error) {
	log := h.Log.With().Str("func", "processRobloxUserInfo").Logger()
	cfg, err := h.OAuthProvider.GetOAuthConfig(ctx, constants.AuthRecipeMethodRoblox)
	if err != nil {
		log.Debug().Err(err).Msg("Error getting oauth config")
		return nil, false, fmt.Errorf("error getting oauth config: %s", err.Error())
	}
	// Roblox is a confidential client (client_secret set); PKCE is optional and
	// no code_challenge is sent at login, so no code_verifier is replayed here.
	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to exchange code for token")
		return nil, false, fmt.Errorf("invalid roblox exchange code: %s", err.Error())
	}

	userInfoURL := constants.RobloxUserInfoURL
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create roblox user info request")
		return nil, false, fmt.Errorf("error creating roblox user info request: %s", err.Error())
	}
	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", oauth2Token.AccessToken)},
	response, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to request roblox user info")
		return nil, false, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read roblox user info response body")
		return nil, false, fmt.Errorf("failed to read roblox response body: %s", err.Error())
	}

	if response.StatusCode >= 400 {
		log.Debug().Err(err).Str("body", string(body)).Msg("Failed to request roblox user info")
		return nil, false, fmt.Errorf("failed to request roblox user info: %s", string(body))
	}

	userRawData := make(map[string]interface{})
	if err := json.Unmarshal(body, &userRawData); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal roblox user info")
		return nil, false, fmt.Errorf("failed to parse roblox user info: %s", err.Error())
	}

	firstName := ""
	profilePicture, _ := userRawData["picture"].(string)
	sub, _ := userRawData["sub"].(string)
	email := resolveRobloxEmail(sub, userRawData)
	// Roblox's userinfo is OIDC-standard, so a real address comes with
	// `email_verified`. The synthetic fallback (the default config, which does
	// not request the `email` scope) is trusted by construction: reserved
	// non-routable domain keyed by the permanent `sub`.
	emailVerified := claimTruthy(userRawData["email_verified"])
	if sub != "" && email == robloxSyntheticEmail(sub) {
		emailVerified = true
	}
	user := &schemas.User{
		GivenName:  &firstName,
		FamilyName: &lastName,
		Email:      &email,
	}

	return user, emailVerified, nil
}

// resolveRobloxEmail prefers the real email Roblox returns; falls back to a
package token

import (
	"crypto/subtle"
	"strconv"

	"github.com/authorizerdev/authorizer/internal/metrics"
)

const (
	// adminSecretMaxFailedAttempts / adminSecretLockoutWindowSeconds throttle
	// online guessing of the admin secret.
	//
	// The admin secret is the single highest-privilege credential in the system,
	// and until now the only thing standing between an attacker and unlimited
	// guesses was the shared 30rps request limiter — the same budget ordinary
	// traffic gets. The budget here is deliberately much tighter and much
	// longer-lived than the login lockout: nobody legitimately mistypes an
	// admin secret dozens of times, and unlike a user account there is no
	// self-service recovery an attacker could grief by tripping it.
	adminSecretMaxFailedAttempts    = 10
	adminSecretLockoutWindowSeconds = int64(15 * 60)
	// adminSecretLockoutPrefix namespaces the counter. Keyed by client IP:
	// there is only one admin secret, so a per-principal key would be a single
	// global counter that any attacker could use to lock out every operator.
	adminSecretLockoutPrefix = "admin_secret_failed_attempts:"
	// adminSecretLockoutUnknownIP is the bucket for callers whose address we
	// could not determine. It is deliberately a NAMED bucket rather than the
	// empty string: keying on "" silently merges every unidentifiable caller
	// into one counter, which is how a pure-gRPC deployment (no forwarded
	// headers, no peer address) turns 10 wrong guesses into an outage for every
	// admin client at once. Callers should ensure this is never needed —
	// transport.MetaFromGRPC falls back to the gRPC peer address for exactly
	// that reason — but if it is, the shared bucket must at least be visible in
	// the key rather than looking like a real client.
	adminSecretLockoutUnknownIP = "unknown"
)

// VerifyAdminSecret is the single gate for every admin-secret comparison —
// AdminLogin's cookie-establishing check and the x-authorizer-admin-secret
// header path both route through it, so neither can be brute-forced while the
// other is throttled.
//
// FAILED attempts are counted, not all attempts. The counter is read before the
// comparison and incremented only when the comparison fails. The alternative —
// increment-then-check, which the login and OTP lockouts use — is right for
// those because they gate a human typing a password a few times a minute, but
// wrong here: this same function runs on EVERY request that authenticates with
// the x-authorizer-admin-secret header, so counting successes means an
// integration issuing more than adminSecretMaxFailedAttempts concurrent admin
// calls from one address gets 401s while presenting the CORRECT secret. The
// concurrency argument for increment-then-check does not apply either: it exists
// to stop parallel requests reading one stale pre-increment count and all
// passing, which matters when each parallel request is an independent guess at a
// short secret. Here every parallel request carries the same operator-chosen
// secret, and overshooting the budget by the in-flight count costs an attacker
// nothing they did not already have.
//
// Returns (valid, locked). A locked caller never reaches the comparison at all,
// so the lockout cannot itself be used as a timing oracle for the secret.
//
// What this does NOT protect against, stated plainly so it is not mistaken for a
// boundary:
//
//   - clientIP comes from utils.GetIP / RequestMetadata, which prefer the
//     X-Real-Ip and X-Forwarded-For request headers. On a deployment that is not
//     behind a proxy that overwrites them, those are attacker-controlled: a
//     guesser rotates the header per request and never fills a bucket. Fixing
//     that needs trusted-proxy configuration this server does not yet have.
//   - a distributed guesser gets adminSecretMaxFailedAttempts per source
//     regardless.
//
// It is defence in depth that makes naive online guessing expensive, not a
// substitute for a high-entropy AdminSecret. The entropy is the control.
func (p *provider) VerifyAdminSecret(clientIP, candidate string) (valid bool, locked bool) {
	// An unconfigured secret must never authenticate anything, empty candidate
	// included.
	if p.config.AdminSecret == "" || candidate == "" {
		return false, false
	}

	// The throttle is defence in depth around the comparison, never a
	// precondition for it. A provider built without a memory store (unit tests,
	// and any future wiring that omits it) must still authenticate correctly
	// rather than nil-panic in a request path — an unrecovered panic here would
	// take down the whole process, turning a missing dependency into an outage.
	if p.dependencies == nil || p.dependencies.MemoryStoreProvider == nil {
		return subtle.ConstantTimeCompare([]byte(candidate), []byte(p.config.AdminSecret)) == 1, false
	}

	if clientIP == "" {
		clientIP = adminSecretLockoutUnknownIP
	}
	lockKey := adminSecretLockoutPrefix + clientIP
	// Fail open on a store fault: an outage must not lock every operator out of
	// their own admin console. Same stance as the login/OTP paths.
	if spent, err := p.dependencies.MemoryStoreProvider.GetCache(lockKey); err != nil {
		p.dependencies.Log.Debug().Err(err).Msg("Failed to read admin-secret failed-attempt counter")
	} else if attempts, _ := strconv.ParseInt(spent, 10, 64); attempts >= adminSecretMaxFailedAttempts {
		metrics.RecordSecurityEvent("admin_secret_locked", "admin_auth")
		p.dependencies.Log.Warn().Int64("attempts", attempts).Str("ip", clientIP).Msg("Admin secret verification locked: too many failed attempts")
		return false, true
	}

	if subtle.ConstantTimeCompare([]byte(candidate), []byte(p.config.AdminSecret)) != 1 {
		if _, err := p.dependencies.MemoryStoreProvider.IncrementCache(lockKey, adminSecretLockoutWindowSeconds); err != nil {
			p.dependencies.Log.Debug().Err(err).Msg("Failed to increment admin-secret failed-attempt counter")
		}
		return false, false
	}

	// Correct secret: clear the budget so a legitimate operator who fat-fingered
	// it a few times starts fresh.
	if err := p.dependencies.MemoryStoreProvider.DeleteCacheByPrefix(lockKey); err != nil {
		p.dependencies.Log.Debug().Err(err).Msg("Failed to reset admin-secret failed-attempt counter")
	}
	return true, false
}
package http_handlers

import (
	"context"

	"github.com/authorizerdev/authorizer/internal/authorization/engine"
	"github.com/authorizerdev/authorizer/internal/constants"
	"github.com/authorizerdev/authorizer/internal/storage"
	"github.com/authorizerdev/authorizer/internal/storage/schemas"
)

// accountHasState reports whether a user account holds state that would be
// silently destroyed by deleting it, and names the first thing found.
//
// This exists to bound the pre-hijack guard in OAuthCallbackHandler. That guard
// deletes an *unverified* pre-existing account rather than linking a federated
// identity to it, which is correct for its intended target — an attacker who
// signed up with someone else's address and never verified it, squatting to
// intercept their later social login. A squatter's account is empty by
// definition: created seconds ago, never used.
//
// A real account is not. StorageProvider.DeleteUser cascades to every
// user-keyed table (schemas.UserOwnedCollections, #749), so nothing is left
// dangling — but "cleanly destroyed" is not "safe to destroy". The replacement
// account gets a fresh id, so deleting a real account silently drops its org
// memberships, its enrolled authenticators and passkeys, and its federated
// identities, and the user gets none of it back. Two things are worse than
// merely lost:
//
//   - FGA grants are NOT covered by that cascade. The tuple store lives outside
//     StorageProvider, and the purge that admin _delete_user runs
//     (service.purgeFgaTuplesForUser) is in the service layer, which this
//     handler's direct StorageProvider.DeleteUser call does not go through. So
//     `user:<dead-id>` grants persist forever while the new account inherits
//     none of them.
//   - the whole thing is triggered by an UNAUTHENTICATED OAuth callback. Nobody
//     proved they own the account being destroyed; they merely presented a
//     provider assertion for the same address.
//
// So an account carrying any of this is not a squatter, and must never be
// deleted to resolve an email collision. The caller refuses the login instead,
// which is recoverable; deletion is not.
//
// Fail-closed by design: a storage or FGA fault reports "has state" (the second
// return names it), because the safe answer when we cannot tell is to refuse
// rather than destroy.
func (h *httpProvider) accountHasState(ctx context.Context, user *schemas.User) (bool, string) {
	if user == nil || user.ID == "" {
		return false, ""
	}

	if memberships, _, err := h.StorageProvider.ListOrgMembershipsByUser(ctx, user.ID, nil); err != nil {
		return true, "org membership lookup failed"
	} else if len(memberships) > 0 {
		return true, "org membership"
	}

	if creds, err := h.StorageProvider.ListWebauthnCredentialsByUserID(ctx, user.ID); err != nil {
		return true, "passkey lookup failed"
	} else if len(creds) > 0 {
		return true, "passkey"
	}

	for _, authenticatorType := range []string{
		constants.EnvKeyTOTPAuthenticator,
		constants.EnvKeyEmailOTPAuthenticator,
		constants.EnvKeySMSOTPAuthenticator,
	} {
		// "Not enrolled" is the normal case and every backend reports it as an
		// error (the not-found contract), so it must be told apart from "the
		// lookup failed" — swallowing both would make a transient storage fault
		// look like an empty account and delete somebody's MFA enrollment.
		a, err := h.StorageProvider.GetAuthenticatorDetailsByUserId(ctx, user.ID, authenticatorType)
		switch {
		case storage.IsNotFound(err):
			continue
		case err != nil:
			return true, "authenticator lookup failed"
		case a != nil:
			return true, "enrolled authenticator"
		}
	}

	// FGA grants are the ones the user notices losing and the ones no cascade
	// would ever clean up. Only ask the engine when one is configured.
	if h.AuthzEngine != nil {
		res, err := h.AuthzEngine.ReadTuples(ctx, engine.ReadTuplesFilter{
			User:     "user:" + user.ID,
			PageSize: 1,
		})
		if err != nil {
			return true, "fga tuple lookup failed"
		}
		if res != nil && len(res.Tuples) > 0 {
			return true, "fga grant"
		}
	}

	return false, ""
}
