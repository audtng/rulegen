package main

	}
	m.Version.Major = data[0]
	m.Version.Minor = data[1]
	cookieLength := int(data[2])
	if len(data) < cookieLength+3 {
		return errBufferTooSmall
	}
	m.Cookie = make([]byte, cookieLength)
