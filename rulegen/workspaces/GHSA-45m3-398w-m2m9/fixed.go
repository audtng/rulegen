package main

	"io"
	"net/http"
	"os"
	"sync"
	"time"

	authTypes "github.com/OliveTin/OliveTin/internal/auth/authpublic"

type OAuth2Handler struct {
	cfg                 *config.Config
	mu                  sync.RWMutex
	registeredStates    map[string]*oauth2State
	registeredProviders map[string]*oauth2.Config
}
		return
	}

	h.mu.Lock()
	h.registeredStates[state] = &oauth2State{
		providerConfig: provider,
		providerName:   providerName,
		Username:       "",
	}
	h.mu.Unlock()

	h.setOAuthCallbackCookie(w, r, "olivetin-sid-oauth", state)

		return nil, state, false
	}

	h.mu.RLock()
	registeredState, ok := h.registeredStates[state]
	h.mu.RUnlock()
	if !ok {
		log.Errorf("State not found in server: %v", state)
		http.Error(w, "State not found in server", http.StatusBadRequest)
	userInfoClient := h.createUserInfoClient(ctx, registeredState.providerConfig, tok, clientSettings)
	userinfo := getUserInfo(h.cfg, userInfoClient, providerConfig)

	h.mu.Lock()
	h.registeredStates[state].Username = userinfo.Username
	h.registeredStates[state].Usergroup = h.computeUsergroup(userinfo, providerConfig)
	h.mu.Unlock()

	http.Redirect(w, r, "/", http.StatusFound)
}
	return stringVal
}

func (h *OAuth2Handler) lookupOAuth2UserByState(state string) (*authTypes.AuthenticatedUser, bool) {
	h.mu.RLock()
	serverState, found := h.registeredStates[state]
	if !found {
		h.mu.RUnlock()
		return nil, false
	}
	user := &authTypes.AuthenticatedUser{
		Username:      serverState.Username,
		UsergroupLine: serverState.Usergroup,
		Provider:      "oauth2",
		SID:           state,
	}
	h.mu.RUnlock()
	return user, true
}

func (h *OAuth2Handler) CheckUserFromOAuth2Cookie(context *authTypes.AuthCheckingContext) *authTypes.AuthenticatedUser {
	cookie, err := context.Request.Cookie("olivetin-sid-oauth")
	if err != nil || cookie.Value == "" {
		return nil
	}

	user, found := h.lookupOAuth2UserByState(cookie.Value)
	if !found {
		log.WithFields(log.Fields{
			"sid":      cookie.Value,
			"provider": "oauth2",
		}).Warnf("Stale session")
		return nil
	}
	return user
}
