package main

		return fail(req, fmt.Errorf("cannot parse arguments: Binding or Action is nil"))
	}

	filterToDefinedArgumentsOnly(req)
	mangleInvalidArgumentValues(req)

	if hasExec(req) {
}

func handleShellBranch(req *ExecutionRequest) bool {
	if hasWebhookTag(req) {
		return fail(req, fmt.Errorf("webhooks cannot use Shell execution; use exec instead. See https://docs.olivetin.app/action_execution/shellvsexec.html"))
	}
	if err := checkShellArgumentSafety(req.Binding.Action); err != nil {
		return fail(req, err)
	}
	}
}

func filterToDefinedArgumentsOnly(req *ExecutionRequest) {
	definedNames := make(map[string]struct{})
	for _, arg := range req.Binding.Action.Arguments {
		definedNames[arg.Name] = struct{}{}
	}
	filtered := make(map[string]string)
	for k, v := range req.Arguments {
		if _, ok := definedNames[k]; ok || strings.HasPrefix(k, "ot_") {
			filtered[k] = v
		}
	}
	req.Arguments = filtered
}

func hasWebhookTag(req *ExecutionRequest) bool {
	for _, tag := range req.Tags {
		if tag == "webhook" {
			return true
		}
	}
	return false
}

func injectSystemArgs(req *ExecutionRequest) {
	req.Arguments["ot_executionTrackingId"] = req.TrackingID
	req.Arguments["ot_username"] = req.AuthenticatedUser.Username
		return
	}

	definedArgs := filterToDefinedArguments(args, action)
	req := &executor.ExecutionRequest{
		Binding:           binding,
		Cfg:               h.cfg,
		Tags:              []string{"webhook"},
		Arguments:         definedArgs,
		AuthenticatedUser: auth.UserFromSystem(h.cfg, "webhook"),
	}

	h.executor.ExecRequest(req)
}

func filterToDefinedArguments(args map[string]string, action *config.Action) map[string]string {
	definedNames := make(map[string]struct{})
	for _, arg := range action.Arguments {
		definedNames[arg.Name] = struct{}{}
	}
	filtered := make(map[string]string)
	for k, v := range args {
		if _, ok := definedNames[k]; ok {
			filtered[k] = v
		}
	}
	return filtered
}
	if action.Shell == "" {
		return nil
	}
	unsafe := map[string]struct{}{"url": {}, "email": {}, "raw_string_multiline": {}, "very_dangerous_raw_string": {}, "password": {}}
	for _, arg := range action.Arguments {
		if _, bad := unsafe[arg.Type]; bad {
			return fmt.Errorf("unsafe argument type '%s' cannot be used with Shell execution. Use 'exec' instead. See https://docs.olivetin.app/action_execution/shellvsexec.html", arg.Type)
