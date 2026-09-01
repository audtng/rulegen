package main

	}
	// extract files
	err = ex.Extract(ctx, in, func(ctx context.Context, f archives.FileInfo) error {
		remote := f.NameInArchive
		// Strip leading "./" from archive paths. Tar files created with
		// relative paths (e.g. "tar -czf archive.tar.gz .") use "./" prefixed
		// entries. Without stripping, rclone encodes the "." as a full-width
		// dot character creating a spurious directory. We only strip "./"
		// specifically to avoid enabling path traversal attacks via "../".
		remote = strings.TrimPrefix(remote, "./")
		// If the entry was exactly "./" (the root dir), skip it
		if remote == "" && f.IsDir() {
			return nil
		}
		if dstDir != "" {
			remote = path.Join(dstDir, remote)
		}
		// check if file should be extracted
		if !fi.Include(remote, f.Size(), f.ModTime(), fs.Metadata{}) {

	return err
}
