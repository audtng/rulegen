package main

	"github.com/abiosoft/lineprefix"
	"github.com/fatih/color"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"

	kcpclientset "github.com/kcp-dev/sdk/client/clientset/versioned/cluster"
	kcptestingserver "github.com/kcp-dev/sdk/testing/server"

	"github.com/kcp-dev/kcp/cmd/test-server/helpers"
)

func startCacheServer(ctx context.Context, logDirPath, workingDir, hostIP string, syntheticDelay time.Duration) (<-chan error, string, error) {
	cyan := color.New(color.BgHiCyan, color.FgHiWhite).SprintFunc()
	inverse := color.New(color.BgHiWhite, color.FgHiCyan).SprintFunc()
	out := lineprefix.New(
		lineprefix.Color(color.New(color.FgHiWhite)),
	)
	cacheWorkingDir := filepath.Join(workingDir, ".kcp-cache")
	cachePort := 8012
	workdir, commandLine := kcptestingserver.Command("cache-server", "cache")
	commandLine = append(
		"--embedded-etcd-peer-port=8011",
		fmt.Sprintf("--secure-port=%d", cachePort),
		fmt.Sprintf("--synthetic-delay=%s", syntheticDelay.String()),
	)
	fmt.Fprintf(out, "running: %v\n", strings.Join(commandLine, " "))
	cmd := exec.CommandContext(ctx, commandLine[0], commandLine[1:]...) //nolint:gosec
			continue
		}

		if _, err := os.Stat(cacheKubeconfigPath); os.IsNotExist(err) {
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
				Contexts: map[string]*clientcmdapi.Context{
					"cache": {
						Cluster: "cache",
					},
				},
				CurrentContext: "cache",
			}
			if err := clientcmd.WriteToFile(cacheServerKubeConfig, cacheKubeconfigPath); err != nil {
				return nil, "", err
			}
		}

		cacheServerKubeConfig, err := clientcmd.LoadFromFile(cacheKubeconfigPath)
		if err != nil {
			return nil, "", err
		}
		cacheClientConfig := clientcmd.NewNonInteractiveClientConfig(*cacheServerKubeConfig, "cache", nil, nil)
		cacheClientRestConfig, err := cacheClientConfig.ClientConfig()
		if err != nil {
			return nil, "", err
