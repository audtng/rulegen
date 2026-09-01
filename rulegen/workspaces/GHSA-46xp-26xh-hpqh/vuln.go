package main

type nonOpManager struct {
}

func (no *nonOpManager) UnsafeSetFileOwnership(file string) error {
	return nil
}

func (no *nonOpManager) SetFileOwnership(file *safepath.Path) error {
	return nil
}

type OwnershipManager struct {
	user string
}
					// Start the VirtualMachineInstance with the PVC attached
					vmi = newVMI(pvcName)

					vmi = libvmops.RunVMIAndExpectLaunch(vmi, 180)

					By(checkingVMInstanceConsoleOut)
					Expect(console.LoginToAlpine(vmi)).To(Succeed())
				},
					Entry("[test_id:3130]with Disk PVC", newRandomVMIWithPVC, true),
					Entry("[test_id:3131]with CDRom PVC", newRandomVMIWithCDRom, true),
					Entry("hostpath disk image file not owned by qemu", newRandomVMIWithPVC, false),
				)
			})

	hdc.lessPVCSpaceToleration = toleration
}

func (hdc DiskImgCreator) Create(vmi *v1.VirtualMachineInstance) error {
	for _, volume := range vmi.Spec.Volumes {
		if hostDisk := volume.VolumeSource.HostDisk; shouldMountHostDisk(hostDisk) {
			if err := hdc.mountHostDiskAndSetOwnership(vmi, volume.Name, hostDisk); err != nil {
	}

	if fileNotExists {
		if err := hdc.handleRequestedSizeAndCreateSparseRaw(vmi, diskDir, filepath.Base(hostDisk.Path), hostDisk); err != nil {
			return err
		}

		if err != nil {
			return err
		}
	}
	// Change file ownership to the qemu user.
	if err := ephemeraldiskutils.DefaultOwnershipManager.SetFileOwnership(diskPath); err != nil {
		log.Log.Reason(err).Errorf("Couldn't set Ownership on %s: %v", diskPath, err)
		return err
	}
	return nil
}
