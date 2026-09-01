package main

		log.Debugf("[creating structure] Created bucket %d, old ID was %d", bucket.ID, oldID)
	}

	// project_view_id is intentionally not writable through Bucket.Update
	// (blocks cross-tenant relocation, GHSA-569v). The importer legitimately
	// remaps buckets onto same-project views, so persist that column directly.
	persistBucketView := func(b *models.Bucket) error {
		_, err := s.Where("id = ?", b.ID).Cols("project_view_id").Update(b)
		return err
	}

	// Create all views, create default views if we don't have any
	viewsByOldIDs := make(map[int64]*models.ProjectView, len(oldViews))
	if len(oldViews) > 0 {
			}

			bucket.ProjectViewID = newView.ID
			err = persistBucketView(bucket)
			if err != nil {
				return
			}
			if view.ViewKind == models.ProjectViewKindKanban {
				for _, b := range bucketsByOldID {
					b.ProjectViewID = view.ID
					err = persistBucketView(b)
					if err != nil {
						return
					}
			"title",
			"limit",
			"position",
		).
		Update(b)
	return
