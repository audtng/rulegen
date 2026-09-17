package main

package proxy

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"github.com/distribution/distribution/v3/internal/client/auth"
	"github.com/distribution/distribution/v3/internal/client/auth/challenge"
	"github.com/distribution/distribution/v3/internal/dcontext"
	"golang.org/x/net/publicsuffix"
)

const challengeHeader = "Docker-Distribution-Api-Version"
func getAuthURLs(remoteURL string) ([]string, error) {
	authURLs := []string{}

	remote, err := url.Parse(remoteURL)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(remoteURL + "/v2/")
	if err != nil {
		return nil, err
	defer resp.Body.Close()

	for _, c := range challenge.ResponseChallenges(resp) {
		if strings.EqualFold(c.Scheme, "bearer") && realmAllowed(remote, c.Parameters["realm"]) {
			authURLs = append(authURLs, c.Parameters["realm"])
		}
	}
	return authURLs, nil
}

func realmAllowed(remote *url.URL, realm string) bool {
	realmURL, err := url.Parse(realm)
	if err != nil {
		return false
	}
	if realmURL.Host == "" || remote == nil || remote.Host == "" {
		return false
	}

	if strings.EqualFold(remote.Host, realmURL.Host) {
		return true
	}

	remoteHost := strings.ToLower(remote.Hostname())
	realmHost := strings.ToLower(realmURL.Hostname())
	if remoteHost == "" || realmHost == "" {
		return false
	}

	if isLiteralOrLocal(remoteHost) || isLiteralOrLocal(realmHost) {
		return false
	}

	return strings.EqualFold(registrableDomain(remoteHost), registrableDomain(realmHost))
}

func isLiteralOrLocal(host string) bool {
	if host == "localhost" {
		return true
	}

	return net.ParseIP(host) != nil
}

func registrableDomain(host string) string {
	domain, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return ""
	}

	return domain
}

func ping(manager challenge.Manager, endpoint, versionHeader string) error {
	resp, err := http.Get(endpoint)
	if err != nil {
package challenge

import (
	"net/http"
	"net/url"
)

// FilteringManager decorates another Manager and drops challenges that do not
// satisfy the configured predicate.
type FilteringManager struct {
	base Manager
	keep func(Challenge) bool
}

// NewFilteringManager returns a Manager that delegates storage to base and
// filters challenges on reads. If keep is nil, the base manager is returned.
func NewFilteringManager(base Manager, keep func(Challenge) bool) Manager {
	if keep == nil {
		return base
	}

	return FilteringManager{
		base: base,
		keep: keep,
	}
}

func (m FilteringManager) GetChallenges(endpoint url.URL) ([]Challenge, error) {
	challenges, err := m.base.GetChallenges(endpoint)
	if err != nil {
		return nil, err
	}

	filtered := make([]Challenge, 0, len(challenges))
	for _, c := range challenges {
		if m.keep(c) {
			filtered = append(filtered, c)
		}
	}

	return filtered, nil
}

func (m FilteringManager) AddResponse(resp *http.Response) error {
	return m.base.AddResponse(resp)
}
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

}

func (r *remoteAuthChallenger) challengeManager() challenge.Manager {
	return challenge.NewFilteringManager(r.cm, func(c challenge.Challenge) bool {
		return !strings.EqualFold(c.Scheme, "bearer") || realmAllowed(&r.remoteURL, c.Parameters["realm"])
	})
}

// tryEstablishChallenges will attempt to get a challenge type for the upstream if none currently exist
