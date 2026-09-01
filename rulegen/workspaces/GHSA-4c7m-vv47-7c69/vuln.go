package main

	return nil
}

func MakeEmailPrimary(email *EmailAddress) error {
	has, err := x.Get(email)
	if err != nil {
		return err
		return errors.EmailNotFound{Email: email.Email}
	}

	if !email.IsActivated {
		return errors.EmailNotVerified{Email: email.Email}
	}
		return nil, err
	}

	info := bindataFileInfo{name: "css/gogs.min.css", size: 64378, mode: os.FileMode(0644), modTime: time.Unix(1584214336, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xd9, 0x49, 0xa9, 0x99, 0x79, 0x58, 0x26, 0xec, 0xaa, 0x9, 0x5a, 0x24, 0x6, 0x69, 0x2e, 0xe0, 0x3a, 0xb1, 0x53, 0xc4, 0x42, 0x72, 0x4d, 0xe0, 0x67, 0x6d, 0xae, 0x6c, 0x8f, 0xc4, 0x27, 0x27}}
	return a, nil
}
		return nil, err
	}

	info := bindataFileInfo{name: "css/gogs.min.css.map", size: 22926, mode: os.FileMode(0644), modTime: time.Unix(1584214336, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x46, 0x89, 0xb2, 0x95, 0x91, 0xfb, 0x5c, 0xda, 0xff, 0x63, 0x54, 0xc5, 0x91, 0xbf, 0x7a, 0x5a, 0xb5, 0x3d, 0xf, 0xf, 0x84, 0x41, 0x2d, 0xc3, 0x18, 0xf5, 0x74, 0xd7, 0xa9, 0x84, 0x70, 0xce}}
	return a, nil
}

	// Make emailaddress primary.
	if c.Query("_method") == "PRIMARY" {
		if err := db.MakeEmailPrimary(&db.EmailAddress{ID: c.QueryInt64("id")}); err != nil {
			c.ServerError("MakeEmailPrimary", err)
			return
		}
