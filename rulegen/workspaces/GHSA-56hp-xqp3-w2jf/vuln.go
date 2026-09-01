package main

	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	InsecureSkipTLSverify bool   // --insecure-skip-verify
	Keyring               string // --keyring
	Password              string // --password
	RepoURL               string // --repo
	Username              string // --username
	Verify                bool   // --verify
		Keyring: c.Keyring,
		Getters: getter.All(settings),
		Options: []getter.Option{
			getter.WithBasicAuth(c.Username, c.Password),
			getter.WithTLSClientConfig(c.CertFile, c.KeyFile, c.CaFile),
			getter.WithInsecureSkipVerifyTLS(c.InsecureSkipTLSverify),
		},
		dl.Verify = downloader.VerifyAlways
	}
	if c.RepoURL != "" {
		chartURL, err := repo.FindChartInAuthAndTLSRepoURL(c.RepoURL, c.Username, c.Password, name, version,
			c.CertFile, c.KeyFile, c.CaFile, c.InsecureSkipTLSverify, getter.All(settings))
		if err != nil {
			return "", err
		}
		name = chartURL
	}

	if err := os.MkdirAll(settings.RepositoryCache, 0755); err != nil {
	"crypto/tls"
	"io"
	"net/http"

	"github.com/pkg/errors"

		req.Header.Set("User-Agent", g.opts.userAgent)
	}

	if g.opts.username != "" && g.opts.password != "" {
		req.SetBasicAuth(g.opts.username, g.opts.password)
	}

	client, err := g.httpClient()

		// Any failure to resolve/download a chart should fail:
		// https://github.com/helm/helm/issues/1439
		churl, username, password, err := m.findChartURL(dep.Name, dep.Version, dep.Repository, repos)
		if err != nil {
			saveError = errors.Wrapf(err, "could not find %s", churl)
			break
			Getters:          m.Getters,
			Options: []getter.Option{
				getter.WithBasicAuth(username, password),
			},
		}

// repoURL is the repository to search
//
// If it finds a URL that is "relative", it will prepend the repoURL.
func (m *Manager) findChartURL(name, version, repoURL string, repos map[string]*repo.ChartRepository) (url, username, password string, err error) {
	if strings.HasPrefix(repoURL, "oci://") {
		return fmt.Sprintf("%s/%s:%s", repoURL, name, version), "", "", nil
	}

	for _, cr := range repos {
			}
			username = cr.Config.Username
			password = cr.Config.Password
			return
		}
	}
	url, err = repo.FindChartInRepoURL(repoURL, name, version, "", "", "", m.Getters)
	if err == nil {
		return url, username, password, err
	}
	err = errors.Errorf("chart %s not found in %s: %s", name, repoURL, err)
	return url, username, password, err
}

// findEntryByName finds an entry in the chart repository whose name matches the given name.
