package main


	router.SetHTTPHandler(handlerHTTP)

	defaultTLSConf, err := m.tlsManager.Get(traefiktls.DefaultTLSStoreName, traefiktls.DefaultTLSConfigName)
	if err != nil {
		log.FromContext(ctx).Errorf("Error during the build of the default TLS configuration: %v", err)
	}

	// Keyed by domain. The source of truth for doing SNI checking, and for what TLS
	// options will actually be used for the connection.
	// As soon as there's (at least) two different tlsOptions found for the same domain,
	// we set the value to the default TLS conf.
	tlsOptionsForHost := map[string]string{}

	// Keyed by domain, then by options reference.
	// As opposed to tlsOptionsForHost, it keeps track of all the (different) TLS
	// options that occur for a given host name, so that later on we can set relevant
	// errors and logging for all the routers concerned (i.e. wrongly configured).
		}

		if len(domains) == 0 {
			// Extra Host(*) rule, for HTTPS routers with no Host rule, and for requests for
			// which the SNI does not match _any_ of the other existing routers Host. This is
			// only about choosing the TLS configuration. The actual routing will be done
			// further on by the HTTPS handler. See examples below.
			router.AddHTTPTLSConfig("*", defaultTLSConf)

			// The server name (from a Host(SNI) rule) is the only parameter (available in
			// HTTP routing rules) on which we can map a TLS config, because it is the only one
			// accessible before decryption (we obtain it during the ClientHello). Therefore,
			// when a router has no Host rule, it does not make any sense to specify some TLS
			// options. Consequently, when it comes to deciding what TLS config will be used,
			// for a request that will match an HTTPS router with no Host rule, the result will
			// depend on the _others_ existing routers (their Host rule, to be precise), and
			// the TLS options associated with them, even though they don't match the incoming
			// request. Consider the following examples:

			//	# conf1
			//	httpRouter1:
			//	httpRouter2:
			//		rule: Host("foo.com") && PathPrefix("/bar")
			//		tlsoptions: myTLSOptions
			//	# When a request for "/foo" comes, even though it won't be routed by
			//	httpRouter2, if its SNI is set to foo.com, myTLSOptions will be used for the TLS
			//	connection. Otherwise, it will fallback to the default TLS config.
			logger.Warnf("No domain found in rule %v, the TLS options applied for this router will depend on the SNI of each request", routerHTTPConfig.Rule)
		}

		tlsConf, err := m.tlsManager.Get(traefiktls.DefaultTLSStoreName, tlsOptionsName)
		if err != nil {
			routerHTTPConfig.AddError(err, true)
			logger.Error(err)
			continue
		}

		for _, domain := range domains {

	sniCheck := snicheck.New(tlsOptionsForHost, handlerHTTPS)

	router.SetHTTPSHandler(sniCheck, defaultTLSConf)

	logger := log.FromContext(ctx)
				break
			}

			logger.Debugf("Adding route for %s with TLS options %s", hostSNI, optionsName)

			router.AddHTTPTLSConfig(hostSNI, config)
		} else {
			routers := make([]string, 0, len(tlsConfigs))
			for _, v := range tlsConfigs {
				configsHTTP[v.routerName].AddError(fmt.Errorf("found different TLS options for routers on the same host %v, so using the default TLS options instead", hostSNI), false)
				routers = append(routers, v.routerName)
			}

			logger.Warnf("Found different TLS options for routers on the same host %v, so using the default TLS options instead for these routers: %#v", hostSNI, routers)

			router.AddHTTPTLSConfig(hostSNI, defaultTLSConf)
		}
	}

	for routerName, routerConfig := range configs {
		ctxRouter := log.With(provider.AddInContext(ctx, routerName), log.Str(log.RouterName, routerName))
		logger := log.FromContext(ctxRouter)
			continue
		}

		handler, err := m.buildTCPHandler(ctxRouter, routerConfig)
		if err != nil {
			routerConfig.AddError(err, true)
			logger.Error(err)
			continue
		}

		domains, err := tcpmuxer.ParseHostSNI(routerConfig.Rule)
		if err != nil {
			routerErr := fmt.Errorf("invalid rule: %q , %w", routerConfig.Rule, err)
			logger.Error(routerErr)
		}

		if routerConfig.TLS == nil {
			logger.Debugf("Adding route for %q", routerConfig.Rule)
			if err := router.AddRoute(routerConfig.Rule, routerConfig.Priority, handler); err != nil {

		if routerConfig.TLS.Passthrough {
			logger.Debugf("Adding Passthrough route for %q", routerConfig.Rule)
			if err := router.AddRouteTLS(routerConfig.Rule, routerConfig.Priority, handler, nil); err != nil {
				routerConfig.AddError(err, true)
				logger.Error(err)
			}
		tlsConf, err := m.tlsManager.Get(traefiktls.DefaultTLSStoreName, tlsOptionsName)
		if err != nil {
			routerConfig.AddError(err, true)
			logger.Error(err)
			continue
		}

		//		rule: HostSNI(foo.com) && ClientIP(IP2)
		//		tlsOption: tlsTwo
		// i.e. same HostSNI but different tlsOptions
		// This is only applicable if the muxer can decide about the routing _before_
		// telling the client about the tlsConf (i.e. before the TLS HandShake). This seems
		// to be the case so far with the existing matchers (HostSNI, and ClientIP), so
		// it's all good. Otherwise, we would have to do as for HTTPS, i.e. disallow
		// different TLS configs for the same HostSNIs.

		logger.Debugf("Adding TLS route for %q", routerConfig.Rule)
		if err := router.AddRouteTLS(routerConfig.Rule, routerConfig.Priority, handler, tlsConf); err != nil {
			routerConfig.AddError(err, true)
			logger.Error(err)
		}
	}

	return router, nil
}

func (m *Manager) buildTCPHandler(ctx context.Context, router *runtime.TCPRouterInfo) (tcp.Handler, error) {
	muxerHTTPS tcpmuxer.Muxer

	// Forwarder handlers.
	// Handles all HTTP requests.
	httpForwarder tcp.Handler
	// Handles (indirectly through muxerHTTPS, or directly) all HTTPS requests.
	httpsForwarder tcp.Handler

	// Neither is used directly, but they are held here, and recreated on config
	// reload, so that they can be passed to the Switcher at the end of the config
	// reload phase.
	httpHandler  http.Handler
	httpsHandler http.Handler

	// TLS configs.
	httpsTLSConfig    *tls.Config            // default TLS config
	hostHTTPTLSConfig map[string]*tls.Config // TLS configs keyed by SNI
}


// ServeTCP forwards the connection to the right TCP/HTTP handler.
func (r *Router) ServeTCP(conn tcp.WriteCloser) {
	// Handling Non-TLS TCP connection early if there is neither HTTP(S) nor TLS
	// routers on the entryPoint, and if there is at least one non-TLS TCP router.
	// In the case of a non-TLS TCP client (that does not "send" first), we would
	// block forever on clientHelloInfo, which is why we want to detect and
	// handle that case first and foremost.
	if r.muxerTCP.HasRoutes() && !r.muxerTCPTLS.HasRoutes() && !r.muxerHTTPS.HasRoutes() {
		connData, err := tcpmuxer.NewConnData("", conn, nil)
		if err != nil {
	// (wrapped inside the returned handler) requested for the given HostSNI.
	handlerHTTPS, catchAllHTTPS := r.muxerHTTPS.Match(connData)
	if handlerHTTPS != nil && !catchAllHTTPS {
		// In order not to depart from the behavior in 2.6, we only allow an HTTPS router
		// to take precedence over a TCP-TLS router if it is _not_ an HostSNI(*) router (so
		// basically any router that has a specific HostSNI based rule).
		handlerHTTPS.ServeTCP(r.GetConn(conn, hello.peeked))
		return
	}
		return
	}

	// needed to handle 404s for HTTPS, as well as all non-Host (e.g. PathPrefix) matches.
	if r.httpsForwarder != nil {
		r.httpsForwarder.ServeTCP(r.GetConn(conn, hello.peeked))
		return
	return r.muxerTCP.AddRoute(rule, priority, target)
}

// AddRouteTLS defines a handler for a given rule and sets the matching tlsConfig.
func (r *Router) AddRouteTLS(rule string, priority int, target tcp.Handler, config *tls.Config) error {
	// TLS PassThrough
	if config == nil {
		return r.muxerTCPTLS.AddRoute(rule, priority, target)
	}

	return r.muxerTCPTLS.AddRoute(rule, priority, &tcp.TLSHandler{
		Next:   target,
		Config: config,
	})
}

// AddHTTPTLSConfig defines a handler for a given sniHost and sets the matching tlsConfig.
func (r *Router) AddHTTPTLSConfig(sniHost string, config *tls.Config) {
	if r.hostHTTPTLSConfig == nil {
	r.httpForwarder = handler
}

// SetHTTPSForwarder sets the tcp handler that will forward the TLS connections to an http handler.
func (r *Router) SetHTTPSForwarder(handler tcp.Handler) {
	for sniHost, tlsConf := range r.hostHTTPTLSConfig {
		// muxerHTTPS only contains single HostSNI rules (and no other kind of rules),
		// so there's no need for specifying a priority for them.
		err := r.muxerHTTPS.AddRoute("HostSNI(`"+sniHost+"`)", 0, &tcp.TLSHandler{
			Next:   handler,
			Config: tlsConf,
		})
		if err != nil {
			log.WithoutContext().Errorf("Error while adding route for host: %v", err)
		}
	}

	r.httpsForwarder = &tcp.TLSHandler{
		Next:   handler,
		Config: r.httpsTLSConfig,

// Conn is a connection proxy that handles Peeked bytes.
type Conn struct {
	// Peeked are the bytes that have been read from Conn for the
	// purposes of route matching, but have not yet been consumed
	// by Read calls. It set to nil by Read when fully consumed.
	Peeked []byte

	// Conn is the underlying connection.
	// It can be type asserted against *net.TCPConn or other types
	// as needed. It should not be read from directly unless
	// Peeked is nil.
	tcp.WriteCloser
}

		return nil, err
	}

	// No valid TLS record has a type of 0x80, however SSLv2 handshakes
	// start with a uint16 length where the MSB is set and the first record
	// is always < 256 bytes long. Therefore typ == 0x80 strongly suggests
	// an SSLv2 client.
	const recordTypeSSLv2 = 0x80
	const recordTypeHandshake = 0x16
	if hdr[0] != recordTypeHandshake {
		if hdr[0] == recordTypeSSLv2 {
			// we consider SSLv2 as TLS and it will be refused by real TLS handshake.
			return &clientHello{
				isTLS:  true,
				peeked: getPeeked(br),
	m.lock.RLock()
	defer m.lock.RUnlock()

	var tlsConfig *tls.Config
	var err error

	sniStrict := false
	config, ok := m.configs[configName]
	if ok {
		sniStrict = config.SniStrict
		tlsConfig, err = buildTLSConfig(config)
	} else {
		err = fmt.Errorf("unknown TLS options: %s", configName)
	}
	if err != nil {
		tlsConfig = &tls.Config{}
	}

	store := m.getStore(storeName)
		err = fmt.Errorf("TLS store %s not found", storeName)
	}
	acmeTLSStore := m.getStore(tlsalpn01.ACMETLS1Protocol)
	if acmeTLSStore == nil {
		err = fmt.Errorf("ACME TLS store %s not found", tlsalpn01.ACMETLS1Protocol)
	}

			certificate := acmeTLSStore.GetBestCertificate(clientHello)
			if certificate == nil {
				log.WithoutContext().Debugf("TLS: no certificate for TLSALPN challenge: %s", domainToCheck)
				// We want the user to eventually get the (alertUnrecognizedName) "unrecognized
				// name" error.
				// Unfortunately, if we returned an error here, since we can't use
				// the unexported error (errNoCertificates) that our caller (config.getCertificate
				// in crypto/tls) uses as a sentinel, it would report an (alertInternalError)
				// "internal error" instead of an alertUnrecognizedName.
				// Which is why we return no error, and we let the caller detect that there's
				// actually no certificate, and fall back into the flow that will report
				// the desired error.
				// https://cs.opensource.google/go/go/+/dev.boringcrypto.go1.17:src/crypto/tls/common.go;l=1058
				return nil, nil
			}
