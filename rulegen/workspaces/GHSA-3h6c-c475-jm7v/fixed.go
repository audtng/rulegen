package main

	ReverseProxyAuthEmail = sec.Key("REVERSE_PROXY_AUTHENTICATION_EMAIL").MustString("X-WEBAUTH-EMAIL")
	MinPasswordLength = sec.Key("MIN_PASSWORD_LENGTH").MustInt(6)
	ImportLocalPaths = sec.Key("IMPORT_LOCAL_PATHS").MustBool(false)
	DisableGitHooks = sec.Key("DISABLE_GIT_HOOKS").MustBool(true)
	OnlyAllowPushIfGiteaEnvironmentSet = sec.Key("ONLY_ALLOW_PUSH_IF_GITEA_ENVIRONMENT_SET").MustBool(true)
	PasswordHashAlgo = sec.Key("PASSWORD_HASH_ALGO").MustString("argon2")
	CSRFCookieHTTPOnly = sec.Key("CSRF_COOKIE_HTTP_ONLY").MustBool(true)
