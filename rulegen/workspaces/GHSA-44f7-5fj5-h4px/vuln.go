package main


type AzureWIProviderFactory struct{} //nolint:revive // ignore linter to have unique type name
type azureWIAuthProvider struct {
	aadToken confidential.AuthResult
	tenantID string
	clientID string
}

type azureWIAuthProviderConf struct {
	Name     string `json:"name"`
	ClientID string `json:"clientID,omitempty"`
}

const (
		}
	}

	// retrieve an AAD Access token
	token, err := azureauth.GetAADAccessToken(context.Background(), tenant, clientID, AADResource)
	if err != nil {
	}

	return &azureWIAuthProvider{
		aadToken: token,
		tenantID: tenant,
		clientID: clientID,
	}, nil
}

		return provider.AuthConfig{}, re.ErrorCodeHostNameInvalid.WithComponentType(re.AuthProvider)
	}

	// need to refresh AAD token if it's expired
	if time.Now().Add(time.Minute * 5).After(d.aadToken.ExpiresOn) {
		newToken, err := azureauth.GetAADAccessToken(ctx, d.tenantID, d.clientID, AADResource)
	identityToken azcore.AccessToken
	clientID      string
	tenantID      string
}

type azureManagedIdentityAuthProviderConf struct {
	Name     string `json:"name"`
	ClientID string `json:"clientID"`
}

const (
			return nil, re.ErrorCodeEnvNotSet.WithDetail("AZURE_CLIENT_ID environment variable is empty").WithComponentType(re.AuthProvider)
		}
	}
	if err != nil {
		return nil, err
	}
	// retrieve an AAD Access token
	token, err := getManagedIdentityToken(context.Background(), client)
	if err != nil {
		identityToken: token,
		clientID:      client,
		tenantID:      tenant,
	}, nil
}

		return provider.AuthConfig{}, err
	}

	// need to refresh AAD token if it's expired
	if time.Now().Add(time.Minute * 5).After(d.identityToken.ExpiresOn) {
		newToken, err := getManagedIdentityToken(ctx, d.clientID)
