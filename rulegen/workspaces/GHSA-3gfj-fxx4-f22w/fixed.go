package main

	return fmt.Sprintf("Cannot write tuple '%s'. Reason: %s", i.TupleKey, i.Reason)
}

func ValidateUser(tk *openfgapb.TupleKey) error {
	if !IsValidUser(tk.GetUser()) {
		return &InvalidTupleError{Reason: "the 'user' field is invalid", TupleKey: tk}

func validateRelationTypeRestrictions(model *openfgapb.AuthorizationModel) error {
	t := New(model)
	allTupleToUsersetDefinitions := t.GetAllTupleToUsersetsDefinitions()

	for objectType := range t.typeDefinitions {
		relations, err := t.GetRelations(objectType)
					if _, err := t.GetRelation(relatedObjectType, relatedRelation); err != nil {
						return InvalidRelationTypeError(objectType, name, relatedObjectType, relatedRelation)
					}

					// you cannot specify a userset if the relation is being used in a `x from y` definition (in the `y` part)
					for _, arrayOfTtus := range allTupleToUsersetDefinitions[objectType] {
						for _, tupleToUserSetDef := range arrayOfTtus {
							if tupleToUserSetDef.Tupleset.Relation == name {
								return &InvalidRelationError{ObjectType: objectType, Relation: name}
							}
						}
					}
				}
			}
		}

	return fmt.Errorf("the relation type '%s' on '%s' in object type '%s' is not valid", relationType, relation, objectType)
}

// GetAllTupleToUsersetsDefinitions returns a map where the key is the object type and the value
// is another map where key=relationName, value=list of tuple to usersets declared in that relation
func (t *TypeSystem) GetAllTupleToUsersetsDefinitions() map[string]map[string][]*openfgapb.TupleToUserset {
	response := make(map[string]map[string][]*openfgapb.TupleToUserset, 0)
	for typeName, typeDef := range t.GetTypeDefinitions() {
		response[typeName] = make(map[string][]*openfgapb.TupleToUserset, 0)
		for relationName, relationDef := range typeDef.GetRelations() {
			ttus := make([]*openfgapb.TupleToUserset, 0)
			response[typeName][relationName] = t.getAllTupleToUsersetsDefinitions(relationDef, &ttus)
		}
	}
	return response
}

func (t *TypeSystem) getAllTupleToUsersetsDefinitions(relationDef *openfgapb.Userset, resp *[]*openfgapb.TupleToUserset) []*openfgapb.TupleToUserset {
	if relationDef.GetTupleToUserset() != nil {
		*resp = append(*resp, relationDef.GetTupleToUserset())
	}
	if relationDef.GetUnion() != nil {
		for _, child := range relationDef.GetUnion().GetChild() {
			t.getAllTupleToUsersetsDefinitions(child, resp)
		}
	}
	if relationDef.GetIntersection() != nil {
		for _, child := range relationDef.GetIntersection().GetChild() {
			t.getAllTupleToUsersetsDefinitions(child, resp)
		}
	}
	if relationDef.GetDifference() != nil {
		t.getAllTupleToUsersetsDefinitions(relationDef.GetDifference().GetBase(), resp)
		t.getAllTupleToUsersetsDefinitions(relationDef.GetDifference().GetSubtract(), resp)
	}
	return *resp
}
			},
			err: NonAssignableRelationError("document", "reader"),
		},
		{
			name: "userset specified as allowed type, but the relation is used in a TTU rewrite",
			model: &openfgapb.AuthorizationModel{
				SchemaVersion: SchemaVersion1_1,
				TypeDefinitions: []*openfgapb.TypeDefinition{
					{
						Type: "user",
					},
					{
						Type: "folder",
						Relations: map[string]*openfgapb.Userset{
							"member": This(),
						},
						Metadata: &openfgapb.Metadata{
							Relations: map[string]*openfgapb.RelationMetadata{
								"member": {
									DirectlyRelatedUserTypes: []*openfgapb.RelationReference{
										{
											Type: "user",
										},
									},
								},
							},
						},
					},
					{
						Type: "document",
						Relations: map[string]*openfgapb.Userset{
							"parent":   This(),
							"can_view": TupleToUserset("parent", "member"),
						},
						Metadata: &openfgapb.Metadata{
							Relations: map[string]*openfgapb.RelationMetadata{
								"parent": {
									DirectlyRelatedUserTypes: []*openfgapb.RelationReference{
										{
											Type:     "folder",
											Relation: "member", //this isn't allowed
										},
									},
								},
							},
						},
					},
				},
			},
			err: &InvalidRelationError{ObjectType: "document", Relation: "parent"},
		},
		{
			name: "userset specified as allowed type, but the relation is used in a TTU rewrite included in a union",
			model: &openfgapb.AuthorizationModel{
				SchemaVersion: SchemaVersion1_1,
				TypeDefinitions: []*openfgapb.TypeDefinition{
					{
						Type: "user",
					},
					{
						Type: "folder",
						Relations: map[string]*openfgapb.Userset{
							"parent": This(),
							"viewer": This(),
						},
						Metadata: &openfgapb.Metadata{
							Relations: map[string]*openfgapb.RelationMetadata{
								"parent": {
									DirectlyRelatedUserTypes: []*openfgapb.RelationReference{
										{
											Type: "folder",
										},
									},
								},
								"viewer": {
									DirectlyRelatedUserTypes: []*openfgapb.RelationReference{
										{
											Type: "user",
										},
									},
								},
							},
						},
					},
					{
						Type: "document",
						Relations: map[string]*openfgapb.Userset{
							"parent": This(),
							"viewer": Union(TupleToUserset("parent", "viewer"), This()),
						},
						Metadata: &openfgapb.Metadata{
							Relations: map[string]*openfgapb.RelationMetadata{
								"parent": {
									DirectlyRelatedUserTypes: []*openfgapb.RelationReference{
										{
											Type:     "folder",
											Relation: "parent", // this isn't allowed
										},
									},
								},
								"viewer": {
									DirectlyRelatedUserTypes: []*openfgapb.RelationReference{
										{
											Type:     "folder",
											Relation: "parent",
										},
										{
											Type: "user",
										},
									},
								},
							},
						},
					},
				},
			},
			err: &InvalidRelationError{ObjectType: "document", Relation: "parent"},
		},
	}

	for _, test := range tests {
	}, nil
}

// getTypeDefinitionRelationUsersets validates a tuple and returns the userset corresponding to the "object" and "relation"
func (query *CheckQuery) getTypeDefinitionRelationUsersets(ctx context.Context, rc *resolutionContext) (*openfgapb.Userset, error) {
	ctx, span := query.tracer.Start(ctx, "getTypeDefinitionRelationUsersets")
	defer span.End()
	return nil
}

// Given this auth model:
//
//	type document
//	  relations
//	    define parent as self
//	    define viewer as reader from parent
//	type folder
//	  relations
//	    define reader as self
//
// and this rc.tk:
//
//	document:budget#viewer@anne
//
// and these tuples:
//
//	folder:budgets#reader@anne
//	document:budget#parent@folder:budget
//
// resolveTupleToUserset first finds all the entities that are related to "document:budget" via the "parent" relation
// and then, for each of those (in this case "folder:budgets"), checks the rc.tk.User (anne) against the "reader" relation of that entity
func (query *CheckQuery) resolveTupleToUserset(
	ctx context.Context,
	rc *resolutionContext,
		relation = rc.tk.GetRelation()
	}

	findTK := tupleUtils.NewTupleKey(rc.tk.GetObject(), relation, "") //findTk=document:budget#parent@

	tracer := rc.tracer.AppendTupleToUserset().AppendString(tupleUtils.ToObjectRelationString(findTK.GetObject(), relation))
	iter, err := rc.read(ctx, query.datastore, findTK)
			break // the user was resolved already, avoid launching extra lookups
		}

		userObj, userRel := tupleUtils.SplitObjectRelation(tuple.GetUser()) // userObj=folder:budgets, userRel=""

		objectType, _ := tupleUtils.SplitObject(rc.tk.GetObject())

		if userObj == Wildcard {

			query.logger.WarnWithContext(
				ctx,
				fmt.Sprintf("unexpected wildcard evaluated on tupleset relation '%s#%s'", objectType, relation),
				zap.String("store_id", rc.store),
				zap.String("authorization_model_id", rc.modelID),
				zap.String("object_type", objectType),
			)
		}

		if tupleUtils.UserSet == tupleUtils.GetUserTypeFromUser(tuple.GetUser()) {
			query.logger.WarnWithContext(
				ctx,
				fmt.Sprintf("unexpected userset evaluated on tupleset relation '%s#%s'", objectType, relation),
				zap.String("store_id", rc.store),
				zap.String("authorization_model_id", rc.modelID),
				zap.String("object_type", objectType),
			)

			return serverErrors.InvalidTuple(
				fmt.Sprintf("unexpected userset evaluated on relation '%s#%s'", tupleUtils.GetType(rc.tk.GetObject()), relation),
				tupleUtils.NewTupleKey(tuple.GetObject(), relation, tuple.GetUser()),
			)
		}

		if !tupleUtils.IsValidObject(userObj) {
			continue // TupleToUserset tuplesets should be of the form 'objectType:id' or 'objectType:id#relation' but are not guaranteed to be because it is neither a user or userset
		}

		usersetRel := node.TupleToUserset.GetComputedUserset().GetRelation() //reader

		if userRel == "" {
			userRel = usersetRel // userRel=reader
		}

		tupleKey := &openfgapb.TupleKey{
			Object:   userObj,         //folder:budgets
			Relation: userRel,         //reader
			User:     rc.tk.GetUser(), //anne
		}
		tracer := tracer.AppendString(tupleUtils.ToObjectRelationString(userObj, userRel))
		nestedRC := rc.fork(tupleKey, tracer, false)
		go func(c chan<- *chanResolveResult) {
			defer wg.Done()

			userset, err := query.getTypeDefinitionRelationUsersets(ctx, nestedRC) // folder:budgets#reader
			if err == nil {
				err = query.resolveNode(ctx, nestedRC, userset, typesys)
			}
			continue
		}

		userType := tupleUtils.GetUserTypeFromUser(user)
		if userType == tupleUtils.UserSet {
			return nil, serverErrors.InvalidTuple(
				fmt.Sprintf("unexpected userset evaluated on relation '%s#%s'", objectType, tupleset),
				tupleUtils.NewTupleKey(targetObject, tupleset, user),
			)
		}

		tObject, tRelation := tupleUtils.SplitObjectRelation(user)
		// We only proceed in the case that tRelation == userset.GetComputedUserset().GetRelation().
		// tRelation may be empty, and in this case, we set it to userset.GetComputedUserset().GetRelation().
			tRelation = userset.GetComputedUserset().GetRelation()
		}

		cs := &openfgapb.TupleKey{
			Object:   tObject,
			Relation: tRelation,
			return serverErrors.HandleTupleValidateError(&tupleUtils.IndirectWriteError{Reason: IndirectWriteErrorReason, TupleKey: tk})
		}

		if err := c.validateNoUsersetForRelationReferencedInTupleset(authModel, tk); err != nil {
			return err
		}

		if err := c.validateTypesForTuple(authModel, tk); err != nil {
			return err
		}
	return nil
}

func (c *WriteCommand) validateNoUsersetForRelationReferencedInTupleset(authModel *openfgapb.AuthorizationModel, tk *openfgapb.TupleKey) error {
	if !tupleUtils.IsObjectRelation(tk.GetUser()) {
		return nil
	}

	objType := tupleUtils.GetType(tk.GetObject())

	// at this point we know tk.User is a userset
	// if tk.Relation is used in a `x from y` definition (in the `y` part), throw an error
	ts := typesystem.New(authModel)
	for _, arrayOfTtus := range ts.GetAllTupleToUsersetsDefinitions()[objType] {
		for _, tupleToUserSetDef := range arrayOfTtus {
			if tupleToUserSetDef.Tupleset.Relation == tk.Relation {
				return serverErrors.InvalidTuple(fmt.Sprintf("Userset '%s' is not allowed to have relation '%s' with '%s'", tk.User, tk.Relation, tk.Object), tk)
			}
		}
	}

	return nil
}

// validateTypesForTuple makes sure that when writing a tuple, the types are compatible.
// 1. If the tuple is of the form (user=person:bob, relation=reader, object=doc:budget), then the type "doc", relation "reader" must allow type "person".
// 2. If the tuple is of the form (user=group:abc#member, relation=reader, object=doc:budget), then the type "doc", relation "reader" must allow type "group", relation "member".
			},
		},
		{
			name: "Error if userset encountered in tupleset relation of a TTU definition",
			typeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "repo",
											Relation: "manager",
										},
										ComputedUserset: &openfgapb.ObjectRelation{
											Relation: "repo_admin",
										},
									}}},
					},
				},
				{
					Type: "group",
					Relations: map[string]*openfgapb.Userset{
						"member": typesystem.This(),
					},
				},
			},
			tuples: []*openfgapb.TupleKey{
				tuple.NewTupleKey("repo:openfga/canaveral", "manager", "org:openfga#repo_admin"),
				tuple.NewTupleKey("org:openfga", "repo_admin", "group:eng#member"),
				tuple.NewTupleKey("group:eng", "member", "github|jose@openfga"),
			},
			resolveNodeLimit: defaultResolveNodeLimit,
			request: &openfgapb.CheckRequest{
				TupleKey: tuple.NewTupleKey("repo:openfga/canaveral", "admin", "github|jose@openfga"),
				Trace:    true,
			},
			err: serverErrors.InvalidTuple(
				"unexpected userset evaluated on relation 'repo#manager'",
				tuple.NewTupleKey("repo:openfga/canaveral", "manager", "org:openfga#repo_admin"),
			),
		},
		{
			name: "ExecuteCanResolveRecursiveComputedUserSets",
				errors.New("unexpected rewrite on relation 'document#parent'"),
			),
		},
		{
			//type org
			//	relations
			//		define viewer as self
			//		define can_view as viewer
			//type document
			//	relations
			//		define parent as self
			//		define viewer as viewer from parent
			name:             "Fails if expanding the computed userset of a tupleToUserset rewrite",
			resolveNodeLimit: defaultResolveNodeLimit,
			typeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "document",
					Relations: map[string]*openfgapb.Userset{
						"parent": typesystem.This(),
						"viewer": typesystem.TupleToUserset("parent", "viewer"),
					},
				},
				{
					Type: "org",
					Relations: map[string]*openfgapb.Userset{
						"viewer":   typesystem.This(),
						"can_view": typesystem.ComputedUserset("viewer"),
					},
				},
			},
			tuples: []*openfgapb.TupleKey{
				tuple.NewTupleKey("org:x", "viewer", "org:y"),
				tuple.NewTupleKey("document:1", "parent", "org:y#can_view"),
				tuple.NewTupleKey("document:1", "parent", "org:z#can_view"), //not relevant
			},
			request: &openfgapb.CheckRequest{
				TupleKey:         tuple.NewTupleKey("document:1", "viewer", "org:y"),
				ContextualTuples: &openfgapb.ContextualTupleKeys{},
			},
			err: serverErrors.InvalidTuple(
				"unexpected userset evaluated on relation 'document#parent'",
				tuple.NewTupleKey("document:1", "parent", "org:y#can_view"),
			),
		},
		{
			//type org
			//	relations
			//		define viewer as self
			//		define can_view as viewer
			//type document
			//	relations
			//		define parent as self
			//		define viewer as viewer from parent
			name:             "Fails if expanding the computed userset of a tupleToUserset rewrite",
			resolveNodeLimit: defaultResolveNodeLimit,
			typeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "document",
					Relations: map[string]*openfgapb.Userset{
						"parent": typesystem.This(),
						"viewer": typesystem.Union(
							typesystem.This(),
							typesystem.TupleToUserset("parent", "viewer"),
						),
					},
				},
			},
			tuples: []*openfgapb.TupleKey{
				tuple.NewTupleKey("document:1", "parent", "document:2#viewer"),
				tuple.NewTupleKey("document:2", "viewer", "jon"),
			},
			request: &openfgapb.CheckRequest{
				TupleKey:         tuple.NewTupleKey("document:1", "viewer", "org:y"),
				ContextualTuples: &openfgapb.ContextualTupleKeys{},
			},
			err: serverErrors.InvalidTuple(
				"unexpected userset evaluated on relation 'document#parent'",
				tuple.NewTupleKey("document:1", "parent", "document:2#viewer"),
			),
		},
	}

	ctx := context.Background()
				},
			},
		},
		{
			name: "Tuple involving userset that is not involved in TTU rewrite",
			typeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "document",
					Relations: map[string]*openfgapb.Userset{
						"parent": typesystem.This(),
						"editor": typesystem.This(),
					},
				},
			},
			tuples: []*openfgapb.TupleKey{
				tuple.NewTupleKey("document:1", "parent", "document:2#editor"),
			},
			request: &openfgapb.ExpandRequest{
				TupleKey: tuple.NewTupleKey("document:1", "parent", ""),
			},
			expected: &openfgapb.ExpandResponse{
				Tree: &openfgapb.UsersetTree{
					Root: &openfgapb.UsersetTree_Node{
						Name: "document:1#parent",
						Value: &openfgapb.UsersetTree_Node_Leaf{
							Leaf: &openfgapb.UsersetTree_Leaf{
								Value: &openfgapb.UsersetTree_Leaf_Users{
									Users: &openfgapb.UsersetTree_Users{
										Users: []string{"document:2#editor"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	require := require.New(t)
				errors.Errorf("unexpected rewrite on relation '%s#%s'", "document", "parent"),
			),
		},
		{
			name: "Tuple involving userset returns error if it is referenced in a TTU rewrite",
			typeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "folder",
					Relations: map[string]*openfgapb.Userset{
						"viewer": typesystem.This(),
					},
				},
				{
					Type: "document",
					Relations: map[string]*openfgapb.Userset{
						"parent": typesystem.This(),
						"editor": typesystem.This(),
						"viewer": typesystem.TupleToUserset("parent", "viewer"),
					},
				},
			},
			tuples: []*openfgapb.TupleKey{
				tuple.NewTupleKey("document:1", "parent", "document:2#editor"),
			},
			request: &openfgapb.ExpandRequest{
				TupleKey: tuple.NewTupleKey("document:1", "viewer", ""),
			},
			expected: serverErrors.InvalidTuple(
				"unexpected userset evaluated on relation 'document#parent'",
				tuple.NewTupleKey("document:1", "parent", "document:2#editor"),
			),
		},
		{
			name: "Tuple involving userset returns error if same ComputedUserset involved in TTU rewrite",
			typeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "document",
					Relations: map[string]*openfgapb.Userset{
						"parent": typesystem.This(),
						"viewer": typesystem.Union(
							typesystem.This(),
							typesystem.TupleToUserset("parent", "viewer"),
						),
					},
				},
			},
			tuples: []*openfgapb.TupleKey{
				tuple.NewTupleKey("document:1", "parent", "document:2#viewer"),
				tuple.NewTupleKey("document:2", "viewer", "jon"),
			},
			request: &openfgapb.ExpandRequest{
				TupleKey: tuple.NewTupleKey("document:1", "viewer", ""),
			},
			expected: serverErrors.InvalidTuple(
				"unexpected userset evaluated on relation 'document#parent'",
				tuple.NewTupleKey("document:1", "parent", "document:2#viewer"),
			),
		},
	}

	require := require.New(t)
		},
		err: serverErrors.NewInternalError("invalid authorization model", errors.New("invalid authorization model")),
	},

	{
		_name: "Write fails if a. schema version is 1.0 b. user is a userset c. relation is referenced in a tupleset of a tupleToUserset relation",
		model: &openfgapb.AuthorizationModel{
			Id:            ulid.Make().String(),
			SchemaVersion: typesystem.SchemaVersion1_0,
			TypeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "folder",
					Relations: map[string]*openfgapb.Userset{
						"owner": typesystem.This(),
						"admin": typesystem.This(),
					},
				},
				{
					Type: "document",
					Relations: map[string]*openfgapb.Userset{
						"parent":   typesystem.This(),
						"can_view": typesystem.TupleToUserset("parent", "owner"), //owner from parent
					},
				},
			},
		},
		request: &openfgapb.WriteRequest{
			Writes: &openfgapb.TupleKeys{TupleKeys: []*openfgapb.TupleKey{
				tuple.NewTupleKey("document:budget", "parent", "folder:budgets#admin"),
			}},
		},
		err: serverErrors.InvalidTuple("Userset 'folder:budgets#admin' is not allowed to have relation 'parent' with 'document:budget'",
			tuple.NewTupleKey("document:budget", "parent", "folder:budgets#admin"),
		),
	},
	{
		_name: "Write fails if a. schema version is 1.0 b. user is a userset c. relation is referenced in a tupleset of a tupleToUserset relation (defined as union)",
		model: &openfgapb.AuthorizationModel{
			Id:            ulid.Make().String(),
			SchemaVersion: typesystem.SchemaVersion1_0,
			TypeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "folder",
					Relations: map[string]*openfgapb.Userset{
						"owner": typesystem.This(),
						"admin": typesystem.This(),
					},
				},
				{
					Type: "document",
					Relations: map[string]*openfgapb.Userset{
						"parent": typesystem.This(),
						"can_view": typesystem.Union(
							typesystem.TupleToUserset("parent", "owner"), //owner from parent
							typesystem.TupleToUserset("parent", "admin"), //admin from parent
						),
					},
				},
			},
		},
		request: &openfgapb.WriteRequest{
			Writes: &openfgapb.TupleKeys{TupleKeys: []*openfgapb.TupleKey{
				tuple.NewTupleKey("document:budget", "parent", "folder:budgets#admin"),
			}},
		},
		err: serverErrors.InvalidTuple("Userset 'folder:budgets#admin' is not allowed to have relation 'parent' with 'document:budget'",
			tuple.NewTupleKey("document:budget", "parent", "folder:budgets#admin"),
		),
	},
	{
		_name: "Write fails if a. schema version is 1.0 b. user is a userset c. relation is referenced in a tupleset of a tupleToUserset relation (defined as intersection)",
		model: &openfgapb.AuthorizationModel{
			Id:            ulid.Make().String(),
			SchemaVersion: typesystem.SchemaVersion1_0,
			TypeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "folder",
					Relations: map[string]*openfgapb.Userset{
						"owner": typesystem.This(),
						"admin": typesystem.This(),
					},
				},
				{
					Type: "document",
					Relations: map[string]*openfgapb.Userset{
						"parent": typesystem.This(),
						"can_view": typesystem.Intersection(
							typesystem.TupleToUserset("parent", "owner"), //owner from parent
							typesystem.TupleToUserset("parent", "admin"), //admin from parent
						),
					},
				},
			},
		},
		request: &openfgapb.WriteRequest{
			Writes: &openfgapb.TupleKeys{TupleKeys: []*openfgapb.TupleKey{
				tuple.NewTupleKey("document:budget", "parent", "folder:budgets#admin"),
			}},
		},
		err: serverErrors.InvalidTuple("Userset 'folder:budgets#admin' is not allowed to have relation 'parent' with 'document:budget'",
			tuple.NewTupleKey("document:budget", "parent", "folder:budgets#admin"),
		),
	},
	{
		_name: "Write fails if a. schema version is 1.0 b. user is a userset c. relation is referenced in a tupleset of a tupleToUserset relation (defined as difference)",
		model: &openfgapb.AuthorizationModel{
			Id:            ulid.Make().String(),
			SchemaVersion: typesystem.SchemaVersion1_0,
			TypeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "folder",
					Relations: map[string]*openfgapb.Userset{
						"owner": typesystem.This(),
						"admin": typesystem.This(),
					},
				},
				{
					Type: "document",
					Relations: map[string]*openfgapb.Userset{
						"parent": typesystem.This(),
						"can_view": typesystem.Difference(
							typesystem.TupleToUserset("parent", "owner"), //owner from parent
							typesystem.TupleToUserset("parent", "admin"), //admin from parent
						),
					},
				},
			},
		},
		request: &openfgapb.WriteRequest{
			Writes: &openfgapb.TupleKeys{TupleKeys: []*openfgapb.TupleKey{
				tuple.NewTupleKey("document:budget", "parent", "folder:budgets#admin"),
			}},
		},
		err: serverErrors.InvalidTuple("Userset 'folder:budgets#admin' is not allowed to have relation 'parent' with 'document:budget'",
			tuple.NewTupleKey("document:budget", "parent", "folder:budgets#admin"),
		),
	},
	{
		_name: "Write succeeds if a. schema version is 1.0 b. user is a userset c. relation is referenced in a tupleset of a tupleToUserset relation of another type",
		model: &openfgapb.AuthorizationModel{
			Id:            ulid.Make().String(),
			SchemaVersion: typesystem.SchemaVersion1_0,
			TypeDefinitions: []*openfgapb.TypeDefinition{
				{
					Type: "folder",
					Relations: map[string]*openfgapb.Userset{
						"owner": typesystem.This(),
						"parent": typesystem.Union( // let's confuse the code. if this were defined in 'document' type, it would fail
							typesystem.TupleToUserset("parent", "owner"),
						),
					},
				},
				{
					Type: "document",
					Relations: map[string]*openfgapb.Userset{
						"owner":  typesystem.This(),
						"parent": typesystem.This(),
					},
				},
			},
		},
		request: &openfgapb.WriteRequest{
			Writes: &openfgapb.TupleKeys{TupleKeys: []*openfgapb.TupleKey{
				tuple.NewTupleKey("document:budget", "parent", "folder:budgets#admin"),
			}},
		},
	},
}

func TestWriteCommand(t *testing.T, datastore storage.OpenFGADatastore) {
		err     error
	}{
		{
			name: "succeeds with a simple model",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			},
		},
		{
			name: "succeeds with a complex model",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: "somestoreid",
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.ExceededEntityLimit("type definitions in an authorization model", datastore.MaxTypesInTypeDefinition()),
		},
		{
			name: "succeeds with empty relations",
			request: &openfgapb.WriteAuthorizationModelRequest{
				TypeDefinitions: []*openfgapb.TypeDefinition{
					{
			},
		},
		{
			name: "succeeds with zero length relations",
			request: &openfgapb.WriteAuthorizationModelRequest{
				TypeDefinitions: []*openfgapb.TypeDefinition{
					{
			},
		},
		{
			name: "fails if the same type appears twice",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(typesystem.ErrDuplicateTypes),
		},
		{
			name: "fails if a relation is not defined",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "fails if unknown relation in computed userset definition",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "fails if unknown relation in tuple to userset definition (computed userset component)",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "", Relation: "owner"}),
		},
		{
			name: "fails if unknown relation in union",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "fails if unknown relation in difference base argument",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "fails if unknown relation in difference subtract argument",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "fails if unknown relation in tuple to userset definition (tupleset component)",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "fails if unknown relation in computed userset",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "bar", Relation: "writer"}),
		},
		{
			name: "fails if unknown relation in intersection",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.RelationUndefinedError{ObjectType: "repo", Relation: "owner"}),
		},
		{
			name: "fails if difference includes same relation twice",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{ObjectType: "repo", Relation: "viewer"}),
		},
		{
			name: "fails if union includes same relation twice",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{ObjectType: "repo", Relation: "viewer"}),
		},
		{
			name: "fails if intersection includes same relation twice",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{ObjectType: "repo", Relation: "viewer"}),
		},
		{
			name: "Success if Union Rewrite Contains Repeated Definitions",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			},
		},
		{
			name: "Success if Intersection Rewrite Contains Repeated Definitions",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			},
		},
		{
			name: "Success if Exclusion Rewrite Contains Repeated Definitions",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			},
		},
		{
			name: "Fails if Tupleset relation involves ComputedUserset rewrite",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			),
		},
		{
			name: "Fails if Tupleset relation involves Union rewrite",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			),
		},
		{
			name: "Fails if Tupleset relation involves Intersection rewrite",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			),
		},
		{
			name: "Fails if Tupleset relation involves Exclusion rewrite",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
			),
		},
		{
			name: "Fails if Tupleset relation involves TupleToUserset rewrite",
			request: &openfgapb.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: []*openfgapb.TypeDefinition{
	openfgapb "go.buf.build/openfga/go/openfga/api/openfga/v1"
)

// ValidateTuple checks whether a tuple is valid. If it is not, returns as error.
// If it is, it grabs the type of the tuple's object, and the tuple's relation, and returns its corresponding userset.
func ValidateTuple(ctx context.Context, backend storage.TypeDefinitionReadBackend, store, authorizationModelID string, tk *openfgapb.TupleKey) (*openfgapb.Userset, error) {
	if err := tuple.ValidateUser(tk); err != nil {
		return nil, err
	return ValidateObjectsRelations(ctx, backend, store, authorizationModelID, tk)
}

// ValidateObjectsRelations checks whether a tuple's object and relations are valid. If they are not, returns an error.
// If they are, it grabs the type of the tuple's object, and the tuple's relation, and returns its corresponding userset.
func ValidateObjectsRelations(ctx context.Context, backend storage.TypeDefinitionReadBackend, store, modelID string, t *openfgapb.TupleKey) (*openfgapb.Userset, error) {
	if !tuple.IsValidRelation(t.GetRelation()) {
		return nil, &tuple.InvalidTupleError{Reason: "invalid relation", TupleKey: t}
