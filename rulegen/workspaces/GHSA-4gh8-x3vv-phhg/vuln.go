package main

// Copyright (c) 2018-2019, Sylabs Inc. All rights reserved.
// Copyright (c) 2018, Divya Cote <divya.cote@gmail.com> All rights reserved.
// Copyright (c) 2017, SingularityWare, LLC. All rights reserved.
// Copyright (c) 2017, Yannick Cote <yhcote@gmail.com> All rights reserved.

// New creates a new empty SIF file.
func New(file string) error {
	cinfo := sif.CreateInfo{
		Pathname:   file,
		Launchstr:  sif.HdrLaunch,
		Sifversion: sif.HdrVersion,
		ID:         uuid.NewV4(),
	}

	_, err := sif.CreateContainer(cinfo)
	if err != nil {
		return err
	}

	return nil
}

// AddOptions contains the options when adding a section to a SIF file.
