package main

)

const noSrvCertMessage = "No server certificate, server is not yet ready to receive traffic"

var (
	cipherSuites         = tls.CipherSuites()
	return tlsConfig
}

func SetupTLSWithCertManager(caManager ClientCAManager, certManager certificate.Manager, clientAuth tls.ClientAuthType, clusterConfig *virtconfig.ClusterConfig) *tls.Config {
	tlsConfig := &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (certificate *tls.Certificate, err error) {
			cert := certManager.Current()
				Certificates: []tls.Certificate{*cert},
				ClientCAs:    clientCAPool,
				ClientAuth:   clientAuth,
			}

			config.BuildNameToCertificate()

type mockCAManager struct {
	caBundle []byte
}

type mockCertManager struct {
	return nil, fmt.Errorf("not implemented")
}

var _ = Describe("TLS", func() {

	var caManager kvtls.ClientCAManager
	var certmanagers map[string]certificate.Manager
	var clusterConfig *virtconfig.ClusterConfig
	var kubeVirtStore cache.Store
		}),
	)

	DescribeTable("should verify self-signed client and server certificates", func(serverSecret, clientSecret string, errStr string) {
		serverTLSConfig := kvtls.SetupTLSWithCertManager(caManager, certmanagers[serverSecret], tls.RequireAndVerifyClientCert, clusterConfig)
		clientTLSConfig := kvtls.SetupTLSForVirtHandlerClients(caManager, certmanagers[clientSecret], false)
		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			components.VirtHandlerServerCertSecretName,
			components.VirtHandlerCertSecretName,
			"",
		),
		Entry(
			"fail if client uses an invalid certificate",
			components.VirtHandlerServerCertSecretName,
			components.VirtHandlerServerCertSecretName,
			"remote error: tls: bad certificate",
		),
		Entry(
			"fail if server uses an invalid certificate",
			components.VirtHandlerCertSecretName,
			components.VirtHandlerCertSecretName,
			"x509: certificate specifies an incompatible key usage",
		),
	)

	})
}

func (app *virtAPIApp) setupTLS(k8sCAManager kvtls.ClientCAManager, kubevirtCAManager kvtls.ClientCAManager) {

	// A VerifyClientCertIfGiven request means we're not guaranteed
	// a client has been authenticated unless they provide a peer
