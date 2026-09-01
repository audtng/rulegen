package main


import (
	"fmt"
	"slices"

	"github.com/go-vela/cli/action"
	"github.com/go-vela/cli/action/secret"
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_COMMAND", "SECRET_COMMAND"},
			Name:    internal.FlagSecretCommands,
			Aliases: []string{"c"},
			Usage:   "enable a secret to be used for a step with commands (default is false for shared secrets)",
			Value:   "true",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SUBSTITUTION", "SECRET_SUBSTITUTION"},
			Name:    internal.FlagSecretSubstitution,
			Aliases: []string{"s"},
			Usage:   "enable a secret to be substituted (default is false for shared secrets)",
			Value:   "true",
		},
		&cli.StringFlag{
	},
	CustomHelpTemplate: fmt.Sprintf(`%s
EXAMPLES:
   1. Add a repository secret.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --value bar
   2. Add a repository secret and disallow usage in commands.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --value bar --commands false
   3. Add an organization secret.
     $ {{.HelpName}} --secret.engine native --secret.type org --org MyOrg --name foo --value bar
   4. Add a shared secret.
     $ {{.HelpName}} --secret.engine native --secret.type shared --org MyOrg --team octokitties --name foo --value bar
   5. Add a repository secret with all event types enabled.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --value bar --event comment --event deployment --event pull_request --event push --event tag
   6. Add a repository secret with an image whitelist.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --value bar --image alpine --image golang:* --image postgres:latest
   7. Add a secret with value from a file.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --value @secret.txt
   8. Add a repository secret with json output.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --value bar --output json
   9. Add a secret or secrets from a file.
     $ {{.HelpName}} --file secret.yml
  10. Add a secret when config or environment variables are set.
     $ {{.HelpName}} --org MyOrg --repo MyRepo --name foo --value bar

DOCUMENTATION:

	//
	// https://pkg.go.dev/github.com/go-vela/cli/action/secret?tab=doc#Config
	s := &secret.Config{
		Action: internal.ActionAdd,
		Engine: c.String(internal.FlagSecretEngine),
		Type:   c.String(internal.FlagSecretType),
		Org:    c.String(internal.FlagOrg),
		Repo:   c.String(internal.FlagRepo),
		Team:   c.String("team"),
		Name:   c.String("name"),
		Value:  c.String("value"),
		Images: c.StringSlice("image"),
		Events: c.StringSlice("event"),
		File:   c.String("file"),
		Output: c.String(internal.FlagOutput),
	}

	// check if allow_command and allow_substitution are provided
	// if they are not, server will not update the fields
	if slices.Contains(c.FlagNames(), internal.FlagSecretCommands) {
		val := c.Bool(internal.FlagSecretCommands)
		s.AllowCommand = &val
	}

	if slices.Contains(c.FlagNames(), internal.FlagSecretSubstitution) {
		val := c.Bool(internal.FlagSecretSubstitution)
		s.AllowSubstitution = &val
	}

	// validate secret configuration

import (
	"fmt"
	"slices"

	"github.com/go-vela/cli/action"
	"github.com/go-vela/cli/action/secret"
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_COMMAND", "SECRET_COMMAND"},
			Name:    internal.FlagSecretCommands,
			Aliases: []string{"c"},
			Usage:   "enable a secret to be used for a step with commands",
			Value:   "true",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SUBSTITUTION", "SECRET_SUBSTITUTION"},
			Name:    internal.FlagSecretSubstitution,
			Aliases: []string{"s"},
			Usage:   "enable a secret to be substituted",
			Value:   "true",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_FILE", "SECRET_FILE"},
			Name:    "file",
	},
	CustomHelpTemplate: fmt.Sprintf(`%s
EXAMPLES:
   1. Update a repository secret.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --value bar
   2. Update a repository secret and disallow usage in commands.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --value bar --commands false
   3. Update an organization secret.
     $ {{.HelpName}} --secret.engine native --secret.type org --org MyOrg --name foo --value bar
   4. Update a shared secret.
     $ {{.HelpName}} --secret.engine native --secret.type shared --org MyOrg --team octokitties --name foo --value bar
   5. Update a repository secret with all event types enabled.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --event comment --event deployment --event pull_request --event push --event tag
   6. Update a repository secret with an image whitelist.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --image alpine --image golang:* --image postgres:latest
   7. Update a secret with value from a file.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --value @secret.txt
   8. Update a repository secret with json output.
     $ {{.HelpName}} --secret.engine native --secret.type repo --org MyOrg --repo MyRepo --name foo --value bar --output json
   9. Update a secret or secrets from a file.
     $ {{.HelpName}} --file secret.yml
  10. Update a secret when config or environment variables are set.
     $ {{.HelpName}} --org MyOrg --repo MyRepo --name foo --value bar

DOCUMENTATION:

	//
	// https://pkg.go.dev/github.com/go-vela/cli/action/secret?tab=doc#Config
	s := &secret.Config{
		Action: internal.ActionUpdate,
		Engine: c.String(internal.FlagSecretEngine),
		Type:   c.String(internal.FlagSecretType),
		Org:    c.String(internal.FlagOrg),
		Repo:   c.String(internal.FlagRepo),
		Team:   c.String("team"),
		Name:   c.String("name"),
		Value:  c.String("value"),
		Images: c.StringSlice("image"),
		Events: c.StringSlice("event"),
		File:   c.String("file"),
		Output: c.String(internal.FlagOutput),
	}

	// check if allow_command and allow_substitution are provided
	// if they are not, server will not update the fields
	if slices.Contains(c.FlagNames(), internal.FlagSecretCommands) {
		val := c.Bool(internal.FlagSecretCommands)
		s.AllowCommand = &val
	}

	if slices.Contains(c.FlagNames(), internal.FlagSecretSubstitution) {
		val := c.Bool(internal.FlagSecretSubstitution)
		s.AllowSubstitution = &val
	}

	// validate secret configuration
	//
	// https://pkg.go.dev/github.com/go-vela/types/library?tab=doc#Secret
	s := &library.Secret{
		Type:              &c.Type,
		Org:               &c.Org,
		Repo:              &c.Repo,
		Team:              &c.Team,
		Name:              &c.Name,
		Value:             &c.Value,
		Images:            &c.Images,
		Events:            &c.Events,
		AllowCommand:      c.AllowCommand,
		AllowSubstitution: c.AllowSubstitution,
	}

	logrus.Tracef("adding secret %s/%s/%s/%s/%s", c.Engine, c.Type, c.Org, name, c.Name)
			//
			// https://pkg.go.dev/github.com/go-vela/cli/action/secret?tab=doc#Config
			s := &Config{
				Action:            "add",
				Engine:            f.Metadata.Engine,
				Type:              s.GetType(),
				Org:               s.GetOrg(),
				Repo:              s.GetRepo(),
				Team:              s.GetTeam(),
				Name:              s.GetName(),
				Value:             s.GetValue(),
				Images:            s.GetImages(),
				Events:            s.GetEvents(),
				AllowCommand:      s.AllowCommand,
				AllowSubstitution: s.AllowSubstitution,
				Output:            c.Output,
			}

			// validate secret configuration
