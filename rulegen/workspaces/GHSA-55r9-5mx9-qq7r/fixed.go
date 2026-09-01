package main

	response.WriteHeader(http.StatusAccepted)
}

// canMount checks if a user has read permission on cached blobs with this specific digest.
// returns true if the user have permission to copy blob from cache.
func canMount(userAc *reqCtx.UserAccessControl, imgStore storageTypes.ImageStore, digest godigest.Digest,
) (bool, error) {
	canMount := true

	// authz enabled
	if userAc != nil {
		canMount = false

		repos, err := imgStore.GetAllDedupeReposCandidates(digest)
		if err != nil {
			// first write
			return false, err
		}

		if len(repos) == 0 {
			canMount = false
		}

		// check if user can read any repo which contain this blob
		for _, repo := range repos {
			if userAc.Can(constants.ReadPermission, repo) {
				canMount = true
			}
		}
	}

	return canMount, nil
}

// CheckBlob godoc
// @Summary Check image blob/layer
// @Description Check an image's blob/layer given a digest

	digest := godigest.Digest(digestStr)

	userAc, err := reqCtx.UserAcFromContext(request.Context())
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)

		return
	}

	userCanMount, err := canMount(userAc, imgStore, digest)
	if err != nil {
		rh.c.Log.Error().Err(err).Msg("unexpected error")
	}

	var blen int64

	if userCanMount {
		ok, blen, err = imgStore.CheckBlob(name, digest)
	} else {
		var lockLatency time.Time

		imgStore.RLock(&lockLatency)
		defer imgStore.RUnlock(&lockLatency)

		ok, blen, _, err = imgStore.StatBlob(name, digest)
	}

	if err != nil {
		details := zerr.GetDetails(err)
		if errors.Is(err, zerr.ErrBadBlobDigest) { //nolint:gocritic // errorslint conflicts with gocritic:IfElseChain
		}

		mountDigest := godigest.Digest(mountDigests[0])

		userAc, err := reqCtx.UserAcFromContext(request.Context())
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)

			return
		}

		userCanMount, err := canMount(userAc, imgStore, mountDigest)
		if err != nil {
			rh.c.Log.Error().Err(err).Msg("unexpected error")
		}

		// zot does not support cross mounting directly and do a workaround creating using hard link.
		// check blob looks for actual path (name+mountDigests[0]) first then look for cache and
		// if found in cache, will do hard link and if fails we will start new upload.
		if userCanMount {
			_, _, err = imgStore.CheckBlob(name, mountDigest)
		}

		if err != nil || !userCanMount {
			upload, err := imgStore.NewBlobUpload(name)
			if err != nil {
				details := zerr.GetDetails(err)
	return nil
}

func (d *BoltDBDriver) GetAllBlobs(digest godigest.Digest) ([]string, error) {
	var blobPath strings.Builder

	blobPaths := []string{}

	if err := d.db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte(constants.BlobsCache))
		if root == nil {
			// this is a serious failure
			err := zerr.ErrCacheRootBucket
			d.log.Error().Err(err).Msg("failed to access root bucket")

			return err
		}

		bucket := root.Bucket([]byte(digest.String()))
		if bucket != nil {
			origin := bucket.Bucket([]byte(constants.OriginalBucket))
			blobPath.Write(d.getOne(origin))
			originBlob := blobPath.String()

			blobPaths = append(blobPaths, originBlob)

			deduped := bucket.Bucket([]byte(constants.DuplicatesBucket))
			if deduped != nil {
				cursor := deduped.Cursor()

				for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
					var blobPath strings.Builder

					blobPath.Write(k)

					duplicateBlob := blobPath.String()

					if duplicateBlob != originBlob {
						blobPaths = append(blobPaths, duplicateBlob)
					}
				}

				return nil
			}
		}

		return zerr.ErrCacheMiss
	}); err != nil {
		return nil, err
	}

	return blobPaths, nil
}

func (d *BoltDBDriver) GetBlob(digest godigest.Digest) (string, error) {
	var blobPath strings.Builder

	return out.OriginalBlobPath, nil
}

func (d *DynamoDBDriver) GetAllBlobs(digest godigest.Digest) ([]string, error) {
	blobPaths := []string{}

	resp, err := d.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"Digest": &types.AttributeValueMemberS{Value: digest.String()},
		},
	})
	if err != nil {
		d.log.Error().Err(err).Str("tableName", d.tableName).Msg("failed to get blob")

		return nil, err
	}

	out := Blob{}

	if resp.Item == nil {
		d.log.Debug().Err(zerr.ErrCacheMiss).Str("digest", string(digest)).Msg("failed to find blob in cache")

		return nil, zerr.ErrCacheMiss
	}

	_ = attributevalue.UnmarshalMap(resp.Item, &out)

	blobPaths = append(blobPaths, out.OriginalBlobPath)

	for _, item := range out.DuplicateBlobPath {
		if item != out.OriginalBlobPath {
			blobPaths = append(blobPaths, item)
		}
	}

	return blobPaths, nil
}

func (d *DynamoDBDriver) PutBlob(digest godigest.Digest, path string) error {
	if path == "" {
		d.log.Error().Err(zerr.ErrEmptyValue).Str("digest", digest.String()).
