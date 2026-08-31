package main

	return fmt.Sprintf("Cannot write tuple '%s'. Reason: %s", i.TupleKey, i.Reason)
}

// ValidateUser returns whether the user is valid.  If not, return error
func ValidateUser(tk *openfgapb.TupleKey) error {
	if !IsValidUser(tk.GetUser()) {
		return &InvalidTupleError{Reason: "the 'user' field is invalid", TupleKey: tk}

func validateRelationTypeRestrictions(model *openfgapb.AuthorizationModel) error {
	t := New(model)

	for objectType := range t.typeDefinitions {
		relations, err := t.GetRelations(objectType)
					if _, err := t.GetRelation(relatedObjectType, relatedRelation); err != nil {
						return InvalidRelationTypeError(objectType, name, relatedObjectType, relatedRelation)
					}
				}
			}
		}

	return fmt.Errorf("the relation type '%s' on '%s' in object type '%s' is not valid", relationType, relation, objectType)
}
			},
			err: NonAssignableRelationError("document", "reader"),
		},
	}

	for _, test := range tests {
	}, nil
}

func (query *CheckQuery) getTypeDefinitionRelationUsersets(ctx context.Context, rc *resolutionContext) (*openfgapb.Userset, error) {
	ctx, span := query.tracer.Start(ctx, "getTypeDefinitionRelationUsersets")
	defer span.End()
	return nil
}

func (query *CheckQuery) resolveTupleToUserset(
	ctx context.Context,
	rc *resolutionContext,
		relation = rc.tk.GetRelation()
	}

	findTK := tupleUtils.NewTupleKey(rc.tk.GetObject(), relation, "")

	tracer := rc.tracer.AppendTupleToUserset().AppendString(tupleUtils.ToObjectRelationString(findTK.GetObject(), relation))
	iter, err := rc.read(ctx, query.datastore, findTK)
			break // the user was resolved already, avoid launching extra lookups
		}

		userObj, userRel := tupleUtils.SplitObjectRelation(tuple.GetUser())

		if userObj == Wildcard {
			objectType, _ := tupleUtils.SplitObject(rc.tk.GetObject())

			query.logger.WarnWithContext(
				ctx,
				fmt.Sprintf("unexpected wildcard evaluated on tupleset relation '%s'", relation),
				zap.String("store_id", rc.store),
				zap.String("authorization_model_id", rc.modelID),
				zap.String("object_type", objectType),
			)
		}

		if !tupleUtils.IsValidObject(userObj) {
			continue // TupleToUserset tuplesets should be of the form 'objectType:id' or 'objectType:id#relation' but are not guaranteed to be because it is neither a user or userset
		}

		usersetRel := node.TupleToUserset.GetComputedUserset().GetRelation()

		// userRel may be empty, and in this case we set it to usersetRel.
		if userRel == "" {
			userRel = usersetRel
		}
		// We only proceed in the case that userRel == usersetRel (=node.TupleToUserset.GetComputedUserset().GetRelation()).
		if userRel != usersetRel {
			continue
		}

		tupleKey := &openfgapb.TupleKey{
			// user from previous lookup
			Object:   userObj,
			Relation: userRel,
			// original tk user
			User: rc.tk.GetUser(),
		}
		tracer := tracer.AppendString(tupleUtils.ToObjectRelationString(userObj, userRel))
		nestedRC := rc.fork(tupleKey, tracer, false)
		go func(c chan<- *chanResolveResult) {
			defer wg.Done()

			userset, err := query.getTypeDefinitionRelationUsersets(ctx, nestedRC)
			if err == nil {
				err = query.resolveNode(ctx, nestedRC, userset, typesys)
			}
			continue
		}

		tObject, tRelation := tupleUtils.SplitObjectRelation(user)
		// We only proceed in the case that tRelation == userset.GetComputedUserset().GetRelation().
		// tRelation may be empty, and in this case, we set it to userset.GetComputedUserset().GetRelation().
			tRelation = userset.GetComputedUserset().GetRelation()
		}

		if tRelation != userset.GetComputedUserset().GetRelation() {
			continue
		}

		cs := &openfgapb.TupleKey{
			Object:   tObject,
			Relation: tRelation,
			return serverErrors.HandleTupleValidateError(&tupleUtils.IndirectWriteError{Reason: IndirectWriteErrorReason, TupleKey: tk})
		}

		if err := c.validateTypesForTuple(authModel, tk); err != nil {
			return err
		}
	return nil
}

// validateTypesForTuple makes sure that when writing a tuple, the types are compatible.
// 1. If the tuple is of the form (user=person:bob, relation=reader, object=doc:budget), then the type "doc", relation "reader" must allow type "person".
// 2. If the tuple is of the form (user=group:abc#member, relation=reader, object=doc:budget), then the type "doc", relation "reader" must allow type "group", relation "member".
			},
		},
		{
			name: "ExecuteReturnsAllowedForTupleToUserset",
			typeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "repo",
											Relation: "manager",
										},
										ComputedUserset: &openfgapb.ObjectRelation{
											Object:   "$TUPLE_USERSET_OBJECT",
											Relation: "repo_admin",
										},
									}}},
					},
				},
				{
					Type: "org",
					Relations: map[string]*openfgapb.Userset{
						// implicit direct?
						"repo_admin": {},
					},
				},
			},
			tuples: []*openfgapb.TupleKey{
				tuple.NewTupleKey("repo:openfga/canaveral", "manager", "org:openfga#repo_admin"),
				tuple.NewTupleKey("org:openfga", "repo_admin", "github|jose@openfga"),
			},
			resolveNodeLimit: defaultResolveNodeLimit,
			request: &openfgapb.CheckRequest{
				TupleKey: tuple.NewTupleKey("repo:openfga/canaveral", "admin", "github|jose@openfga"),
				Trace:    true,
			},
			response: &openfgapb.CheckResponse{
				Allowed:    true,
				Resolution: ".union.1(tuple-to-userset).repo:openfga/canaveral#manager.org:openfga#repo_admin.(direct).",
			},
		},
		{
			name: "ExecuteCanResolveRecursiveComputedUserSets",
				errors.New("unexpected rewrite on relation 'document#parent'"),
			),
		},
	}

	ctx := context.Background()
				},
			},
		},
	}

	require := require.New(t)
				errors.Errorf("unexpected rewrite on relation '%s#%s'", "document", "parent"),
			),
		},
	}

	require := require.New(t)
		},
		err: serverErrors.NewInternalError("invalid authorization model", errors.New("invalid authorization model")),
	},
}

func TestWriteCommand(t *testing.T, datastore storage.OpenFGADatastore) {
		err     error
	}{
		{
			name: "succeeds",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			},
		},
		{
			name: "succeeds part II",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: "somestoreid",
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.ExceededEntityLimit("type definitions in an authorization model", datastore.MaxTypesInTypeDefinition()),
		},
		{
			name: "empty relations is valid",
			request: &openfgapb.WriteAuthorizationModelRequest{
				TypeDefinitions: []*openfgapb.TypeDefinition{
					{
			},
		},
		{
			name: "zero length relations is valid",
			request: &openfgapb.WriteAuthorizationModelRequest{
				TypeDefinitions: []*openfgapb.TypeDefinition{
					{
			},
		},
		{
			name: "ExecuteWriteFailsIfSameTypeTwice",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(typesystem.ErrDuplicateTypes),
		},
		{
			name: "ExecuteWriteFailsIfEmptyRewrites",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "ExecuteWriteFailsIfUnknownRelationInComputedUserset",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "ExecuteWriteFailsIfUnknownRelationInTupleToUserset",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "", Relation: "owner"}),
		},
		{
			name: "ExecuteWriteFailsIfUnknownRelationInUnion",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "ExecuteWriteFailsIfUnknownRelationInDifferenceBaseArgument",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "ExecuteWriteFailsIfUnknownRelationInDifferenceSubtractArgument",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "ExecuteWriteFailsIfUnknownRelationInTupleToUsersetTupleset",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "ExecuteWriteFailsIfUnknownRelationInTupleToUsersetComputedUserset",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
					{
						Type: "repo",
						Relations: map[string]*openfgapb.Userset{
							"writer": {
								Userset: &openfgapb.Userset_This{},
							},
							"viewer": {
								Userset: &openfgapb.Userset_TupleToUserset{
									TupleToUserset: &openfgapb.TupleToUserset{
										Tupleset: &openfgapb.ObjectRelation{
											Relation: "writer",
										},
										ComputedUserset: &openfgapb.ObjectRelation{
											Relation: "owner",
										},
									},
								},
							},
						},
					},
				},
			},
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "", Relation: "owner"}),
		},
		{
			name: "ExecuteWriteFailsIfTupleToUsersetReferencesUnknownRelation",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "bar", Relation: "writer"}),
		},
		{
			name: "ExecuteWriteFailsIfUnknownRelationInIntersection",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "ExecuteWriteFailsIfDifferenceIncludesSameRelationTwice",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{ObjectType: "repo", Relation: "viewer"}),
		},
		{
			name: "ExecuteWriteFailsIfUnionIncludesSameRelationTwice",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{ObjectType: "repo", Relation: "viewer"}),
		},
		{
			name: "ExecuteWriteFailsIfIntersectionIncludesSameRelationTwice",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{ObjectType: "repo", Relation: "viewer"}),
		},
		{
			name: "Union Rewrite Contains Repeated Definitions",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			},
		},
		{
			name: "Intersection Rewrite Contains Repeated Definitions",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			},
		},
		{
			name: "Exclusion Rewrite Contains Repeated Definitions",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			},
		},
		{
			name: "Tupleset relation involves ComputedUserset rewrite",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			),
		},
		{
			name: "Tupleset relation involves Union rewrite",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			),
		},
		{
			name: "Tupleset relation involves Intersection rewrite",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			),
		},
		{
			name: "Tupleset relation involves Exclusion rewrite",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			),
		},
		{
			name: "Tupleset relation involves TupleToUserset rewrite",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
	openfgapb "go.buf.build/openfga/go/openfga/api/openfga/v1"
)

// ValidateTuple returns whether a *openfgapb.TupleKey is valid
func ValidateTuple(ctx context.Context, backend storage.TypeDefinitionReadBackend, store, authorizationModelID string, tk *openfgapb.TupleKey) (*openfgapb.Userset, error) {
	if err := tuple.ValidateUser(tk); err != nil {
		return nil, err
	return ValidateObjectsRelations(ctx, backend, store, authorizationModelID, tk)
}

// ValidateObjectsRelations returns whether a tuple's object and relations are valid
func ValidateObjectsRelations(ctx context.Context, backend storage.TypeDefinitionReadBackend, store, modelID string, t *openfgapb.TupleKey) (*openfgapb.Userset, error) {
	if !tuple.IsValidRelation(t.GetRelation()) {
		return nil, &tuple.InvalidTupleError{Reason: "invalid relation", TupleKey: t}
