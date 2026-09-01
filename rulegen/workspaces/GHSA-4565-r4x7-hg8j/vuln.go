package main

}

// ChangeCollaborationAccessMode sets new access mode for the collaboration.
func (r *Repository) ChangeCollaborationAccessMode(userID int64, mode AccessMode) error {
	// Discard invalid input
	if mode <= AccessModeNone || mode > AccessModeOwner {
		return nil
	}

	}

	if form.Permission != nil {
		if err := c.Repo.Repository.ChangeCollaborationAccessMode(collaborator.ID, database.ParseAccessMode(*form.Permission)); err != nil {
			c.Error(err, "change collaboration access mode")
			return
		}

func ChangeCollaborationAccessMode(c *context.Context) {
	if err := c.Repo.Repository.ChangeCollaborationAccessMode(
		c.QueryInt64("uid"),
		database.AccessMode(c.QueryInt("mode"))); err != nil {
		log.Error("ChangeCollaborationAccessMode: %v", err)
