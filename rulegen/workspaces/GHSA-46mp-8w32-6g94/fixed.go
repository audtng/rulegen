package main

}

func matchSignatures(signatures []oci.Signature, subject, subjectRegExp, issuer, issuerRegExp string, extensions map[string]string) error {
	if subject == "" && issuer == "" && subjectRegExp == "" && issuerRegExp == "" && len(extensions) == 0 {
		return nil
	}

