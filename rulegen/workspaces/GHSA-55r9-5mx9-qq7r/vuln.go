package main

	response.WriteHeader(http.StatusAccepted)
}

// CheckBlob godoc
// @Summary Check image blob/layer
// @Description Check an image's blob/layer given a digest

	digest := godigest.Digest(digestStr)

	ok, blen, err := imgStore.CheckBlob(name, digest)
	if err != nil {
		details := zerr.GetDetails(err)
		if errors.Is(err, zerr.ErrBadBlobDigest) { //nolint:gocritic // errorslint conflicts with gocritic:IfElseChain
		}

		mountDigest := godigest.Digest(mountDigests[0])
		// zot does not support cross mounting directly and do a workaround creating using hard link.
		// check blob looks for actual path (name+mountDigests[0]) first then look for cache and
		// if found in cache, will do hard link and if fails we will start new upload.
		_, _, err := imgStore.CheckBlob(name, mountDigest)
		if err != nil {
			upload, err := imgStore.NewBlobUpload(name)
			if err != nil {
				details := zerr.GetDetails(err)
	return nil
}

func (d *BoltDBDriver) GetBlob(digest godigest.Digest) (string, error) {
	var blobPath strings.Builder

	return out.OriginalBlobPath, nil
}

func (d *DynamoDBDriver) PutBlob(digest godigest.Digest, path string) error {
	if path == "" {
		d.log.Error().Err(zerr.ErrEmptyValue).Str("digest", digest.String()).
