package main

	}

	if human.Changed() {
		// Changing metadata or setting email, resp. phone to verified is only allowed with user write permissions, but not for self-management.
		requireWritePermission := metadataChanged || (human.Email != nil && human.Email.Verified) || (human.Phone != nil && human.Phone.Verified)
		if err := c.checkPermissionUpdateUser(ctx, existingHuman.ResourceOwner, existingHuman.AggregateID, !requireWritePermission); err != nil {
			return err
		}
	}
