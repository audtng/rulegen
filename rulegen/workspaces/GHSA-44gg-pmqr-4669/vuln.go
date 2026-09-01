package main

	return entrypoint, args
}

func parseSecurityOpt(container *Container, config *runconfig.Config) error {
	var (
		label_opts []string
		err        error
	)

	for _, opt := range config.SecurityOpt {
		}
		switch con[0] {
		case "label":
			label_opts = append(label_opts, con[1])
		case "apparmor":
			container.AppArmorProfile = con[1]
		default:
		}
	}

	container.ProcessLabel, container.MountLabel, err = label.InitLabels(label_opts)
	return err
}

		execCommands:    newExecStore(),
	}
	container.root = daemon.containerRoot(container.ID)
	err = parseSecurityOpt(container, config)
	return container, err
}

}

func (daemon *Daemon) setHostConfig(container *Container, hostConfig *runconfig.HostConfig) error {
	// Validate the HostConfig binds. Make sure that:
	// the source exists
	for _, bind := range hostConfig.Binds {
	Entrypoint      []string
	NetworkDisabled bool
	OnBuild         []string
	SecurityOpt     []string
}

func ContainerConfigFromJob(job *engine.Job) *Config {
	}
	job.GetenvJson("ExposedPorts", &config.ExposedPorts)
	job.GetenvJson("Volumes", &config.Volumes)
	config.SecurityOpt = job.GetenvList("SecurityOpt")
	if PortSpecs := job.GetenvList("PortSpecs"); PortSpecs != nil {
		config.PortSpecs = PortSpecs
	}
