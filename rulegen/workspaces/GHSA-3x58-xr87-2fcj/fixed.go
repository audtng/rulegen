package main


		case html.StartTagToken:

			mostRecentlyStartedToken = normaliseElementName(token.Data)

			aps, ok := p.elsAndAttrs[token.Data]
			if !ok {

		case html.EndTagToken:

			if mostRecentlyStartedToken == normaliseElementName(token.Data) {
				mostRecentlyStartedToken = ""
			}


			if !skipElementContent {
				switch mostRecentlyStartedToken {
				case `script`:
					// not encouraged, but if a policy allows JavaScript we
					// should not HTML escape it as that would break the output
					buff.WriteString(token.Data)
				case `style`:
					// not encouraged, but if a policy allows CSS styles we
					// should not HTML escape it as that would break the output
					buff.WriteString(token.Data)
	}
	return aps, matched
}


// normaliseElementName takes a HTML element like <script> which is user input
// and returns a lower case version of it that is immune to UTF-8 to ASCII
// conversion tricks (like the use of upper case cyrillic i scrİpt which a
// strings.ToLower would convert to script). Instead this func will preserve
// all non-ASCII as their escaped equivalent, i.e. \u0130 which reveals the
// characters when lower cased
func normaliseElementName(str string) string {
	// that useful QuoteToASCII put quote marks at the start and end
	// so those are trimmed off
	return strings.TrimSuffix(
		strings.TrimPrefix(
			strings.ToLower(
				strconv.QuoteToASCII(str),
			),
			`"`),
		`"`,
	)
}
