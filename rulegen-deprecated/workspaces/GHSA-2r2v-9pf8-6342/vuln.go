package main

	"encoding/base64"
	"errors"
	"fmt"
	"github.com/h44z/wg-portal/internal/app"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"path"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	evbus "github.com/vardius/message-bus"

	"github.com/h44z/wg-portal/internal/config"
	"github.com/h44z/wg-portal/internal/domain"
)

type UserManager interface {
	oauthAuthenticators map[string]domain.OauthAuthenticator
	ldapAuthenticators  map[string]domain.LdapAuthenticator

	users UserManager
}

func NewAuthenticator(cfg *config.Auth, bus evbus.MessageBus, users UserManager) (*Authenticator, error) {
	a := &Authenticator{
		cfg:   cfg,
		bus:   bus,
		users: users,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
}

func (a *Authenticator) setupExternalAuthProviders(ctx context.Context) error {
	extUrl, err := url.Parse(a.cfg.CallbackUrlPrefix)
	if err != nil {
		return fmt.Errorf("failed to parse external url: %w", err)
	}
		authProviders = append(authProviders, domain.LoginProviderInfo{
			Identifier:  providerId,
			Name:        providerName,
			ProviderUrl: fmt.Sprintf("%s/%s/init", a.cfg.CallbackUrlPrefix, providerId),
			CallbackUrl: fmt.Sprintf("%s/%s/callback", a.cfg.CallbackUrlPrefix, providerId),
		})
	}

	return user, nil
}

func (a *Authenticator) passwordAuthentication(ctx context.Context, identifier domain.UserIdentifier, password string) (*domain.User, error) {
	ctx = domain.SetUserInfo(ctx, domain.SystemAdminContextUserInfo()) // switch to admin user context to check if user exists

	var ldapUserInfo *domain.AuthenticatorUserInfo
	var ldapProvider domain.LdapAuthenticator
	}

	if !userInDatabase {
		user, err := a.processUserInfo(ctx, ldapUserInfo, domain.UserSourceLdap, ldapProvider.GetName(), ldapProvider.RegistrationEnabled())
		if err != nil {
			return nil, fmt.Errorf("unable to process user information: %w", err)
		}

// region oauth authentication

func (a *Authenticator) OauthLoginStep1(_ context.Context, providerId string) (authCodeUrl, state, nonce string, err error) {
	oauthProvider, ok := a.oauthAuthenticators[providerId]
	if !ok {
		return "", "", "", fmt.Errorf("missing oauth provider %s", providerId)
		return nil, fmt.Errorf("unable to parse user information: %w", err)
	}

	ctx = domain.SetUserInfo(ctx, domain.SystemAdminContextUserInfo()) // switch to admin user context to check if user exists
	user, err := a.processUserInfo(ctx, userInfo, domain.UserSourceOauth, oauthProvider.GetName(), oauthProvider.RegistrationEnabled())
	if err != nil {
		return nil, fmt.Errorf("unable to process user information: %w", err)
	}
	return user, nil
}

func (a *Authenticator) processUserInfo(ctx context.Context, userInfo *domain.AuthenticatorUserInfo, source domain.UserSource, provider string, withReg bool) (*domain.User, error) {
	// Search user in backend
	user, err := a.users.GetUser(ctx, userInfo.Identifier)
	switch {
	return user, nil
}

func (a *Authenticator) registerNewUser(ctx context.Context, userInfo *domain.AuthenticatorUserInfo, source domain.UserSource, provider string) (*domain.User, error) {
	// convert user info to domain.User
	user := &domain.User{
		Identifier:   userInfo.Identifier,
)

type Auth struct {
	OpenIDConnect     []OpenIDConnectProvider `yaml:"oidc"`
	OAuth             []OAuthProvider         `yaml:"oauth"`
	Ldap              []LdapProvider          `yaml:"ldap"`
	CallbackUrlPrefix string                  `yaml:"callback_url_prefix"`
}

type BaseFields struct {

type OauthFields struct {
	BaseFields `yaml:",inline"`
	IsAdmin    string `yaml:"is_admin"`
}

type LdapFields struct {
	// DisplayName is shown to the user on the login page. If it is empty, ProviderName will be displayed.
	DisplayName string `yaml:"display_name"`

	BaseUrl string `yaml:"base_url"`

	// ClientID is the application's ID.
	ClientID string `yaml:"client_id"`

	TokenURL    string `yaml:"token_url"`
	UserInfoURL string `yaml:"user_info_url"`

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string `yaml:"redirect_url"`

	// Scope specifies optional requested permissions.
	Scopes []string `yaml:"scopes"`

	userManager, err := users.NewUserManager(cfg, eventBus, database, database)
	internal.AssertNoError(err)

	authenticator, err := auth.NewAuthenticator(&cfg.Auth, eventBus, userManager)
	internal.AssertNoError(err)

	wireGuardManager, err := wireguard.NewWireGuardManager(cfg, eventBus, wireGuard, wgQuick, database)
