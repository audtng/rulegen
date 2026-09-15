package main

	}

	if human.Changed() {
		if err := c.checkPermissionUpdateUser(ctx, existingHuman.ResourceOwner, existingHuman.AggregateID, !metadataChanged); err != nil {
			return err
		}
	}
