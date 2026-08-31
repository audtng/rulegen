package main

				linterErr = multierr.Append(linterErr, err)
			}
		}
		if err := l.lintCommands(config, container, area); err != nil {
			linterErr = multierr.Append(linterErr, err)
		}
	}
	return nil
}

func (l *Linter) lintCommands(config *WorkflowConfig, c *types.Container, field string) error {
	if len(c.Commands) == 0 {
		return nil
	}
	if len(c.Settings) != 0 {
		var keys []string
		for key := range c.Settings {
			keys = append(keys, key)
		}
		return newLinterError(fmt.Sprintf("Cannot configure both commands and custom attributes %v", keys), config.File, fmt.Sprintf("%s.%s", field, c.Name), false)
	}
	return nil
}
}

func (c *Container) IsPlugin() bool {
	return len(c.Commands) == 0 && len(c.Entrypoint) == 0
}

func (c *Container) IsTrustedCloneImage() bool {
		return err
	}

	return os.WriteFile(output, data, 0644)
}
