package main

/*
Copyright 2026 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"fmt"

	"github.com/spf13/pflag"

	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/request/anonymous"
	"k8s.io/apiserver/pkg/authentication/request/union"
	"k8s.io/apiserver/pkg/authentication/request/x509"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
)

type Authentication struct {
	ClientCAFile string

	// EmbeddedAuthenticator is an optional authenticator delegating to the
	// parent server's authentication chain. Set by the shard server when
	// the cache server runs embedded.
	EmbeddedAuthenticator authenticator.Request
}

func NewAuthentication() *Authentication {
	return &Authentication{}
}

func (o *Authentication) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ClientCAFile, "client-ca-file", o.ClientCAFile, "Path to a PEM-encoded certificate bundle. If set, any request presenting a client certificate signed by one of the authorities in the bundle is authenticated with an identity corresponding to the CommonName of the client certificate.")
}

func (o *Authentication) Validate() []error {
	return nil
}

func (o *Authentication) ApplyTo(authenticationInfo *genericapiserver.AuthenticationInfo, servingInfo *genericapiserver.SecureServingInfo) error {
	if o.ClientCAFile == "" && o.EmbeddedAuthenticator == nil {
		// This validation cannot happen in .Validate because these
		// options may be set by the shard embedding the cache server.
		// For the standalone cache server it doesn't matter if it's
		// validated here or in .Validate.
		return fmt.Errorf("either --client-ca-file or an embedded authenticator must be configured")
	}

	var authenticators []authenticator.Request

	if o.ClientCAFile != "" {
		caProvider, err := dynamiccertificates.NewDynamicCAContentFromFile("client-ca", o.ClientCAFile)
		if err != nil {
			return fmt.Errorf("unable to load client CA file %q: %w", o.ClientCAFile, err)
		}

		if err := authenticationInfo.ApplyClientCert(caProvider, servingInfo); err != nil {
			return fmt.Errorf("unable to apply client cert: %w", err)
		}

		authenticators = append(authenticators, x509.NewDynamic(caProvider.VerifyOptions, x509.CommonNameUserConversion))
	}
	if o.EmbeddedAuthenticator != nil {
		authenticators = append(authenticators, o.EmbeddedAuthenticator)
	}
	authenticators = append(authenticators, anonymous.NewAuthenticator(nil))

	authenticationInfo.Authenticator = union.New(authenticators...)

	return nil
}
	"github.com/abiosoft/lineprefix"
	"github.com/fatih/color"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"

	kcpclientset "github.com/kcp-dev/sdk/client/clientset/versioned/cluster"
	kcptestingserver "github.com/kcp-dev/sdk/testing/server"
	"github.com/kcp-dev/sdk/testing/third_party/library-go/crypto"

	"github.com/kcp-dev/kcp/cmd/test-server/helpers"
)

func startCacheServer(ctx context.Context, logDirPath, workingDir, hostIP string, syntheticDelay time.Duration, clientCA *crypto.CA, clientCAPath string) (<-chan error, string, error) {
	cyan := color.New(color.BgHiCyan, color.FgHiWhite).SprintFunc()
	inverse := color.New(color.BgHiWhite, color.FgHiCyan).SprintFunc()
	out := lineprefix.New(
		lineprefix.Color(color.New(color.FgHiWhite)),
	)
	cacheWorkingDir := filepath.Join(workingDir, ".kcp-cache")

	// Generate a client certificate for accessing the cache server.
	cacheClientCert := filepath.Join(cacheWorkingDir, "cache-client.crt")
	cacheClientKey := filepath.Join(cacheWorkingDir, "cache-client.key")
	if err := os.MkdirAll(cacheWorkingDir, 0755); err != nil {
		return nil, "", err
	}
	_, err := clientCA.MakeClientCertificate(cacheClientCert, cacheClientKey,
		&user.DefaultInfo{Name: "cache-client", Groups: []string{"system:masters"}}, 365)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create cache client cert: %w", err)
	}

	// Use absolute paths for the client cert/key so they resolve correctly
	// regardless of how the kubeconfig is loaded. ClientConfigLoadingRules
	// resolves relative paths relative to the kubeconfig file, while
	// LoadFromFile + NewNonInteractiveClientConfig uses them as-is from CWD.
	absCacheClientCert, err := filepath.Abs(cacheClientCert)
	if err != nil {
		return nil, "", fmt.Errorf("failed to resolve absolute path for cache client cert: %w", err)
	}
	absCacheClientKey, err := filepath.Abs(cacheClientKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to resolve absolute path for cache client key: %w", err)
	}

	cachePort := 8012
	workdir, commandLine := kcptestingserver.Command("cache-server", "cache")
	commandLine = append(
		"--embedded-etcd-peer-port=8011",
		fmt.Sprintf("--secure-port=%d", cachePort),
		fmt.Sprintf("--synthetic-delay=%s", syntheticDelay.String()),
		fmt.Sprintf("--client-ca-file=%s", clientCAPath),
	)
	fmt.Fprintf(out, "running: %v\n", strings.Join(commandLine, " "))
	cmd := exec.CommandContext(ctx, commandLine[0], commandLine[1:]...) //nolint:gosec
			continue
		}

		cacheServerCert, err := os.ReadFile(filepath.Join(cacheWorkingDir, "apiserver.crt"))
		if err != nil {
			return nil, "", err
		}
		cacheServerKubeConfig := clientcmdapi.Config{
			Clusters: map[string]*clientcmdapi.Cluster{
				"cache": {
					Server:                   fmt.Sprintf("https://localhost:%d", cachePort),
					CertificateAuthorityData: cacheServerCert,
				},
			},
			AuthInfos: map[string]*clientcmdapi.AuthInfo{
				"cache": {
					ClientCertificate: absCacheClientCert,
					ClientKey:         absCacheClientKey,
				},
			},
			Contexts: map[string]*clientcmdapi.Context{
				"cache": {
					Cluster:  "cache",
					AuthInfo: "cache",
				},
			},
			CurrentContext: "cache",
		}
		if err := clientcmd.WriteToFile(cacheServerKubeConfig, cacheKubeconfigPath); err != nil {
			return nil, "", err
		}

		loadedKubeConfig, err := clientcmd.LoadFromFile(cacheKubeconfigPath)
		if err != nil {
			return nil, "", err
		}
		cacheClientConfig := clientcmd.NewNonInteractiveClientConfig(*loadedKubeConfig, "cache", nil, nil)
		cacheClientRestConfig, err := cacheClientConfig.ClientConfig()
		if err != nil {
			return nil, "", err
/*
Copyright 2026 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"github.com/spf13/pflag"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/authorization/authorizerfactory"
	"k8s.io/apiserver/pkg/authorization/path"
	"k8s.io/apiserver/pkg/authorization/union"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

type Authorization struct {
	AlwaysAllowPaths  []string
	AlwaysAllowGroups []string
}

func NewAuthorization() *Authorization {
	return &Authorization{
		AlwaysAllowPaths:  []string{"/healthz", "/readyz", "/livez"},
		AlwaysAllowGroups: []string{user.SystemPrivilegedGroup},
	}
}

func (o *Authorization) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.AlwaysAllowPaths, "authorization-always-allow-paths", o.AlwaysAllowPaths,
		"A list of HTTP paths to skip during authorization, i.e. these are authorized without contacting the 'core' kubernetes server.")
}

func (o *Authorization) Validate() []error {
	return nil
}

func (o *Authorization) ApplyTo(config *genericapiserver.Config) error {
	var authorizers []authorizer.Authorizer

	if len(o.AlwaysAllowGroups) > 0 {
		authorizers = append(authorizers, authorizerfactory.NewPrivilegedGroups(o.AlwaysAllowGroups...))
	}

	if len(o.AlwaysAllowPaths) > 0 {
		a, err := path.NewAuthorizer(o.AlwaysAllowPaths)
		if err != nil {
			return err
		}
		authorizers = append(authorizers, a)
	}

	config.Authorization.Authorizer = union.New(authorizers...)
	return nil
}
