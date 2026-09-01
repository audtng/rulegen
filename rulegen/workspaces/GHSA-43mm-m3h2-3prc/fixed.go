package main

// MethodJSONAuth is used to identify json auth.
const MethodJSONAuth settings.AuthMethod = "json"

// dummyHash is used to prevent user enumeration timing attacks.
// It MUST be a valid bcrypt hash.
const dummyHash = "$2a$10$O4mEMeOL/nit6zqe.WQXauLRbRlzb3IgLHsa26Pf0N/GiU9b.wK1m"

type jsonCred struct {
	Password  string `json:"password"`
	Username  string `json:"username"`
	}

	u, err := usr.Get(srv.Root, cred.Username)

	hash := dummyHash
	if err == nil {
		hash = u.Password
	}

	if !users.CheckPwd(cred.Password, hash) {
		return nil, os.ErrPermission
	}

	if err != nil {
		return nil, os.ErrPermission
	}

