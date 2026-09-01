package main

type nonOpManager struct {
}

func (no *nonOpManager) UnsafeSetFileOwnership(_ string) error {
	return nil
}

func (no *nonOpManager) SetFileOwnership(_ *safepath.Path) error {
	return nil
}

func MockDefaultOwnershipManagerWithFailure() {
	DefaultOwnershipManager = &failureManager{}
}

type failureManager struct {
}

func (no *failureManager) UnsafeSetFileOwnership(_ string) error {
	panic("unexpected call to UnsafeSetFileOwnership")
}

func (no *failureManager) SetFileOwnership(_ *safepath.Path) error {
	panic("unexpected call to SetFileOwnership")
}

type OwnershipManager struct {
	user string
}
					// Start the VirtualMachineInstance with the PVC attached
					vmi = newVMI(pvcName)

					if imageOwnedByQEMU {
						vmi = libvmops.RunVMIAndExpectLaunch(vmi, 180)

						By(checkingVMInstanceConsoleOut)
						Expect(console.LoginToAlpine(vmi)).To(Succeed())
					} else {
						By("Starting a VirtualMachineInstance")
						createdVMI := libvmops.RunVMIAndExpectScheduling(vmi, 60)

						By(fmt.Sprintf("Checking that VirtualMachineInstance start failed: starting at %v", time.Now()))
						ctx, cancel := context.WithCancel(context.Background())
						defer cancel()
						event := watcher.New(createdVMI).Timeout(60*time.Second).SinceWatchedObjectResourceVersion().WaitFor(ctx, watcher.WarningEvent, "SyncFailed")
						Expect(event.Message).To(ContainSubstring("Could not open '/var/run/kubevirt-private/vmi-disks/disk0/disk.img': Permission denied"), "VMI should not be started")
					}
				},
					Entry("[test_id:3130]with Disk PVC", newRandomVMIWithPVC, true),
					Entry("[test_id:3131]with CDRom PVC", newRandomVMIWithCDRom, true),
					Entry("unless hostpath disk image file not owned by qemu", newRandomVMIWithPVC, false),
				)
			})

	hdc.lessPVCSpaceToleration = toleration
}

func (hdc *DiskImgCreator) Create(vmi *v1.VirtualMachineInstance) error {
	for _, volume := range vmi.Spec.Volumes {
		if hostDisk := volume.VolumeSource.HostDisk; shouldMountHostDisk(hostDisk) {
			if err := hdc.mountHostDiskAndSetOwnership(vmi, volume.Name, hostDisk); err != nil {
	}

	if fileNotExists {
		if err = hdc.handleRequestedSizeAndCreateSparseRaw(vmi, diskDir, filepath.Base(hostDisk.Path), hostDisk); err != nil {
			return err
		}

		if err != nil {
			return err
		}
		// Change file ownership to the qemu user.
		if err = ephemeraldiskutils.DefaultOwnershipManager.SetFileOwnership(diskPath); err != nil {
			log.Log.Reason(err).Errorf("Couldn't set Ownership on %s: %v", diskPath, err)
			return err
		}
	}
	return nil
}
