package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	{name: "db_nodes_autoinc", stage: patchPreDaemonStorage, run: patchDBNodesAutoInc},
	{name: "clustering_server_cert_trust", stage: patchPreDaemonStorage, run: patchClusteringServerCertTrust},
	{name: "dnsmasq_entries_include_device_name", stage: patchPostDaemonStorage, run: patchDnsmasqEntriesIncludeDeviceName},
}

type patch struct {
	return nil
}

// Patches end here

// Here are a couple of legacy patches that were originally in
// VolumePostHook function returned from a storage action that should be run later to complete the action.
type VolumePostHook func(vol Volume) error

// BaseDirectories maps volume types to the expected directories.
var BaseDirectories = map[VolumeType][]string{
	VolumeTypeContainer: {"containers", "containers-snapshots"},
	VolumeTypeCustom:    {"custom", "custom-snapshots"},
	VolumeTypeImage:     {"images"},
	VolumeTypeVM:        {"virtual-machines", "virtual-machines-snapshots"},
}

// Volume represents a storage volume, and provides functions to mount and unmount it.
	poolMountPath := GetPoolMountPath(poolName)

	for _, volType := range d.Info().VolumeTypes {
		if len(BaseDirectories[volType]) < 1 {
			return nil, fmt.Errorf("Cannot get base directory name for volume type %q", volType)
		}

		volTypePath := filepath.Join(poolMountPath, BaseDirectories[volType][0])
		ents, err := ioutil.ReadDir(volTypePath)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to list directory %q for volume type %q", volTypePath, volType)
		}
