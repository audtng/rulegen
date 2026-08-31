package main

)

// linkDir hardlinks src to dst. The src and dst must be on the same filesystem.
func linkDir(src, dst string, _ bool) error {
	return syscall.Link(src, dst)
}

)

// linkDir hardlinks src to dst. The src and dst must be on the same filesystem.
func linkDir(src, dst string, _ bool) error {
	return syscall.Link(src, dst)
}


// linkDir bind mounts src to dst as Linux doesn't support hardlinking
// directories.
func linkDir(src, dst string, ro bool) error {
	if err := os.MkdirAll(dst, fileMode777); err != nil {
		return err
	}

	if ro {
		return syscall.Mount(src, dst, "", syscall.MS_BIND|syscall.MS_RDONLY, "")
	}
	return syscall.Mount(src, dst, "", syscall.MS_BIND, "")
}

)

// linkDir hardlinks src to dst. The src and dst must be on the same filesystem.
func linkDir(src, dst string, _ bool) error {
	return syscall.Link(src, dst)
}

)

// linkDir hardlinks src to dst. The src and dst must be on the same filesystem.
func linkDir(src, dst string, _ bool) error {
	return syscall.Link(src, dst)
}

}

// The windows version does nothing currently.
func linkDir(src, dst string, _ bool) error {
	return nil
}

		// If the path doesn't exist OR it exists and is empty, link it
		empty, _ := pathEmpty(t.SharedTaskDir)
		if !pathExists(t.SharedTaskDir) || empty {
			if err := linkDir(t.SharedAllocDir, t.SharedTaskDir, false); err != nil {
				return fmt.Errorf("Failed to mount shared directory for task: %w", err)
			}
			if err := linkDir(t.LogDir, filepath.Join(t.SharedTaskDir, "logs"), true); err != nil {
				return fmt.Errorf("Failed to mount shared directory for task: %w", err)
			}

		}
	}


	// Check if the directory has the shared alloc mounted.
	if pathExists(t.SharedTaskDir) {
		if err := unlinkDir(filepath.Join(t.SharedTaskDir, "logs")); err != nil {
			mErr = multierror.Append(mErr,
				fmt.Errorf("failed to unmount logs dir %q: %w",
					filepath.Join(t.SharedTaskDir, "logs"), err))
		}

		if err := unlinkDir(t.SharedTaskDir); err != nil {
			mErr = multierror.Append(mErr,
				fmt.Errorf("failed to unmount shared alloc dir %q: %w", t.SharedTaskDir, err))

/*
Package fifo implements functions to create and open a fifo for inter-process
communication in an OS agnostic way. A few assumptions should be made when using
this package. First, New() must always be called before Open(). Second Open()
returns an io.ReadWriteCloser that is only connected with the io.ReadWriteCloser
returned from New().

On Unix, all exported functions use os.Root under the hood to avoid chasing
symlinks out of their parent directory. On Windows, this is unnecessary because
named pipes exist in their own namespace and not the filesystem.
*/
package fifo
	require := require.New(t)
	var path string

	rootDir := t.TempDir()

	if runtime.GOOS == "windows" {
		path = "//./pipe/fifo"
	} else {
		path = filepath.Join(rootDir, "fifo")
	}

	readerOpenFn, err := CreateAndRead(path)
	require := require.New(t)
	var path string

	rootDir := t.TempDir()

	if runtime.GOOS == "windows" {
		path = "//./pipe/" + uuid.Generate()[:4]
	} else {
		path = filepath.Join(rootDir, "fifo")
	}

	readerOpenFn, err := CreateAndRead(path)
// SPDX-License-Identifier: BUSL-1.1

//go:build !windows

package fifo

	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CreateAndRead creates a fifo at the given path, and returns an open function
// for reading. For compatibility with windows, the fifo must not exist
// already.
//
// It returns a reader open function that may block until a writer opens
// so it's advised to run it in a goroutine different from reader goroutine
func CreateAndRead(path string) (func() (io.ReadCloser, error), error) {
	// create first
	if err := mkfifo(path, 0600); err != nil {
		return nil, fmt.Errorf("error creating fifo %v: %w", path, err)
	}

	return func() (io.ReadCloser, error) {
}

func OpenReader(path string) (io.ReadCloser, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("error opening fifo parent directory %q: %w", dir, err)
	}
	defer root.Close()

	// also uses O_NOFOLLOW under the hood
	f, err := root.OpenFile(base, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("error opening reader at %s: %w", path, err)
	}
	return f, nil
}

// OpenWriter opens a fifo file for writer, assuming it already exists, returns io.WriteCloser
func OpenWriter(path string) (io.WriteCloser, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("error opening fifo parent directory %q: %w", dir, err)
	}
	defer root.Close()

	// also uses O_NOFOLLOW under the hood
	f, err := root.OpenFile(base, os.O_WRONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("error opening writer at %s: %w", path, err)
	}
	return f, nil
}

// Remove a fifo that already exists at a given path
func Remove(path string) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("error opening root dir %q: %w", dir, err)
	}
	defer root.Close()

	return root.Remove(base)
}

func IsClosedErr(err error) bool {
	}
	return false
}
// Copyright IBM Corp. 2015, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !linux && !freebsd && !netbsd && !openbsd && !windows

package fifo

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func mkfifo(path string, mode uint32) (err error) {
	// macOS doesn't support mkfifoat
	err = unix.Mkfifo(path, mode)
	if err != nil {
		return fmt.Errorf("error creating fifo: %w", err)
	}
	return nil
}
// Copyright IBM Corp. 2015, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build linux || freebsd || netbsd || openbsd

package fifo

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func mkfifo(path string, mode uint32) (err error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("error opening fifo parent directory %q: %v", dir, err)
	}
	defer root.Close()

	parent, err := root.Open(".")
	if err != nil {
		return fmt.Errorf("error getting file handle to fifo parent directory %q: %v", dir, err)
	}
	defer parent.Close()

	// os.Root doesn't support creating a FIFO, so we need to drop to the
	// syscall and grab the parent's FD
	err = unix.Mkfifoat(int(parent.Fd()), base, mode)
	if err != nil {
		return fmt.Errorf("error creating fifo: %w", err)
	}
	return nil
}
	allocDirBind := fmt.Sprintf("%s:%s", task.TaskDir().SharedAllocDir, task.Env[taskenv.AllocDir])
	taskLocalBind := fmt.Sprintf("%s:%s", task.TaskDir().LocalDir, task.Env[taskenv.TaskLocalDir])
	secretDirBind := fmt.Sprintf("%s:%s", task.TaskDir().SecretsDir, task.Env[taskenv.SecretsDir])
	selinuxLabel := d.config.Volumes.SelinuxLabel

	binds := []string{allocDirBind, taskLocalBind, secretDirBind}

	logsROFlag := "ro"
	if selinuxLabel != "" {
		// Apply SELinux Label to each built-in bind
		for i := range binds {
			binds[i] = fmt.Sprintf("%s:%s", binds[i], selinuxLabel)
		}
		logsROFlag = "ro," + selinuxLabel
	}
	allocLogsDirBind := fmt.Sprintf("%s/logs:%s/logs:%s", task.TaskDir().SharedAllocDir, task.Env[taskenv.AllocDir], logsROFlag)
	binds = append(binds, allocLogsDirBind)

	for _, userbind := range driverConfig.Volumes {
		// This assumes host OS = docker container OS.
		"alloc:/alloc:z",
		"redis-demo/local:/local:z",
		"redis-demo/secrets:/secrets:z",
		"alloc/logs:/alloc/logs:ro,z",
		"/etc/ssl/certs:/etc/ssl/certs:ro,z",
		"/var/www:/srv/www:z",
	}, cc.Host.Binds)
		`alloc:c:/alloc`,
		`redis-demo\local:c:/local`,
		`redis-demo\secrets:c:/secrets`,
		`alloc/logs:c:/alloc/logs:ro`,
		`c:\etc\ssl\certs:c:/etc/ssl/certs`,
		`c:\var\www:c:/srv/www`,
	}, cc.Host.Binds)
