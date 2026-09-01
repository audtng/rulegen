package main

	"github.com/containers/podman/v4/pkg/signal"
	"github.com/containers/storage/pkg/idtools"
	stypes "github.com/containers/storage/types"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
	return up[0], up[1]
}

// ParseRegistryCreds takes a credentials string in the form USERNAME:PASSWORD
// and returns a DockerAuthConfig
func ParseRegistryCreds(creds string) (*types.DockerAuthConfig, error) {
	"github.com/containers/podman/v4/pkg/auth"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/storage/pkg/fileutils"
	"github.com/containers/storage/pkg/ioutils"
	"github.com/containers/storage/pkg/regexp"
		stdout = options.Out
	}

	excludes := options.Excludes
	if len(excludes) == 0 {
		excludes, err = parseDockerignore(options.ContextDirectory)
		if err != nil {
			return nil, err
		}
	}

	contextDir, err = filepath.Abs(options.ContextDirectory)
	if err != nil {
		logrus.Errorf("Cannot find absolute path of %v: %v", options.ContextDirectory, err)
		if strings.HasPrefix(containerfile, contextDir+string(filepath.Separator)) {
			containerfile = strings.TrimPrefix(containerfile, contextDir+string(filepath.Separator))
			dontexcludes = append(dontexcludes, "!"+containerfile)
		} else {
			// If Containerfile does not exist, assume it is in context directory and do Not add to tarfile
			if _, err := os.Lstat(containerfile); err != nil {
					return nil, err
				}
				containerfile = c
			} else {
				// If Containerfile does exist and not in the context directory, add it to the tarfile
				tarContent = append(tarContent, containerfile)
		}
		newContainerFiles = append(newContainerFiles, filepath.ToSlash(containerfile))
	}
	if len(newContainerFiles) > 0 {
		cFileJSON, err := json.Marshal(newContainerFiles)
		if err != nil {
		params.Set("dockerfile", string(cFileJSON))
	}

	// build secrets are usually absolute host path or relative to context dir on host
	// in any case move secret to current context and ship the tar.
	if secrets := options.CommonBuildOpts.Secrets; len(secrets) > 0 {
	})
	return rc, nil
}

func parseDockerignore(root string) ([]string, error) {
	ignore, err := os.ReadFile(filepath.Join(root, ".containerignore"))
	if err != nil {
		var dockerIgnoreErr error
		ignore, dockerIgnoreErr = os.ReadFile(filepath.Join(root, ".dockerignore"))
		if dockerIgnoreErr != nil && !os.IsNotExist(dockerIgnoreErr) {
			return nil, err
		}
	}
	rawexcludes := strings.Split(string(ignore), "\n")
	excludes := make([]string, 0, len(rawexcludes))
	for _, e := range rawexcludes {
		if len(e) == 0 || e[0] == '#' {
			continue
		}
		excludes = append(excludes, e)
	}
	return excludes, nil
}
	"github.com/containers/podman/v4/pkg/auth"
	"github.com/containers/podman/v4/pkg/channel"
	"github.com/containers/podman/v4/pkg/rootless"
	"github.com/containers/storage/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/gorilla/schema"
	reporter := channel.NewWriter(make(chan []byte))
	defer reporter.Close()

	runtime := r.Context().Value(api.RuntimeKey).(*libpod.Runtime)
	buildOptions := buildahDefine.BuildOptions{
		AddCapabilities:         addCaps,
		From:                           fromImage,
		IDMappingOptions:               &idMappingOptions,
		IgnoreUnrecognizedInstructions: query.Ignore,
		Isolation:                      isolation,
		Jobs:                           &jobs,
		Labels:                         labels,
