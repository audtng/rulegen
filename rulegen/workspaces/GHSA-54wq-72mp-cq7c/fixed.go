package main

	"io/fs"
	"log"
	"net"
	"net/mail"
	"os"
	"regexp"
	"strconv"
				break
			}

			match := extractAndValidateAddress(mailFromRE, args)
			if match == nil {
				s.writef("501 5.5.4 Syntax error in parameters or arguments (invalid FROM parameter)")
			} else {
				break
			}

			match := extractAndValidateAddress(rcptToRE, args)
			if match == nil {
				s.writef("501 5.5.4 Syntax error in parameters or arguments (invalid TO parameter)")
			} else {

	return authenticated, err
}

// Extract and validate email address from a regex match.
// This ensures that only RFC 5322 email addresses are accepted (if set).
func extractAndValidateAddress(re *regexp.Regexp, args string) []string {
	match := re.FindStringSubmatch(args)
	if match == nil || strings.Contains(match[1], " ") {
		return nil
	}

	// first argument will be the email address, validate it if not empty
	if match[1] != "" {
		// fmt.Println("Validating email address:", match[1])
		_, err := mail.ParseAddress(match[1])
		if err != nil {
			return nil
		}
	}

	return match
}
