package main

			`<meta http-equiv="refresh" content="{{"asd: 123"}}">`,
			`<meta http-equiv="refresh" content="asd: 123">`,
		},
		{
			"meta content url with whitespace before equals",
			`<meta http-equiv="refresh" content="0;url ={{"javascript:alert(1)"}}">`,
			`<meta http-equiv="refresh" content="0;url =#ZgotmplZ">`,
		},
		{
			"meta content url with tab before equals",
			"<meta http-equiv=\"refresh\" content=\"0;url\t={{\"javascript:alert(1)\"}}\">",
			"<meta http-equiv=\"refresh\" content=\"0;url\t=#ZgotmplZ\">",
		},
		{
			"meta content url with space after equals",
			`<meta http-equiv="refresh" content="0;url= {{"javascript:alert(1)"}}">`,
			`<meta http-equiv="refresh" content="0;url= #ZgotmplZ">`,
		},
		{
			"meta content url with whitespace both sides of equals",
			"<meta http-equiv=\"refresh\" content=\"0;url \t= {{\"javascript:alert(1)\"}}\">",
			"<meta http-equiv=\"refresh\" content=\"0;url \t= #ZgotmplZ\">",
		},
	}

	for _, test := range tests {

// tMetaContent is the context transition function for the meta content attribute state.
func tMetaContent(c context, s []byte) (context, int) {
	for i := range len(s) {
		if i+3 <= len(s)-1 && bytes.EqualFold(s[i:i+3], []byte("url")) {
			if j := eatWhiteSpace(s, i+3); j < len(s) && s[j] == '=' {
				c.state = stateMetaContentURL
				return c, j + 1
			}
		}
	}
	return c, len(s)

// tMetaContentURL is the context transition function for the "url=" part of a meta content attribute state.
func tMetaContentURL(c context, s []byte) (context, int) {
	for i := range len(s) {
		if s[i] == ';' {
			c.state = stateMetaContent
			return c, i + 1
