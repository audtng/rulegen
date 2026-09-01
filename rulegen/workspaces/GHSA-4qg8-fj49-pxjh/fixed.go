package main

	if req.TSAPolicyOID == "" {
		oidInts = nil
	} else {
		// 128 is the max number of sub-identifiers per
		// https://datatracker.ietf.org/doc/html/rfc2578#section-3.5
		if c := strings.Count(req.TSAPolicyOID, "."); c > 128 {
			return nil, excesssivelyLongOID, fmt.Errorf("oid has %d sub identifiers, expected 128", c)
		}
		for _, v := range strings.SplitN(req.TSAPolicyOID, ".", 129) {
			i, _ := strconv.Atoi(v)
			oidInts = append(oidInts, i)
		}

func getContentType(r *http.Request) (string, error) {
	contentTypeHeader := r.Header.Get("Content-Type")
	if strings.Count(contentTypeHeader, "application/") != 1 {
		return "", errors.New("content-type header should specify application only once")
	}
	splitHeader := strings.SplitN(contentTypeHeader, "application/", 2)
	if len(splitHeader) != 2 {
		return "", errors.New("expected header value to be split into two pieces")
	}

const (
	failedToGenerateTimestampResponse        = "Error generating timestamp response"
	excesssivelyLongOID                      = "OID should be comprised of at most 128 components"
	WeakHashAlgorithmTimestampRequest        = "Weak hash algorithm in timestamp request"
	InconsistentDigestLengthTimestampRequest = "Message digest has incorrect length for specified algorithm"
)
