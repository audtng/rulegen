package main

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2024 Matthew Penner

//go:build unix

package ufs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/sys/unix"
)

// UnixFS is a filesystem that uses the unix package to make io calls.
//
// This is used for proper sand-boxing and full control over the exact syscalls
// being performed.
type UnixFS struct {
	// basePath is the base path for file operations to take place in.
	basePath string

	// dirfd holds the file descriptor of BasePath and is used to ensure
	// operations are restricted into descendants of BasePath.
	dirfd atomic.Int64

	// useOpenat2 controls whether the `openat2` syscall is used instead of the
	// older `openat` syscall.
	useOpenat2 bool
}

// NewUnixFS creates a new sandboxed unix filesystem. BasePath is used as the
// sandbox path, operations on BasePath itself are not allowed, but any
// operations on its descendants are. Symlinks pointing outside BasePath are
// checked and prevented from enabling an escape in a non-raceable manor.
func NewUnixFS(basePath string, useOpenat2 bool) (*UnixFS, error) {
	basePath = strings.TrimSuffix(basePath, "/")
	// We don't need Openat2, if we are given a basePath that is already unsafe
	// I give up on trying to sandbox it.
	dirfd, err := unix.Openat(AT_EMPTY_PATH, basePath, O_DIRECTORY|O_RDONLY, 0)
	if err != nil {
		return nil, convertErrorType(err)
	}

	fs := &UnixFS{
		basePath:   basePath,
		useOpenat2: useOpenat2,
	}
	fs.dirfd.Store(int64(dirfd))
	return fs, nil
}

// BasePath returns the base path of the UnixFS sandbox, file operations
// pointing outside this path are prohibited and will be blocked by all
// operations implemented by UnixFS.
func (fs *UnixFS) BasePath() string {
	return fs.basePath
}

// Close releases the file descriptor used to sandbox operations within the
// base path of the filesystem.
func (fs *UnixFS) Close() error {
	// Once closed, change dirfd to something invalid to detect when it has been
	// closed.
	defer func() {
		fs.dirfd.Store(-1)
	}()
	return unix.Close(int(fs.dirfd.Load()))
}

// Chmod changes the mode of the named file to mode.
//
// If the file is a symbolic link, it changes the mode of the link's target.
// If there is an error, it will be of type *PathError.
//
// A different subset of the mode bits are used, depending on the
// operating system.
//
// On Unix, the mode's permission bits, ModeSetuid, ModeSetgid, and
// ModeSticky are used.
//
// On Windows, only the 0200 bit (owner writable) of mode is used; it
// controls whether the file's read-only attribute is set or cleared.
// The other bits are currently unused. For compatibility with Go 1.12
// and earlier, use a non-zero mode. Use mode 0400 for a read-only
// file and 0600 for a readable+writable file.
//
// On Plan 9, the mode's permission bits, ModeAppend, ModeExclusive,
// and ModeTemporary are used.
func (fs *UnixFS) Chmod(name string, mode FileMode) error {
	dirfd, name, closeFd, err := fs.safePath(name)
	defer closeFd()
	if err != nil {
		return err
	}
	return convertErrorType(unix.Fchmodat(dirfd, name, uint32(mode), 0))
}

// Chown changes the numeric uid and gid of the named file.
//
// If the file is a symbolic link, it changes the uid and gid of the link's target.
// A uid or gid of -1 means to not change that value.
// If there is an error, it will be of type *PathError.
//
// On Windows or Plan 9, Chown always returns the syscall.EWINDOWS or
// EPLAN9 error, wrapped in *PathError.
func (fs *UnixFS) Chown(name string, uid, gid int) error {
	return fs.fchown(name, uid, gid, 0)
}

// Lchown changes the numeric uid and gid of the named file.
//
// If the file is a symbolic link, it changes the uid and gid of the link itself.
// If there is an error, it will be of type *PathError.
//
// On Windows, it always returns the syscall.EWINDOWS error, wrapped
// in *PathError.
func (fs *UnixFS) Lchown(name string, uid, gid int) error {
	// With AT_SYMLINK_NOFOLLOW, Fchownat acts like Lchown but allows us to
	// pass a dirfd.
	return fs.fchown(name, uid, gid, AT_SYMLINK_NOFOLLOW)
}

// fchown is a re-usable Fchownat syscall used by Chown and Lchown.
func (fs *UnixFS) fchown(name string, uid, gid, flags int) error {
	dirfd, name, closeFd, err := fs.safePath(name)
	defer closeFd()
	if err != nil {
		return err
	}
	return convertErrorType(unix.Fchownat(dirfd, name, uid, gid, flags))
}

// Chownat is like Chown but allows passing an existing directory file
// descriptor rather than needing to resolve one.
func (fs *UnixFS) Chownat(dirfd int, name string, uid, gid int) error {
	return convertErrorType(unix.Fchownat(dirfd, name, uid, gid, 0))
}

// Lchownat is like Lchown but allows passing an existing directory file
// descriptor rather than needing to resolve one.
func (fs *UnixFS) Lchownat(dirfd int, name string, uid, gid int) error {
	return convertErrorType(unix.Fchownat(dirfd, name, uid, gid, AT_SYMLINK_NOFOLLOW))
}

// Chtimes changes the access and modification times of the named
// file, similar to the Unix utime() or utimes() functions.
//
// The underlying filesystem may truncate or round the values to a
// less precise time unit.
//
// If there is an error, it will be of type *PathError.
func (fs *UnixFS) Chtimes(name string, atime, mtime time.Time) error {
	dirfd, name, closeFd, err := fs.safePath(name)
	defer closeFd()
	if err != nil {
		return err
	}
	return fs.Chtimesat(dirfd, name, atime, mtime)
}

// Chtimesat is like Chtimes but allows passing an existing directory file
// descriptor rather than needing to resolve one.
func (fs *UnixFS) Chtimesat(dirfd int, name string, atime, mtime time.Time) error {
	var utimes [2]unix.Timespec
	set := func(i int, t time.Time) {
		if t.IsZero() {
			utimes[i] = unix.Timespec{Sec: unix.UTIME_OMIT, Nsec: unix.UTIME_OMIT}
		} else {
			utimes[i] = unix.NsecToTimespec(t.UnixNano())
		}
	}
	set(0, atime)
	set(1, mtime)
	// This does support `AT_SYMLINK_NOFOLLOW` as well if needed.
	if err := unix.UtimesNanoAt(dirfd, name, utimes[0:], 0); err != nil {
		return convertErrorType(&PathError{Op: "chtimes", Path: name, Err: err})
	}
	return nil
}

// Create creates or truncates the named file. If the file already exists,
// it is truncated.
//
// If the file does not exist, it is created with mode 0666
// (before umask). If successful, methods on the returned File can
// be used for I/O; the associated file descriptor has mode O_RDWR.
// If there is an error, it will be of type *PathError.
func (fs *UnixFS) Create(name string) (File, error) {
	return fs.OpenFile(name, O_CREATE|O_WRONLY|O_TRUNC, 0o644)
}

// Mkdir creates a new directory with the specified name and permission
// bits (before umask).
//
// If there is an error, it will be of type *PathError.
func (fs *UnixFS) Mkdir(name string, mode FileMode) error {
	dirfd, name, closeFd, err := fs.safePath(name)
	defer closeFd()
	if err != nil {
		return err
	}
	return fs.Mkdirat(dirfd, name, mode)
}

func (fs *UnixFS) Mkdirat(dirfd int, name string, mode FileMode) error {
	return convertErrorType(unix.Mkdirat(dirfd, name, uint32(mode)))
}

// MkdirAll creates a directory named path, along with any necessary
// parents, and returns nil, or else returns an error.
//
// The permission bits perm (before umask) are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
func (fs *UnixFS) MkdirAll(name string, mode FileMode) error {
	// Ensure name is somewhat clean before continuing.
	name, err := fs.unsafePath(name)
	if err != nil {
		return err
	}
	return fs.mkdirAll(name, mode)
}

// Open opens the named file for reading.
//
// If successful, methods on the returned file can be used for reading; the
// associated file descriptor has mode O_RDONLY.
//
// If there is an error, it will be of type *PathError.
func (fs *UnixFS) Open(name string) (File, error) {
	return fs.OpenFile(name, O_RDONLY, 0)
}

// OpenFile is the generalized open call; most users will use Open
// or Create instead. It opens the named file with specified flag
// (O_RDONLY etc.).
//
// If the file does not exist, and the O_CREATE flag
// is passed, it is created with mode perm (before umask). If successful,
// methods on the returned File can be used for I/O.
//
// If there is an error, it will be of type *PathError.
func (fs *UnixFS) OpenFile(name string, flag int, mode FileMode) (File, error) {
	fd, err := fs.openFile(name, flag, mode)
	if err != nil {
		return nil, err
	}
	// Do not close `fd` here, it is passed to a file that needs the fd, the
	// caller of this function is responsible for calling Close() on the File
	// to release the file descriptor.
	return os.NewFile(uintptr(fd), name), nil
}

func (fs *UnixFS) openFile(name string, flag int, mode FileMode) (int, error) {
	dirfd, name, closeFd, err := fs.safePath(name)
	defer closeFd()
	if err != nil {
		return 0, err
	}
	return fs.openat(dirfd, name, flag, mode)
}

func (fs *UnixFS) OpenFileat(dirfd int, name string, flag int, mode FileMode) (File, error) {
	fd, err := fs.openat(dirfd, name, flag, mode)
	if err != nil {
		return nil, err
	}
	// Do not close `fd` here, it is passed to a file that needs the fd, the
	// caller of this function is responsible for calling Close() on the File
	// to release the file descriptor.
	return os.NewFile(uintptr(fd), name), nil
}

// ReadDir reads the named directory,
//
// returning all its directory entries sorted by filename.
// If an error occurs reading the directory, ReadDir returns the entries it
// was able to read before the error, along with the error.
func (fs *UnixFS) ReadDir(path string) ([]DirEntry, error) {
	dirfd, name, closeFd, err := fs.safePath(path)
	defer closeFd()
	if err != nil {
		return nil, err
	}
	fd, err := fs.openat(dirfd, name, O_DIRECTORY|O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(fd)
	return fs.readDir(fd, name, nil)
}

// RemoveStat is a combination of Stat and Remove, it is used to more
// efficiently remove a file when the caller needs to stat it before
// removing it.
//
// This optimized function exists for our QuotaFS abstraction, which needs
// to track writes to a filesystem. When removing a file, the QuotaFS needs
// to know if the entry is a file and if so, how large it is. Because we
// need to Stat a file in order to get its mode and size, we will already
// know if the entry needs to be removed by using Unlink or Rmdir. The
// standard `Remove` method just tries both Unlink and Rmdir (in that order)
// as it ends up usually being faster and more efficient than calling Stat +
// the proper operation in the first place.
func (fs *UnixFS) RemoveStat(name string) (FileInfo, error) {
	dirfd, name, closeFd, err := fs.safePath(name)
	defer closeFd()
	if err != nil {
		return nil, err
	}

	// Lstat name, we use Lstat as Unlink doesn't care about symlinks.
	s, err := fs.Lstatat(dirfd, name)
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		err = fs.unlinkat(dirfd, name, AT_REMOVEDIR) // Rmdir
	} else {
		err = fs.unlinkat(dirfd, name, 0)
	}
	if err != nil {
		return s, convertErrorType(&PathError{Op: "remove", Path: name, Err: err})
	}
	return s, nil
}

// Remove removes the named file or (empty) directory.
//
// If there is an error, it will be of type *PathError.
func (fs *UnixFS) Remove(name string) error {
	dirfd, name, closeFd, err := fs.safePath(name)
	defer closeFd()
	if err != nil {
		return err
	}

	// Prevent trying to Remove the base directory.
	if name == "." {
		return &PathError{
			Op:   "remove",
			Path: name,
			Err:  ErrBadPathResolution,
		}
	}

	// System call interface forces us to know
	// whether name is a file or directory.
	// Try both: it is cheaper on average than
	// doing a Stat plus the right one.
	err = fs.unlinkat(dirfd, name, 0)
	if err == nil {
		return nil
	}
	err1 := fs.unlinkat(dirfd, name, AT_REMOVEDIR) // Rmdir
	if err1 == nil {
		return nil
	}

	// Both failed: figure out which error to return.
	// OS X and Linux differ on whether unlink(dir)
	// returns EISDIR, so can't use that. However,
	// both agree that rmdir(file) returns ENOTDIR,
	// so we can use that to decide which error is real.
	// Rmdir might also return ENOTDIR if given a bad
	// file path, like /etc/passwd/foo, but in that case,
	// both errors will be ENOTDIR, so it's okay to
	// use the error from unlink.
	if err1 != unix.ENOTDIR {
		err = err1
	}
	return convertErrorType(&PathError{Op: "remove", Path: name, Err: err})
}

// RemoveAll removes path and any children it contains.
//
// It removes everything it can but returns the first error
// it encounters. If the path does not exist, RemoveAll
// returns nil (no error).
//
// If there is an error, it will be of type *PathError.
func (fs *UnixFS) RemoveAll(name string) error {
	name, err := fs.unsafePath(name)
	if err != nil {
		return err
	}
	// While removeAll internally checks this, I want to make sure we check it
	// and return the proper error so our tests can ensure that this will never
	// be a possibility.
	if name == "." {
		return &PathError{
			Op:   "removeall",
			Path: name,
			Err:  ErrBadPathResolution,
		}
	}
	return fs.removeAll(name)
}

func (fs *UnixFS) unlinkat(dirfd int, name string, flags int) error {
	return ignoringEINTR(func() error {
		return unix.Unlinkat(dirfd, name, flags)
	})
}

// Rename renames (moves) oldpath to newpath.
//
// If newpath already exists and is not a directory, Rename replaces it.
// OS-specific restrictions may apply when oldpath and newpath are in different directories.
// Even within the same directory, on non-Unix platforms Rename is not an atomic operation.
//
// If there is an error, it will be of type *LinkError.
func (fs *UnixFS) Rename(oldpath, newpath string) error {
	// Simple case: both paths are the same.
	if oldpath == newpath {
		return nil
	}

	olddirfd, oldname, closeFd, err := fs.safePath(oldpath)
	defer closeFd()
	if err != nil {
		return err
	}
	// Ensure that we are not trying to rename the base directory itself.
	// While unix.Renameat ends up throwing a "device or resource busy" error,
	// that doesn't mean we are protecting the system properly.
	if oldname == "." {
		return convertErrorType(&PathError{
			Op:   "rename",
			Path: oldname,
			Err:  ErrBadPathResolution,
		})
	}
	// Stat the old target to return proper errors.
	if _, err := fs.Lstatat(olddirfd, oldname); err != nil {
		return err
	}

	newdirfd, newname, closeFd2, err := fs.safePath(newpath)
	if err != nil {
		closeFd2()
		if !errors.Is(err, ErrNotExist) {
			return convertErrorType(err)
		}
		var pathErr *PathError
		if !errors.As(err, &pathErr) {
			return convertErrorType(err)
		}
		if err := fs.MkdirAll(pathErr.Path, 0o755); err != nil {
			return err
		}
		newdirfd, newname, closeFd2, err = fs.safePath(newpath)
		defer closeFd2()
		if err != nil {
			return err
		}
	} else {
		defer closeFd2()
	}

	// Ensure that we are not trying to rename the base directory itself.
	// While unix.Renameat ends up throwing a "device or resource busy" error,
	// that doesn't mean we are protecting the system properly.
	if newname == "." {
		return convertErrorType(&PathError{
			Op:   "rename",
			Path: newname,
			Err:  ErrBadPathResolution,
		})
	}
	// Stat the new target to return proper errors.
	_, err = fs.Lstatat(newdirfd, newname)
	switch {
	case err == nil:
		return convertErrorType(&PathError{
			Op:   "rename",
			Path: newname,
			Err:  ErrExist,
		})
	case !errors.Is(err, ErrNotExist):
		return err
	}
	return unix.Renameat(olddirfd, oldname, newdirfd, newname)
}

// Stat returns a FileInfo describing the named file.
//
// If there is an error, it will be of type *PathError.
func (fs *UnixFS) Stat(name string) (FileInfo, error) {
	return fs.fstat(name, 0)
}

// Statat is like Stat but allows passing an existing directory file
// descriptor rather than needing to resolve one.
func (fs *UnixFS) Statat(dirfd int, name string) (FileInfo, error) {
	return fs.fstatat(dirfd, name, 0)
}

// Lstat returns a FileInfo describing the named file.
//
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link. Lstat makes no attempt to follow the link.
//
// If there is an error, it will be of type *PathError.
func (fs *UnixFS) Lstat(name string) (FileInfo, error) {
	return fs.fstat(name, AT_SYMLINK_NOFOLLOW)
}

// Lstatat is like Lstat but allows passing an existing directory file
// descriptor rather than needing to resolve one.
func (fs *UnixFS) Lstatat(dirfd int, name string) (FileInfo, error) {
	return fs.fstatat(dirfd, name, AT_SYMLINK_NOFOLLOW)
}

func (fs *UnixFS) fstat(name string, flags int) (FileInfo, error) {
	dirfd, name, closeFd, err := fs.safePath(name)
	defer closeFd()
	if err != nil {
		return nil, err
	}
	return fs.fstatat(dirfd, name, flags)
}

func (fs *UnixFS) fstatat(dirfd int, name string, flags int) (FileInfo, error) {
	var s fileStat
	if err := ignoringEINTR(func() error {
		return unix.Fstatat(dirfd, name, &s.sys, flags)
	}); err != nil {
		return nil, &PathError{Op: "stat", Path: name, Err: err}
	}
	fillFileStatFromSys(&s, name)
	return &s, nil
}

// Symlink creates newname as a symbolic link to oldname.
//
// On Windows, a symlink to a non-existent oldname creates a file symlink;
// if oldname is later created as a directory the symlink will not work.
//
// If there is an error, it will be of type *LinkError.
func (fs *UnixFS) Symlink(oldpath, newpath string) error {
	dirfd, newpath, closeFd, err := fs.safePath(newpath)
	defer closeFd()
	if err != nil {
		return err
	}
	if err := ignoringEINTR(func() error {
		// We aren't concerned with oldpath here as a symlink can point anywhere
		// it wants.
		return unix.Symlinkat(oldpath, dirfd, newpath)
	}); err != nil {
		return &LinkError{Op: "symlink", Old: oldpath, New: newpath, Err: err}
	}
	return nil
}

// Touch will attempt to open a file for reading and/or writing. If the file
// does not exist it will be created, and any missing parent directories will
// also be created. The opened file may be truncated, only if `flag` has
// O_TRUNC set.
func (fs *UnixFS) Touch(path string, flag int, mode FileMode) (File, error) {
	if flag&O_CREATE == 0 {
		flag |= O_CREATE
	}
	dirfd, name, closeFd, err := fs.safePath(path)
	defer closeFd()
	if err == nil {
		return fs.OpenFileat(dirfd, name, flag, mode)
	}
	if !errors.Is(err, ErrNotExist) {
		return nil, err
	}
	var pathErr *PathError
	if !errors.As(err, &pathErr) {
		return nil, err
	}
	if err := fs.MkdirAll(pathErr.Path, 0o755); err != nil {
		return nil, err
	}
	// Try to open the file one more time after creating its parent directories.
	return fs.OpenFile(path, flag, mode)
}

// WalkDir walks the file tree rooted at root, calling fn for each file or
// directory in the tree, including root.
//
// All errors that arise visiting files and directories are filtered by fn:
// see the [WalkDirFunc] documentation for details.
//
// The files are walked in lexical order, which makes the output deterministic
// but requires WalkDir to read an entire directory into memory before proceeding
// to walk that directory.
//
// WalkDir does not follow symbolic links found in directories,
// but if root itself is a symbolic link, its target will be walked.
func (fs *UnixFS) WalkDir(root string, fn WalkDirFunc) error {
	return WalkDir(fs, root, fn)
}

// openat is a wrapper around both unix.Openat and unix.Openat2. If the UnixFS
// was configured to enable openat2 support, unix.Openat2 will be used instead
// of unix.Openat due to having better security properties for our use-case.
func (fs *UnixFS) openat(dirfd int, name string, flag int, mode FileMode) (int, error) {
	if flag&O_NOFOLLOW == 0 {
		flag |= O_NOFOLLOW
	}

	var fd int
	for {
		var err error
		if fs.useOpenat2 {
			fd, err = fs._openat2(dirfd, name, uint64(flag), uint64(syscallMode(mode)))
		} else {
			fd, err = fs._openat(dirfd, name, flag, uint32(syscallMode(mode)))
		}
		if err == nil {
			break
		}
		// We have to check EINTR here, per issues https://go.dev/issue/11180 and https://go.dev/issue/39237.
		if err == unix.EINTR {
			continue
		}
		return 0, convertErrorType(err)
	}

	// If we are not using openat2, do additional path checking. This assumes
	// that openat2 is using `RESOLVE_BENEATH` to avoid the same security
	// issue.
	if !fs.useOpenat2 {
		var finalPath string
		finalPath, err := filepath.EvalSymlinks(filepath.Join("/proc/self/fd/", strconv.Itoa(dirfd)))
		if err != nil {
			return fd, convertErrorType(err)
		}
		if err != nil {
			if !errors.Is(err, ErrNotExist) {
				return fd, fmt.Errorf("failed to evaluate symlink: %w", convertErrorType(err))
			}

			// The target of one of the symlinks (EvalSymlinks is recursive)
			// does not exist. So get the path that does not exist and use
			// that for further validation instead.
			var pErr *PathError
			if ok := errors.As(err, &pErr); !ok {
				return fd, fmt.Errorf("failed to evaluate symlink: %w", convertErrorType(err))
			}
			finalPath = pErr.Path
		}

		// Check if the path is within our root.
		if !fs.unsafeIsPathInsideOfBase(finalPath) {
			return fd, convertErrorType(&PathError{
				Op:   "openat",
				Path: name,
				Err:  ErrBadPathResolution,
			})
		}
	}
	return fd, nil
}

// _openat is a wrapper around unix.Openat. This method should never be directly
// called, use `openat` instead.
func (fs *UnixFS) _openat(dirfd int, name string, flag int, mode uint32) (int, error) {
	// Ensure the O_CLOEXEC flag is set.
	// Go sets this in the os package, but since we are directly using unix
	// we need to set it ourselves.
	if flag&O_CLOEXEC == 0 {
		flag |= O_CLOEXEC
	}
	// O_LARGEFILE is set by Openat for us automatically.
	fd, err := unix.Openat(dirfd, name, flag, mode)
	switch {
	case err == nil:
		return fd, nil
	case err == unix.EINTR:
		return 0, err
	case err == unix.EAGAIN:
		return 0, err
	default:
		return 0, &PathError{Op: "openat", Path: name, Err: err}
	}
}

// _openat2 is a wonderful syscall that supersedes the `openat` syscall. It has
// improved validation and security characteristics that weren't available or
// considered when `openat` was originally implemented. As such, it is only
// present in Kernel 5.6 and above.
//
// This method should never be directly called, use `openat` instead.
func (fs *UnixFS) _openat2(dirfd int, name string, flag uint64, mode uint64) (int, error) {
	// Ensure the O_CLOEXEC flag is set.
	// Go sets this when using the os package, but since we are directly using
	// the unix package we need to set it ourselves.
	if flag&O_CLOEXEC == 0 {
		flag |= O_CLOEXEC
	}
	// Ensure the O_LARGEFILE flag is set.
	// Go sets this for unix.Open, unix.Openat, but not unix.Openat2.
	if flag&O_LARGEFILE == 0 {
		flag |= O_LARGEFILE
	}
	fd, err := unix.Openat2(dirfd, name, &unix.OpenHow{
		Flags: flag,
		Mode:  mode,
		// This is the bread and butter of preventing a symlink escape, without
		// this option, we have to handle path validation fully on our own.
		//
		// This is why using Openat2 over Openat is preferred if available.
		Resolve: unix.RESOLVE_BENEATH,
	})
	switch {
	case err == nil:
		return fd, nil
	case err == unix.EINTR:
		return 0, err
	case err == unix.EAGAIN:
		return 0, err
	default:
		return 0, &PathError{Op: "openat2", Path: name, Err: err}
	}
}

func (fs *UnixFS) SafePath(path string) (int, string, func(), error) {
	return fs.safePath(path)
}

func (fs *UnixFS) safePath(path string) (dirfd int, file string, closeFd func(), err error) {
	// Default closeFd to a NO-OP.
	closeFd = func() {}

	// Use unsafePath to clean the path and strip BasePath if path is absolute.
	var name string
	name, err = fs.unsafePath(path)
	if err != nil {
		return
	}

	// Check if dirfd was closed, this will happen if (*UnixFS).Close()
	// was called.
	fsDirfd := int(fs.dirfd.Load())
	if fsDirfd == -1 {
		err = ErrClosed
		return
	}

	// Split the parent from the last element in the path, this gives us the
	// "file name" and the full path to its parent.
	var dir string
	dir, file = filepath.Split(name)
	// If dir is empty then name is not nested.
	if dir == "" {
		// We don't need to set closeFd here as it will default to a NO-OP and
		// `fs.dirfd` is re-used until the filesystem is no-longer needed.
		dirfd = fsDirfd

		// Return dirfd, name, an empty closeFd func, and no error
		return
	}

	// Dir will usually contain a trailing slash as filepath.Split doesn't
	// trim slashes.
	dir = strings.TrimSuffix(dir, "/")
	dirfd, err = fs.openat(fsDirfd, dir, O_DIRECTORY|O_RDONLY, 0)
	if dirfd != 0 {
		// Set closeFd to close the newly opened directory file descriptor.
		closeFd = func() { _ = unix.Close(dirfd) }
	}

	// Return dirfd, name, the closeFd func, and err
	return
}

// unsafePath prefixes the given path and prefixes it with the filesystem's
// base path, cleaning the result. The path returned by this function may not
// be inside the filesystem's base path, additional checks are required to
// safely use paths returned by this function.
func (fs *UnixFS) unsafePath(path string) (string, error) {
	// Calling filepath.Clean on the joined directory will resolve it to the
	// absolute path, removing any ../ type of resolution arguments, and leaving
	// us with a direct path link.
	//
	// This will also trim the existing root path off the beginning of the path
	// passed to the function since that can get a bit messy.
	r := filepath.Clean(filepath.Join(fs.basePath, strings.TrimPrefix(path, fs.basePath)))

	if fs.unsafeIsPathInsideOfBase(r) {
		// This is kinda ironic isn't it.
		// We do this as we are operating with dirfds and `*at` syscalls which
		// behave differently if given an absolute path.
		//
		// First trim the BasePath, then trim any leading slashes.
		r = strings.TrimPrefix(strings.TrimPrefix(r, fs.basePath), "/")
		// If the path is empty then return "." as the path is pointing to the
		// root.
		if r == "" {
			return ".", nil
		}
		return r, nil
	}

	return "", &PathError{
		Op:   "safePath",
		Path: path,
		Err:  ErrBadPathResolution,
	}
}

// unsafeIsPathInsideOfBase checks if the given path is inside the filesystem's
// base path.
func (fs *UnixFS) unsafeIsPathInsideOfBase(path string) bool {
	return strings.HasPrefix(
		strings.TrimSuffix(path, "/")+"/",
		fs.basePath+"/",
	)
}
package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"emperror.dev/errors"
	"github.com/apex/log"
	"github.com/gabriel-vasile/mimetype"
	ignore "github.com/sabhiram/go-gitignore"

	"github.com/pterodactyl/wings/config"
	"github.com/pterodactyl/wings/internal/ufs"
)

// TODO: detect and enable
var useOpenat2 = true

type Filesystem struct {
	unixFS *ufs.Quota

	mu                sync.RWMutex
	lastLookupTime    *usageLookupTime
	lookupInProgress  atomic.Bool
	diskCheckInterval time.Duration
	denylist          *ignore.GitIgnore

	isTest bool
}

// New creates a new Filesystem instance for a given server.
func New(root string, size int64, denylist []string) (*Filesystem, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	unixFS, err := ufs.NewUnixFS(root, config.UseOpenat2())
	if err != nil {
		return nil, err
	}
	quota := ufs.NewQuota(unixFS, size)

	return &Filesystem{
		unixFS: quota,

		diskCheckInterval: time.Duration(config.Get().System.DiskCheckInterval),
		lastLookupTime:    &usageLookupTime{},
		denylist:          ignore.CompileIgnoreLines(denylist...),
	}, nil
}

// Path returns the root path for the Filesystem instance.
func (fs *Filesystem) Path() string {
	return fs.unixFS.BasePath()
}

// ReadDir reads directory entries.
func (fs *Filesystem) ReadDir(path string) ([]ufs.DirEntry, error) {
	return fs.unixFS.ReadDir(path)
}

// ReadDirStat is like ReadDir except that it returns FileInfo for each entry
// instead of just a DirEntry.
func (fs *Filesystem) ReadDirStat(path string) ([]ufs.FileInfo, error) {
	return ufs.ReadDirMap(fs.unixFS.UnixFS, path, func(e ufs.DirEntry) (ufs.FileInfo, error) {
		return e.Info()
	})
}

// File returns a reader for a file instance as well as the stat information.
func (fs *Filesystem) File(p string) (ufs.File, Stat, error) {
	f, err := fs.unixFS.Open(p)
	if err != nil {
		return nil, Stat{}, err
	}
	st, err := statFromFile(f)
	if err != nil {
		_ = f.Close()
		return nil, Stat{}, err
	}
	return f, st, nil
}

func (fs *Filesystem) UnixFS() *ufs.UnixFS {
	return fs.unixFS.UnixFS
}

// Touch acts by creating the given file and path on the disk if it is not present
// already. If  it is present, the file is opened using the defaults which will truncate
// the contents. The opened file is then returned to the caller.
func (fs *Filesystem) Touch(p string, flag int) (ufs.File, error) {
	return fs.unixFS.Touch(p, flag, 0o644)
}

// Writefile writes a file to the system. If the file does not already exist one
// will be created. This will also properly recalculate the disk space used by
// the server when writing new files or modifying existing ones.
//
// DEPRECATED: use `Write` instead.
func (fs *Filesystem) Writefile(p string, r io.Reader) error {
	var currentSize int64
	st, err := fs.unixFS.Stat(p)
	if err != nil && !errors.Is(err, ufs.ErrNotExist) {
		return errors.Wrap(err, "server/filesystem: writefile: failed to stat file")
	} else if err == nil {
		if st.IsDir() {
			// TODO: resolved
			return errors.WithStack(&Error{code: ErrCodeIsDirectory, resolved: ""})
		}
		currentSize = st.Size()
	}

	// Touch the file and return the handle to it at this point. This will
	// create or truncate the file, and create any necessary parent directories
	// if they are missing.
	file, err := fs.unixFS.Touch(p, ufs.O_RDWR|ufs.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("error touching file: %w", err)
	}
	defer file.Close()

	// Do not use CopyBuffer here, it is wasteful as the file implements
	// io.ReaderFrom, which causes it to not use the buffer anyways.
	n, err := io.Copy(file, r)

	// Adjust the disk usage to account for the old size and the new size of the file.
	fs.unixFS.Add(n - currentSize)

	if err := fs.chownFile(p); err != nil {
		return fmt.Errorf("error chowning file: %w", err)
	}
	// Return the error from io.Copy.
	return err
}

func (fs *Filesystem) Write(p string, r io.Reader, newSize int64, mode ufs.FileMode) error {
	var currentSize int64
	st, err := fs.unixFS.Stat(p)
	if err != nil && !errors.Is(err, ufs.ErrNotExist) {
		return errors.Wrap(err, "server/filesystem: writefile: failed to stat file")
	} else if err == nil {
		if st.IsDir() {
			// TODO: resolved
			return errors.WithStack(&Error{code: ErrCodeIsDirectory, resolved: ""})
		}
		currentSize = st.Size()
	}

	// Check that the new size we're writing to the disk can fit. If there is currently
	// a file we'll subtract that current file size from the size of the buffer to determine
	// the amount of new data we're writing (or amount we're removing if smaller).
	if err := fs.HasSpaceFor(newSize - currentSize); err != nil {
		return err
	}

	// Touch the file and return the handle to it at this point. This will
	// create or truncate the file, and create any necessary parent directories
	// if they are missing.
	file, err := fs.unixFS.Touch(p, ufs.O_RDWR|ufs.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer file.Close()

	// Do not use CopyBuffer here, it is wasteful as the file implements
	// io.ReaderFrom, which causes it to not use the buffer anyways.
	n, err := io.Copy(file, io.LimitReader(r, newSize))

	// Adjust the disk usage to account for the old size and the new size of the file.
	fs.unixFS.Add(n - currentSize)

	if err := fs.chownFile(p); err != nil {
		return err
	}
	// Return the error from io.Copy.
	return err
}

// CreateDirectory creates a new directory (name) at a specified path (p) for
// the server.
func (fs *Filesystem) CreateDirectory(name string, p string) error {
	return fs.unixFS.MkdirAll(filepath.Join(p, name), 0o755)
}

func (fs *Filesystem) Rename(oldpath, newpath string) error {
	return fs.unixFS.Rename(oldpath, newpath)
}

func (fs *Filesystem) Symlink(oldpath, newpath string) error {
	return fs.unixFS.Symlink(oldpath, newpath)
}

func (fs *Filesystem) chownFile(name string) error {
	if fs.isTest {
		return nil
	}

	uid := config.Get().System.User.Uid
	gid := config.Get().System.User.Gid
	return fs.unixFS.Lchown(name, uid, gid)
}

// Chown recursively iterates over a file or directory and sets the permissions on all of the
// underlying files. Iterate over all of the files and directories. If it is a file just
// go ahead and perform the chown operation. Otherwise dig deeper into the directory until
// we've run out of directories to dig into.
func (fs *Filesystem) Chown(p string) error {
	if fs.isTest {
		return nil
	}
	uid := config.Get().System.User.Uid
	gid := config.Get().System.User.Gid

	dirfd, name, closeFd, err := fs.unixFS.SafePath(p)
	defer closeFd()
	if err != nil {
		return err
	}

	// Start by just chowning the initial path that we received.
	if err := fs.unixFS.Lchownat(dirfd, name, uid, gid); err != nil {
		return errors.Wrap(err, "server/filesystem: chown: failed to chown path")
	}

	// If this is not a directory we can now return from the function, there is nothing
	// left that we need to do.
	if st, err := fs.unixFS.Lstatat(dirfd, name); err != nil || !st.IsDir() {
		return nil
	}

	// This walker is probably some of the most efficient code in Wings. It has
	// an internally re-used buffer for listing directory entries and doesn't
	// need to check if every individual path it touches is safe as the code
	// doesn't traverse symlinks, is immune to symlink timing attacks, and
	// gives us a dirfd and file name to make a direct syscall with.
	if err := fs.unixFS.WalkDirat(dirfd, name, func(dirfd int, name, _ string, info ufs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err := fs.unixFS.Lchownat(dirfd, name, uid, gid); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("server/filesystem: chown: failed to chown during walk function: %w", err)
	}
	return nil
}

func (fs *Filesystem) Chmod(path string, mode ufs.FileMode) error {
	return fs.unixFS.Chmod(path, mode)
}

// Begin looping up to 50 times to try and create a unique copy file name. This will take
// an input of "file.txt" and generate "file copy.txt". If that name is already taken, it will
// then try to write "file copy 2.txt" and so on, until reaching 50 loops. At that point we
// Could probably make this more efficient by checking if there are any files matching the copy
// pattern, and trying to find the highest number and then incrementing it by one rather than
// looping endlessly.
func (fs *Filesystem) findCopySuffix(dirfd int, name, extension string) (string, error) {
	var i int
	suffix := " copy"

		n := name + suffix + extension
		// If we stat the file and it does not exist that means we're good to create the copy. If it
		// does exist, we'll just continue to the next loop and try again.
		if _, err := fs.unixFS.Lstatat(dirfd, n); err != nil {
			if !errors.Is(err, ufs.ErrNotExist) {
				return "", err
			}
			break
		}

	return name + suffix + extension, nil
}

// Copy copies a given file to the same location and appends a suffix to the
// file to indicate that it has been copied.
func (fs *Filesystem) Copy(p string) error {
	dirfd, name, closeFd, err := fs.unixFS.SafePath(p)
	defer closeFd()
	if err != nil {
		return err
	}
	source, err := fs.unixFS.OpenFileat(dirfd, name, ufs.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer source.Close()
	info, err := source.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() || !info.Mode().IsRegular() {
		// If this is a directory or not a regular file, just throw a not-exist error
		// since anything calling this function should understand what that means.
		return ufs.ErrNotExist
	}
	currentSize := info.Size()

	// Check that copying this file wouldn't put the server over its limit.
	if err := fs.HasSpaceFor(currentSize); err != nil {
		return err
	}

	base := info.Name()
	extension := filepath.Ext(base)
	baseName := strings.TrimSuffix(base, extension)

	// Ensure that ".tar" is also counted as apart of the file extension.
	// There might be a better way to handle this for other double file extensions,
	// but this is a good workaround for now.
	if strings.HasSuffix(baseName, ".tar") {
		extension = ".tar" + extension
		baseName = strings.TrimSuffix(baseName, ".tar")
	}

	newName, err := fs.findCopySuffix(dirfd, baseName, extension)
	if err != nil {
		return err
	}
	dst, err := fs.unixFS.OpenFileat(dirfd, newName, ufs.O_WRONLY|ufs.O_CREATE, info.Mode())
	if err != nil {
		return err
	}

	// Do not use CopyBuffer here, it is wasteful as the file implements
	// io.ReaderFrom, which causes it to not use the buffer anyways.
	n, err := io.Copy(dst, io.LimitReader(source, currentSize))
	fs.unixFS.Add(n)

	if !fs.isTest {
		if err := fs.unixFS.Lchownat(dirfd, newName, config.Get().System.User.Uid, config.Get().System.User.Gid); err != nil {
			return err
		}
	}
	// Return the error from io.Copy.
	return err
}

// TruncateRootDirectory removes _all_ files and directories from a server's
	if err := os.Mkdir(fs.Path(), 0o755); err != nil {
		return err
	}
	_ = fs.unixFS.Close()
	unixFS, err := ufs.NewUnixFS(fs.Path(), config.UseOpenat2())
	if err != nil {
		return err
	}
	var limit int64
	if fs.isTest {
		limit = 0
	} else {
		limit = fs.unixFS.Limit()
	}
	fs.unixFS = ufs.NewQuota(unixFS, limit)
	return nil
}

// Delete removes a file or folder from the system. Prevents the user from
// accidentally (or maliciously) removing their root server data directory.
func (fs *Filesystem) Delete(p string) error {
	return fs.unixFS.RemoveAll(p)
}

//type fileOpener struct {
//	fs   *Filesystem
//	busy uint
//}
//
//// Attempts to open a given file up to "attempts" number of times, using a backoff. If the file
//// cannot be opened because of a "text file busy" error, we will attempt until the number of attempts
//// has been exhaused, at which point we will abort with an error.
//func (fo *fileOpener) open(path string, flags int, perm ufs.FileMode) (ufs.File, error) {
//	for {
//		f, err := fo.fs.unixFS.OpenFile(path, flags, perm)
//
//		// If there is an error because the text file is busy, go ahead and sleep for a few
//		// hundred milliseconds and then try again up to three times before just returning the
//		// error back to the caller.
//		//
//		// Based on code from: https://github.com/golang/go/issues/22220#issuecomment-336458122
//		if err != nil && fo.busy < 3 && strings.Contains(err.Error(), "text file busy") {
//			time.Sleep(100 * time.Millisecond << fo.busy)
//			fo.busy++
//			continue
//		}
//
//		return f, err
//	}
//}

// ListDirectory lists the contents of a given directory and returns stat
// information about each file and folder within it.
func (fs *Filesystem) ListDirectory(p string) ([]Stat, error) {
	// Read entries from the path on the filesystem, using the mapped reader, so
	// we can map the DirEntry slice into a Stat slice with mimetype information.
	out, err := ufs.ReadDirMap(fs.unixFS.UnixFS, p, func(e ufs.DirEntry) (Stat, error) {
		info, err := e.Info()
		if err != nil {
			return Stat{}, err
		}

		var d string
		if e.Type().IsDir() {
			d = "inode/directory"
		} else {
			d = "application/octet-stream"
		}
		var m *mimetype.MIME
		if e.Type().IsRegular() {
			// TODO: I should probably find a better way to do this.
			eO := e.(interface {
				Open() (ufs.File, error)
			})
			f, err := eO.Open()
			if err != nil {
				return Stat{}, err
			}
			m, err = mimetype.DetectReader(f)
			if err != nil {
				log.Error(err.Error())
			}
			_ = f.Close()
		}

		st := Stat{FileInfo: info, Mimetype: d}
		if m != nil {
			st.Mimetype = m.String()
		}
		return st, nil
	})
	if err != nil {
		return nil, err
	}

	// Sort entries alphabetically.
	slices.SortStableFunc(out, func(a, b Stat) int {
		switch {
		case a.Name() == b.Name():
			return 0
		case a.Name() > b.Name():
			return 1
		default:
			return -1
		}
	})

	// Sort folders before other file types.
	slices.SortStableFunc(out, func(a, b Stat) int {
		switch {
		case a.IsDir() && b.IsDir():
			return 0
		case a.IsDir():
			return 1
		default:
			return -1
		}
	})

	return out, nil
}

func (fs *Filesystem) Chtimes(path string, atime, mtime time.Time) error {
	if fs.isTest {
		return nil
	}
	return fs.unixFS.Chtimes(path, atime, mtime)
}
// SPDX-License-Identifier: BSD-2-Clause

// Some code in this file was derived from https://github.com/karrick/godirwalk.

//go:build unix

package ufs

import (
	"bytes"
	iofs "io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"unsafe"

	"golang.org/x/sys/unix"
)

type WalkDiratFunc func(dirfd int, name, relative string, d DirEntry, err error) error

func (fs *UnixFS) WalkDirat(dirfd int, name string, fn WalkDiratFunc) error {
	if dirfd == 0 {
		// TODO: proper validation, ideally a dedicated function.
		dirfd = int(fs.dirfd.Load())
	}
	info, err := fs.Lstatat(dirfd, name)
	if err != nil {
		err = fn(dirfd, name, name, nil, err)
	} else {
		b := newScratchBuffer()
		err = fs.walkDir(b, dirfd, name, name, iofs.FileInfoToDirEntry(info), fn)
	}
	if err == SkipDir || err == SkipAll {
		return nil
	}
	return err
}

func (fs *UnixFS) walkDir(b []byte, parentfd int, name, relative string, d DirEntry, walkDirFn WalkDiratFunc) error {
	if err := walkDirFn(parentfd, name, relative, d, nil); err != nil || !d.IsDir() {
		if err == SkipDir && d.IsDir() {
			// Successfully skipped directory.
			err = nil
		}
		return err
	}

	dirfd, err := fs.openat(parentfd, name, O_DIRECTORY|O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer unix.Close(dirfd)

	dirs, err := fs.readDir(dirfd, name, b)
	if err != nil {
		// Second call, to report ReadDir error.
		err = walkDirFn(dirfd, name, relative, d, err)
		if err != nil {
			if err == SkipDir && d.IsDir() {
				err = nil
			}
			return err
		}
	}

	for _, d1 := range dirs {
		if err := fs.walkDir(b, dirfd, d1.Name(), path.Join(relative, d1.Name()), d1, walkDirFn); err != nil {
			if err == SkipDir {
				break
			}
			return err
		}
	}
	return nil
}

// ReadDirMap .
// TODO: document
func ReadDirMap[T any](fs *UnixFS, path string, fn func(DirEntry) (T, error)) ([]T, error) {
	dirfd, name, closeFd, err := fs.safePath(path)
	defer closeFd()
	if err != nil {
		return nil, err
	}
	fd, err := fs.openat(dirfd, name, O_DIRECTORY|O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(fd)

	entries, err := fs.readDir(fd, ".", nil)
	if err != nil {
		return nil, err
	}

	out := make([]T, len(entries))
	for i, e := range entries {
		idx := i
		e := e
		v, err := fn(e)
		if err != nil {
			return nil, err
		}
		out[idx] = v
	}
	return out, nil
}

// nameOffset is a compile time constant
const nameOffset = int(unsafe.Offsetof(unix.Dirent{}.Name))

func nameFromDirent(de *unix.Dirent) (name []byte) {
	// Because this GOOS' syscall.Dirent does not provide a field that specifies
	// the name length, this function must first calculate the max possible name
	// length, and then search for the NULL byte.
	ml := int(de.Reclen) - nameOffset

	// Convert syscall.Dirent.Name, which is array of int8, to []byte, by
	// overwriting Cap, Len, and Data slice header fields to the max possible
	// name length computed above, and finding the terminating NULL byte.
	//
	// TODO: is there an alternative to the deprecated SliceHeader?
	// SliceHeader was mainly deprecated due to it being misused for avoiding
	// allocations when converting a byte slice to a string, ref;
	// https://go.dev/issue/53003
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&name))
	sh.Cap = ml
	sh.Len = ml
	sh.Data = uintptr(unsafe.Pointer(&de.Name[0]))

	if index := bytes.IndexByte(name, 0); index >= 0 {
		// Found NULL byte; set slice's cap and len accordingly.
		sh.Cap = index
		sh.Len = index
		return
	}

	// NOTE: This branch is not expected, but included for defensive
	// programming, and provides a hard stop on the name based on the structure
	// field array size.
	sh.Cap = len(de.Name)
	sh.Len = sh.Cap
	return
}

// modeTypeFromDirent converts a syscall defined constant, which is in purview
// of OS, to a constant defined by Go, assumed by this project to be stable.
//
// When the syscall constant is not recognized, this function falls back to a
// Stat on the file system.
func (fs *UnixFS) modeTypeFromDirent(fd int, de *unix.Dirent, osDirname, osBasename string) (FileMode, error) {
	switch de.Type {
	case unix.DT_REG:
		return 0, nil
	case unix.DT_DIR:
		return ModeDir, nil
	case unix.DT_LNK:
		return ModeSymlink, nil
	case unix.DT_CHR:
		return ModeDevice | ModeCharDevice, nil
	case unix.DT_BLK:
		return ModeDevice, nil
	case unix.DT_FIFO:
		return ModeNamedPipe, nil
	case unix.DT_SOCK:
		return ModeSocket, nil
	default:
		// If syscall returned unknown type (e.g., DT_UNKNOWN, DT_WHT), then
		// resolve actual mode by reading file information.
		return fs.modeType(fd, filepath.Join(osDirname, osBasename))
	}
}

// modeType returns the mode type of the file system entry identified by
// osPathname by calling os.LStat function, to intentionally not follow symbolic
// links.
//
// Even though os.LStat provides all file mode bits, we want to ensure same
// values returned to caller regardless of whether we obtained file mode bits
// from syscall or stat call. Therefore, mask out the additional file mode bits
// that are provided by stat but not by the syscall, so users can rely on their
// values.
func (fs *UnixFS) modeType(dirfd int, name string) (os.FileMode, error) {
	fi, err := fs.Lstatat(dirfd, name)
	if err == nil {
		return fi.Mode() & ModeType, nil
	}
	return 0, err
}

var minimumScratchBufferSize = os.Getpagesize()

func newScratchBuffer() []byte {
	return make([]byte, minimumScratchBufferSize)
}

func (fs *UnixFS) readDir(fd int, name string, b []byte) ([]DirEntry, error) {
	scratchBuffer := b
	if scratchBuffer == nil || len(scratchBuffer) < minimumScratchBufferSize {
		scratchBuffer = newScratchBuffer()
	}

	var entries []DirEntry
	var workBuffer []byte

	var sde unix.Dirent
	for {
		if len(workBuffer) == 0 {
			n, err := unix.Getdents(fd, scratchBuffer)
			if err != nil {
				if err == unix.EINTR {
					continue
				}
				return nil, convertErrorType(err)
			}
			if n <= 0 {
				// end of directory: normal exit
				return entries, nil
			}
			workBuffer = scratchBuffer[:n] // trim work buffer to number of bytes read
		}

		// "Go is like C, except that you just put `unsafe` all over the place".
		copy((*[unsafe.Sizeof(unix.Dirent{})]byte)(unsafe.Pointer(&sde))[:], workBuffer)
		workBuffer = workBuffer[sde.Reclen:] // advance buffer for next iteration through loop

		if sde.Ino == 0 {
			continue // inode set to 0 indicates an entry that was marked as deleted
		}

		nameSlice := nameFromDirent(&sde)
		nameLength := len(nameSlice)

		if nameLength == 0 || (nameSlice[0] == '.' && (nameLength == 1 || (nameLength == 2 && nameSlice[1] == '.'))) {
			continue
		}

		childName := string(nameSlice)
		mt, err := fs.modeTypeFromDirent(fd, &sde, name, childName)
		if err != nil {
			return nil, convertErrorType(err)
		}
		entries = append(entries, &dirent{name: childName, path: name, modeType: mt, dirfd: fd, fs: fs})
	}
}

// dirent stores the name and file system mode type of discovered file system
// entries.
type dirent struct {
	name     string
	path     string
	modeType FileMode

	dirfd int
	fs    *UnixFS
}

func (de dirent) Name() string {
	return de.name
}

func (de dirent) IsDir() bool {
	return de.modeType&ModeDir != 0
}

func (de dirent) Type() FileMode {
	return de.modeType
}

func (de dirent) Info() (FileInfo, error) {
	if de.fs == nil {
		return nil, nil
	}
	return de.fs.Lstatat(de.dirfd, de.name)
}

func (de dirent) Open() (File, error) {
	if de.fs == nil {
		return nil, nil
	}
	return de.fs.OpenFileat(de.dirfd, de.name, O_RDONLY, 0)
}

// reset releases memory held by entry err and name, and resets mode type to 0.
func (de *dirent) reset() {
	de.name = ""
	de.path = ""
	de.modeType = 0
}
