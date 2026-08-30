package main

		return nil, false
	}

	trustedCIDRs, _ := c.engine.prepareTrustedCIDRs()
	c.engine.trustedCIDRs = trustedCIDRs
	if c.engine.trustedCIDRs != nil {
		for _, cidr := range c.engine.trustedCIDRs {
			if cidr.Contains(remoteIP) {
	assert.True(t, c.IsAborted())
}

func TestContextClientIP(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", nil)

	resetContextForClientIPTests(c)

	// Legacy tests (validating that the defaults don't break the

	// No trusted proxies
	c.engine.TrustedProxies = []string{}
	c.engine.RemoteIPHeaders = []string{"X-Forwarded-For"}
	assert.Equal(t, "40.40.40.40", c.ClientIP())

	// Last proxy is trusted, but the RemoteAddr is not
	c.engine.TrustedProxies = []string{"30.30.30.30"}
	assert.Equal(t, "40.40.40.40", c.ClientIP())

	// Only trust RemoteAddr
	c.engine.TrustedProxies = []string{"40.40.40.40"}
	assert.Equal(t, "20.20.20.20", c.ClientIP())

	// All steps are trusted
	c.engine.TrustedProxies = []string{"40.40.40.40", "30.30.30.30", "20.20.20.20"}
	assert.Equal(t, "20.20.20.20", c.ClientIP())

	// Use CIDR
	c.engine.TrustedProxies = []string{"40.40.25.25/16", "30.30.30.30"}
	assert.Equal(t, "20.20.20.20", c.ClientIP())

	// Use hostname that resolves to all the proxies
	c.engine.TrustedProxies = []string{"foo"}
	assert.Equal(t, "40.40.40.40", c.ClientIP())

	// Use hostname that returns an error
	c.engine.TrustedProxies = []string{"bar"}
	assert.Equal(t, "40.40.40.40", c.ClientIP())

	// X-Forwarded-For has a non-IP element
	c.engine.TrustedProxies = []string{"40.40.40.40"}
	c.Request.Header.Set("X-Forwarded-For", " blah ")
	assert.Equal(t, "40.40.40.40", c.ClientIP())

	// happen, but we should test it to make sure we handle it
	// gracefully.
	c.engine.TrustedProxies = []string{"baz"}
	c.Request.Header.Set("X-Forwarded-For", " 30.30.30.30 ")
	assert.Equal(t, "40.40.40.40", c.ClientIP())

	c.engine.TrustedProxies = []string{"40.40.40.40"}
	c.Request.Header.Del("X-Forwarded-For")
	c.engine.RemoteIPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}
	assert.Equal(t, "10.10.10.10", c.ClientIP())
	buffer := new(bytes.Buffer)

	router := New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Output: buffer,
		Formatter: func(param LogFormatterParams) string {
