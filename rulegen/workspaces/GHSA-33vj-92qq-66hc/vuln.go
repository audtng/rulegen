package main

	"time"

	crmetadata "github.com/checkpoint-restore/checkpointctl/lib"
	"github.com/checkpoint-restore/go-criu/v7/stats"
	"github.com/checkpoint-restore/go-criu/v7/utils"
	"github.com/containerd/containerd/api/types/runc/options"
	"github.com/containerd/containerd/v2/client"
	"github.com/opencontainers/image-spec/identity"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// checkIfCheckpointOCIImage returns checks if the input refers to a checkpoint image.
// It returns the StorageImageID of the image the input resolves to, nil otherwise.
func (c *criService) checkIfCheckpointOCIImage(ctx context.Context, input string) (string, error) {
		}(archiveFile)

		filter := archive.WithFilter(func(hdr *tar.Header) (bool, error) {
			// The checkpoint archive is unpacked twice if using a tar file directly.
			// The first time only the metadata files are relevant to prepare the
			// restore operation. This filter function ignores the large parts of
		}
	}

	if createAnnotations != nil {
		// The hash also needs to be update or Kubernetes thinks the container needs to be restarted
		_, ok1 := createAnnotations["io.kubernetes.container.hash"]
		_, ok2 := originalAnnotations["io.kubernetes.container.hash"]

		if ok1 && ok2 {
			originalAnnotations["io.kubernetes.container.hash"] = createAnnotations["io.kubernetes.container.hash"]
		}

		// The restart count also needs to be correctly updated
		_, ok1 = createAnnotations["io.kubernetes.container.restartCount"]
		_, ok2 = originalAnnotations["io.kubernetes.container.restartCount"]

		if ok1 && ok2 {
			originalAnnotations["io.kubernetes.container.restartCount"] = createAnnotations["io.kubernetes.container.restartCount"]
		}
	}

	// Pulling the image the checkpoint is based on. This is a bit different
	// than automatic image pulling. The checkpoint image is not automatically
	if _, err := reference.ParseAnyReference(config.RootfsImageName); err != nil {
		return "", fmt.Errorf("error parsing reference: %q is not a valid repository/tag %v", config.RootfsImageName, err)
	}
	tagImage, err := c.client.ImageService().Get(ctx, config.RootfsImageRef)
	if err != nil {
		return "", fmt.Errorf("failed to get checkpoint base image %s: %w", config.RootfsImageRef, err)
	}
	// Second step is to tag the image with the same tag it used to have
	// during checkpointing. For the error that the image NAME:TAG already
	// exists is ignored. It could happen that NAME:TAG now belongs to
	// another NAME@DIGEST than during checkpointing and the restore will
	// happen on another image.
	// TODO: handle if NAME:TAG points to a different NAME@DIGEST
	tagImage.Name = config.RootfsImageName
	_, err = c.client.ImageService().Create(ctx, tagImage)
	if err != nil {
		if !errdefs.IsAlreadyExists(err) {
			return "", fmt.Errorf("failed to tag checkpoint base image %s with %s: %w", config.RootfsImageRef, config.RootfsImageName, err)
		}
	}

	var image imagestore.Image
	for i := 1; i < 500; i++ {
		return "", err
	}

	if restoreStorageImageID != "" {
		if err := fs.CopyDir(containerRootDir, mountPoint); err != nil {
			return "", err
		}
		if err := mount.UnmountAll(mountPoint, 0); err != nil {
			return "", err
		}
	} else {
		// unpack the checkpoint archive
		filter := archive.WithFilter(func(hdr *tar.Header) (bool, error) {
			excludePatterns := []string{
				crmetadata.ConfigDumpFile,
				crmetadata.SpecDumpFile,

		// Start from the beginning of the checkpoint archive
		archiveFile.Seek(0, 0)
		_, err = archive.Apply(ctx, containerRootDir, archiveFile, []archive.ApplyOpt{filter}...)

		if err != nil {
			return "", fmt.Errorf("unpacking of checkpoint archive %s failed: %w", containerRootDir, err)
		}
	}
	log.G(ctx).Debugf("Unpacked checkpoint in %s", containerRootDir)

	// Restore container log file (if it exists)
	containerLog := filepath.Join(containerRootDir, "container.log")
	_, err = c.os.Stat(containerLog)
	if err == nil {
		if err := c.os.CopyFile(containerLog, meta.LogPath, 0600); err != nil {
			return "", fmt.Errorf("restoring container log file %s failed: %w", containerLog, err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task for container %q: %w", r.GetContainerId(), err)
	}
	img, err := task.Checkpoint(ctx, []client.CheckpointTaskOpts{withCheckpointOpts(i.Runtime.Name, c.getContainerRootDir(r.GetContainerId()))}...)
	if err != nil {
		return nil, fmt.Errorf("checkpointing container %q failed: %w", r.GetContainerId(), err)
	}
		return nil, fmt.Errorf("failed to unmarshall blob into checkpoint data OCI index: %w", err)
	}

	cpPath := filepath.Join(c.getContainerRootDir(r.GetContainerId()), "ctrd-checkpoint")
	if err := os.MkdirAll(cpPath, 0o700); err != nil {
		return nil, err
	}
	defer os.RemoveAll(cpPath)

	// This internal containerd file is used by checkpointctl for
	// checkpoint archive analysis.
	if err := c.os.CopyFile(
		filepath.Join(c.getContainerRootDir(r.GetContainerId()), crmetadata.StatusFile),
		filepath.Join(cpPath, crmetadata.StatusFile),
		0o600,
	); err != nil {
		return nil, err
	}

	// This file is created by CRIU and includes timing analysis.
	// Also used by checkpointctl
	if err := c.os.CopyFile(
		filepath.Join(c.getContainerRootDir(r.GetContainerId()), stats.StatsDump),
		filepath.Join(cpPath, stats.StatsDump),
		0o600,
	); err != nil {
		return nil, err
	}

	// The log file created by CRIU. This file could be missing.
	// Let's ignore errors if the file is missing.
	if err := c.os.CopyFile(
		filepath.Join(c.getContainerRootDir(r.GetContainerId()), crmetadata.DumpLogFile),
		filepath.Join(cpPath, crmetadata.DumpLogFile),
		0o600,
	); err != nil {
		if !errors.Is(errors.Unwrap(err), os.ErrNotExist) {
			return nil, err
		}
	}

	// Save the existing container log file
	_, err = c.os.Stat(criContainerStatus.GetStatus().GetLogPath())

	return nil
}
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/containerd/log"
//
// It returns:
//  1. New options that include new "lowedir=..." mount option.
//  2. "Clean up" function -- it should be called as a defer one before
//     checking for error, because if do the second and avoid calling "clean up",
//     you're going to have "dirty" setup -- there's no guarantee that those
//     temporary mount points for lowedirs will be cleaned properly.
//  3. Error -- nil if everything's fine, otherwise an error.
func prepareIDMappedOverlay(usernsFd int, options []string) ([]string, func(), error) {
	lowerIdx, lowerDirs := findOverlayLowerdirs(options)
		return options, nil, fmt.Errorf("failed to parse overlay lowerdir's from given options")
	}

	tempRemountsLocation, err := os.MkdirTemp(tempMountLocation, "ovl-idmapped")
	if err != nil {
		return options, nil, fmt.Errorf("failed to create temporary overlay lowerdir mount location: %w", err)
	}

	tmpLowerdirs, idMapCleanUp, err := doPrepareIDMappedOverlay(tempRemountsLocation, lowerDirs, usernsFd)
	if err != nil {
		return options, idMapCleanUp, fmt.Errorf("failed to create idmapped mount: %w", err)
	}
				userNsCleanUp func()
			)
			options, userNsCleanUp, err = prepareIDMappedOverlay(int(usernsFd.Fd()), options)
			defer userNsCleanUp()

			if err != nil {
				return fmt.Errorf("failed to prepare idmapped overlay: %w", err)
			}
			// To not meet concurrency issues while using the same lowedirs
			// for different containers, replace them by temporary directories,
			if optionsSize(options) >= pagesize-512 {
	}

	var flags int
	for flag := range unprivilegedFlags {
		if int(statfs.Flags)&flag == flag {
			flags |= flag
		}
	return flags, nil
}

func doPrepareIDMappedOverlay(tempRemountsLocation string, lowerDirs []string, usernsFd int) ([]string, func(), error) {
	tmpLowerDirs := make([]string, 0, len(lowerDirs))

	cleanUp := func() {
		for _, lowerDir := range tmpLowerDirs {
			if err := unix.Unmount(lowerDir, 0); err != nil {
				log.L.WithError(err).Warnf("failed to unmount temp lowerdir %s", lowerDir)
				continue
			}
			// Using os.Remove() so if it's not empty, we don't delete files in the
			// rootfs.
			if err := os.Remove(lowerDir); err != nil {
				log.L.WithError(err).Warnf("failed to remove temporary overlay lowerdir")
			}
		}

		// This dir should be empty now. Otherwise, we don't do anything.
		if err := os.Remove(tempRemountsLocation); err != nil {
			log.L.WithError(err).Infof("failed to remove temporary overlay dir")
		}
	}
	for i, lowerDir := range lowerDirs {
		tmpLowerDir := filepath.Join(tempRemountsLocation, strconv.Itoa(i))
		tmpLowerDirs = append(tmpLowerDirs, tmpLowerDir)

		if err := os.MkdirAll(tmpLowerDir, 0700); err != nil {
			return nil, cleanUp, fmt.Errorf("failed to create temporary dir: %w", err)
		}
		if err := IDMapMountWithAttrs(lowerDir, tmpLowerDir, usernsFd, unix.MOUNT_ATTR_RDONLY, 0); err != nil {
			return nil, cleanUp, err
		}
	}
	return tmpLowerDirs, cleanUp, nil
}

// parseMountOptions takes fstab style mount options and parses them for

func mountAt(chdir string, source, target, fstype string, flags uintptr, data string) error {
	if chdir == "" {
		return unix.Mount(source, target, fstype, flags, data)
	}

	ch := make(chan error, 1)
			ch <- err
			return
		}

		ch <- unix.Mount(source, target, fstype, flags, data)
	}()
	return <-ch
}
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
			return nil
		}

		parts := strings.Split(userstr, ":")
		switch len(parts) {
		case 1:
			v, err := strconv.Atoi(parts[0])
			if err != nil || v < minUserID || v > maxUserID {
				// if we cannot parse as an int32 then try to see if it is a username
				return WithUsername(userstr)(ctx, client, c, s)
			}
			return WithUserID(uint32(v))(ctx, client, c, s)
		case 2:
			var (
			)
			var uid, gid uint32
			v, err := strconv.Atoi(parts[0])
			if err != nil || v < minUserID || v > maxUserID {
				username = parts[0]
			} else {
				uid = uint32(v)
			}
			v, err = strconv.Atoi(parts[1])
			if err != nil || v < minGroupID || v > maxGroupID {
				groupname = parts[1]
			} else {
				gid = uint32(v)
			}
			if err != nil {
				return err
			}
			ugroups, groupErr := user.ParseGroupFile(gpath)
			if groupErr != nil && !os.IsNotExist(groupErr) {
				return groupErr
			}
	if err != nil {
		return user.User{}, err
	}
	users, err := user.ParsePasswdFileFilter(ppath, filter)
	if err != nil {
		return user.User{}, err
	}
	if err != nil {
		return 0, err
	}
	groups, err := user.ParseGroupFileFilter(gpath, filter)
	if err != nil {
		return 0, err
	}
	if err != nil {
		return []uint32{}, err
	}
	groups, err := user.ParseGroupFileFilter(gpath, filter)
	if err != nil {
		return []uint32{}, err
	}
		return nil
	}
}
