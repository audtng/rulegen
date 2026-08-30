package main

		if err != nil {
			return ctx, nil, nil, nopTeardown, err
		}
		return outCtx, appClient.(*grpc.ClientConn), nil, nopTeardown, nil
	}

	outCtx = p.telemetryFn(outCtx)
	outCtx = metadata.AppendToOutgoingContext(outCtx, invokev1.CallerIDHeader, p.appID, invokev1.CalleeIDHeader, target.id)

	appMetadataToken := security.GetAppToken()
	if appMetadataToken != "" {
		outCtx = metadata.AppendToOutgoingContext(outCtx, securityConsts.APITokenHeader, appMetadataToken)
	}

	pt := &grpcProxy.ProxyTarget{
		ID:        target.id,
		Namespace: target.namespace,
			}, nil
		})

		ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{diagnostics.GRPCProxyAppIDKey: []string{"a"}})
		proxy := p.(*proxy)
		_, conn, _, teardown, err := proxy.intercept(ctx, "/test")
		defer teardown(true)

		require.NoError(t, err)
		assert.NotNil(t, conn)
		assert.Equal(t, "a", conn.Target())
	})

	t.Run("proxy to a remote app", func(t *testing.T) {
		assert.Equal(t, "b", md["a"][0])
		assert.Equal(t, "a", md[invokev1.CallerIDHeader][0])
		assert.Equal(t, "b", md[invokev1.CalleeIDHeader][0])
		assert.Equal(t, "token1", md[securityConsts.APITokenHeader][0])
	})

	t.Run("access policies applied", func(t *testing.T) {
		"DAPR_COMPONENTS_SOCKETS_FOLDER", socket.Directory(),
	))
}
