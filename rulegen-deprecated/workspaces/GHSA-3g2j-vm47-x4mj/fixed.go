package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	{name: "db_nodes_autoinc", stage: patchPreDaemonStorage, run: patchDBNodesAutoInc},
	{name: "clustering_server_cert_trust", stage: patchPreDaemonStorage, run: patchClusteringServerCertTrust},
	{name: "dnsmasq_entries_include_device_name", stage: patchPostDaemonStorage, run: patchDnsmasqEntriesIncludeDeviceName},
	{name: "pool_fix_default_permissions", stage: patchPostDaemonStorage, run: patchDefaultStoragePermissions},
}

type patch struct {
	return nil
}

// patchDefaultStoragePermissions re-applies the default modes to all storage pools.
func patchDefaultStoragePermissions(_ string, d *Daemon) error {
	pools, err := d.cluster.GetStoragePoolNames()
	if err != nil {
		// Skip the rest of the patch if no storage pools were found.
		if errors.Is(err, db.ErrNoSuchObject) {
			return nil
		}

		return fmt.Errorf("Failed getting storage pool names: %w", err)
	}

	for _, pool := range pools {
		for _, volEntry := range storageDrivers.BaseDirectories {
			for _, volDir := range volEntry.Paths {
				path := storageDrivers.GetPoolMountPath(pool) + "/" + volDir

				err := os.Chmod(path, volEntry.Mode)
				if err != nil && !errors.Is(err, fs.ErrNotExist) {
					return fmt.Errorf("Failed to set directory mode %q: %w", path, err)
				}
			}
		}
	}

	return nil
}

// Patches end here

// Here are a couple of legacy patches that were originally in
// VolumePostHook function returned from a storage action that should be run later to complete the action.
type VolumePostHook func(vol Volume) error

type baseDirectory struct {
	Paths []string
	Mode  os.FileMode
}

// BaseDirectories maps volume types to the expected directories.
var BaseDirectories = map[VolumeType]baseDirectory{
	VolumeTypeContainer: {Paths: []string{"containers", "containers-snapshots"}, Mode: 0o711}, // Containers may be run as non-root, so 0700 won't work, however as containers have their own sub-directory with correct ownership that is 0100 this is OK.
	VolumeTypeCustom:    {Paths: []string{"custom", "custom-snapshots"}, Mode: 0o700},
	VolumeTypeImage:     {Paths: []string{"images"}, Mode: 0o700},
	VolumeTypeVM:        {Paths: []string{"virtual-machines", "virtual-machines-snapshots"}, Mode: 0o700},
}

// Volume represents a storage volume, and provides functions to mount and unmount it.
	poolMountPath := GetPoolMountPath(poolName)

	for _, volType := range d.Info().VolumeTypes {
		if len(BaseDirectories[volType].Paths) < 1 {
			return nil, fmt.Errorf("Cannot get base directory name for volume type %q", volType)
		}

		volTypePath := filepath.Join(poolMountPath, BaseDirectories[volType].Paths[0])
		ents, err := os.ReadDir(volTypePath)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to list directory %q for volume type %q", volTypePath, volType)
		}
