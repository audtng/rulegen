package main

		// This means that we now have to automatically allow `runtime/default`
		// if a user specifies `docker/default` and vice versa in an SCC.
		if s.runtimeDefaultAllowed &&
			(p == v1.DeprecatedSeccompProfileDockerDefault ||
				p == v1.SeccompProfileRuntimeDefault) {
			return nil
		}
	}
