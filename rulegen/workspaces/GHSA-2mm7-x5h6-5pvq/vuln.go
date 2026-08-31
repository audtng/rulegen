package main

		}
	}
	if ec.Privileged {
		if p.Capabilities == nil {
			p.Capabilities = &specs.LinuxCapabilities{}
		}
		p.Capabilities.Bounding = caps.GetAllCapabilities()
		p.Capabilities.Permitted = p.Capabilities.Bounding
		p.Capabilities.Inheritable = p.Capabilities.Bounding
		p.Capabilities.Effective = p.Capabilities.Bounding
	}
	if apparmor.HostSupports() {
		var appArmorProfile string
		Version: specs.Version,
		Process: &specs.Process{
			Capabilities: &specs.LinuxCapabilities{
				Bounding:    caps.DefaultCapabilities(),
				Permitted:   caps.DefaultCapabilities(),
				Inheritable: caps.DefaultCapabilities(),
				Effective:   caps.DefaultCapabilities(),
			},
		},
		Root: &specs.Root{},
var deviceCgroupRuleRegex = regexp.MustCompile("^([acb]) ([0-9]+|\\*):([0-9]+|\\*) ([rwm]{1,3})$")

// SetCapabilities sets the provided capabilities on the spec
// All capabilities are added if privileged is true
func SetCapabilities(s *specs.Spec, caplist []string) error {
	s.Process.Capabilities.Effective = caplist
	s.Process.Capabilities.Bounding = caplist
	s.Process.Capabilities.Permitted = caplist
	s.Process.Capabilities.Inheritable = caplist
	// setUser has already been executed here
	// if non root drop capabilities in the way execve does
	if s.Process.User.UID != 0 {
		s.Process.Capabilities.Effective = []string{}
		s.Process.Capabilities.Permitted = []string{}
	}
	return nil
}
