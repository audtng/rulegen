package main

		return fail(req, fmt.Errorf("cannot parse arguments: Binding or Action is nil"))
	}

	mangleInvalidArgumentValues(req)

	if hasExec(req) {
}

func handleShellBranch(req *ExecutionRequest) bool {
	if err := checkShellArgumentSafety(req.Binding.Action); err != nil {
		return fail(req, err)
	}
	}
}

func injectSystemArgs(req *ExecutionRequest) {
	req.Arguments["ot_executionTrackingId"] = req.TrackingID
	req.Arguments["ot_username"] = req.AuthenticatedUser.Username
		return
	}

	req := &executor.ExecutionRequest{
		Binding:           binding,
		Cfg:               h.cfg,
		Tags:              []string{"webhook"},
		Arguments:         args,
		AuthenticatedUser: auth.UserFromSystem(h.cfg, "webhook"),
	}

	h.executor.ExecRequest(req)
}
	if action.Shell == "" {
		return nil
	}
	unsafe := map[string]struct{}{"url": {}, "email": {}, "raw_string_multiline": {}, "very_dangerous_raw_string": {}}
	for _, arg := range action.Arguments {
		if _, bad := unsafe[arg.Type]; bad {
			return fmt.Errorf("unsafe argument type '%s' cannot be used with Shell execution. Use 'exec' instead. See https://docs.olivetin.app/action_execution/shellvsexec.html", arg.Type)
