package main

			logger.ErrorContext(p.ctx, "write to auth file",
				slog.Any("err", err),
			)
		}

		return c.OpenVPNPluginFuncError
	case management.ClientAuthPending:
		pendingRespCh, err := p.managementClient.RegisterPendingPoller(currentClientID)
		if err != nil {
	}
}

// TestPluginDenyNonWebAuthClient verifies that a client not supporting webauth
// receives OPENVPN_PLUGIN_FUNC_ERROR so that OpenVPN rejects the connection.
// Previously the plugin returned OPENVPN_PLUGIN_FUNC_SUCCESS on denial, which
// caused OpenVPN to let the client in despite the management-interface deny.
func TestPluginDenyNonWebAuthClient(t *testing.T) {
	t.Parallel()

	unixSocket, err := nettest.LocalPath()
	require.NoError(t, err)

	passwordFile, err := os.CreateTemp(t.TempDir(), "openvpn-auth-oauth2-test-")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = passwordFile.Close()
	})

	_, err = passwordFile.WriteString("password")
	require.NoError(t, err)

	argv, cStrings := testutil.CreateCStringArray([]string{"openvpn-auth-oauth2", "unix://" + unixSocket, passwordFile.Name()})

	t.Cleanup(func() {
		testutil.FreeCStringArray(argv, cStrings)
	})

	openArgs := &c.OpenVPNPluginArgsOpenIn{
		Callbacks: testutil.Callbacks(),
		Argv:      argv,
	}
	openRet := &c.OpenVPNPluginArgsOpenReturn{}

	status := PluginOpenV3(PluginStructVerMin, openArgs, openRet)
	require.Equal(t, c.OpenVPNPluginFuncSuccess, status)

	t.Cleanup(func() {
		PluginCloseV1(openRet.Handle)
	})

	args := &c.OpenVPNPluginArgsFuncIn{
		Handle: openRet.Handle,
		Argv:   argv,
		Type:   c.OpenVPNPluginUp,
	}
	ret := &c.OpenVPNPluginArgsFuncReturn{}

	status = PluginFuncV3(PluginStructVerMin, args, ret)
	require.Equal(t, c.OpenVPNPluginFuncSuccess, status)

	logger := testsuite.NewTestLogger()

	// clientListener must not be closed, because it is used by the httpClientListener.
	clientListener, err := nettest.NewLocalListener("tcp")
	require.NoError(t, err)

	_, resourceServerURL, clientCredentials, err := testutils.SetupResourceServer(t, clientListener, logger.Logger, nil)
	require.NoError(t, err)

	conf := config.Defaults
	conf.OpenVPN.Addr = types.URL{URL: &url.URL{Scheme: "unix", Path: unixSocket}}
	conf.OpenVPN.Password = "password"
	conf.OAuth2.OpenVPNUsernameClaim = testsuite.SubjectClaim
	conf.HTTP.BaseURL = types.URL{URL: &url.URL{Scheme: "http", Host: clientListener.Addr().String()}}
	conf.HTTP.Secret = testsuite.Secret
	conf.OAuth2.Issuer = resourceServerURL
	conf.OAuth2.Nonce = true
	conf.OAuth2.RefreshNonce = config.OAuth2RefreshNonceEmpty
	conf.OAuth2.Client.ID = clientCredentials.ID
	conf.OAuth2.Client.Secret = clientCredentials.Secret

	tokenStorage := tokenstorage.NewInMemory(testsuite.Secret, time.Hour)
	_, openVPNClient := testutils.SetupOpenVPNOAuth2Clients(t.Context(), t, conf, logger.Logger, http.DefaultClient, tokenStorage)

	errOpenVPNClientCh := make(chan error, 1)

	go func(errCh chan<- error) {
		errCh <- openVPNClient.Connect(t.Context())
	}(errOpenVPNClientCh)

	select {
	case err := <-errOpenVPNClientCh:
		require.NoError(t, err)
	default:
	}

	time.Sleep(50 * time.Millisecond)

	clientContextPtr := PluginClientConstructorV1(openRet.Handle)
	require.NotNil(t, clientContextPtr)

	authControlFile, err := os.CreateTemp(t.TempDir(), "auth_control_file")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, authControlFile.Close())
	})

	// Client without webauth support (no IV_SSO=webauth)
	envp, envCStrings := testutil.CreateCStringArray([]string{
		"n_clients=0",
		"password=",
		"session_id=SESSIONID",
		"untrusted_port=17016",
		"untrusted_ip=192.168.65.1",
		"common_name=user@example.com",
		"username=",
		"session_state=Initial",
		"auth_control_file=" + authControlFile.Name(),
	})

	t.Cleanup(func() {
		testutil.FreeCStringArray(envp, envCStrings)
	})

	args = &c.OpenVPNPluginArgsFuncIn{
		Handle:           openRet.Handle,
		Argv:             argv,
		Envp:             envp,
		Type:             c.OpenVPNPluginAuthUserPassVerify,
		PerClientContext: clientContextPtr,
	}
	ret = &c.OpenVPNPluginArgsFuncReturn{}

	// A client without webauth support must be denied: the plugin must return ERROR.
	status = PluginFuncV3(PluginStructVerMin, args, ret)
	require.Equal(t, c.OpenVPNPluginFuncError, status, logger.String())

	// The auth control file must reflect the denial.
	data, err := os.ReadFile(authControlFile.Name())
	require.NoError(t, err)
	require.Equal(t, "0", string(data))

	PluginClientDestructorV1(args.Handle, clientContextPtr)
}

func TestPluginOpenV3_InvalidArgs(t *testing.T) {
	t.Parallel()

