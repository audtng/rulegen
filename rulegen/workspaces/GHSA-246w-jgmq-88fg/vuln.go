package main

	require.NoError(t, err)

	argv, cStrings := testutil.CreateCStringArray([]string{"openvpn-auth-oauth2", "unix://" + unixSocket, passwordFile.Name()})
	t.Cleanup(func() {
		testutil.FreeCStringArray(argv, cStrings)
	})
