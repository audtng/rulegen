package main

	return entrypoint, args
}

func parseSecurityOpt(container *Container, config *runconfig.HostConfig) error {
	var (
		labelOpts []string
		err       error
	)

	for _, opt := range config.SecurityOpt {
		}
		switch con[0] {
		case "label":
			labelOpts = append(labelOpts, con[1])
		case "apparmor":
			container.AppArmorProfile = con[1]
		default:
		}
	}

	container.ProcessLabel, container.MountLabel, err = label.InitLabels(labelOpts)
	return err
}

		execCommands:    newExecStore(),
	}
	container.root = daemon.containerRoot(container.ID)
	return container, err
}

}

func (daemon *Daemon) setHostConfig(container *Container, hostConfig *runconfig.HostConfig) error {
	if err := parseSecurityOpt(container, hostConfig); err != nil {
		return err
	}
	// Validate the HostConfig binds. Make sure that:
	// the source exists
	for _, bind := range hostConfig.Binds {
	Entrypoint      []string
	NetworkDisabled bool
	OnBuild         []string
}

func ContainerConfigFromJob(job *engine.Job) *Config {
	}
	job.GetenvJson("ExposedPorts", &config.ExposedPorts)
	job.GetenvJson("Volumes", &config.Volumes)
	if PortSpecs := job.GetenvList("PortSpecs"); PortSpecs != nil {
		config.PortSpecs = PortSpecs
	}
