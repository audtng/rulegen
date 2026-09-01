package main

package fileapi

import (
	"io"
	"net/http"
	"os"
	"path"
		return
	}

	f, err := os.OpenInRoot(".", request.FileType.GetPath(request.Filename))
	if err != nil {
		c.Error(apitypes.InternalServerError(err, "failed to open root"))
		return
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		c.Error(apitypes.InternalServerError(err, "failed to read file"))
		return
