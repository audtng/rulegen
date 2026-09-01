package main

// MethodJSONAuth is used to identify json auth.
const MethodJSONAuth settings.AuthMethod = "json"

type jsonCred struct {
	Password  string `json:"password"`
	Username  string `json:"username"`
	}

	u, err := usr.Get(srv.Root, cred.Username)
	if err != nil || !users.CheckPwd(cred.Password, u.Password) {
		return nil, os.ErrPermission
	}

