package main

var appliedVersion string
var initMu sync.Mutex

const UITranslationKeyForErrorMessage = "invalidSamlAttrs"

func InitializeSamlServiceProvider(configToSet *v3.SamlConfig, name string) error {

	initMu.Lock()
	userPrincipal, groupPrincipals, err = s.getSamlPrincipals(config, samlData)
	if err != nil {
		log.Error(err)
		// UI uses this translation key to get the error message
		http.Redirect(w, r, redirectURL+"errorCode=422&errorMsg="+UITranslationKeyForErrorMessage, http.StatusFound)
		return
	}
	allowedPrincipals := config.AllowedPrincipalIDs
