package main

	daemonCommand.Flags().Duration("tick", 3*time.Second, "Tick for polling events")
	daemonCommand.Flags().Int("vsock-port", 0, "Use vsock server instead a UNIX socket")
	daemonCommand.Flags().String("virtio-port", "", "Use virtio server instead a UNIX socket")
	daemonCommand.Flags().Int("socket-owner", 0, "UID of the main user that owns the UNIX socket (0 for root)")
	return daemonCommand
}

	if err != nil {
		return err
	}
	socketOwner, err := cmd.Flags().GetInt("socket-owner")
	if err != nil {
		return err
	}
	if tick == 0 {
		return errors.New("tick must be specified")
	}
		if err != nil {
			return err
		}
		// The daemon runs as root (for iptables), but the host connects to the
		// socket as the main user over an SSH local-forward. Hand the socket to
		// that user so 0600 restricts access to the main user rather than root.
		if socketOwner > 0 {
			if err := os.Chown(socket, socketOwner, -1); err != nil {
				return err
			}
		}
		if err := os.Chmod(socket, 0o600); err != nil {
			return err
		}
		l = socketL
	installSystemdCommand.Flags().Bool("guestagent-updated", false, "Indicate that the guest agent has been updated")
	installSystemdCommand.Flags().Int("vsock-port", 0, "Use vsock server on specified port")
	installSystemdCommand.Flags().String("virtio-port", "", "Use virtio server instead a UNIX socket")
	installSystemdCommand.Flags().Int("socket-owner", 0, "UID of the main user that owns the UNIX socket (0 for root)")
	return installSystemdCommand
}

	if err != nil {
		return err
	}
	socketOwner, err := cmd.Flags().GetInt("socket-owner")
	if err != nil {
		return err
	}
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}
	unit, err := generateSystemdUnit(vsockPort, virtioPort, socketOwner, debug)
	if err != nil {
		return err
	}
//go:embed lima-guestagent.TEMPLATE.service
var systemdUnitTemplate string

func generateSystemdUnit(vsockPort int, virtioPort string, socketOwner int, debug bool) ([]byte, error) {
	selfExeAbs, err := os.Executable()
	if err != nil {
		return nil, err
	if virtioPort != "" {
		args = append(args, fmt.Sprintf("--virtio-port %s", virtioPort))
	}
	if socketOwner > 0 {
		args = append(args, fmt.Sprintf("--socket-owner %d", socketOwner))
	}
	if debug {
		args = append(args, "--debug")
	}
