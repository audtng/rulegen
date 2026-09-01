package main

	"github.com/containers/podman/v4/pkg/signal"
	"github.com/containers/storage/pkg/idtools"
	stypes "github.com/containers/storage/types"
	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
	return up[0], up[1]
}

// Takes build context and validates `.containerignore` or `.dockerignore`
// if they are symlink outside of buildcontext. Returns list of files to be
// excluded and resolved path to the ignore files inside build context or error
func ParseDockerignore(containerfiles []string, root string) ([]string, string, error) {
	ignoreFile := ""
	path, err := securejoin.SecureJoin(root, ".containerignore")
	if err != nil {
		return nil, ignoreFile, err
	}
	// set resolved ignore file so imagebuildah
	// does not attempts to re-resolve it
	ignoreFile = path
	ignore, err := os.ReadFile(path)
	if err != nil {
		var dockerIgnoreErr error
		path, symlinkErr := securejoin.SecureJoin(root, ".dockerignore")
		if symlinkErr != nil {
			return nil, ignoreFile, symlinkErr
		}
		// set resolved ignore file so imagebuildah
		// does not attempts to re-resolve it
		ignoreFile = path
		ignore, dockerIgnoreErr = os.ReadFile(path)
		if os.IsNotExist(dockerIgnoreErr) {
			// In this case either ignorefile was not found
			// or it is a symlink to unexpected file in such
			// case manually set ignorefile to `/dev/null` so
			// internally imagebuildah does not attempts to re-resolve
			// this invalid symlink and instead reads a blank file.
			ignoreFile = "/dev/null"
		}
		// after https://github.com/containers/buildah/pull/4239 build supports
		// <Containerfile>.containerignore or <Containerfile>.dockerignore as ignore file
		// so remote must support parsing that.
		if dockerIgnoreErr != nil {
			for _, containerfile := range containerfiles {
				if _, err := os.Stat(filepath.Join(root, containerfile+".containerignore")); err == nil {
					path, symlinkErr = securejoin.SecureJoin(root, containerfile+".containerignore")
					if symlinkErr == nil {
						ignoreFile = path
						ignore, dockerIgnoreErr = os.ReadFile(path)
					}
				}
				if _, err := os.Stat(filepath.Join(root, containerfile+".dockerignore")); err == nil {
					path, symlinkErr = securejoin.SecureJoin(root, containerfile+".dockerignore")
					if symlinkErr == nil {
						ignoreFile = path
						ignore, dockerIgnoreErr = os.ReadFile(path)
					}
				}
				if dockerIgnoreErr == nil {
					break
				}
			}
		}
		if dockerIgnoreErr != nil && !os.IsNotExist(dockerIgnoreErr) {
			return nil, ignoreFile, err
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
	return excludes, ignoreFile, nil
}

// ParseRegistryCreds takes a credentials string in the form USERNAME:PASSWORD
// and returns a DockerAuthConfig
func ParseRegistryCreds(creds string) (*types.DockerAuthConfig, error) {
	"github.com/containers/podman/v4/pkg/auth"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/util"
	"github.com/containers/storage/pkg/fileutils"
	"github.com/containers/storage/pkg/ioutils"
	"github.com/containers/storage/pkg/regexp"
		stdout = options.Out
	}

	contextDir, err = filepath.Abs(options.ContextDirectory)
	if err != nil {
		logrus.Errorf("Cannot find absolute path of %v: %v", options.ContextDirectory, err)
		if strings.HasPrefix(containerfile, contextDir+string(filepath.Separator)) {
			containerfile = strings.TrimPrefix(containerfile, contextDir+string(filepath.Separator))
			dontexcludes = append(dontexcludes, "!"+containerfile)
			dontexcludes = append(dontexcludes, "!"+containerfile+".dockerignore")
			dontexcludes = append(dontexcludes, "!"+containerfile+".containerignore")
		} else {
			// If Containerfile does not exist, assume it is in context directory and do Not add to tarfile
			if _, err := os.Lstat(containerfile); err != nil {
					return nil, err
				}
				containerfile = c
				dontexcludes = append(dontexcludes, "!"+containerfile)
				dontexcludes = append(dontexcludes, "!"+containerfile+".dockerignore")
				dontexcludes = append(dontexcludes, "!"+containerfile+".containerignore")
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

	excludes := options.Excludes
	if len(excludes) == 0 {
		excludes, _, err = util.ParseDockerignore(newContainerFiles, options.ContextDirectory)
		if err != nil {
			return nil, err
		}
	}

	// build secrets are usually absolute host path or relative to context dir on host
	// in any case move secret to current context and ship the tar.
	if secrets := options.CommonBuildOpts.Secrets; len(secrets) > 0 {
	})
	return rc, nil
}
	"github.com/containers/podman/v4/pkg/auth"
	"github.com/containers/podman/v4/pkg/channel"
	"github.com/containers/podman/v4/pkg/rootless"
	"github.com/containers/podman/v4/pkg/util"
	"github.com/containers/storage/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/gorilla/schema"
	reporter := channel.NewWriter(make(chan []byte))
	defer reporter.Close()

	_, ignoreFile, err := util.ParseDockerignore(containerFiles, contextDirectory)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, fmt.Errorf("processing ignore file: %w", err))
		return
	}

	runtime := r.Context().Value(api.RuntimeKey).(*libpod.Runtime)
	buildOptions := buildahDefine.BuildOptions{
		AddCapabilities:         addCaps,
		From:                           fromImage,
		IDMappingOptions:               &idMappingOptions,
		IgnoreUnrecognizedInstructions: query.Ignore,
		IgnoreFile:                     ignoreFile,
		Isolation:                      isolation,
		Jobs:                           &jobs,
		Labels:                         labels,
