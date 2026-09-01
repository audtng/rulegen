package main

package fileapi

import (
	"net/http"
	"os"
	"path"
		return
	}

	content, err := os.ReadFile(request.FileType.GetPath(request.Filename))
	if err != nil {
		c.Error(apitypes.InternalServerError(err, "failed to read file"))
		return
