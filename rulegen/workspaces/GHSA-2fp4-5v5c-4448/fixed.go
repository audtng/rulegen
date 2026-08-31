package main

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

	var playlist playlistp.Playlist
	if playlistPath != "" {
		if pl, err := c.playlistStore.Read(playlistPath); err == nil && pl != nil {
			playlist = *pl
		}
	}

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	}
}

func TestPlaylistTraversalDenied(t *testing.T) {
	t.Parallel()
	f := newFixture(t)
	writePrivatePlaylist(t, f)

	encode := func(parts ...string) string {
		return playlistIDEncode(filepath.Join(parts...)).String()
	}
	altID := fmt.Sprint(f.alt.ID)

	cases := []struct {
		name    string
		handler handlerSubsonic
		id      string
	}{
		{"get_cross_user", f.contr.ServeGetPlaylist, encode(altID, "..", "1", "private.m3u")},
		{"delete_cross_user", f.contr.ServeDeletePlaylist, encode(altID, "..", "1", "shared.m3u")},
		{"create_escapes_base", f.contr.ServeCreateOrUpdatePlaylist, encode("..", "injected.m3u")},
	}
	for _, tc := range cases {
		body := f.query(t, tc.handler, f.alt, url.Values{"id": {tc.id}, "name": {"x"}})
		var sub spec.SubsonicResponse
		if err := json.Unmarshal([]byte(body), &sub); err != nil {
			t.Fatalf("%s: unmarshal: %v", tc.name, err)
		}
		if sub.Response.Status != "failed" || sub.Response.Error == nil {
			t.Fatalf("%s: expected failure, got: %s", tc.name, body)
		}
	}

	if _, err := f.contr.playlistStore.Read(filepath.Join("1", "shared.m3u")); err != nil {
		t.Fatalf("shared playlist deleted via traversal: %v", err)
	}
	escaped := filepath.Join(f.contr.playlistStore.BasePath(), "..", "injected.m3u")
	if _, err := os.Stat(escaped); !os.IsNotExist(err) {
		t.Fatalf("file written outside playlists dir: stat err=%v", err)
	}
}

func writePrivatePlaylist(t *testing.T, f *fixture) string {
	t.Helper()
	relPath := filepath.Join("1", "private.m3u")
