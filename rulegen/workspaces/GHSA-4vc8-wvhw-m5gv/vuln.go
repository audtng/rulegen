package main

	healthHandler := srv.monitoredHandler(http.HandlerFunc(srv.healthHandler), "health")
	logStreamHandler := srv.monitoredHandler(newLogStreamEndpointHandler(httpCtxt), "logstream")
	embeddedCLIHandler := srv.monitoredHandler(newEmbeddedCLIHandler(httpCtxt), "commands")
	var debuglogAuth httpcontext.CompositeAuthorizer = []authentication.Authorizer{
		tagKindAuthorizer{names.MachineTagKind, names.ControllerAgentTagKind},
		controllerAdminAuthorizer{
			controllerTag: systemState.ControllerTag(),
		},
		modelPermissionAuthorizer{
			perm: permission.ReadAccess,
		},
		GetHandler:  modelCharmsHandler.ServeGet,
	}, "charms")
	var modelCharmsUploadAuthorizer httpcontext.CompositeAuthorizer = []authentication.Authorizer{
		controllerAdminAuthorizer{
			controllerTag: systemState.ControllerTag(),
		},
		modelPermissionAuthorizer{
			perm: permission.WriteAccess,
		},
		ctxt:          httpCtxt,
		stateAuthFunc: httpCtxt.stateForRequestAuthenticatedUser,
	}, "tools")
	modelToolsUploadAuthorizer := tagKindAuthorizer{names.UserTagKind}
	modelToolsDownloadHandler := srv.monitoredHandler(newToolsDownloadHandler(httpCtxt), "tools")
	resourcesHandler := srv.monitoredHandler(&ResourcesHandler{
		StateAuthFunc: func(req *http.Request, tagKinds ...string) (ResourcesBackend, state.PoolHelper, names.Tag,
		},
	}, "units")

	controllerAdminAuthorizer := controllerAdminAuthorizer{
		controllerTag: systemState.ControllerTag(),
	}
	migrateCharmsHandler := &charmsHandler{
		ctxt:          httpCtxt,
		dataDir:       srv.dataDir,
	}, {
		pattern:    "/tools",
		handler:    modelToolsUploadHandler,
		authorizer: modelToolsUploadAuthorizer,
	}, {
		pattern:         "/tools/:version",
		handler:         modelToolsDownloadHandler,
