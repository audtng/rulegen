package main

	"strings"
	"time"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/rules"
	"github.com/spf13/afero"
)

var (
		}
	}

	// regular file
	if file != nil && !file.IsSymlink {
		return file, nil
	return file, nil
}

// Checksum checksums a given File for a given User, using a specific
// algorithm. The checksums data is saved on File object.
func (i *FileInfo) Checksum(algo string) error {
		isSymlink, isInvalidLink := false, false
		if IsSymlink(f.Mode()) {
			isSymlink = true
			// It's a symbolic link. We try to follow it. The scoped filesystem
			// refuses to dereference a link whose target escapes the scope
			// (permission error); such a link is omitted from the listing
			// entirely so it cannot leak the target's metadata. Any other
			// failure means a broken link, which we surface as an invalid link
			// rather than the target's information.
			info, err := i.Fs.Stat(fPath)
			switch {
			case err == nil:
				f = info
			case errors.Is(err, os.ErrPermission):
				continue
			default:
				isInvalidLink = true
			}
		}
	"github.com/spf13/afero"
)

func TestScopedFs(t *testing.T) {
	t.Run("path inside scope is allowed", func(t *testing.T) {
		scope := t.TempDir()
		if err := os.WriteFile(filepath.Join(scope, "file.txt"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		fs := NewScopedFs(afero.NewOsFs(), scope)

		if _, err := fs.Stat("/file.txt"); err != nil {
			t.Fatalf("expected in-scope file to be accessible, got %v", err)
		}
	})

	t.Run("new file inside scope can be created", func(t *testing.T) {
		scope := t.TempDir()
		fs := NewScopedFs(afero.NewOsFs(), scope)

		f, err := fs.OpenFile("/does-not-exist-yet.txt", os.O_RDWR|os.O_CREATE, 0o644)
		if err != nil {
			t.Fatalf("expected to create a new in-scope file, got %v", err)
		}
		_ = f.Close()
	})

	// Regression for #5975: when the scope resolves to the filesystem root,
	// root+separator used to be "//", which no path matched, so every write
	// was rejected with os.ErrPermission (HTTP 403).
	t.Run("filesystem root scope allows access", func(t *testing.T) {
		f := filepath.Join(t.TempDir(), "file.txt")
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		fs := NewScopedFs(afero.NewOsFs(), "/")

		if _, err := fs.Stat(f); err != nil {
			t.Fatalf("expected a path under root scope to be accessible, got %v", err)
		}
	})

	t.Run("escaping symlink to a sibling is rejected", func(t *testing.T) {
		base := t.TempDir()
		scope := filepath.Join(base, "srv")
		sibling := filepath.Join(base, "srvother")
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(sibling, "secret.txt"), []byte("secret"), 0o644); err != nil {
			t.Fatal(err)
		}
		// A symlink lexically inside the scope pointing at a sibling directory
		// must not be followed for reads or stats.
		if err := os.Symlink(sibling, filepath.Join(scope, "escape")); err != nil {
			t.Skipf("cannot create symlink: %v", err)
		}
		fs := NewScopedFs(afero.NewOsFs(), scope)

		if _, err := fs.Stat("/escape"); !os.IsPermission(err) {
			t.Fatalf("expected stat of escaping symlink to be rejected, got %v", err)
		}
		if _, err := fs.Open("/escape/secret.txt"); !os.IsPermission(err) {
			t.Fatalf("expected read through escaping symlink to be rejected, got %v", err)
		}
	})

		if err := os.Symlink(filepath.Join(scope, "real"), filepath.Join(scope, "link")); err != nil {
			t.Skipf("cannot create symlink: %v", err)
		}
		fs := NewScopedFs(afero.NewOsFs(), scope)

		if _, err := fs.Stat("/link/f.txt"); err != nil {
			t.Fatalf("expected in-scope symlink target to be accessible, got %v", err)
		}
	})
}
	}

	// Filesystem scoped to the shared directory, as a public share would be.
	bfs := NewScopedFs(afero.NewOsFs(), filepath.Join(scope, "shared"))

	if _, err := stat(&FileOptions{Fs: bfs, Path: "/link/secret.txt"}); !os.IsPermission(err) {
		t.Fatalf("expected permission error for linked-ancestor escape, got %v", err)
package files

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
)

// ScopedFs is an afero.Fs that confines every operation to a base directory and
// refuses to follow a symbolic link whose on-disk target resolves outside that
// base. It wraps an *afero.BasePathFs — which already provides the lexical
// confinement — and adds a per-operation scope check on every call that would
// dereference a symlink at the OS layer (open, stat, lstat, chmod, …).
type ScopedFs struct {
	base *afero.BasePathFs
}

var (
	_ afero.Fs      = (*ScopedFs)(nil)
	_ afero.Lstater = (*ScopedFs)(nil)
)

func NewScopedFs(source afero.Fs, path string) *ScopedFs {
	if s, ok := source.(*ScopedFs); ok {
		source = s.base
	}
	return &ScopedFs{base: afero.NewBasePathFs(source, path).(*afero.BasePathFs)}
}

// Base returns the underlying *afero.BasePathFs.
func (s *ScopedFs) Base() *afero.BasePathFs { return s.base }

// guard returns an error if name's on-disk target resolves outside the scope.
func (s *ScopedFs) guard(name string) error {
	ok, err := s.within(name)
	if err != nil {
		return err
	}
	if !ok {
		return os.ErrPermission
	}
	return nil
}

// within reports whether the on-disk target of p — after resolving any symbolic
// links — stays within the scoped root. It exists to stop a symlink that lives
// lexically inside the scope but points outside it from being followed for
// reads, writes, or shares.
//
// Paths that do not exist yet (e.g. a brand-new file being created) are
// validated against their nearest existing ancestor, so legitimate new files
// are always allowed.
//
// Note: a dangling symlink whose target does not yet exist resolves to its
// containing directory and is therefore allowed; writing through such a link
// could still create a file outside the scope. This is treated as best-effort
// and relies on rejecting existing escaping symlinks, which covers the
// disclosure and overwrite vectors.
func (s *ScopedFs) within(p string) (bool, error) {
	root, err := filepath.EvalSymlinks(afero.FullBaseFsPath(s.base, "/"))
	if err != nil {
		return false, err
	}

	target := afero.FullBaseFsPath(s.base, p)
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
	// "/srvother" is not treated as being inside "/srv". When root is itself the
	// filesystem boundary (e.g. "/"), it already ends in a separator, so avoid
	// producing "//" — which no path would match — and accept any path under it.
	prefix := root
	if !strings.HasSuffix(prefix, string(filepath.Separator)) {
		prefix += string(filepath.Separator)
	}

	return resolved == root || strings.HasPrefix(resolved, prefix), nil
}

func (s *ScopedFs) Create(name string) (afero.File, error) {
	if err := s.guard(name); err != nil {
		return nil, err
	}
	return s.base.Create(name)
}

func (s *ScopedFs) Mkdir(name string, perm os.FileMode) error {
	if err := s.guard(name); err != nil {
		return err
	}
	return s.base.Mkdir(name, perm)
}

func (s *ScopedFs) MkdirAll(path string, perm os.FileMode) error {
	if err := s.guard(path); err != nil {
		return err
	}
	return s.base.MkdirAll(path, perm)
}

func (s *ScopedFs) Open(name string) (afero.File, error) {
	if err := s.guard(name); err != nil {
		return nil, err
	}
	return s.base.Open(name)
}

func (s *ScopedFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if err := s.guard(name); err != nil {
		return nil, err
	}
	return s.base.OpenFile(name, flag, perm)
}

func (s *ScopedFs) Remove(name string) error {
	return s.base.Remove(name)
}

func (s *ScopedFs) RemoveAll(path string) error {
	return s.base.RemoveAll(path)
}

func (s *ScopedFs) Rename(oldname, newname string) error {
	if err := s.guard(oldname); err != nil {
		return err
	}
	if err := s.guard(newname); err != nil {
		return err
	}
	return s.base.Rename(oldname, newname)
}

func (s *ScopedFs) Stat(name string) (os.FileInfo, error) {
	if err := s.guard(name); err != nil {
		return nil, err
	}
	return s.base.Stat(name)
}

func (s *ScopedFs) Name() string { return "ScopedFs" }

func (s *ScopedFs) Chmod(name string, mode os.FileMode) error {
	if err := s.guard(name); err != nil {
		return err
	}
	return s.base.Chmod(name, mode)
}

func (s *ScopedFs) Chown(name string, uid, gid int) error {
	if err := s.guard(name); err != nil {
		return err
	}
	return s.base.Chown(name, uid, gid)
}

func (s *ScopedFs) Chtimes(name string, atime, mtime time.Time) error {
	if err := s.guard(name); err != nil {
		return err
	}
	return s.base.Chtimes(name, atime, mtime)
}

func (s *ScopedFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	if err := s.guard(name); err != nil {
		return nil, false, err
	}
	return s.base.LstatIfPossible(name)
}
	"path/filepath"
	"testing"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/spf13/afero"
)

		t.Skipf("cannot create symlink: %v", err)
	}

	afs := files.NewScopedFs(afero.NewOsFs(), scope)

	err := Copy(afs, "/srcdir", "/dstdir", 0o644, 0o755)
	if err == nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	afs := files.NewScopedFs(afero.NewOsFs(), scope)

	if err := Copy(afs, "/srcdir", "/dstdir", 0o644, 0o755); err != nil {
		t.Fatalf("expected copy of an in-scope symlink to succeed, got: %v", err)
import (
	"errors"
	"io/fs"

	"github.com/spf13/afero"
)

// CopyDir copies a directory from source to dest and all
// of its sub-directories. It doesn't stop if it finds an error
// during the copy. Returns an error if any.
func CopyDir(afs afero.Fs, source, dest string, fileMode, dirMode fs.FileMode) error {
	// Get properties of source.
	srcinfo, err := afs.Stat(source)
	if err != nil {
	"path/filepath"

	"github.com/spf13/afero"
)

// MoveFile moves file from src to dst.
// CopyFile copies a file from source to dest and returns
// an error if any.
func CopyFile(afs afero.Fs, source, dest string, fileMode, dirMode fs.FileMode) error {
	// Open the source file.
	src, err := afs.Open(source)
	if err != nil {
	"path/filepath"
	"strings"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/share"
	"golang.org/x/crypto/bcrypt"
)

var withHashFile = func(fn handleFunc) handleFunc {
			filePath = ifPath
		}

		// set fs root to the shared file/folder. ScopedFs (not a bare
		// BasePathFs) so the share is also symlink-confined: a link inside the
		// shared subtree that points elsewhere in the owner's scope — outside
		// the share — must not be followed.
		d.user.Fs = files.NewScopedFs(d.user.Fs, basePath)

		// the filesystem is now rebased onto basePath, so paths handed to the
		// rule checker are relative to it. Resolve them back to the user's
	"testing"

	"github.com/asdine/storm/v3"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/share"
	"github.com/filebrowser/filebrowser/v2/storage"
	"github.com/filebrowser/filebrowser/v2/storage/bolt"
	"github.com/filebrowser/filebrowser/v2/users"
	"github.com/spf13/afero"
)

// symlinkShareStorage builds a storage whose single user is rooted at a real
	}
	st.Users = &customFSUser{
		Store: st.Users,
		fs:    files.NewScopedFs(afero.NewOsFs(), scope),
	}
	return st
}
	"testing"

	"github.com/asdine/storm/v3"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/rules"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/share"
	"github.com/filebrowser/filebrowser/v2/storage/bolt"
	"github.com/filebrowser/filebrowser/v2/users"
	"github.com/spf13/afero"
)

func TestPublicShareHandlerAuthentication(t *testing.T) {
				t.Fatalf("failed to save settings: %v", err)
			}

			fs := files.NewScopedFs(afero.NewOsFs(), t.TempDir())
			if err := fs.MkdirAll("/projects/private", 0o755); err != nil {
				t.Fatalf("failed to create private dir: %v", err)
			}
	if err != nil {
		return nil, err
	}
	// Mirror production (users.User init), where a user's filesystem is always a
	// scoped, symlink-confining ScopedFs rather than a bare afero.Fs.
	user.Fs = files.NewScopedFs(cu.fs, "/")

	return user, nil
}
	"path/filepath"
	"strings"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/fileutils"
	"github.com/filebrowser/filebrowser/v2/users"
	"github.com/mholt/archives"
)

func slashClean(name string) string {
		return nil, nil
	}

	info, err := d.user.Fs.Stat(path)
	if err != nil {
		return nil, err
	"strings"
	"time"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/fileutils"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/spf13/afero"
)

var resourceGetHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
			return http.StatusForbidden, nil
		}

		err = checkParent(src, dst)
		if err != nil {
			return http.StatusBadRequest, err
}

func writeFile(afs afero.Fs, dst string, in io.Reader, fileMode, dirMode fs.FileMode) (os.FileInfo, error) {
	dir, _ := path.Split(dst)
	err := afs.MkdirAll(dir, dirMode)
	if err != nil {
	"strings"
	"time"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/share"
	"golang.org/x/crypto/bcrypt"
)

func withPermShare(fn handleFunc) handleFunc {
	// Only allow sharing paths that currently exist. Otherwise a share could be
	// created for a non-existent path and would silently start exposing
	// whatever file later appears there.
	//
	// d.user.Fs is scoped, so Stat also refuses to follow a symlink whose target
	// escapes the user's scope: that returns a permission error here and so
	// blocks creating a share that points out of scope.
	if _, err := d.user.Fs.Stat(r.URL.Path); err != nil {
		return errToStatus(err), err
	}

	var s *share.Link
	var body share.CreateBody
	if r.Body != nil {
	"strings"
	"time"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/spf13/afero"
)

// keepUploadActive periodically touches the cache entry to prevent eviction during transfer
			return http.StatusForbidden, nil
		}

		file, err := files.NewFileInfo(&files.FileOptions{
			Fs:         d.user.Fs,
			Path:       r.URL.Path,
			return http.StatusUnsupportedMediaType, nil
		}

		uploadOffset, err := getUploadOffset(r)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("invalid upload offset")
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/storage/bolt"
	"github.com/filebrowser/filebrowser/v2/users"
	}
	st.Users = &customFSUser{
		Store: st.Users,
		fs:    files.NewScopedFs(afero.NewOsFs(), userScope),
	}

	// Forge a valid auth token for user ID 1.
import (
	"path/filepath"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/rules"
	"github.com/spf13/afero"
)

// ViewMode describes a view mode.

// User describes a user.
type User struct {
	ID                    uint            `storm:"id,increment" json:"id"`
	Username              string          `storm:"unique" json:"username"`
	Password              string          `json:"password"`
	Scope                 string          `json:"scope"`
	Locale                string          `json:"locale"`
	LockPassword          bool            `json:"lockPassword"`
	ViewMode              ViewMode        `json:"viewMode"`
	SingleClick           bool            `json:"singleClick"`
	RedirectAfterCopyMove bool            `json:"redirectAfterCopyMove"`
	Perm                  Permissions     `json:"perm"`
	Commands              []string        `json:"commands"`
	Sorting               files.Sorting   `json:"sorting"`
	Fs                    *files.ScopedFs `json:"-" yaml:"-"`
	Rules                 []rules.Rule    `json:"rules"`
	HideDotfiles          bool            `json:"hideDotfiles"`
	DateFormat            bool            `json:"dateFormat"`
	AceEditorTheme        string          `json:"aceEditorTheme"`
}

// GetRules implements rules.Provider.
	if u.Fs == nil {
		scope := u.Scope
		scope = filepath.Join(baseScope, filepath.Join("/", scope))
		u.Fs = files.NewScopedFs(afero.NewOsFs(), scope)
	}

	return nil

// FullPath gets the full path for a user's relative path.
func (u *User) FullPath(path string) string {
	return afero.FullBaseFsPath(u.Fs.Base(), path)
}
