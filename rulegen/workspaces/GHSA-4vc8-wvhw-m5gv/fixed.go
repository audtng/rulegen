package main

	healthHandler := srv.monitoredHandler(http.HandlerFunc(srv.healthHandler), "health")
	logStreamHandler := srv.monitoredHandler(newLogStreamEndpointHandler(httpCtxt), "logstream")
	embeddedCLIHandler := srv.monitoredHandler(newEmbeddedCLIHandler(httpCtxt), "commands")
	controllerAdminAuthorizer := controllerAdminAuthorizer{
		controllerTag: systemState.ControllerTag(),
	}
	var debuglogAuth httpcontext.CompositeAuthorizer = []authentication.Authorizer{
		tagKindAuthorizer{names.MachineTagKind, names.ControllerAgentTagKind},
		controllerAdminAuthorizer,
		modelPermissionAuthorizer{
			perm: permission.ReadAccess,
		},
		GetHandler:  modelCharmsHandler.ServeGet,
	}, "charms")
	var modelCharmsUploadAuthorizer httpcontext.CompositeAuthorizer = []authentication.Authorizer{
		controllerAdminAuthorizer,
		modelPermissionAuthorizer{
			perm: permission.WriteAccess,
		},
		ctxt:          httpCtxt,
		stateAuthFunc: httpCtxt.stateForRequestAuthenticatedUser,
	}, "tools")
	var modelToolsUploadAuthorizer httpcontext.CompositeAuthorizer = []authentication.Authorizer{
		controllerAdminAuthorizer,
		modelPermissionAuthorizer{
			perm: permission.AdminAccess,
		},
	}
	modelToolsDownloadHandler := srv.monitoredHandler(newToolsDownloadHandler(httpCtxt), "tools")
	resourcesHandler := srv.monitoredHandler(&ResourcesHandler{
		StateAuthFunc: func(req *http.Request, tagKinds ...string) (ResourcesBackend, state.PoolHelper, names.Tag,
		},
	}, "units")

	migrateCharmsHandler := &charmsHandler{
		ctxt:          httpCtxt,
		dataDir:       srv.dataDir,
	}, {
		pattern:    "/tools",
		handler:    modelToolsUploadHandler,
		authorizer: controllerAdminAuthorizer,
	}, {
		pattern:         "/tools/:version",
		handler:         modelToolsDownloadHandler,
