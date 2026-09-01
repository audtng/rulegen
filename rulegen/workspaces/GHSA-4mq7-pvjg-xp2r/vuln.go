package main

	return &AuthenticatorOAuth2Introspection{c: c, logger: l, provider: p, clientMap: make(map[string]*http.Client)}
}

func (a *AuthenticatorOAuth2Introspection) GetID() string {
	return "oauth2_introspection"
}

type Audience []string
	return errUnsupportedType
}

func (a *AuthenticatorOAuth2Introspection) tokenFromCache(config *AuthenticatorOAuth2IntrospectionConfiguration, token string, ss fosite.ScopeStrategy) *AuthenticatorOAuth2IntrospectionResult {
	if !config.Cache.Enabled {
		return nil
		return nil
	}

	i, found := a.tokenCache.Get(token)
	if !found {
		return nil
	}
		return
	}

	if v, err := json.Marshal(i); err != nil {
		return
	} else if a.cacheTTL != nil {
		a.tokenCache.SetWithTTL(token, v, 1, *a.cacheTTL)
	} else {
		a.tokenCache.Set(token, v, 1)
	}
}

		if err != nil {
			return errors.WithStack(err)
		}
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode == http.StatusTooManyRequests {
			return errors.WithStack(helper.NewErrTooManyRequestsWithHeaders(resp))
		}
	}

	if len(i.TokenUse) > 0 && i.TokenUse != "access_token" {
		return errors.WithStack(helper.ErrForbidden.WithReason(fmt.Sprintf("Use of introspected token is not an access token but \"%s\"", i.TokenUse)))
	}

	if !i.Active {
	}

	for _, audience := range cf.Audience {
		if !slices.Contains([]string(i.Audience), audience) {
			return errors.WithStack(helper.ErrForbidden.WithReason(fmt.Sprintf("Token audience is not intended for target audience %s", audience)))
		}
	}

import (
	"bytes"
	"crypto/md5" //nolint:gosec
	"encoding/json"
	"fmt"
	"net/http"
	return "hydrator"
}

func (a *MutatorHydrator) cacheKey(config *MutatorHydratorConfig, session string) string {
	return fmt.Sprintf("%s|%x", config.Api.URL, md5.Sum([]byte(session))) //nolint:gosec
}

func (a *MutatorHydrator) hydrateFromCache(config *MutatorHydratorConfig, session string) (*authn.AuthenticationSession, bool) {
	if !config.Cache.Enabled {
		return nil, false
	}
	return item.Copy(), true
}

func (a *MutatorHydrator) hydrateToCache(config *MutatorHydratorConfig, key string, session *authn.AuthenticationSession) {
	if !config.Cache.Enabled {
		return
	}

	if a.hydrateCache.SetWithTTL(a.cacheKey(config, key), session.Copy(), 0, config.Cache.ttl) {
		a.d.Logger().Debug("Cache reject item")
	}
}
		return err
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(session); err != nil {
		return errors.WithStack(err)
	}

	encodedSession := b.String()
	if cacheSession, ok := a.hydrateFromCache(cfg, encodedSession); ok {
		*session = *cacheSession
		return nil
	} else if _, err := url.ParseRequestURI(cfg.Api.URL); err != nil {
		return errors.New(ErrInvalidAPIURL)
	}
	req, err := http.NewRequest("POST", cfg.Api.URL, &b)
	if err != nil {
		return errors.WithStack(err)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close() //nolint:errcheck

	switch res.StatusCode {
	case http.StatusOK:
