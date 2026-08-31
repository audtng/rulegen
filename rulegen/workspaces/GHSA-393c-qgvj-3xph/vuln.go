package main

	return rel, nil
}

// GetLFSLockByID returns release by given id.
func GetLFSLockByID(ctx context.Context, id int64) (*LFSLock, error) {
	lock := new(LFSLock)
	has, err := db.GetEngine(ctx).ID(id).Get(lock)
	if err != nil {
		return nil, err
	} else if !has {
// DeleteLFSLockByID deletes a lock by given ID.
func DeleteLFSLockByID(ctx context.Context, id int64, repo *repo_model.Repository, u *user_model.User, force bool) (*LFSLock, error) {
	return db.WithTx2(ctx, func(ctx context.Context) (*LFSLock, error) {
		lock, err := GetLFSLockByID(ctx, id)
		if err != nil {
			return nil, err
		}
			})
			return
		}
		lock, err := git_model.GetLFSLockByID(ctx, v)
		if err != nil && !git_model.IsErrLFSLockNotExist(err) {
			log.Error("Unable to get lock with ID[%s]: Error: %v", v, err)
		}
