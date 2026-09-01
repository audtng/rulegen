package main

	}
	m.Version.Major = data[0]
	m.Version.Minor = data[1]
	cookieLength := data[2]
	if len(data) < (int(cookieLength) + 3) {
		return errBufferTooSmall
	}
	m.Cookie = make([]byte, cookieLength)
