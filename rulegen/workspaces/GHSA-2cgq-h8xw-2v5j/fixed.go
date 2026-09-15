package main

					},
				},
			},
			PotentiallyUnsafeConfigAnnotations: []string{
				"bundle",
				"org.systemd.property.", // prefix form
				"org.criu.config",
			},
		}

		if seccomp.Enabled {
