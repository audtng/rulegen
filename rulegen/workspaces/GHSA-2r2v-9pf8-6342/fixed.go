package main

	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/h44z/wg-portal/internal/app"
	"github.com/h44z/wg-portal/internal/config"
	"github.com/h44z/wg-portal/internal/domain"
	"github.com/sirupsen/logrus"
	evbus "github.com/vardius/message-bus"
)

type UserManager interface {
	oauthAuthenticators map[string]domain.OauthAuthenticator
	ldapAuthenticators  map[string]domain.LdapAuthenticator

	// URL prefix for the callback endpoints, this is a combination of the external URL and the API prefix
	callbackUrlPrefix string

	users UserManager
}

func NewAuthenticator(cfg *config.Auth, extUrl string, bus evbus.MessageBus, users UserManager) (
	*Authenticator,
	error,
) {
	a := &Authenticator{
		cfg:               cfg,
		bus:               bus,
		users:             users,
		callbackUrlPrefix: fmt.Sprintf("%s/api/v0", extUrl),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
}

func (a *Authenticator) setupExternalAuthProviders(ctx context.Context) error {
	extUrl, err := url.Parse(a.callbackUrlPrefix)
	if err != nil {
		return fmt.Errorf("failed to parse external url: %w", err)
	}
		authProviders = append(authProviders, domain.LoginProviderInfo{
			Identifier:  providerId,
			Name:        providerName,
			ProviderUrl: fmt.Sprintf("/auth/login/%s/init", providerId),
			CallbackUrl: fmt.Sprintf("/auth/login/%s/callback", providerId),
		})
	}

	return user, nil
}

func (a *Authenticator) passwordAuthentication(
	ctx context.Context,
	identifier domain.UserIdentifier,
	password string,
) (*domain.User, error) {
	ctx = domain.SetUserInfo(ctx,
		domain.SystemAdminContextUserInfo()) // switch to admin user context to check if user exists

	var ldapUserInfo *domain.AuthenticatorUserInfo
	var ldapProvider domain.LdapAuthenticator
	}

	if !userInDatabase {
		user, err := a.processUserInfo(ctx, ldapUserInfo, domain.UserSourceLdap, ldapProvider.GetName(),
			ldapProvider.RegistrationEnabled())
		if err != nil {
			return nil, fmt.Errorf("unable to process user information: %w", err)
		}

// region oauth authentication

func (a *Authenticator) OauthLoginStep1(_ context.Context, providerId string) (
	authCodeUrl, state, nonce string,
	err error,
) {
	oauthProvider, ok := a.oauthAuthenticators[providerId]
	if !ok {
		return "", "", "", fmt.Errorf("missing oauth provider %s", providerId)
		return nil, fmt.Errorf("unable to parse user information: %w", err)
	}

	ctx = domain.SetUserInfo(ctx,
		domain.SystemAdminContextUserInfo()) // switch to admin user context to check if user exists
	user, err := a.processUserInfo(ctx, userInfo, domain.UserSourceOauth, oauthProvider.GetName(),
		oauthProvider.RegistrationEnabled())
	if err != nil {
		return nil, fmt.Errorf("unable to process user information: %w", err)
	}
	return user, nil
}

func (a *Authenticator) processUserInfo(
	ctx context.Context,
	userInfo *domain.AuthenticatorUserInfo,
	source domain.UserSource,
	provider string,
	withReg bool,
) (*domain.User, error) {
	// Search user in backend
	user, err := a.users.GetUser(ctx, userInfo.Identifier)
	switch {
	return user, nil
}

func (a *Authenticator) registerNewUser(
	ctx context.Context,
	userInfo *domain.AuthenticatorUserInfo,
	source domain.UserSource,
	provider string,
) (*domain.User, error) {
	// convert user info to domain.User
	user := &domain.User{
		Identifier:   userInfo.Identifier,
)

type Auth struct {
	OpenIDConnect []OpenIDConnectProvider `yaml:"oidc"`
	OAuth         []OAuthProvider         `yaml:"oauth"`
	Ldap          []LdapProvider          `yaml:"ldap"`
}

type BaseFields struct {

type OauthFields struct {
	BaseFields `yaml:",inline"`
	IsAdmin    string `yaml:"is_admin"` // If the value is "true", the user is an admin.
}

type LdapFields struct {
	// DisplayName is shown to the user on the login page. If it is empty, ProviderName will be displayed.
	DisplayName string `yaml:"display_name"`

	// ClientID is the application's ID.
	ClientID string `yaml:"client_id"`

	TokenURL    string `yaml:"token_url"`
	UserInfoURL string `yaml:"user_info_url"`

	// Scope specifies optional requested permissions.
	Scopes []string `yaml:"scopes"`

	userManager, err := users.NewUserManager(cfg, eventBus, database, database)
	internal.AssertNoError(err)

	authenticator, err := auth.NewAuthenticator(&cfg.Auth, cfg.Web.ExternalUrl, eventBus, userManager)
	internal.AssertNoError(err)

	wireGuardManager, err := wireguard.NewWireGuardManager(cfg, eventBus, wireGuard, wgQuick, database)
