package main

	forwardPath bool
	// templates
	Templates map[string]*pongo2.Template
}

// New returns a Django render engine for Fiber
			LayoutName: "embed",
			Funcmap:    make(map[string]interface{}),
		},
	}
	return engine
}
			LayoutName: "embed",
			Funcmap:    make(map[string]interface{}),
		},
	}
	return engine
}
			LayoutName: "embed",
			Funcmap:    make(map[string]interface{}),
		},
		forwardPath: true,
	}
	return engine
	pongoset := pongo2.NewSet("default", pongoloader)
	// Set template settings
	pongoset.Globals.Update(e.Funcmap)
	// Enable autoescaping
	pongo2.SetAutoescape(true)

	// Loop trough each Directory and register template files
	walkFn := func(path string, info os.FileInfo, err error) error {
	return true
}

// Render will render the template by name
func (e *Engine) Render(out io.Writer, name string, binding interface{}, layout ...string) error {
	if !e.Loaded || e.ShouldReload {
escaping
django/README.md      | 27 +++++++++++++++++++++++----
django/django_test.go | 16 ++++++++++++++++
2 files changed, 39 insertions(+), 4 deletions(-)
	core.Engine
	// forward the base path to the template Engine
	forwardPath bool
	// templates
	Templates map[string]*pongo2.Template
	// set auto escape globally
	AutoEscape bool
}

// New returns a Django render engine for Fiber
			LayoutName: "embed",
			Funcmap:    make(map[string]interface{}),
		},
		AutoEscape: true,
	}
	return engine
}
			LayoutName: "embed",
			Funcmap:    make(map[string]interface{}),
		},
		AutoEscape: true,
	}
	return engine
}
			LayoutName: "embed",
			Funcmap:    make(map[string]interface{}),
		},
		AutoEscape: true,
		forwardPath: true,
	}
	return engine
	// Set template settings
	pongoset.Globals.Update(e.Funcmap)
	// Set autoescaping
	pongo2.SetAutoescape(e.AutoEscape)

	// Loop trough each Directory and register template files
	walkFn := func(path string, info os.FileInfo, err error) error {

// SetAutoEscape sets the auto-escape property of the template engine
func (e *Engine) SetAutoEscape(autoEscape bool) {
	e.AutoEscape = autoEscape
}

// Render will render the template by name
		if bind == nil {
			bind = make(map[string]interface{}, 1)
		}
		bind[e.LayoutName] = parsed
		lay := e.Templates[layout[0]]
		if lay == nil {
			return fmt.Errorf("LayoutName %s does not exist", layout[0])
