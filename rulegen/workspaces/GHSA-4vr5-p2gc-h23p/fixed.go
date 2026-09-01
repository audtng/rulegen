package main

	}
	// extract files
	err = ex.Extract(ctx, in, func(ctx context.Context, f archives.FileInfo) error {
		remote, err := destPath(f.NameInArchive, dstDir)
		if err != nil {
			return err
		}
		// Skip the archive root entry ("./") which has no name of its own.
		if remote == "" {
			return nil
		}
		// check if file should be extracted
		if !fi.Include(remote, f.Size(), f.ModTime(), fs.Metadata{}) {

	return err
}

// destPath maps an archive entry name onto its destination remote within
// dstDir, returning an error if the name is unsafe.
//
// Archive entry names are attacker controlled. A leading "./" is stripped:
// tar archives created with relative paths (e.g. "tar -czf archive.tar.gz .")
// use "./" prefixed entries and, without stripping, rclone would encode the
// "." as a full-width dot character creating a spurious directory.
//
// Any entry with a ".." path component is then rejected to prevent a path
// traversal ("Zip Slip") attack: path.Join collapses "..", so an entry such
// as "../escaped.txt" joined onto "dir" would resolve to "escaped.txt" and be
// written outside the selected destination directory. Both "/" and "\" are
// treated as separators when looking for ".." segments: archive names should
// use "/", but a crafted archive may use "\", which the local backend treats
// as a path separator on Windows.
//
// The returned remote is empty for the archive root entry ("./"), which the
// caller should skip.
func destPath(nameInArchive, dstDir string) (string, error) {
	remote := strings.TrimPrefix(nameInArchive, "./")
	isSeparator := func(r rune) bool { return r == '/' || r == '\\' }
	for _, segment := range strings.FieldsFunc(remote, isSeparator) {
		if segment == ".." {
			return "", fmt.Errorf("refusing to extract archive entry %q with a %q path component", nameInArchive, "..")
		}
	}
	if remote == "" {
		return "", nil
	}
	if dstDir != "" {
		remote = path.Join(dstDir, remote)
	}
	return remote, nil
}
