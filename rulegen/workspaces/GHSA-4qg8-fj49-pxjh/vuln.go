package main

	if req.TSAPolicyOID == "" {
		oidInts = nil
	} else {
		for _, v := range strings.Split(req.TSAPolicyOID, ".") {
			i, _ := strconv.Atoi(v)
			oidInts = append(oidInts, i)
		}

func getContentType(r *http.Request) (string, error) {
	contentTypeHeader := r.Header.Get("Content-Type")
	splitHeader := strings.Split(contentTypeHeader, "application/")
	if len(splitHeader) != 2 {
		return "", errors.New("expected header value to be split into two pieces")
	}

const (
	failedToGenerateTimestampResponse        = "Error generating timestamp response"
	WeakHashAlgorithmTimestampRequest        = "Weak hash algorithm in timestamp request"
	InconsistentDigestLengthTimestampRequest = "Message digest has incorrect length for specified algorithm"
)
