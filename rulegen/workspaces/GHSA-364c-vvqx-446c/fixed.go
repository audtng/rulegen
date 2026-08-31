package main

	fmt.Fprintf(os.Stderr, "\n")
	return nil
}

// ValidFileName checks if a filename is valid
// and returns true only if it all of the characters are either
// 0-9, a-z, A-Z, ., _, -, space, or /
func ValidFileName(fname string) bool {
	for _, r := range fname {
		if !((r >= '0' && r <= '9') ||
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			r == '.' || r == '_' || r == '-' || r == ' ' || r == '/') {
			return false
		}
	}
	return true
}
	c.EmptyFoldersToTransfer = senderInfo.EmptyFoldersToTransfer
	c.TotalNumberFolders = senderInfo.TotalNumberFolders
	c.FilesToTransfer = senderInfo.FilesToTransfer
	for i, fi := range c.FilesToTransfer {
		// Issues #593 - sanitize the sender paths and prevent ".." from being used
		c.FilesToTransfer[i].FolderRemote = filepath.Clean(fi.FolderRemote)
		if strings.Contains(c.FilesToTransfer[i].FolderRemote, "..") {
			return true, fmt.Errorf("invalid path detected: '%s'", fi.FolderRemote)
		}
		// Issues #593 - disallow specific folders like .ssh
		if strings.Contains(c.FilesToTransfer[i].FolderRemote, ".ssh") {
			return true, fmt.Errorf("invalid path detected: '%s'", fi.FolderRemote)
		}

	}
	c.TotalNumberOfContents = 0
	if c.FilesToTransfer != nil {
		c.TotalNumberOfContents += len(c.FilesToTransfer)
		filePath := filepath.Join(destination, f.Name)
		fmt.Fprintf(os.Stderr, "\r\033[2K")
		fmt.Fprintf(os.Stderr, "\rUnzipping file %s", filePath)
		// Issue #593 conceal path traversal vulnerability
		// make sure the filepath does not have ".."
		filePath = filepath.Clean(filePath)
		if strings.Contains(filePath, "..") {
			log.Fatalf("Invalid file path %s\n", filePath)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
src/croc/croc.go   |  4 ++++
src/utils/utils.go | 15 +++++++++++++++
2 files changed, 19 insertions(+)
