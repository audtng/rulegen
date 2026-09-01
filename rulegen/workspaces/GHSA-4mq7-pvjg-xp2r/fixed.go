package main

	return &AuthenticatorOAuth2Introspection{c: c, logger: l, provider: p, clientMap: make(map[string]*http.Client)}
}

func (a *AuthenticatorOAuth2Introspection) GetID() string { return "oauth2_introspection" }

func (a *AuthenticatorOAuth2Introspection) WaitForCache() {
	if a.tokenCache != nil {
		a.tokenCache.Wait()
	}
}

type Audience []string
	return errUnsupportedType
}

func tokenCacheKey(token, endpoint string) string {
	return fmt.Sprintf("%s|%s", token, endpoint)
}

func (a *AuthenticatorOAuth2Introspection) tokenFromCache(config *AuthenticatorOAuth2IntrospectionConfiguration, token string, ss fosite.ScopeStrategy) *AuthenticatorOAuth2IntrospectionResult {
	if !config.Cache.Enabled {
		return nil
		return nil
	}

	key := tokenCacheKey(token, config.IntrospectionURL)
	i, found := a.tokenCache.Get(key)
	if !found {
		return nil
	}
		return
	}

	key := tokenCacheKey(token, config.IntrospectionURL)
	v, err := json.Marshal(i)
	if err != nil {
		return
	}

	if a.cacheTTL != nil {
		a.tokenCache.SetWithTTL(key, v, 1, *a.cacheTTL)
	} else {
		a.tokenCache.Set(key, v, 1)
	}
}

		if err != nil {
			return errors.WithStack(err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode == http.StatusTooManyRequests {
			return errors.WithStack(helper.NewErrTooManyRequestsWithHeaders(resp))
		}
	}

	if i.TokenUse != "" && i.TokenUse != "access_token" {
		return errors.WithStack(helper.ErrForbidden.WithReason(fmt.Sprintf("Use of introspected token is not an access token but %q", i.TokenUse)))
	}

	if !i.Active {
	}

	for _, audience := range cf.Audience {
		if !slices.Contains(i.Audience, audience) {
			return errors.WithStack(helper.ErrForbidden.WithReason(fmt.Sprintf("Token audience is not intended for target audience %s", audience)))
		}
	}

import (
	"bytes"
	//nolint:gosec
	"encoding/json"
	"fmt"
	"net/http"
	return "hydrator"
}

func (a *MutatorHydrator) cacheKey(config *MutatorHydratorConfig, session []byte) string {
	return fmt.Sprintf("%s|%s", session, config.Api.URL)
}

func (a *MutatorHydrator) hydrateFromCache(config *MutatorHydratorConfig, session []byte) (*authn.AuthenticationSession, bool) {
	if !config.Cache.Enabled {
		return nil, false
	}
	return item.Copy(), true
}

func (a *MutatorHydrator) hydrateToCache(config *MutatorHydratorConfig, rawSession []byte, session *authn.AuthenticationSession) {
	if !config.Cache.Enabled {
		return
	}

	if a.hydrateCache.SetWithTTL(a.cacheKey(config, rawSession), session.Copy(), 0, config.Cache.ttl) {
		a.d.Logger().Debug("Cache reject item")
	}
}
		return err
	}

	encodedSession, err := json.Marshal(session)
	if err != nil {
		return errors.WithStack(err)
	}
	if cacheSession, ok := a.hydrateFromCache(cfg, encodedSession); ok {
		*session = *cacheSession
		return nil
	} else if _, err := url.ParseRequestURI(cfg.Api.URL); err != nil {
		return errors.New(ErrInvalidAPIURL)
	}
	req, err := http.NewRequest("POST", cfg.Api.URL, bytes.NewReader(encodedSession))
	if err != nil {
		return errors.WithStack(err)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() { _ = res.Body.Close() }()

	switch res.StatusCode {
	case http.StatusOK:
