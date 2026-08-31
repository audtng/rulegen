package main

var appliedVersion string
var initMu sync.Mutex

func InitializeSamlServiceProvider(configToSet *v3.SamlConfig, name string) error {

	initMu.Lock()
	userPrincipal, groupPrincipals, err = s.getSamlPrincipals(config, samlData)
	if err != nil {
		log.Error(err)
		http.Redirect(w, r, redirectURL+"errorCode=422&errorMsg=Invalid saml attributes", http.StatusFound)
		return
	}
	allowedPrincipals := config.AllowedPrincipalIDs
