package main

	DefaultIdleTimeout = 180 * time.Second

	// DefaultReadTimeout defines the default maximum duration for reading the entire request, including the body.
	DefaultReadTimeout = 5 * time.Second

	// DefaultAcmeCAServer is the default ACME API endpoint.
	DefaultAcmeCAServer = "https://acme-v02.api.letsencrypt.org/directory"
