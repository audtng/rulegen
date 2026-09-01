package main

		log.Debugf("[creating structure] Created bucket %d, old ID was %d", bucket.ID, oldID)
	}

	// Create all views, create default views if we don't have any
	viewsByOldIDs := make(map[int64]*models.ProjectView, len(oldViews))
	if len(oldViews) > 0 {
			}

			bucket.ProjectViewID = newView.ID
			err = bucket.Update(s, user)
			if err != nil {
				return
			}
			if view.ViewKind == models.ProjectViewKindKanban {
				for _, b := range bucketsByOldID {
					b.ProjectViewID = view.ID
					err = b.Update(s, user)
					if err != nil {
						return
					}
			"title",
			"limit",
			"position",
			"project_view_id",
		).
		Update(b)
	return
