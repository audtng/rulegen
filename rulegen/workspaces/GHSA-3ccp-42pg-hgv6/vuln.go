package main

		return
	}

	if fa.forwardBody {
		bodyBytes, err := fa.readBodyBytes(req)
		if errors.Is(err, errBodyTooLarge) {
			logger.Debug().Msgf("Request body is too large, maxBodySize: %d", fa.maxBodySize)
	assert.Equal(t, 1, nextCallCount)
}

func TestForwardAuthForwardBodyEmptyBody(t *testing.T) {
	var serverCallCount int
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
	}

	outReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(outReq)

	assert.Equal(t, "chunk 0\nchunk 1\nchunk 2\n", string(body))
}

type transportManagerMock struct {
	tlsConfig *tls.Config
}
