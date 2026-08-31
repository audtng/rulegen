package main


func (f *Fs) getTusLocationOrRetry(ctx context.Context, resp *http.Response, err error) (bool, string, error) {

	// resp is nil if the HTTP transaction failed before a response
	// was received, eg on connection refused or reset
	if resp != nil {
		switch resp.StatusCode {
		case 201:
			location := resp.Header.Get("Location")
			return false, location, nil
		case 412:
			return false, "", ErrVersionMismatch
		case 413:
			return false, "", ErrLargeUpload
		}
	}

	retry, err := f.shouldRetry(ctx, resp, err)
