package main

	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	InsecureSkipTLSverify bool   // --insecure-skip-verify
	Keyring               string // --keyring
	Password              string // --password
	PassCredentialsAll    bool   // --pass-credentials
	RepoURL               string // --repo
	Username              string // --username
	Verify                bool   // --verify
		Keyring: c.Keyring,
		Getters: getter.All(settings),
		Options: []getter.Option{
			getter.WithPassCredentialsAll(c.PassCredentialsAll),
			getter.WithTLSClientConfig(c.CertFile, c.KeyFile, c.CaFile),
			getter.WithInsecureSkipVerifyTLS(c.InsecureSkipTLSverify),
		},
		dl.Verify = downloader.VerifyAlways
	}
	if c.RepoURL != "" {
		chartURL, err := repo.FindChartInAuthAndTLSAndPassRepoURL(c.RepoURL, c.Username, c.Password, name, version,
			c.CertFile, c.KeyFile, c.CaFile, c.InsecureSkipTLSverify, c.PassCredentialsAll, getter.All(settings))
		if err != nil {
			return "", err
		}
		name = chartURL

		// Only pass the user/pass on when the user has said to or when the
		// location of the chart repo and the chart are the same domain.
		u1, err := url.Parse(c.RepoURL)
		if err != nil {
			return "", err
		}
		u2, err := url.Parse(chartURL)
		if err != nil {
			return "", err
		}

		// Host on URL (returned from url.Parse) contains the port if present.
		// This check ensures credentials are not passed between different
		// services on different ports.
		if c.PassCredentialsAll || (u1.Scheme == u2.Scheme && u1.Host == u2.Host) {
			dl.Options = append(dl.Options, getter.WithBasicAuth(c.Username, c.Password))
		} else {
			dl.Options = append(dl.Options, getter.WithBasicAuth("", ""))
		}
	} else {
		dl.Options = append(dl.Options, getter.WithBasicAuth(c.Username, c.Password))
	}

	if err := os.MkdirAll(settings.RepositoryCache, 0755); err != nil {
	"crypto/tls"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

		req.Header.Set("User-Agent", g.opts.userAgent)
	}

	// Before setting the basic auth credentials, make sure the URL associated
	// with the basic auth is the one being fetched.
	u1, err := url.Parse(g.opts.url)
	if err != nil {
		return buf, errors.Wrap(err, "Unable to parse getter URL")
	}
	u2, err := url.Parse(href)
	if err != nil {
		return buf, errors.Wrap(err, "Unable to parse URL getting from")
	}

	// Host on URL (returned from url.Parse) contains the port if present.
	// This check ensures credentials are not passed between different
	// services on different ports.
	if g.opts.passCredentialsAll || (u1.Scheme == u2.Scheme && u1.Host == u2.Host) {
		if g.opts.username != "" && g.opts.password != "" {
			req.SetBasicAuth(g.opts.username, g.opts.password)
		}
	}

	client, err := g.httpClient()

		// Any failure to resolve/download a chart should fail:
		// https://github.com/helm/helm/issues/1439
		churl, username, password, passcredentialsall, err := m.findChartURL(dep.Name, dep.Version, dep.Repository, repos)
		if err != nil {
			saveError = errors.Wrapf(err, "could not find %s", churl)
			break
			Getters:          m.Getters,
			Options: []getter.Option{
				getter.WithBasicAuth(username, password),
				getter.WithPassCredentialsAll(passcredentialsall),
			},
		}

// repoURL is the repository to search
//
// If it finds a URL that is "relative", it will prepend the repoURL.
func (m *Manager) findChartURL(name, version, repoURL string, repos map[string]*repo.ChartRepository) (url, username, password string, passcredentialsall bool, err error) {
	if strings.HasPrefix(repoURL, "oci://") {
		return fmt.Sprintf("%s/%s:%s", repoURL, name, version), "", "", false, nil
	}

	for _, cr := range repos {
			}
			username = cr.Config.Username
			password = cr.Config.Password
			passcredentialsall = cr.Config.PassCredentialsAll
			return
		}
	}
	url, err = repo.FindChartInRepoURL(repoURL, name, version, "", "", "", m.Getters)
	if err == nil {
		return url, username, password, false, err
	}
	err = errors.Errorf("chart %s not found in %s: %s", name, repoURL, err)
	return url, username, password, false, err
}

// findEntryByName finds an entry in the chart repository whose name matches the given name.
