package main

	fmt.Fprintf(os.Stderr, "\n")
	return nil
}
	c.EmptyFoldersToTransfer = senderInfo.EmptyFoldersToTransfer
	c.TotalNumberFolders = senderInfo.TotalNumberFolders
	c.FilesToTransfer = senderInfo.FilesToTransfer
	c.TotalNumberOfContents = 0
	if c.FilesToTransfer != nil {
		c.TotalNumberOfContents += len(c.FilesToTransfer)
		filePath := filepath.Join(destination, f.Name)
		fmt.Fprintf(os.Stderr, "\r\033[2K")
		fmt.Fprintf(os.Stderr, "\rUnzipping file %s", filePath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
src/croc/croc.go   |  4 ++++
src/utils/utils.go | 15 +++++++++++++++
2 files changed, 19 insertions(+)
