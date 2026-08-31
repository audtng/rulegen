package main

)

const noSrvCertMessage = "No server certificate, server is not yet ready to receive traffic"
const serverNotReadyMsg = "Server is not yet ready to receive traffic"

var (
	cipherSuites         = tls.CipherSuites()
	return tlsConfig
}

func SetupTLSWithCertManager(caManager KubernetesCAManager, certManager certificate.Manager, clientAuth tls.ClientAuthType, clusterConfig *virtconfig.ClusterConfig) *tls.Config {
	tlsConfig := &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (certificate *tls.Certificate, err error) {
			cert := certManager.Current()
				Certificates: []tls.Certificate{*cert},
				ClientCAs:    clientCAPool,
				ClientAuth:   clientAuth,
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					if len(verifiedChains) == 0 || len(verifiedChains[0]) == 0 {
						return nil
					}

					certificate, err := x509.ParseCertificate(rawCerts[0])
					if err != nil {
						return fmt.Errorf("failed to parse peer certificate: %v", err)
					}

					CNs, err := caManager.GetCNs()
					if err != nil {
						log.Log.Reason(err).Error(serverNotReadyMsg)
						return fmt.Errorf(serverNotReadyMsg)
					}

					if len(CNs) == 0 {
						return nil
					}
					for _, CN := range CNs {
						if certificate.Subject.CommonName == CN {
							return nil
						}
					}

					return fmt.Errorf("Common name is invalid")
				},
			}

			config.BuildNameToCertificate()

type mockCAManager struct {
	caBundle []byte
	cns      []string
}

type mockCertManager struct {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCAManager) GetCNs() ([]string, error) {
	return m.cns, nil
}

var _ = Describe("TLS", func() {

	var caManager kvtls.KubernetesCAManager
	var certmanagers map[string]certificate.Manager
	var clusterConfig *virtconfig.ClusterConfig
	var kubeVirtStore cache.Store
		}),
	)

	DescribeTable("should verify self-signed client and server certificates", func(serverSecret, clientSecret, errStr string, cns []string) {
		caManager.(*mockCAManager).cns = cns
		serverTLSConfig := kvtls.SetupTLSWithCertManager(caManager, certmanagers[serverSecret], tls.RequireAndVerifyClientCert, clusterConfig)
		clientTLSConfig := kvtls.SetupTLSForVirtHandlerClients(caManager, certmanagers[clientSecret], false)
		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			components.VirtHandlerServerCertSecretName,
			components.VirtHandlerCertSecretName,
			"",
			[]string{"kubevirt.io:system:client:virt-handler"},
		),
		Entry(
			"connect with proper certificates with no CN auth",
			components.VirtHandlerServerCertSecretName,
			components.VirtHandlerCertSecretName,
			"",
			[]string{},
		),
		Entry(
			"fail if client uses an invalid certificates (CN)",
			components.VirtHandlerServerCertSecretName,
			components.VirtHandlerCertSecretName,
			"remote error: tls: bad certificate",
			[]string{"kubevirt.io:system:clientv2:virt-handler"},
		),
		Entry(
			"fail if client uses an invalid certificate",
			components.VirtHandlerServerCertSecretName,
			components.VirtHandlerServerCertSecretName,
			"remote error: tls: bad certificate",
			[]string{"kubevirt.io:system:client:virt-handler"},
		),
		Entry(
			"fail if server uses an invalid certificate",
			components.VirtHandlerCertSecretName,
			components.VirtHandlerCertSecretName,
			"x509: certificate specifies an incompatible key usage",
			[]string{"kubevirt.io:system:client:virt-handler"},
		),
	)

	})
}

func (app *virtAPIApp) setupTLS(k8sCAManager kvtls.KubernetesCAManager, kubevirtCAManager kvtls.ClientCAManager) {

	// A VerifyClientCertIfGiven request means we're not guaranteed
	// a client has been authenticated unless they provide a peer
