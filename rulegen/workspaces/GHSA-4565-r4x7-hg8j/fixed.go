package main

}

// ChangeCollaborationAccessMode sets new access mode for the collaboration.
// The actor's access mode bounds the new mode, so an actor can never grant a
// level higher than their own.
func (r *Repository) ChangeCollaborationAccessMode(actorMode AccessMode, userID int64, mode AccessMode) error {
	// Collaborators can hold at most admin access.
	if mode <= AccessModeNone || mode > AccessModeAdmin {
		return nil
	}

	// Actors must not grant a level above their own.
	if mode > actorMode {
		return nil
	}

	}

	if form.Permission != nil {
		if err := c.Repo.Repository.ChangeCollaborationAccessMode(c.Repo.AccessMode, collaborator.ID, database.ParseAccessMode(*form.Permission)); err != nil {
			c.Error(err, "change collaboration access mode")
			return
		}

func ChangeCollaborationAccessMode(c *context.Context) {
	if err := c.Repo.Repository.ChangeCollaborationAccessMode(
		c.Repo.AccessMode,
		c.QueryInt64("uid"),
		database.AccessMode(c.QueryInt("mode"))); err != nil {
		log.Error("ChangeCollaborationAccessMode: %v", err)
