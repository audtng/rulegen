package main

	}

	// execute celRoleEntry.AuthProgram
	pbAuth, err := b.runCelProgram(ctx, celRoleEntry, allClaims)
	if err != nil {
		return logical.ErrorResponse("error executing cel program: %s", err.Error()), nil
	}
}

// runCelProgram executes the CelProgram for the celRoleEntry and returns a pb.Auth or error
func (b *jwtAuthBackend) runCelProgram(ctx context.Context, celRoleEntry *celRoleEntry, allClaims map[string]any) (*pb.Auth, error) {
	result, err := b.celEvalProgram(celRoleEntry.CelProgram, allClaims)
	if err != nil {
		return nil, fmt.Errorf("Cel role auth program failed: %w", err)
	}
			if !ok {
				t.Fatalf("Expected jwtAuthBackend, got %T", logicalBackend)
			}
			role, err := b.runCelProgram(context.Background(), &tc.celRole, tc.claims)
			if tc.validateResult != nil {
				tc.validateResult(t, err, role)
			}
		},
		"expiration_leeway": {
			Type: framework.TypeSignedDurationSecond,
			Description: `Duration in seconds of leeway when validating expiration of a token to account for clock skew. 
Defaults to 150 (2.5 minutes) if set to 0 and can be disabled if set to -1.`,
			Default: claimDefaultLeeway,
		},
		"not_before_leeway": {
			Type: framework.TypeSignedDurationSecond,
			Description: `Duration in seconds of leeway when validating not before values of a token to account for clock skew. 
Defaults to 150 (2.5 minutes) if set to 0 and can be disabled if set to -1.`,
			Default: claimDefaultLeeway,
		},
		"clock_skew_leeway": {
			Type: framework.TypeSignedDurationSecond,
			Description: `Duration in seconds of leeway when validating all claims to account for clock skew. 
Defaults to 60 (1 minute) if set to 0 and can be disabled if set to -1.`,
			Default: jwt.DefaultLeeway,
		},
		"bound_audiences": {
			},
			"expiration_leeway": {
				Type: framework.TypeSignedDurationSecond,
				Description: `Duration in seconds of leeway when validating expiration of a token to account for clock skew. 
Defaults to 150 (2.5 minutes) if set to 0 and can be disabled if set to -1.`,
				Default: claimDefaultLeeway,
			},
			"not_before_leeway": {
				Type: framework.TypeSignedDurationSecond,
				Description: `Duration in seconds of leeway when validating not before values of a token to account for clock skew. 
Defaults to 150 (2.5 minutes) if set to 0 and can be disabled if set to -1.`,
				Default: claimDefaultLeeway,
			},
			"clock_skew_leeway": {
				Type: framework.TypeSignedDurationSecond,
				Description: `Duration in seconds of leeway when validating all claims to account for clock skew. 
Defaults to 60 (1 minute) if set to 0 and can be disabled if set to -1.`,
				Default: jwt.DefaultLeeway,
			},
			"bound_audiences": {

func (b *jwtAuthBackend) validateCelProgram(program celhelper.CelProgram) (bool, error) {
	// adding a minimal jwtClaims collection here, for validating usages in CEL expression
	_, err := b.celEvalProgram(program, map[string]any{"sub": "email@example.com", "aud": "audience", "iss": "issuer"})
	if err != nil {
		return false, fmt.Errorf("failed to validate CEL program: %w", err)
	}
	return true, nil
}

func (b *jwtAuthBackend) celEvalProgram(program celhelper.CelProgram, jwtClaims map[string]any) (any, error) {
	env, err := b.celEnv(program)
	if err != nil {
		return nil, err
	// The "request" key allows CEL expressions to access and evaluate against input fields.
	// Additional variables and evaluated results will be added dynamically during processing.
	evaluationData := map[string]interface{}{
		"claims": jwtClaims,
		"now":    time.Now(),
	}

	// Evaluate all variables

	ctx := context.Background()

	testVals := func(caseSensitive bool) {
		// Clear storage
		userList, err := storage.List(ctx, "user/")
			}
		}

		loginReq := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "login/Hermes Conrad",
			Data: map[string]interface{}{
				"password": "hermes",
			},
			Storage:    storage,
			Connection: &logical.Connection{},
		}
		resp, err = b.HandleRequest(ctx, loginReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v", err, resp)
		}
		expected := []string{"grouppolicy", "userpolicy"}
		if !reflect.DeepEqual(expected, resp.Auth.Policies) {
			t.Fatalf("bad: policies: expected: %q, actual: %q", expected, resp.Auth.Policies)
		}
	}

	cleanup, cfg := ldap.PrepareTestContainer(t, "latest")
	"context"
	"errors"
	"fmt"

	"github.com/openbao/openbao/sdk/v2/framework"
	"github.com/openbao/openbao/sdk/v2/helper/cidrutil"
		return nil, errors.New("missing username")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	effectiveUsername, policies, resp, groupNames, err := b.Login(ctx, req, username, password, cfg.UsernameAsAlias)
	if err != nil || (resp != nil && resp.IsError()) {
		return resp, err
func (b *backend) pathLoginAliasLookahead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	if username == "" {
		return nil, errors.New("missing username")
	}

	return &logical.Response{
}

func (b *backend) pathLoginAliasLookahead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	if username == "" {
		return nil, errors.New("missing username")
	}

type MFACachedAuthResponse struct {
	CachedAuth            *logical.Auth
	RequestPath           string
	RequestNSID           string
	RequestNSPath         string
	}

	// First on core
	_, err = c.RegisterAuth(ctx, 0, "auth/github/login", auth, "", true)
	if err != nil {
		t.Fatal(err)
	}

	auth.TokenPolicies[0] = "default"
	_, err = c.RegisterAuth(ctx, 0, "auth/github/login", auth, "", true)
	if err == nil {
		t.Fatal("expected error")
	}

		"-H", "Content-Type: application/json",
		"--data", `{"password": "password"}`,
		"https://" + vaultAddr + ":8200/v1/auth/userpass/login/testing",
	}
	stdout, stderr, retcode, err = curlRunner.RunCmdWithOutput(ctx, curlResult.Container.ID, curlCmd)
	t.Logf("cURL Command: %v\nstdout: %v\nstderr: %v\n", curlCmd, string(stdout), string(stderr))
	var data map[string]interface{}
	err = json.Unmarshal(stdout, &data)
	require.NoError(t, err)

	auth := data["auth"].(map[string]interface{})
	remoteToken := auth["client_token"].(string)

	// Using the remote token locally should fail...
	}

	// MFA validation has passed. Let's generate the token
	resp, err := b.Core.LoginMFACreateToken(ctx, cachedResponseAuth.RequestPath, cachedResponseAuth.CachedAuth, req.Data, !req.IsInlineAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to create a token. error: %v", err)
	}

// LoginMFACreateToken creates a token after the login MFA is validated.
// It also applies the lease quotas on the original login request path.
func (c *Core) LoginMFACreateToken(ctx context.Context, reqPath string, cachedAuth *logical.Auth, loginRequestData map[string]interface{}, persistToken bool) (*logical.Response, error) {
	auth := cachedAuth
	resp := &logical.Response{
		Auth: auth,
		role = reqRole.(string)
	}

	_, resp, err = c.LoginCreateToken(ctx, ns, reqPath, mountPoint, role, resp, persistToken)
	return resp, err
}

	}

	// if user lockout feature is not disabled, check if the user is locked
	if !isUserLockoutDisabled {
		isloginUserLocked, err := c.isUserLocked(ctx, entry, req)
		if err != nil {
			return nil, nil, err
		}
		if isloginUserLocked {
			return nil, nil, logical.ErrPermissionDenied
		}
	}

	// Route the request
	// if routeErr has invalid credentials error, update the userFailedLoginMap
	if routeErr != nil && routeErr == logical.ErrInvalidCredentials {
		if !isUserLockoutDisabled {
			err := c.failedUserLoginProcess(ctx, entry, req)
			if err != nil {
				return nil, nil, err
			}
				// and return MFARequirement only
				respAuth := &MFACachedAuthResponse{
					CachedAuth:            resp.Auth,
					RequestPath:           req.Path,
					RequestNSID:           ns.ID,
					RequestNSPath:         ns.Path,
			role = c.DetermineRoleFromLoginRequest(ctx, req.MountPoint, req.Data)
		}

		_, respTokenCreate, errCreateToken := c.LoginCreateToken(ctx, ns, req.Path, source, role, resp, req.IsInlineAuth)
		if errCreateToken != nil {
			return respTokenCreate, nil, errCreateToken
		}
	// For service tokens on ent it is taken care by registerAuth RPC calls.
	// This update is done as part of registerAuth of RPC calls from standby
	// to active node. This is added there to reduce RPC calls
	if !isUserLockoutDisabled && (auth.TokenType == logical.TokenTypeBatch) {
		loginUserInfoKey := FailedLoginUser{
			aliasName:     auth.Alias.Name,
			mountAccessor: auth.Alias.MountAccessor,
		}

		// We don't need to try to delete the lockedUsers storage entry, since we're
		// processing a login request. If a login attempt is allowed, it means the user is
		// unlocked and we only add storage entry when the user gets locked.
		err = c.LocalUpdateUserFailedLoginInfo(ctx, loginUserInfoKey, nil, true)
		if err != nil {
			return nil, nil, err
		}
// LoginCreateToken creates a token as a result of a login request.
// If MFA is enforced, mfa/validate endpoint calls this functions
// after successful MFA validation to generate the token.
func (c *Core) LoginCreateToken(ctx context.Context, ns *namespace.Namespace, reqPath, mountPoint, role string, resp *logical.Response, isInlineAuth bool) (bool, *logical.Response, error) {
	auth := resp.Auth
	source := strings.TrimPrefix(mountPoint, credentialRoutePrefix)
	source = strings.ReplaceAll(source, "/", "-")
	}

	leaseGenerated := false
	te, err := c.RegisterAuth(ctx, tokenTTL, reqPath, auth, role, !isInlineAuth)
	switch {
	case err == nil:
		if auth.TokenType != logical.TokenTypeBatch {
// failedUserLoginProcess updates the userFailedLoginMap with login count and  last failed
// login time for users with failed login attempt
// If the user gets locked for current login attempt, it updates the storage entry too
func (c *Core) failedUserLoginProcess(ctx context.Context, mountEntry *MountEntry, req *logical.Request) error {
	// get the user lockout configuration for the user
	userLockoutConfiguration := c.getUserLockoutConfiguration(mountEntry)

	// determine the key for userFailedLoginInfo map
	loginUserInfoKey, err := c.getLoginUserInfoKey(ctx, mountEntry, req)
	if err != nil {
		return err
	}

	// get entry from userFailedLoginInfo map for the key
	userFailedLoginInfo := c.LocalGetUserFailedLoginInfo(ctx, loginUserInfoKey)

	// update the last failed login time with current time
	failedLoginInfo := FailedLoginInfo{
	}

	// update the userFailedLoginInfo map (and/or storage) with the updated/new entry
	err = c.LocalUpdateUserFailedLoginInfo(ctx, loginUserInfoKey, &failedLoginInfo, false)
	if err != nil {
		return err
	}
}

// isUserLocked determines if the login request user is locked
func (c *Core) isUserLocked(ctx context.Context, mountEntry *MountEntry, req *logical.Request) (locked bool, err error) {
	// get userFailedLoginInfo map key for login user
	loginUserInfoKey, err := c.getLoginUserInfoKey(ctx, mountEntry, req)
	if err != nil {
		return false, err
	}

	// get entry from userFailedLoginInfo map for the key
		// entry not found in userFailedLoginInfo map, check storage to re-verify
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return false, fmt.Errorf("could not retrieve namespace from context: %w", err)
		}

		view := NamespaceView(c.barrier, ns).SubView(coreLockedUsersPath).SubView(loginUserInfoKey.mountAccessor + "/")
		existingEntry, err := view.Get(ctx, loginUserInfoKey.aliasName)
		if err != nil {
			return false, err
		}

		var lastLoginTime int
		if existingEntry == nil {
			// no storage entry found, user is not locked
			return false, nil
		}

		err = jsonutil.DecodeJSON(existingEntry.Value, &lastLoginTime)
		if err != nil {
			return false, err
		}

		// if time passed from last login time is within lockout duration, the user is locked
		if time.Now().Unix()-int64(lastLoginTime) < int64(userLockoutConfiguration.LockoutDuration.Seconds()) {
			// user locked
			return true, nil
		}

		// else user is not locked. Entry is stale, this will be removed from storage during cleanup

		if isCountOverLockoutThreshold && isWithinLockoutDuration {
			// user locked
			return true, nil
		}
	}
	return false, nil
}

// getUserLockoutConfiguration gets the user lockout configuration for a mount entry
// store, and registers a corresponding token lease to the expiration manager.
// role is the login role used as part of the creation of the token entry. If not
// relevant, can be omitted (by being provided as "").
func (c *Core) RegisterAuth(ctx context.Context, tokenTTL time.Duration, path string, auth *logical.Auth, role string, persistToken bool) (*logical.TokenEntry, error) {
	// We first assign token policies to what was returned from the backend
	// via auth.Policies. Then, we get the full set of policies into
	// auth.Policies from the backend + entity information -- this is not
		// Successful login, remove any entry from userFailedLoginInfo map
		// if it exists. This is done for service tokens (for oss) here.
		// For ent it is taken care by registerAuth RPC calls.
		if auth.Alias != nil {
			loginUserInfoKey := FailedLoginUser{
				aliasName:     auth.Alias.Name,
				mountAccessor: auth.Alias.MountAccessor,
			}

			// We don't need to try to delete the lockedUsers storage entry, since we're
			// processing a login request. If a login attempt is allowed, it means the user is
			// unlocked and we only add storage entry when the user gets locked.
			err = c.LocalUpdateUserFailedLoginInfo(ctx, loginUserInfoKey, nil, true)
			if err != nil {
				return nil, err
			}
