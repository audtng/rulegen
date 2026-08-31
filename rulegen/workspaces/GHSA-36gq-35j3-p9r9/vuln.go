package main

			return err
		}

		switch opts.Format {
		case "json":
			content, err = project.MarshalJSON()
		options.Services = project.ServiceNames()
	}

	var observedState Containers
	observedState, err := s.getContainers(ctx, project.Name, oneOffInclude, true)
	if err != nil {
		return err
	}
