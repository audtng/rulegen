package main

		return
	}

	forwardBody := fa.forwardBody
	// When a CONNECT method has a body with an unknown length we consider the bytes as tunnel data.
	// Therefore, we do not want to forward them to the auth server.
	if req.Method == http.MethodConnect && req.ContentLength < 0 {
		forwardBody = false
	}

	if forwardBody {
		forwardReq.ContentLength = req.ContentLength
		forwardReq.TransferEncoding = req.TransferEncoding

		bodyBytes, err := fa.readBodyBytes(req)
		if errors.Is(err, errBodyTooLarge) {
			logger.Debug().Msgf("Request body is too large, maxBodySize: %d", fa.maxBodySize)
	assert.Equal(t, 1, nextCallCount)
}

func TestForwardAuthDoesNotForwardCONNECTBody(t *testing.T) {
	testCases := []struct {
		desc          string
		contentLength int64
		enableHTTP2   bool
		expectedBody  string
	}{
		{
			desc:          "HTTP/1.1 CONNECT with a chunked body",
			contentLength: -1,
		},
		{
			desc:          "HTTP/2 CONNECT with a chunked body",
			enableHTTP2:   true,
			contentLength: -1,
		},
		{
			desc:          "HTTP/1.1 CONNECT with a fixed content-length body",
			contentLength: 3,
			expectedBody:  "foo",
		},
		{
			desc:          "HTTP/2 CONNECT with a fixed content-length body",
			enableHTTP2:   true,
			contentLength: 3,
			expectedBody:  "foo",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			var serverCallCount int
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				serverCallCount++

				forwardedData, err := io.ReadAll(req.Body)
				require.NoError(t, err)

				assert.Equal(t, test.expectedBody, string(forwardedData))
			}))
			t.Cleanup(server.Close)

			var nextCallCount int
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { nextCallCount++ })

			auth := dynamic.ForwardAuth{
				ForwardBody:           true,
				PreserveRequestMethod: true,
				Address:               server.URL,
			}
			middleware, err := NewForward(t.Context(), next, auth, "authTest")
			require.NoError(t, err)

			ts := httptest.NewUnstartedServer(middleware)
			if test.enableHTTP2 {
				ts.EnableHTTP2 = true
				ts.StartTLS()
			} else {
				ts.Start()
			}
			t.Cleanup(ts.Close)

			// Explicitly set ContentLength so we can cover both fixed-length and unknown-length bodies.
			// For unknown length (ContentLength = -1), the HTTP/1.1 client will use chunked encoding.
			req := testhelpers.MustNewRequest(http.MethodConnect, ts.URL, bytes.NewReader([]byte("foo")))
			req.ContentLength = test.contentLength

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, 1, serverCallCount)
			assert.Equal(t, 1, nextCallCount)
		})
	}
}

func TestForwardAuthForwardBodyEmptyBody(t *testing.T) {
	var serverCallCount int
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
	}

	// The FastProxy does not support CONNECT tunneling,
	// so any CONNECT request is rejected.
	if req.Method == http.MethodConnect {
		rw.WriteHeader(http.StatusNotImplemented)
		return
	}

	outReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(outReq)

	assert.Equal(t, "chunk 0\nchunk 1\nchunk 2\n", string(body))
}

func TestConnectRequest(t *testing.T) {
	var callCount int
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		callCount++
	}))
	t.Cleanup(server.Close)

	builder := NewProxyBuilder(&transportManagerMock{}, static.FastProxyConfig{})

	serverURL, err := url.JoinPath(server.URL)
	require.NoError(t, err)

	proxyHandler, err := builder.Build("", testhelpers.MustParseURL(serverURL), true, true)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodConnect, "/", http.NoBody)
	res := httptest.NewRecorder()

	proxyHandler.ServeHTTP(res, req)

	assert.Equal(t, 0, callCount)
	assert.Equal(t, http.StatusNotImplemented, res.Code)
}

type transportManagerMock struct {
	tlsConfig *tls.Config
}
