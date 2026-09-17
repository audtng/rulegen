package main

	"time"

	crmetadata "github.com/checkpoint-restore/checkpointctl/lib"
	"github.com/checkpoint-restore/go-criu/v7/utils"
	"github.com/containerd/containerd/api/types/runc/options"
	"github.com/containerd/containerd/v2/client"
	"github.com/opencontainers/image-spec/identity"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"

	// TODO: This package import is kept to prevent merge conflicts while integrating multiple
	// branches, specifically because this changes vendoring.
	_ "github.com/checkpoint-restore/go-criu/v7/stats"
)

// copyNoFollow copies the regular file at src to dst without following a symlink
// at the final path component of src.
//
// The checkpoint code reads files (container.log, status, stats-dump, dump.log)
// out of the container state directory, which can contain entries unpacked from a
// checkpoint archive or OCI image. Those entries are externally provided, so they
// are read defensively.
//
// src is first lstat'd (which does not follow a final-component symlink) and must
// be a regular file; non-regular entries are rejected before src is ever opened.
// src is then opened with O_NOFOLLOW as a belt-and-suspenders guard in case the
// entry changes type between the lstat and the open.
func copyNoFollow(src, dst string, perm os.FileMode) error {
	fi, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("refusing to copy %s: not a regular file", src)
	}

	in, err := os.OpenFile(src, os.O_RDONLY|unix.O_NOFOLLOW, 0)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// checkpointArchiveEntryAllowed reports whether a tar entry from a checkpoint
// archive may be unpacked. Legitimate checkpoint archives contain only regular
// files and directories; other entry types (symlinks, hardlinks, device and fifo
// nodes) are not produced by the checkpoint code and are rejected as a hardening
// measure.
func checkpointArchiveEntryAllowed(hdr *tar.Header) bool {
	switch hdr.Typeflag {
	//nolint:staticcheck // TypeRegA is deprecated but we may still receive an external tar with TypeRegA
	case tar.TypeReg, tar.TypeRegA, tar.TypeDir, tar.TypeXGlobalHeader:
		return true
	default:
		return false
	}
}

// assertCheckpointDirSafe verifies that the populated restore directory contains
// only regular files and directories.
//
// The OCI-image restore path copies checkpoint content into the restore dir with
// fs.CopyDir, which (unlike the tar unpack filter) faithfully recreates any
// symlinks and special files present in the image. Restore-time consumers open
// paths under this directory, so non-regular entries are rejected before they run.
func assertCheckpointDirSafe(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// d.Type() reports the entry type without following symlinks.
		if d.IsDir() || d.Type().IsRegular() {
			return nil
		}
		return fmt.Errorf("refusing to restore checkpoint: %s is not a regular file or directory", path)
	})
}

// checkIfCheckpointOCIImage returns checks if the input refers to a checkpoint image.
// It returns the StorageImageID of the image the input resolves to, nil otherwise.
func (c *criService) checkIfCheckpointOCIImage(ctx context.Context, input string) (string, error) {
		}(archiveFile)

		filter := archive.WithFilter(func(hdr *tar.Header) (bool, error) {
			// Reject entry types the checkpoint code never produces (symlinks,
			// hardlinks, device/fifo nodes) so they are not recreated on disk.
			if !checkpointArchiveEntryAllowed(hdr) {
				log.G(ctx).Warnf("Skipping unexpected checkpoint archive entry %q (type %d)", hdr.Name, hdr.Typeflag)
				return false, nil
			}
			// The checkpoint archive is unpacked twice if using a tar file directly.
			// The first time only the metadata files are relevant to prepare the
			// restore operation. This filter function ignores the large parts of
		}
	}

	originalAnnotations = filterAndMergeAnnotations(
		ctx,
		originalAnnotations,
		createAnnotations,
	)

	// Pulling the image the checkpoint is based on. This is a bit different
	// than automatic image pulling. The checkpoint image is not automatically
	if _, err := reference.ParseAnyReference(config.RootfsImageName); err != nil {
		return "", fmt.Errorf("error parsing reference: %q is not a valid repository/tag %v", config.RootfsImageName, err)
	}

	var image imagestore.Image
	for i := 1; i < 500; i++ {
		return "", err
	}

	// Confine all checkpoint content to a dedicated subdirectory of the container
	// state dir instead of unpacking it directly into the state dir, so it cannot
	// collide with containerd's own files there. Create it fresh; RemoveAll unlinks
	// any pre-existing entry without following it.
	restoreDir := filepath.Join(containerRootDir, checkpointRestoreDir)
	if err := os.RemoveAll(restoreDir); err != nil {
		return "", err
	}
	if err := os.Mkdir(restoreDir, 0o700); err != nil {
		return "", err
	}

	if restoreStorageImageID != "" {
		if err := fs.CopyDir(restoreDir, mountPoint); err != nil {
			return "", err
		}
		if err := mount.UnmountAll(mountPoint, 0); err != nil {
			return "", err
		}
		// fs.CopyDir recreates any symlinks/special files from the image; reject
		// them here so restore-time consumers only ever open regular files.
		if err := assertCheckpointDirSafe(restoreDir); err != nil {
			return "", err
		}
	} else {
		// unpack the checkpoint archive
		filter := archive.WithFilter(func(hdr *tar.Header) (bool, error) {
			// Reject entry types the checkpoint code never produces (symlinks,
			// hardlinks, device/fifo nodes) so they are not recreated on disk.
			if !checkpointArchiveEntryAllowed(hdr) {
				log.G(ctx).Warnf("Skipping unexpected checkpoint archive entry %q (type %d)", hdr.Name, hdr.Typeflag)
				return false, nil
			}
			excludePatterns := []string{
				crmetadata.ConfigDumpFile,
				crmetadata.SpecDumpFile,

		// Start from the beginning of the checkpoint archive
		archiveFile.Seek(0, 0)
		_, err = archive.Apply(ctx, restoreDir, archiveFile, []archive.ApplyOpt{filter}...)

		if err != nil {
			return "", fmt.Errorf("unpacking of checkpoint archive %s failed: %w", restoreDir, err)
		}
	}
	log.G(ctx).Debugf("Unpacked checkpoint in %s", restoreDir)

	// Restore container log file (if it exists).
	//
	// container.log was unpacked from a checkpoint archive/OCI image, so it is
	// copied without following a final-component symlink.
	containerLog := filepath.Join(restoreDir, "container.log")
	if err := copyNoFollow(containerLog, meta.LogPath, 0600); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("restoring container log file %s failed: %w", containerLog, err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task for container %q: %w", r.GetContainerId(), err)
	}

	cpPath := filepath.Join(c.getContainerRootDir(container.ID), "ctrd-checkpoint")
	// ctrd-checkpoint may already exist from a prior checkpoint operation. RemoveAll
	// unlinks any existing entry (including a symlink) itself rather than its target,
	// so creating the directory afterwards cannot write through a link.
	if err := os.RemoveAll(cpPath); err != nil {
		return nil, err
	}
	if err := os.Mkdir(cpPath, 0o700); err != nil {
		return nil, err
	}
	defer os.RemoveAll(cpPath)

	// Point CRIU's work directory (where it writes dump.log and stats-dump) at the
	// dedicated, freshly-created checkpoint dir instead of the persistent container
	// state dir. Otherwise checkpoint creation litters those files into the state
	// dir where they are never cleaned up; here they land directly where they are
	// archived from and are removed with cpPath.
	img, err := task.Checkpoint(ctx, []client.CheckpointTaskOpts{withCheckpointOpts(i.Runtime.Name, cpPath)}...)
	if err != nil {
		return nil, fmt.Errorf("checkpointing container %q failed: %w", r.GetContainerId(), err)
	}
		return nil, fmt.Errorf("failed to unmarshall blob into checkpoint data OCI index: %w", err)
	}

	// This internal containerd file is used by checkpointctl for checkpoint archive
	// analysis. It lives in the container state dir, which can hold files from a
	// prior checkpoint operation, so it is read without following symlinks.
	if err := copyNoFollow(
		filepath.Join(c.getContainerRootDir(container.ID), crmetadata.StatusFile),
		filepath.Join(cpPath, crmetadata.StatusFile),
		0o600,
	); err != nil {
		return nil, err
	}

	// dump.log and stats-dump are written directly into cpPath by CRIU via its
	// work directory (see withCheckpointOpts above), so they are already present
	// for archiving and do not need to be copied out of the container state dir.

	// Save the existing container log file
	_, err = c.os.Stat(criContainerStatus.GetStatus().GetLogPath())

	return nil
}

func filterAndMergeAnnotations(
	ctx context.Context,
	checkpointAnnotations map[string]string,
	createAnnotations map[string]string,
) map[string]string {
	result := make(map[string]string)

	for k, v := range checkpointAnnotations {
		if strings.HasPrefix(k, "cdi.k8s.io/") || k == "cdi.k8s.io" {
			log.G(ctx).Warnf("Denying annotation %q in checkpoint restore", k)
			continue
		}
		result[k] = v
	}

	// The hash also needs to be update or Kubernetes thinks the container needs to be restarted
	_, ok1 := createAnnotations["io.kubernetes.container.hash"]
	_, ok2 := result["io.kubernetes.container.hash"]

	if ok1 && ok2 {
		result["io.kubernetes.container.hash"] = createAnnotations["io.kubernetes.container.hash"]
	}

	// The restart count also needs to be correctly updated
	_, ok1 = createAnnotations["io.kubernetes.container.restartCount"]
	_, ok2 = result["io.kubernetes.container.restartCount"]

	if ok1 && ok2 {
		result["io.kubernetes.container.restartCount"] = createAnnotations["io.kubernetes.container.restartCount"]
	}

	return result
}
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/containerd/log"
//
// It returns:
//  1. New options that include new "lowedir=..." mount option.
//  2. "Clean up" function -- it should be called only if no error is returned.
//  3. Error -- nil if everything's fine, otherwise an error.
func prepareIDMappedOverlay(usernsFd int, options []string) ([]string, func(), error) {
	lowerIdx, lowerDirs := findOverlayLowerdirs(options)
		return options, nil, fmt.Errorf("failed to parse overlay lowerdir's from given options")
	}

	tmpLowerdirs, idMapCleanUp, err := doPrepareIDMappedOverlay(tempMountLocation, lowerDirs, usernsFd)
	if err != nil {
		return options, idMapCleanUp, fmt.Errorf("failed to create idmapped mount: %w", err)
	}
				userNsCleanUp func()
			)
			options, userNsCleanUp, err = prepareIDMappedOverlay(int(usernsFd.Fd()), options)
			if err != nil {
				return fmt.Errorf("failed to prepare idmapped overlay: %w", err)
			}
			defer userNsCleanUp()

			// To not meet concurrency issues while using the same lowedirs
			// for different containers, replace them by temporary directories,
			if optionsSize(options) >= pagesize-512 {
	}

	var flags int
	for _, flag := range unprivilegedFlags {
		if int(statfs.Flags)&flag == flag {
			flags |= flag
		}
	return flags, nil
}

func doPrepareIDMappedOverlay(tmpDir string, lowerDirs []string, usernsFd int) (_ []string, _ func(), retErr error) {
	commonDir, err := getCommonDirectory(lowerDirs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to determine common parent: %w", err)
	}

	tempRemountsLocation, err := os.MkdirTemp(tmpDir, "ovl-idmapped")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temporary overlay lowerdir mount location: %w", err)
	}
	cleanDir := func() {
		if err := os.Remove(tempRemountsLocation); err != nil {
			log.L.WithError(err).Infof("failed to remove idmapped directory")
		}
	}
	defer func() {
		if retErr != nil {
			cleanDir()
		}
	}()

	// IDMapMount the directory containing all the layers
	if err := IDMapMountWithAttrs(commonDir, tempRemountsLocation, usernsFd, unix.MOUNT_ATTR_RDONLY, 0); err != nil {
		return nil, nil, err
	}
	cleanMount := func() {
		// Use the Unmount helper that does retries because there can be easily an open fd
		// to the idmapped directory and when containerd forks to create a userns fd (maybe
		// for another container), it will make the mount busy for a few ms.
		err := Unmount(tempRemountsLocation, 0)
		if err != nil {
			log.L.WithError(err).Warnf("failed to unmount idmapped directory %s: %v", tempRemountsLocation, err)
		}
	}
	defer func() {
		if retErr != nil {
			cleanMount()
		}
	}()

	// Build new lower dir paths through the idmapped directory
	tmpLowerDirs := buildIDMappedPaths(lowerDirs, commonDir, tempRemountsLocation)

	cleanup := func() {
		cleanMount()
		cleanDir()
	}
	return tmpLowerDirs, cleanup, nil
}

// getCommonDirectory finds the common directory among the lowerDirs passed in.
// "/" and "." are considered invalid common directories and are treated as error
func getCommonDirectory(lowerDirs []string) (string, error) {
	commonPrefix := longestCommonPrefix(lowerDirs)
	if commonPrefix == "" {
		return "", fmt.Errorf("no common prefix found")
	}

	// Ensure the common prefix ends at a directory boundary
	commonPrefix = path.Dir(commonPrefix)

	if commonPrefix == "." || commonPrefix == "/" {
		return "", fmt.Errorf("invalid common directory: %s", commonPrefix)
	}

	return commonPrefix, nil
}

// buildIDMappedPaths constructs new lower directory paths through an idmapped mount of the commonDir.
// It takes the original lowerDirs, the commonDir of those dirs, and rewrites the paths
// to go through the idMappedDir directory to achieve idmapped lowerdirs ready for overlayfs
func buildIDMappedPaths(lowerDirs []string, commonDir, idMappedDir string) []string {
	tmpLowerDirs := make([]string, 0, len(lowerDirs))

	for _, lowerDir := range lowerDirs {
		relativePath := strings.TrimPrefix(lowerDir, commonDir)
		tmpLowerDirs = append(tmpLowerDirs, filepath.Join(idMappedDir, relativePath))
	}

	return tmpLowerDirs
}

// parseMountOptions takes fstab style mount options and parses them for

func mountAt(chdir string, source, target, fstype string, flags uintptr, data string) error {
	if chdir == "" {
		err := unix.Mount(source, target, fstype, flags, data)
		if err != nil {
			return fmt.Errorf("mount source: %q, target: %q, fstype: %s, flags: %d, data: %q, err: %w", source, target, fstype, flags, data, err)
		}
		return nil
	}

	ch := make(chan error, 1)
			ch <- err
			return
		}
		err := unix.Mount(source, target, fstype, flags, data)
		if err != nil {
			err = fmt.Errorf("mount source: %q, target: %q, fstype: %s, flags: %d, data: %q, err: %w", source, target, fstype, flags, data, err)
		}
		ch <- err
	}()
	return <-ch
}
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
			return nil
		}

		isErrRange := func(err error) bool {
			var numErr *strconv.NumError
			return errors.As(err, &numErr) && numErr.Err == strconv.ErrRange
		}

		parts := strings.Split(userstr, ":")
		switch len(parts) {
		case 1:
			v, err := strconv.Atoi(parts[0])
			if err != nil {
				if isErrRange(err) {
					return fmt.Errorf("invalid USER value %q: uid out of range", userstr)
				}
				// Non-numeric user value; treat it as a username.
				return WithUsername(userstr)(ctx, client, c, s)
			}
			if v < minUserID || v > maxUserID {
				return fmt.Errorf("invalid USER value %q: uid out of range", userstr)
			}
			return WithUserID(uint32(v))(ctx, client, c, s)
		case 2:
			var (
			)
			var uid, gid uint32
			v, err := strconv.Atoi(parts[0])
			if err != nil {
				if isErrRange(err) {
					return fmt.Errorf("invalid USER value %q: uid out of range", userstr)
				}
				username = parts[0]
			} else if v < minUserID || v > maxUserID {
				return fmt.Errorf("invalid USER value %q: uid out of range", userstr)
			} else {
				uid = uint32(v)
			}
			v, err = strconv.Atoi(parts[1])
			if err != nil {
				if isErrRange(err) {
					return fmt.Errorf("invalid USER value %q: gid out of range", userstr)
				}
				groupname = parts[1]
			} else if v < minGroupID || v > maxGroupID {
				return fmt.Errorf("invalid USER value %q: gid out of range", userstr)
			} else {
				gid = uint32(v)
			}
			if err != nil {
				return err
			}
			var ugroups []user.Group
			f, groupErr := openBoundedUserFile(gpath)
			if groupErr == nil {
				ugroups, groupErr = user.ParseGroup(f)
				f.Close()
			}
			if groupErr != nil && !os.IsNotExist(groupErr) {
				return groupErr
			}
	if err != nil {
		return user.User{}, err
	}
	f, err := openBoundedUserFile(ppath)
	if err != nil {
		return user.User{}, err
	}
	defer f.Close()
	users, err := user.ParsePasswdFilter(f, filter)
	if err != nil {
		return user.User{}, err
	}
	if err != nil {
		return 0, err
	}
	f, err := openBoundedUserFile(gpath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	groups, err := user.ParseGroupFilter(f, filter)
	if err != nil {
		return 0, err
	}
	if err != nil {
		return []uint32{}, err
	}
	f, err := openBoundedUserFile(gpath)
	if err != nil {
		return []uint32{}, err
	}
	defer f.Close()
	groups, err := user.ParseGroupFilter(f, filter)
	if err != nil {
		return []uint32{}, err
	}
		return nil
	}
}

// maxUserFileBytes caps how much data is read from any user-database file
// opened via openBoundedUserFile. Real systems keep these files well under
// 1 MiB; 10 MiB is generous headroom while keeping peak memory during
// user.ParsePasswd/ParseGroup bounded to single-digit MiB.
const maxUserFileBytes = 10 << 20

// openBoundedUserFile opens path and returns an io.ReadCloser that errors out
// if more than maxUserFileBytes are read from it. Non-regular sources are
// rejected before opening, so callers never block on FIFOs or device files
// and parsers never consume bytes from them.
//
// openBoundedUserFile does NOT perform any path validation. It does not guard
// against symlink traversal or paths that escape the container rootfs, and it
// follows symlinks both when stat-ing and when opening. Callers are responsible
// for confining path to the intended root beforehand (e.g. via fs.RootPath,
// which resolves every symlink component and re-anchors absolute links to the
// root) and must not pass attacker-controlled, unresolved paths directly.
func openBoundedUserFile(path string) (io.ReadCloser, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &limitedFile{
		Closer: f,
		// Allow one byte past the cap so an overflow surfaces as an
		// error rather than a silent EOF that the parser would treat as
		// a clean end-of-file (and miss any entries past the cap).
		r:    &io.LimitedReader{R: f, N: maxUserFileBytes + 1},
		name: path,
	}, nil
}

// limitedFile is an io.ReadCloser whose Read returns an error once more than
// maxUserFileBytes have been read.
type limitedFile struct {
	io.Closer
	r    *io.LimitedReader
	name string
}

func (l *limitedFile) Read(p []byte) (int, error) {
	n, err := l.r.Read(p)
	if l.r.N == 0 {
		return n, fmt.Errorf("%q exceeds %d bytes", l.name, maxUserFileBytes)
	}
	return n, err
}
