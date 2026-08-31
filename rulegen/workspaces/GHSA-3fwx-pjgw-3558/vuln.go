package main

	if err != nil {
		return err
	}
	return idtools.MkdirAllAndChown(p, 0701, idtools.CurrentIdentity())
}
	}
	ctr.RWLayer = rwLayer

	if err := idtools.MkdirAndChown(ctr.Root, 0701, idtools.CurrentIdentity()); err != nil {
		return nil, err
	}
	if err := idtools.MkdirAndChown(ctr.CheckpointDir(), 0700, idtools.CurrentIdentity()); err != nil {
		return nil, err
	}

	}

	daemonRepo := filepath.Join(config.Root, "containers")
	if err := idtools.MkdirAllAndChown(daemonRepo, 0701, idtools.CurrentIdentity()); err != nil {
		return nil, err
	}

		}
	}

	// if user namespaces are enabled we will create a subtree underneath the specified root
	// with any/all specified remapped root uid/gid options on the daemon creating
	// a new subdirectory with ownership set to the remapped uid/gid (so as to allow
	// `chdir()` to work for containers namespaced to that uid/gid)
	if config.RemappedRoot != "" {
		id := idtools.CurrentIdentity()
		// First make sure the current root dir has the correct perms.
		if err := idtools.MkdirAllAndChown(config.Root, 0701, id); err != nil {
			return errors.Wrapf(err, "could not create or set daemon root permissions: %s", config.Root)
		}

		config.Root = filepath.Join(rootDir, fmt.Sprintf("%d.%d", remappedRoot.UID, remappedRoot.GID))
		logrus.Debugf("Creating user namespaced daemon root: %s", config.Root)
		// Create the root directory if it doesn't exist
		if err := idtools.MkdirAllAndChown(config.Root, 0701, id); err != nil {
			return fmt.Errorf("Cannot create daemon root: %s: %v", config.Root, err)
		}
		// we also need to verify that any pre-existing directories in the path to
	}

	currentID := idtools.CurrentIdentity()
	// Create the root aufs driver dir
	if err := idtools.MkdirAllAndChown(root, 0701, currentID); err != nil {
		return nil, err
	}

	// Populate the dir structure
	for _, p := range paths {
		if err := idtools.MkdirAllAndChown(path.Join(root, p), 0701, currentID); err != nil {
			return nil, err
		}
	}
		return nil, graphdriver.ErrPrerequisites
	}

	if err := idtools.MkdirAllAndChown(home, 0701, idtools.CurrentIdentity()); err != nil {
		return nil, err
	}

	if err != nil {
		return err
	}
	if err := idtools.MkdirAllAndChown(subvolumes, 0701, idtools.CurrentIdentity()); err != nil {
		return err
	}
	if parent == "" {
		return nil, graphdriver.ErrNotSupported
	}

	if err := idtools.MkdirAllAndChown(path.Join(home, linkDir), 0701, idtools.CurrentIdentity()); err != nil {
		return nil, err
	}

	}
	root := idtools.Identity{UID: rootUID, GID: rootGID}

	currentID := idtools.CurrentIdentity()
	if err := idtools.MkdirAllAndChown(path.Dir(dir), 0701, currentID); err != nil {
		return err
	}
	if err := idtools.MkdirAndChown(dir, 0701, currentID); err != nil {
		return err
	}

		return nil
	}

	if err := idtools.MkdirAndChown(path.Join(dir, workDirName), 0701, currentID); err != nil {
		return err
	}

		logrus.WithField("storage-driver", "overlay").Warn(overlayutils.ErrDTypeNotSupported("overlay", backingFs))
	}

	// Create the driver home dir
	if err := idtools.MkdirAllAndChown(home, 0701, idtools.CurrentIdentity()); err != nil {
		return nil, err
	}

	d := &Driver{
		home:          home,
		uidMaps:       uidMaps,
	root := idtools.Identity{UID: rootUID, GID: rootGID}

	currentID := idtools.CurrentIdentity()
	if err := idtools.MkdirAllAndChown(path.Dir(dir), 0701, currentID); err != nil {
		return err
	}
	if err := idtools.MkdirAndChown(dir, 0701, currentID); err != nil {
		return err
	}

		logger.Warn(overlayutils.ErrDTypeNotSupported("overlay2", backingFs))
	}

	if err := idtools.MkdirAllAndChown(path.Join(home, linkDir), 0701, idtools.CurrentIdentity()); err != nil {
		return nil, err
	}

		return err
	}
	root := idtools.Identity{UID: rootUID, GID: rootGID}
	current := idtools.CurrentIdentity()

	if err := idtools.MkdirAllAndChown(path.Dir(dir), 0701, current); err != nil {
		return err
	}
	if err := idtools.MkdirAndChown(dir, 0701, current); err != nil {
		return err
	}

	if err := d.parseOptions(options); err != nil {
		return nil, err
	}

	if err := idtools.MkdirAllAndChown(home, 0701, idtools.CurrentIdentity()); err != nil {
		return nil, err
	}

func (d *Driver) create(id, parent string, size uint64) error {
	dir := d.dir(id)
	rootIDs := d.idMapping.RootPair()
	if err := idtools.MkdirAllAndChown(filepath.Dir(dir), 0701, idtools.CurrentIdentity()); err != nil {
		return err
	}
	if err := idtools.MkdirAndChown(dir, 0755, rootIDs); err != nil {
		return nil, fmt.Errorf("BUG: zfs get all -t filesystem -rHp '%s' should contain '%s'", options.fsName, options.fsName)
	}

	if err := idtools.MkdirAllAndChown(base, 0701, idtools.CurrentIdentity()); err != nil {
		return nil, fmt.Errorf("Failed to create '%s': %v", base, err)
	}

