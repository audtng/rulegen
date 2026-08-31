package main

)

// linkDir hardlinks src to dst. The src and dst must be on the same filesystem.
func linkDir(src, dst string) error {
	return syscall.Link(src, dst)
}

)

// linkDir hardlinks src to dst. The src and dst must be on the same filesystem.
func linkDir(src, dst string) error {
	return syscall.Link(src, dst)
}


// linkDir bind mounts src to dst as Linux doesn't support hardlinking
// directories.
func linkDir(src, dst string) error {
	if err := os.MkdirAll(dst, fileMode777); err != nil {
		return err
	}

	return syscall.Mount(src, dst, "", syscall.MS_BIND, "")
}

)

// linkDir hardlinks src to dst. The src and dst must be on the same filesystem.
func linkDir(src, dst string) error {
	return syscall.Link(src, dst)
}

)

// linkDir hardlinks src to dst. The src and dst must be on the same filesystem.
func linkDir(src, dst string) error {
	return syscall.Link(src, dst)
}

}

// The windows version does nothing currently.
func linkDir(src, dst string) error {
	return nil
}

		// If the path doesn't exist OR it exists and is empty, link it
		empty, _ := pathEmpty(t.SharedTaskDir)
		if !pathExists(t.SharedTaskDir) || empty {
			if err := linkDir(t.SharedAllocDir, t.SharedTaskDir); err != nil {
				return fmt.Errorf("Failed to mount shared directory for task: %w", err)
			}
		}
	}


	// Check if the directory has the shared alloc mounted.
	if pathExists(t.SharedTaskDir) {
		if err := unlinkDir(t.SharedTaskDir); err != nil {
			mErr = multierror.Append(mErr,
				fmt.Errorf("failed to unmount shared alloc dir %q: %w", t.SharedTaskDir, err))

/*
Package fifo implements functions to create and open a fifo for inter-process
communication in an OS agnostic way. A few assumptions should be made when
using this package. First, New() must always be called before Open(). Second
Open() returns an io.ReadWriteCloser that is only connected with the
io.ReadWriteCloser returned from New().
*/
package fifo
	require := require.New(t)
	var path string

	if runtime.GOOS == "windows" {
		path = "//./pipe/fifo"
	} else {
		path = filepath.Join(t.TempDir(), "fifo")
	}

	readerOpenFn, err := CreateAndRead(path)
	require := require.New(t)
	var path string

	if runtime.GOOS == "windows" {
		path = "//./pipe/" + uuid.Generate()[:4]
	} else {
		path = filepath.Join(t.TempDir(), "fifo")
	}

	readerOpenFn, err := CreateAndRead(path)
// SPDX-License-Identifier: BUSL-1.1

//go:build !windows
// +build !windows

package fifo

	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

// CreateAndRead creates a fifo at the given path, and returns an open function for reading.
// For compatibility with windows, the fifo must not exist already.
//
// It returns a reader open function that may block until a writer opens
// so it's advised to run it in a goroutine different from reader goroutine
func CreateAndRead(path string) (func() (io.ReadCloser, error), error) {
	// create first
	if err := mkfifo(path, 0600); err != nil {
		return nil, fmt.Errorf("error creating fifo %v: %v", path, err)
	}

	return func() (io.ReadCloser, error) {
}

func OpenReader(path string) (io.ReadCloser, error) {
	return os.OpenFile(path, unix.O_RDONLY, os.ModeNamedPipe)
}

// OpenWriter opens a fifo file for writer, assuming it already exists, returns io.WriteCloser
func OpenWriter(path string) (io.WriteCloser, error) {
	return os.OpenFile(path, unix.O_WRONLY, os.ModeNamedPipe)
}

// Remove a fifo that already exists at a given path
func Remove(path string) error {
	return os.Remove(path)
}

func IsClosedErr(err error) bool {
	}
	return false
}

func mkfifo(path string, mode uint32) (err error) {
	return unix.Mkfifo(path, mode)
}
	allocDirBind := fmt.Sprintf("%s:%s", task.TaskDir().SharedAllocDir, task.Env[taskenv.AllocDir])
	taskLocalBind := fmt.Sprintf("%s:%s", task.TaskDir().LocalDir, task.Env[taskenv.TaskLocalDir])
	secretDirBind := fmt.Sprintf("%s:%s", task.TaskDir().SecretsDir, task.Env[taskenv.SecretsDir])
	binds := []string{allocDirBind, taskLocalBind, secretDirBind}

	selinuxLabel := d.config.Volumes.SelinuxLabel
	if selinuxLabel != "" {
		// Apply SELinux Label to each built-in bind
		for i := range binds {
			binds[i] = fmt.Sprintf("%s:%s", binds[i], selinuxLabel)
		}
	}

	for _, userbind := range driverConfig.Volumes {
		// This assumes host OS = docker container OS.
		"alloc:/alloc:z",
		"redis-demo/local:/local:z",
		"redis-demo/secrets:/secrets:z",
		"/etc/ssl/certs:/etc/ssl/certs:ro,z",
		"/var/www:/srv/www:z",
	}, cc.Host.Binds)
		`alloc:c:/alloc`,
		`redis-demo\local:c:/local`,
		`redis-demo\secrets:c:/secrets`,
		`c:\etc\ssl\certs:c:/etc/ssl/certs`,
		`c:\var\www:c:/srv/www`,
	}, cc.Host.Binds)
