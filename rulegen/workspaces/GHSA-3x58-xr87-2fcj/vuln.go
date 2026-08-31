package main


		case html.StartTagToken:

			mostRecentlyStartedToken = strings.ToLower(token.Data)

			aps, ok := p.elsAndAttrs[token.Data]
			if !ok {

		case html.EndTagToken:

			if mostRecentlyStartedToken == strings.ToLower(token.Data) {
				mostRecentlyStartedToken = ""
			}


			if !skipElementContent {
				switch mostRecentlyStartedToken {
				case "script":
					// not encouraged, but if a policy allows JavaScript we
					// should not HTML escape it as that would break the output
					buff.WriteString(token.Data)
				case "style":
					// not encouraged, but if a policy allows CSS styles we
					// should not HTML escape it as that would break the output
					buff.WriteString(token.Data)
	}
	return aps, matched
}
