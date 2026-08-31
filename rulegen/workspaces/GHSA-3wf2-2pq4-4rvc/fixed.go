package main

				linterErr = multierr.Append(linterErr, err)
			}
		}
		if err := l.lintSettings(config, container, area); err != nil {
			linterErr = multierr.Append(linterErr, err)
		}
	}
	return nil
}

func (l *Linter) lintSettings(config *WorkflowConfig, c *types.Container, field string) error {
	if len(c.Settings) == 0 {
		return nil
	}
	if len(c.Commands) != 0 {
		return newLinterError("Cannot configure both commands and settings", config.File, fmt.Sprintf("%s.%s", field, c.Name), false)
	}
	if len(c.Entrypoint) != 0 {
		return newLinterError("Cannot configure both entrypoint and settings", config.File, fmt.Sprintf("%s.%s", field, c.Name), false)
	}
	if len(c.Environment) != 0 {
		return newLinterError("Cannot configure both environment and settings", config.File, fmt.Sprintf("%s.%s", field, c.Name), false)
	}
	return nil
}
}

func (c *Container) IsPlugin() bool {
	return len(c.Commands) == 0 &&
		len(c.Entrypoint) == 0 &&
		len(c.Environment) == 0
}

func (c *Container) IsTrustedCloneImage() bool {
		return err
	}

	return os.WriteFile(output, data, 0o644)
}
