package main

			`<meta http-equiv="refresh" content="{{"asd: 123"}}">`,
			`<meta http-equiv="refresh" content="asd: 123">`,
		},
	}

	for _, test := range tests {

// tMetaContent is the context transition function for the meta content attribute state.
func tMetaContent(c context, s []byte) (context, int) {
	for i := 0; i < len(s); i++ {
		if i+3 <= len(s)-1 && bytes.Equal(bytes.ToLower(s[i:i+4]), []byte("url=")) {
			c.state = stateMetaContentURL
			return c, i + 4
		}
	}
	return c, len(s)

// tMetaContentURL is the context transition function for the "url=" part of a meta content attribute state.
func tMetaContentURL(c context, s []byte) (context, int) {
	for i := 0; i < len(s); i++ {
		if s[i] == ';' {
			c.state = stateMetaContent
			return c, i + 1
