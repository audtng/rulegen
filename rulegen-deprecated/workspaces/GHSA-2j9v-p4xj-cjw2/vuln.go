package main

	daemonCommand.Flags().Duration("tick", 3*time.Second, "Tick for polling events")
	daemonCommand.Flags().Int("vsock-port", 0, "Use vsock server instead a UNIX socket")
	daemonCommand.Flags().String("virtio-port", "", "Use virtio server instead a UNIX socket")
	return daemonCommand
}

	if err != nil {
		return err
	}
	if tick == 0 {
		return errors.New("tick must be specified")
	}
		if err != nil {
			return err
		}
		if err := os.Chmod(socket, 0o777); err != nil {
			return err
		}
		l = socketL
	installSystemdCommand.Flags().Bool("guestagent-updated", false, "Indicate that the guest agent has been updated")
	installSystemdCommand.Flags().Int("vsock-port", 0, "Use vsock server on specified port")
	installSystemdCommand.Flags().String("virtio-port", "", "Use virtio server instead a UNIX socket")
	return installSystemdCommand
}

	if err != nil {
		return err
	}
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}
	unit, err := generateSystemdUnit(vsockPort, virtioPort, debug)
	if err != nil {
		return err
	}
//go:embed lima-guestagent.TEMPLATE.service
var systemdUnitTemplate string

func generateSystemdUnit(vsockPort int, virtioPort string, debug bool) ([]byte, error) {
	selfExeAbs, err := os.Executable()
	if err != nil {
		return nil, err
	if virtioPort != "" {
		args = append(args, fmt.Sprintf("--virtio-port %s", virtioPort))
	}
	if debug {
		args = append(args, "--debug")
	}
