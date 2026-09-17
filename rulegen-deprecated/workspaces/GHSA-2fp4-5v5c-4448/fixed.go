package main

		return nil, err
	}

	absPath, err := fileutil.SafeJoin(s.basePath, relPath)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("stat m3u: %w", err)
		return nil, errors.New("path is a directory")
	}

	// derive metadata from the cleaned relative path so traversal segments
	// can't spoof the first-segment user id check.
	cleanRel, _ := filepath.Rel(s.basePath, absPath)

	var playlist Playlist
	playlist.UpdatedAt = stat.ModTime()

	playlist.UserID, err = userIDFromPath(cleanRel)
	if err != nil {
		playlist.UserID = 1
	}

	playlist.Name = strings.TrimSuffix(filepath.Base(cleanRel), filepath.Ext(cleanRel))

	file, err := os.Open(absPath)
	if err != nil {
		return err
	}

	absPath, err := fileutil.SafeJoin(s.basePath, relPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return fmt.Errorf("make m3u base dir: %w", err)
	}
	file, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE, 0o666)
		return err
	}

	absPath, err := fileutil.SafeJoin(s.basePath, relPath)
	if err != nil {
		return err
	}
	return os.Remove(absPath)
}

func firstPathEl(path string) string {
package fileutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrEscapesBase = errors.New("path escapes base")

var nonAlphaNumExpr = regexp.MustCompile("[^a-zA-Z0-9_.]+")

const maxFilenameLength = 200
func HasPrefix(p, prefix string) bool {
	return p == prefix || strings.HasPrefix(p, filepath.Clean(prefix)+string(filepath.Separator))
}

func SafeJoin(base, rel string) (string, error) {
	abs := filepath.Join(base, rel)
	if !HasPrefix(abs, base) {
		return "", fmt.Errorf("%w: %q", ErrEscapesBase, rel)
	}
	return abs, nil
}

	var playlist playlistp.Playlist
	if playlistPath != "" {
		if pl, err := c.playlistStore.Read(playlistPath); err == nil && pl != nil {
			playlist = *pl
		}
	}
