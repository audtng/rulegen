package main

		return nil, err
	}

	absPath := filepath.Join(s.basePath, relPath)
	stat, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("stat m3u: %w", err)
		return nil, errors.New("path is a directory")
	}

	var playlist Playlist
	playlist.UpdatedAt = stat.ModTime()

	playlist.UserID, err = userIDFromPath(relPath)
	if err != nil {
		playlist.UserID = 1
	}

	playlist.Name = strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath))

	file, err := os.Open(absPath)
	if err != nil {
		return err
	}

	absPath := filepath.Join(s.basePath, relPath)
	if err := os.MkdirAll(filepath.Dir(absPath), 0o777); err != nil {
		return fmt.Errorf("make m3u base dir: %w", err)
	}
	file, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE, 0o666)
		return err
	}

	return os.Remove(filepath.Join(s.basePath, relPath))
}

func firstPathEl(path string) string {
package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var nonAlphaNumExpr = regexp.MustCompile("[^a-zA-Z0-9_.]+")

const maxFilenameLength = 200
func HasPrefix(p, prefix string) bool {
	return p == prefix || strings.HasPrefix(p, filepath.Clean(prefix)+string(filepath.Separator))
}

	var playlist playlistp.Playlist
	if playlistPath != "" {
		if pl, err := c.playlistStore.Read(playlistPath); err != nil && pl != nil {
			playlist = *pl
		}
	}
