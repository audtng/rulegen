package main

	"strings"
	"time"

	"github.com/spf13/afero"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/rules"
)

var (
		}
	}

	if file != nil {
		ok, scopeErr := WithinScope(opts.Fs, opts.Path)
		if scopeErr != nil || !ok {
			return nil, os.ErrPermission
		}
	}

	// regular file
	if file != nil && !file.IsSymlink {
		return file, nil
	return file, nil
}

// WithinScope reports whether the on-disk target of p — after resolving any
// symbolic links — stays within the scoped root of fsys. It exists to stop a
// symlink that lives lexically inside a user's scope but points outside it
// from being followed for reads, writes, or shares.
//
// Paths that do not exist yet (e.g. a brand-new file being created) are
// validated against their nearest existing ancestor, so legitimate new files
// are always allowed. For a filesystem that is not scoped with BasePathFs the
// check is a no-op and returns true.
//
// Note: a dangling symlink whose target does not yet exist resolves to its
// containing directory and is therefore allowed; writing through such a link
// could still create a file outside the scope. Callers that create files
// should treat this as best-effort and rely on rejecting existing escaping
// symlinks, which covers the disclosure and overwrite vectors.
func WithinScope(fsys afero.Fs, p string) (bool, error) {
	bfs, ok := fsys.(*afero.BasePathFs)
	if !ok {
		// Not a scoped filesystem; nothing to enforce.
		return true, nil
	}

	root, err := filepath.EvalSymlinks(afero.FullBaseFsPath(bfs, "/"))
	if err != nil {
		return false, err
	}

	target := afero.FullBaseFsPath(bfs, p)
	resolved, err := filepath.EvalSymlinks(target)
	for errors.Is(err, fs.ErrNotExist) {
		parent := filepath.Dir(target)
		if parent == target {
			break
		}
		target = parent
		resolved, err = filepath.EvalSymlinks(target)
	}
	if err != nil {
		return false, err
	}

	// Compare against root with a trailing separator so a sibling like
	// "/srvother" is not treated as being inside "/srv". When root is itself
	// the filesystem boundary (e.g. "/"), it already ends in a separator, so
	// avoid producing "//" — which no path would match — and accept any path
	// under it.
	prefix := root
	if !strings.HasSuffix(prefix, string(filepath.Separator)) {
		prefix += string(filepath.Separator)
	}

	return resolved == root || strings.HasPrefix(resolved, prefix), nil
}

// Checksum checksums a given File for a given User, using a specific
// algorithm. The checksums data is saved on File object.
func (i *FileInfo) Checksum(algo string) error {
		isSymlink, isInvalidLink := false, false
		if IsSymlink(f.Mode()) {
			isSymlink = true
			if ok, scopeErr := WithinScope(i.Fs, fPath); scopeErr != nil || !ok {
				continue
			}
			// It's a symbolic link. We try to follow it. If it doesn't work,
			// we stay with the link information instead of the target's.
			info, err := i.Fs.Stat(fPath)
			if err == nil {
				f = info
			} else {
				isInvalidLink = true
			}
		}
	"github.com/spf13/afero"
)

func TestWithinScope(t *testing.T) {
	t.Run("non-scoped filesystem is a no-op", func(t *testing.T) {
		ok, err := WithinScope(afero.NewOsFs(), "/anything")
		if err != nil || !ok {
			t.Fatalf("expected (true, nil), got (%v, %v)", ok, err)
		}
	})

	t.Run("path inside a nested scope is allowed", func(t *testing.T) {
		scope := t.TempDir()
		if err := os.WriteFile(filepath.Join(scope, "file.txt"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		bfs := afero.NewBasePathFs(afero.NewOsFs(), scope)

		ok, err := WithinScope(bfs, "/file.txt")
		if err != nil || !ok {
			t.Fatalf("expected (true, nil), got (%v, %v)", ok, err)
		}
	})

	t.Run("new file inside scope is allowed", func(t *testing.T) {
		scope := t.TempDir()
		bfs := afero.NewBasePathFs(afero.NewOsFs(), scope)

		ok, err := WithinScope(bfs, "/does-not-exist-yet.txt")
		if err != nil || !ok {
			t.Fatalf("expected (true, nil), got (%v, %v)", ok, err)
		}
	})

	// Regression for #5975: when the scope resolves to the filesystem root,
	// root+separator used to be "//", which no path matched, so every write
	// was rejected with os.ErrPermission (HTTP 403).
	t.Run("filesystem root scope allows writes", func(t *testing.T) {
		f := filepath.Join(t.TempDir(), "file.txt")
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		bfs := afero.NewBasePathFs(afero.NewOsFs(), "/")

		ok, err := WithinScope(bfs, f)
		if err != nil || !ok {
			t.Fatalf("expected (true, nil) for a path under root scope, got (%v, %v)", ok, err)
		}
	})

	t.Run("sibling of a nested scope is rejected", func(t *testing.T) {
		base := t.TempDir()
		scope := filepath.Join(base, "srv")
		sibling := filepath.Join(base, "srvother")
				t.Fatal(err)
			}
		}
		// A symlink lexically inside the scope pointing at a sibling directory
		// must not be followed.
		link := filepath.Join(scope, "escape")
		if err := os.Symlink(sibling, link); err != nil {
			t.Fatal(err)
		}
		bfs := afero.NewBasePathFs(afero.NewOsFs(), scope)

		ok, err := WithinScope(bfs, "/escape")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected escaping symlink to a sibling directory to be rejected")
		}
	})

		if err := os.Symlink(filepath.Join(scope, "real"), filepath.Join(scope, "link")); err != nil {
			t.Skipf("cannot create symlink: %v", err)
		}
		bfs := afero.NewBasePathFs(afero.NewOsFs(), scope)

		ok, err := WithinScope(bfs, "/link/f.txt")
		if err != nil || !ok {
			t.Fatalf("expected (true, nil) for an in-scope symlink target, got (%v, %v)", ok, err)
		}
	})
}
	}

	// Filesystem scoped to the shared directory, as a public share would be.
	bfs := afero.NewBasePathFs(afero.NewOsFs(), filepath.Join(scope, "shared"))

	if _, err := stat(&FileOptions{Fs: bfs, Path: "/link/secret.txt"}); !os.IsPermission(err) {
		t.Fatalf("expected permission error for linked-ancestor escape, got %v", err)
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

		t.Skipf("cannot create symlink: %v", err)
	}

	afs := afero.NewBasePathFs(afero.NewOsFs(), scope)

	err := Copy(afs, "/srcdir", "/dstdir", 0o644, 0o755)
	if err == nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	afs := afero.NewBasePathFs(afero.NewOsFs(), scope)

	if err := Copy(afs, "/srcdir", "/dstdir", 0o644, 0o755); err != nil {
		t.Fatalf("expected copy of an in-scope symlink to succeed, got: %v", err)
import (
	"errors"
	"io/fs"
	"os"

	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/files"
)

// CopyDir copies a directory from source to dest and all
// of its sub-directories. It doesn't stop if it finds an error
// during the copy. Returns an error if any.
func CopyDir(afs afero.Fs, source, dest string, fileMode, dirMode fs.FileMode) error {
	if ok, err := files.WithinScope(afs, source); err != nil || !ok {
		if err != nil {
			return err
		}
		return os.ErrPermission
	}

	// Get properties of source.
	srcinfo, err := afs.Stat(source)
	if err != nil {
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/files"
)

// MoveFile moves file from src to dst.
// CopyFile copies a file from source to dest and returns
// an error if any.
func CopyFile(afs afero.Fs, source, dest string, fileMode, dirMode fs.FileMode) error {
	if ok, err := files.WithinScope(afs, source); err != nil || !ok {
		if err != nil {
			return err
		}
		return os.ErrPermission
	}

	// Open the source file.
	src, err := afs.Open(source)
	if err != nil {
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"golang.org/x/crypto/bcrypt"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/share"
)

var withHashFile = func(fn handleFunc) handleFunc {
			filePath = ifPath
		}

		// set fs root to the shared file/folder
		d.user.Fs = afero.NewBasePathFs(d.user.Fs, basePath)

		// the filesystem is now rebased onto basePath, so paths handed to the
		// rule checker are relative to it. Resolve them back to the user's
	"testing"

	"github.com/asdine/storm/v3"
	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/share"
	"github.com/filebrowser/filebrowser/v2/storage"
	"github.com/filebrowser/filebrowser/v2/storage/bolt"
	"github.com/filebrowser/filebrowser/v2/users"
)

// symlinkShareStorage builds a storage whose single user is rooted at a real
	}
	st.Users = &customFSUser{
		Store: st.Users,
		fs:    afero.NewBasePathFs(afero.NewOsFs(), scope),
	}
	return st
}
	"testing"

	"github.com/asdine/storm/v3"
	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/rules"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/share"
	"github.com/filebrowser/filebrowser/v2/storage/bolt"
	"github.com/filebrowser/filebrowser/v2/users"
)

func TestPublicShareHandlerAuthentication(t *testing.T) {
				t.Fatalf("failed to save settings: %v", err)
			}

			fs := afero.NewBasePathFs(afero.NewOsFs(), t.TempDir())
			if err := fs.MkdirAll("/projects/private", 0o755); err != nil {
				t.Fatalf("failed to create private dir: %v", err)
			}
	if err != nil {
		return nil, err
	}
	user.Fs = cu.fs

	return user, nil
}
	"path/filepath"
	"strings"

	"github.com/mholt/archives"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/fileutils"
	"github.com/filebrowser/filebrowser/v2/users"
)

func slashClean(name string) string {
		return nil, nil
	}

	if ok, err := files.WithinScope(d.user.Fs, path); err != nil || !ok {
		return nil, nil
	}

	info, err := d.user.Fs.Stat(path)
	if err != nil {
		return nil, err
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
	"github.com/spf13/afero"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/fileutils"
)

var resourceGetHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
			return http.StatusForbidden, nil
		}

		for _, p := range []string{src, dst} {
			if ok, scopeErr := files.WithinScope(d.user.Fs, p); scopeErr != nil || !ok {
				if scopeErr != nil {
					return errToStatus(scopeErr), scopeErr
				}
				return http.StatusForbidden, nil
			}
		}

		err = checkParent(src, dst)
		if err != nil {
			return http.StatusBadRequest, err
}

func writeFile(afs afero.Fs, dst string, in io.Reader, fileMode, dirMode fs.FileMode) (os.FileInfo, error) {
	if ok, err := files.WithinScope(afs, dst); err != nil || !ok {
		return nil, os.ErrPermission
	}

	dir, _ := path.Split(dst)
	err := afs.MkdirAll(dir, dirMode)
	if err != nil {
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/share"
)

func withPermShare(fn handleFunc) handleFunc {
	// Only allow sharing paths that currently exist. Otherwise a share could be
	// created for a non-existent path and would silently start exposing
	// whatever file later appears there.
	if _, err := d.user.Fs.Stat(r.URL.Path); err != nil {
		return errToStatus(err), err
	}

	// Refuse to create a share whose on-disk target escapes the user's scope
	// (e.g. via a symlink), mirroring the read/write guards. The public serve
	// path already rejects these, but blocking creation avoids dangling shares.
	if ok, err := files.WithinScope(d.user.Fs, r.URL.Path); err != nil || !ok {
		if err != nil {
			return errToStatus(err), err
		}
		return http.StatusForbidden, nil
	}

	var s *share.Link
	var body share.CreateBody
	if r.Body != nil {
	"strings"
	"time"

	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/files"
)

// keepUploadActive periodically touches the cache entry to prevent eviction during transfer
			return http.StatusForbidden, nil
		}

		if ok, scopeErr := files.WithinScope(d.user.Fs, r.URL.Path); scopeErr != nil || !ok {
			return http.StatusForbidden, nil
		}

		file, err := files.NewFileInfo(&files.FileOptions{
			Fs:         d.user.Fs,
			Path:       r.URL.Path,
			return http.StatusUnsupportedMediaType, nil
		}

		if ok, scopeErr := files.WithinScope(d.user.Fs, r.URL.Path); scopeErr != nil || !ok {
			return http.StatusForbidden, nil
		}

		uploadOffset, err := getUploadOffset(r)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("invalid upload offset")
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/storage/bolt"
	"github.com/filebrowser/filebrowser/v2/users"
	}
	st.Users = &customFSUser{
		Store: st.Users,
		fs:    afero.NewBasePathFs(afero.NewOsFs(), userScope),
	}

	// Forge a valid auth token for user ID 1.
import (
	"path/filepath"

	"github.com/spf13/afero"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/rules"
)

// ViewMode describes a view mode.

// User describes a user.
type User struct {
	ID                    uint          `storm:"id,increment" json:"id"`
	Username              string        `storm:"unique" json:"username"`
	Password              string        `json:"password"`
	Scope                 string        `json:"scope"`
	Locale                string        `json:"locale"`
	LockPassword          bool          `json:"lockPassword"`
	ViewMode              ViewMode      `json:"viewMode"`
	SingleClick           bool          `json:"singleClick"`
	RedirectAfterCopyMove bool          `json:"redirectAfterCopyMove"`
	Perm                  Permissions   `json:"perm"`
	Commands              []string      `json:"commands"`
	Sorting               files.Sorting `json:"sorting"`
	Fs                    afero.Fs      `json:"-" yaml:"-"`
	Rules                 []rules.Rule  `json:"rules"`
	HideDotfiles          bool          `json:"hideDotfiles"`
	DateFormat            bool          `json:"dateFormat"`
	AceEditorTheme        string        `json:"aceEditorTheme"`
}

// GetRules implements rules.Provider.
	if u.Fs == nil {
		scope := u.Scope
		scope = filepath.Join(baseScope, filepath.Join("/", scope))
		u.Fs = afero.NewBasePathFs(afero.NewOsFs(), scope)
	}

	return nil

// FullPath gets the full path for a user's relative path.
func (u *User) FullPath(path string) string {
	return afero.FullBaseFsPath(u.Fs.(*afero.BasePathFs), path)
}
