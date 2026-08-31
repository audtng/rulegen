package main

		}
	}
	if ec.Privileged {
		p.Capabilities = &specs.LinuxCapabilities{
			Bounding:  caps.GetAllCapabilities(),
			Permitted: caps.GetAllCapabilities(),
			Effective: caps.GetAllCapabilities(),
		}
	}
	if apparmor.HostSupports() {
		var appArmorProfile string
		Version: specs.Version,
		Process: &specs.Process{
			Capabilities: &specs.LinuxCapabilities{
				Bounding:  caps.DefaultCapabilities(),
				Permitted: caps.DefaultCapabilities(),
				Effective: caps.DefaultCapabilities(),
			},
		},
		Root: &specs.Root{},
var deviceCgroupRuleRegex = regexp.MustCompile("^([acb]) ([0-9]+|\\*):([0-9]+|\\*) ([rwm]{1,3})$")

// SetCapabilities sets the provided capabilities on the spec
// All capabilities are added if privileged is true.
func SetCapabilities(s *specs.Spec, caplist []string) error {
	// setUser has already been executed here
	if s.Process.User.UID == 0 {
		s.Process.Capabilities = &specs.LinuxCapabilities{
			Effective: caplist,
			Bounding:  caplist,
			Permitted: caplist,
		}
	} else {
		// Do not set Effective and Permitted capabilities for non-root users,
		// to match what execve does.
		s.Process.Capabilities = &specs.LinuxCapabilities{
			Bounding: caplist,
		}
	}
	return nil
}
