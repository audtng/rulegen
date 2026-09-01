package main

	"io"
	"net/http"
	"os"
	"time"

	authTypes "github.com/OliveTin/OliveTin/internal/auth/authpublic"

type OAuth2Handler struct {
	cfg                 *config.Config
	registeredStates    map[string]*oauth2State
	registeredProviders map[string]*oauth2.Config
}
		return
	}

	h.registeredStates[state] = &oauth2State{
		providerConfig: provider,
		providerName:   providerName,
		Username:       "",
	}

	h.setOAuthCallbackCookie(w, r, "olivetin-sid-oauth", state)

		return nil, state, false
	}

	registeredState, ok := h.registeredStates[state]
	if !ok {
		log.Errorf("State not found in server: %v", state)
		http.Error(w, "State not found in server", http.StatusBadRequest)
	userInfoClient := h.createUserInfoClient(ctx, registeredState.providerConfig, tok, clientSettings)
	userinfo := getUserInfo(h.cfg, userInfoClient, providerConfig)

	h.registeredStates[state].Username = userinfo.Username
	h.registeredStates[state].Usergroup = h.computeUsergroup(userinfo, providerConfig)

	http.Redirect(w, r, "/", http.StatusFound)
}
	return stringVal
}

func (h *OAuth2Handler) CheckUserFromOAuth2Cookie(context *authTypes.AuthCheckingContext) *authTypes.AuthenticatedUser {
	cookie, err := context.Request.Cookie("olivetin-sid-oauth")

	user := &authTypes.AuthenticatedUser{}

	if err != nil {
		return nil
	}

	if cookie.Value == "" {
		return nil
	}

	serverState, found := h.registeredStates[cookie.Value]

	if !found {
		log.WithFields(log.Fields{
			"sid":      cookie.Value,
			"provider": "oauth2",
		}).Warnf("Stale session")

		return nil
	}

	user.Username = serverState.Username
	user.UsergroupLine = serverState.Usergroup
	user.Provider = "oauth2"
	user.SID = cookie.Value

	return user
}
