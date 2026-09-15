package main

)

// loadPayloads loads the input payloads from a map to a data map
func (generator *PayloadGenerator) loadPayloads(payloads map[string]interface{}, templatePath, templateDirectory string, sandbox bool) (map[string][]string, error) {
	loadedPayloads := make(map[string][]string)

	for name, payload := range payloads {
			if len(elements) >= 2 {
				loadedPayloads[name] = elements
			} else {
				if sandbox {
					pt = filepath.Clean(pt)
					templatePathDir := filepath.Dir(templatePath)
					if !(templatePathDir != "/" && strings.HasPrefix(pt, templatePathDir)) && !strings.HasPrefix(pt, templateDirectory) {
						return nil, errors.New("denied payload file path specified")
					}
	ClientCAFile string
	// Deprecated: Use ZTLS library
	ZTLS bool
	// Sandbox enables sandboxed nuclei template execution
	Sandbox bool
	// ShowMatchLine enables display of match line number
	ShowMatchLine bool
	// EnablePprof enables exposing pprof runtime information with a webserver.
}

// New creates a new generator structure for payload generation
func New(payloads map[string]interface{}, attackType AttackType, templatePath string, sandbox bool, catalog catalog.Catalog, customAttackType string) (*PayloadGenerator, error) {
	if attackType.String() == "" {
		attackType = BatteringRamAttack
	}
		return nil, err
	}

	compiled, err := generator.loadPayloads(payloadsFinal, templatePath, config.DefaultConfig.TemplatesDirectory, sandbox)
	if err != nil {
		return nil, err
	}
