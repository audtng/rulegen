package main

		count, err := db.GetEngine(ctx).Count(&LFSMetaObject{Pointer: lfs.Pointer{Oid: oid}})
		return count > 0, err
	}
	// LFS objects are repository code content, so authorization must require
	// Code-unit access; other unit accesses (e.g. Issues) must not authorize
	// reuse of an existing LFS object across repositories.
	cond := repo_model.AccessibleRepositoryCondition(user, unit.TypeCode)
	count, err := db.GetEngine(ctx).Where(cond).Join("INNER", "repository", "`lfs_meta_object`.repository_id = `repository`.id").Count(&LFSMetaObject{Pointer: lfs.Pointer{Oid: oid}})
	return count > 0, err
}
			newMetas := make([]*LFSMetaObject, 0, len(metas))
			cond := builder.In(
				"`lfs_meta_object`.repository_id",
				builder.Select("`repository`.id").From("repository").Where(repo_model.AccessibleRepositoryCondition(user, unit.TypeCode)),
			)
			if err := db.GetEngine(ctx).Cols("oid").Where(cond).In("oid", oids...).GroupBy("oid").Find(&newMetas); err != nil {
				return err
