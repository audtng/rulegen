package main

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("internal/graph/check")
	defaultMaxConcurrentReadsForCheck = math.MaxUint32
)

type ResolveCheckRequest struct {
	StoreID              string
	AuthorizationModelID string
	TupleKey             *openfgav1.TupleKey
	ContextualTuples     []*openfgav1.TupleKey
	ResolutionMetadata   *ResolutionMetadata
}

type ResolveCheckResponse struct {
		return nil, fmt.Errorf("relation '%s' undefined for object type '%s'", relation, objectType)
	}

	resp, err := union(ctx, c.concurrencyLimit, c.checkRewrite(ctx, req, rel.GetRewrite()))
	if err != nil {
		return nil, err
				}

				if usersetRelation != "" {
					handlers = append(handlers, c.dispatch(
						ctx,
						&ResolveCheckRequest{
							StoreID:              storeID,
							AuthorizationModelID: req.GetAuthorizationModelID(),
							TupleKey:             tuple.NewTupleKey(usersetObject, usersetRelation, tk.GetUser()),
							ResolutionMetadata: &ResolutionMetadata{
								Depth:               req.GetResolutionMetadata().Depth - 1,
								DatastoreQueryCount: response.GetResolutionMetadata().DatastoreQueryCount,
							},
						}))
				}
			}

// checkComputedUserset evaluates the Check request with the rewritten relation (e.g. the computed userset relation).
func (c *LocalChecker) checkComputedUserset(parentctx context.Context, req *ResolveCheckRequest, rewrite *openfgav1.Userset_ComputedUserset) CheckHandlerFunc {
	return func(ctx context.Context) (*ResolveCheckResponse, error) {
		ctx, span := tracer.Start(ctx, "checkComputedUserset")
		defer span.End()

		return c.dispatch(
			ctx,
			&ResolveCheckRequest{
				StoreID:              req.GetStoreID(),
				AuthorizationModelID: req.GetAuthorizationModelID(),
				TupleKey: tuple.NewTupleKey(
					req.TupleKey.GetObject(),
					rewrite.ComputedUserset.GetRelation(),
					req.TupleKey.GetUser(),
				),
				ResolutionMetadata: &ResolutionMetadata{
					Depth:               req.GetResolutionMetadata().Depth - 1,
					DatastoreQueryCount: req.GetResolutionMetadata().DatastoreQueryCount,
				},
			})(ctx)
	}
}
				}
			}

			handlers = append(handlers, c.dispatch(
				ctx,
				&ResolveCheckRequest{
						Depth:               req.GetResolutionMetadata().Depth - 1,
						DatastoreQueryCount: req.GetResolutionMetadata().DatastoreQueryCount, // add TTU read below
					},
				}))
		}

		})
	}
}
				ContextualTuples: req.GetContextualTuples().GetTupleKeys(),
			}, reverseExpandResultsChan, resolutionMetadata)
			if err != nil {
				resultsChan <- ListObjectsResult{Err: err}
			}

					},
				})
				if err != nil {
					resultsChan <- ListObjectsResult{Err: err}
					return
				}
				if errors.Is(result.Err, serverErrors.AuthorizationModelResolutionTooComplex) {
					return nil, result.Err
				}
				return nil, serverErrors.HandleError("", result.Err)
			}


	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/internal/graph"
	serverErrors "github.com/openfga/openfga/pkg/server/errors"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/storage/storagewrappers"
	"github.com/openfga/openfga/pkg/tuple"
		ctx = graph.ContextWithResolutionDepth(ctx, 0)
	} else {
		if depth >= c.resolveNodeLimit {
			return serverErrors.AuthorizationModelResolutionTooComplex
		}

		ctx = graph.ContextWithResolutionDepth(ctx, depth+1)
		},
	})
	if err != nil {
		if errors.Is(err, graph.ErrResolutionDepthExceeded) {
			return nil, serverErrors.AuthorizationModelResolutionTooComplex
		}

	}
}

// This test ensures that when the data storage fails, ListObjects v0 throws an error
func TestListObjects_Unoptimized_UnhappyPaths(t *testing.T) {
	ctx := context.Background()
	store := ulid.Make().String()
	modelID := ulid.Make().String()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDatastore := mockstorage.NewMockOpenFGADatastore(mockController)

	mockDatastore.EXPECT().ReadAuthorizationModel(gomock.Any(), store, modelID).AnyTimes().Return(&openfgav1.AuthorizationModel{
		SchemaVersion: typesystem.SchemaVersion1_1,
		TypeDefinitions: parser.MustParse(`
		type user

		type repo
		  relations
		    define allowed: [user] as self
		    define viewer: [user] as self and allowed
		`),
	}, nil)
	mockDatastore.EXPECT().ReadStartingWithUser(gomock.Any(), store, gomock.Any()).AnyTimes().Return(nil, errors.New("error reading from storage"))

	s := MustNewServerWithOpts(
		WithDatastore(mockDatastore),
	)

	t.Run("error_listing_objects_from_storage_in_non-streaming_version", func(t *testing.T) {
		res, err := s.ListObjects(ctx, &openfgav1.ListObjectsRequest{
			StoreId:              store,
			AuthorizationModelId: modelID,
			Type:                 "repo",
			Relation:             "viewer",
			User:                 "user:bob",
		})

		require.Nil(t, res)
		require.ErrorIs(t, err, serverErrors.NewInternalError("", errors.New("error reading from storage")))
	})

	t.Run("error_listing_objects_from_storage_in_streaming_version", func(t *testing.T) {
		err := s.StreamedListObjects(&openfgav1.StreamedListObjectsRequest{
			StoreId:              store,
			AuthorizationModelId: modelID,
			Type:                 "repo",
			Relation:             "viewer",
			User:                 "user:bob",
		}, NewMockStreamServer())

		require.ErrorIs(t, err, serverErrors.NewInternalError("", errors.New("error reading from storage")))
	})
}

// This test ensures that when the data storage fails for known eror, ListObjects v0 throws the correct error
func TestListObjects_Unoptimized_UnhappyPaths_Known_Error(t *testing.T) {
	ctx := context.Background()
	store := ulid.Make().String()
	modelID := ulid.Make().String()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDatastore := mockstorage.NewMockOpenFGADatastore(mockController)

	mockDatastore.EXPECT().ReadAuthorizationModel(gomock.Any(), store, modelID).AnyTimes().Return(&openfgav1.AuthorizationModel{
		SchemaVersion: typesystem.SchemaVersion1_1,
		TypeDefinitions: parser.MustParse(`
		type user

		type repo
		  relations
		    define allowed: [user] as self
		    define viewer: [user] as self and allowed
		`),
	}, nil)
	mockDatastore.EXPECT().ReadStartingWithUser(gomock.Any(), store, gomock.Any()).AnyTimes().Return(nil, serverErrors.AuthorizationModelResolutionTooComplex)

	s := MustNewServerWithOpts(
		WithDatastore(mockDatastore),
	)

	t.Run("error_listing_objects_from_storage_in_non-streaming_version", func(t *testing.T) {
		res, err := s.ListObjects(ctx, &openfgav1.ListObjectsRequest{
			StoreId:              store,
			AuthorizationModelId: modelID,
			Type:                 "repo",
			Relation:             "viewer",
			User:                 "user:bob",
		})

		require.Nil(t, res)
		require.ErrorIs(t, err, serverErrors.AuthorizationModelResolutionTooComplex)
	})

	t.Run("error_listing_objects_from_storage_in_streaming_version", func(t *testing.T) {
		err := s.StreamedListObjects(&openfgav1.StreamedListObjectsRequest{
			StoreId:              store,
			AuthorizationModelId: modelID,
			Type:                 "repo",
			Relation:             "viewer",
			User:                 "user:bob",
		}, NewMockStreamServer())

		require.ErrorIs(t, err, serverErrors.AuthorizationModelResolutionTooComplex)
	})
}

// This test ensures that when the data storage fails, ListObjects v1 throws an error
func TestListObjects_UnhappyPaths(t *testing.T) {
	ctx := context.Background()
	store := ulid.Make().String()
	modelID := ulid.Make().String()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDatastore := mockstorage.NewMockOpenFGADatastore(mockController)

	mockDatastore.EXPECT().ReadAuthorizationModel(gomock.Any(), store, modelID).AnyTimes().Return(&openfgav1.AuthorizationModel{
		SchemaVersion: typesystem.SchemaVersion1_1,
		TypeDefinitions: []*openfgav1.TypeDefinition{
			{
				Type: "user",
			},
			{
				Type: "document",
				Relations: map[string]*openfgav1.Userset{
					"viewer": typesystem.This(),
				},
				Metadata: &openfgav1.Metadata{
					Relations: map[string]*openfgav1.RelationMetadata{
						"viewer": {
							DirectlyRelatedUserTypes: []*openfgav1.RelationReference{
								typesystem.DirectRelationReference("user", ""),
								typesystem.WildcardRelationReference("user"),
							},
						},
					},
				},
			},
		},
	}, nil)
	mockDatastore.EXPECT().ReadStartingWithUser(gomock.Any(), store, storage.ReadStartingWithUserFilter{
		ObjectType: "document",
		Relation:   "viewer",
		UserFilter: []*openfgav1.ObjectRelation{
			{Object: "user:*"},
			{Object: "user:bob"},
		}}).AnyTimes().Return(nil, errors.New("error reading from storage"))

	s := MustNewServerWithOpts(
		WithDatastore(mockDatastore),
	)

	t.Run("error_listing_objects_from_storage_in_non-streaming_version", func(t *testing.T) {
		res, err := s.ListObjects(ctx, &openfgav1.ListObjectsRequest{
			StoreId:              store,
			AuthorizationModelId: modelID,
			Type:                 "document",
			Relation:             "viewer",
			User:                 "user:bob",
		})

		require.Nil(t, res)
		require.ErrorIs(t, err, serverErrors.NewInternalError("", errors.New("error reading from storage")))
	})

	t.Run("error_listing_objects_from_storage_in_streaming_version", func(t *testing.T) {
		err := s.StreamedListObjects(&openfgav1.StreamedListObjectsRequest{
			StoreId:              store,
			AuthorizationModelId: modelID,
			Type:                 "document",
			Relation:             "viewer",
			User:                 "user:bob",
		}, NewMockStreamServer())

		require.ErrorIs(t, err, serverErrors.NewInternalError("", errors.New("error reading from storage")))
	})
}

// This test ensures that when the data storage fails with known errors, ListObjects v1 throws an error
func TestListObjects_UnhappyPaths_Known_Error(t *testing.T) {
	ctx := context.Background()
	store := ulid.Make().String()
	modelID := ulid.Make().String()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDatastore := mockstorage.NewMockOpenFGADatastore(mockController)

	mockDatastore.EXPECT().ReadAuthorizationModel(gomock.Any(), store, modelID).AnyTimes().Return(&openfgav1.AuthorizationModel{
		SchemaVersion: typesystem.SchemaVersion1_1,
		TypeDefinitions: []*openfgav1.TypeDefinition{
			{
				Type: "user",
			},
			{
				Type: "document",
				Relations: map[string]*openfgav1.Userset{
					"viewer": typesystem.This(),
				},
				Metadata: &openfgav1.Metadata{
					Relations: map[string]*openfgav1.RelationMetadata{
						"viewer": {
							DirectlyRelatedUserTypes: []*openfgav1.RelationReference{
								typesystem.DirectRelationReference("user", ""),
								typesystem.WildcardRelationReference("user"),
							},
						},
					},
				},
			},
		},
	}, nil)
	mockDatastore.EXPECT().ReadStartingWithUser(gomock.Any(), store, storage.ReadStartingWithUserFilter{
		ObjectType: "document",
		Relation:   "viewer",
		UserFilter: []*openfgav1.ObjectRelation{
			{Object: "user:*"},
			{Object: "user:bob"},
		}}).AnyTimes().Return(nil, serverErrors.AuthorizationModelResolutionTooComplex)

	s := MustNewServerWithOpts(
		WithDatastore(mockDatastore),
	)

	t.Run("error_listing_objects_from_storage_in_non-streaming_version", func(t *testing.T) {
		res, err := s.ListObjects(ctx, &openfgav1.ListObjectsRequest{
			StoreId:              store,
			AuthorizationModelId: modelID,
			Type:                 "document",
			Relation:             "viewer",
			User:                 "user:bob",
		})

		require.Nil(t, res)
		require.ErrorIs(t, err, serverErrors.AuthorizationModelResolutionTooComplex)
	})

	t.Run("error_listing_objects_from_storage_in_streaming_version", func(t *testing.T) {
		err := s.StreamedListObjects(&openfgav1.StreamedListObjectsRequest{
			StoreId:              store,
			AuthorizationModelId: modelID,
			Type:                 "document",
			Relation:             "viewer",
			User:                 "user:bob",
		}, NewMockStreamServer())

		require.ErrorIs(t, err, serverErrors.AuthorizationModelResolutionTooComplex)
	})
}

	parser "github.com/craigpastro/openfga-dsl-parser/v2"
	"github.com/oklog/ulid/v2"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/server/commands/reverseexpand"
	serverErrors "github.com/openfga/openfga/pkg/server/errors"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/tuple"
	"github.com/openfga/openfga/pkg/typesystem"
				tuple.NewTupleKey("folder:folder2", "parent", "folder:folder1"),
				tuple.NewTupleKey("folder:folder3", "parent", "folder:folder2"),
			},
			expectedError:        serverErrors.AuthorizationModelResolutionTooComplex,
			expectedDSQueryCount: 0,
		},
		{

import (
	"context"
	"errors"
	"fmt"
	"testing"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/logger"
	"github.com/openfga/openfga/pkg/server/commands"
	serverErrors "github.com/openfga/openfga/pkg/server/errors"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/typesystem"
	"github.com/stretchr/testify/require"
)

func WriteAuthorizationModelTest(t *testing.T, datastore storage.OpenFGADatastore) {
		name          string
		request       *openfgav1.WriteAuthorizationModelRequest
		allowSchema10 bool
		err           error
	}{
		{
			name: "fails_if_too_many_types",
				SchemaVersion:   typesystem.SchemaVersion1_1,
			},
			allowSchema10: false,
			err:           serverErrors.ExceededEntityLimit("type definitions in an authorization model", datastore.MaxTypesPerAuthorizationModel()),
		},
		{
			name: "fails_if_a_relation_is_not_defined",
				SchemaVersion: typesystem.SchemaVersion1_1,
			},
			allowSchema10: false,
			err:           serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{ObjectType: "repo", Relation: "owner", Cause: typesystem.ErrInvalidUsersetRewrite}),
		},
		{
			name: "Fails_if_type_info_metadata_is_omitted_in_1.1_model",
				},
			},
			allowSchema10: false,
			err: serverErrors.InvalidAuthorizationModelInput(
				errors.New("the assignable relation 'reader' in object type 'document' must contain at least one relation type"),
			),
		},
		{
			name: "Fails_if_writing_1_0_model_because_it_will_be_interpreted_as_1_1",
				},
			},
			allowSchema10: true,
			err:           serverErrors.InvalidAuthorizationModelInput(typesystem.AssignableRelationError("document", "reader")),
		},
		{
			name: "Works_if_no_schema_version",
				`),
				SchemaVersion: typesystem.SchemaVersion1_1,
			},
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{
				ObjectType: "document",
				Relation:   "viewer",
				Cause:      typesystem.ErrNoEntrypoints},
			),
		},
		{
			name: "self_referencing_type_restriction_without_entrypoint_2",
				`),
				SchemaVersion: typesystem.SchemaVersion1_1,
			},
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{
				ObjectType: "document",
				Relation:   "viewer",
				Cause:      typesystem.ErrNoEntrypoints,
			}),
		},
		{
			name: "self_referencing_type_restriction_without_entrypoint_3",
				`),
				SchemaVersion: typesystem.SchemaVersion1_1,
			},
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{
				ObjectType: "document",
				Relation:   "viewer",
				Cause:      typesystem.ErrNoEntrypoints,
			}),
		},
		{
			name: "rewritten_relation_in_intersection_unresolvable",
				`),
				SchemaVersion: typesystem.SchemaVersion1_1,
			},
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{
				ObjectType: "document",
				Relation:   "action1",
				Cause:      typesystem.ErrNoEntryPointsLoop,
			}),
		},
		{
			name: "direct_relationship_with_entrypoint",
				`),
			},
		},

		{
			name: "rewritten_relation_in_exclusion_unresolvable",
			request: &openfgav1.WriteAuthorizationModelRequest{
				    define action3 as admin but not action1
				`),
			},
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{
				ObjectType: "document",
				Relation:   "action1",
				Cause:      typesystem.ErrNoEntryPointsLoop,
			}),
		},
		{
			name: "no_entrypoint_3a",
				    define editor: [user] as self
				`),
			},
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{
				ObjectType: "document",
				Relation:   "viewer",
				Cause:      typesystem.ErrNoEntrypoints,
			}),
		},

		{
			name: "no_entrypoint_3b",
			request: &openfgav1.WriteAuthorizationModelRequest{
				    define editor: [user] as self
				`),
			},
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{
				ObjectType: "document",
				Relation:   "viewer",
				Cause:      typesystem.ErrNoEntrypoints,
			}),
		},
		{
			name: "no_entrypoint_4",
				    define viewer as editor from parent
				`),
			},
			err: serverErrors.InvalidAuthorizationModelInput(&typesystem.InvalidRelationError{
				ObjectType: "document",
				Relation:   "editor",
				Cause:      typesystem.ErrNoEntrypoints,
			}),
		},
		{
			name: "self_referencing_type_restriction_with_entrypoint_1",
					},
				},
			},
			err: serverErrors.InvalidAuthorizationModelInput(
				fmt.Errorf("the type name of a type definition cannot be an empty string"),
			),
		},
		{
			name: "relation_name_is_empty_string",
					},
				},
			},
			err: serverErrors.InvalidAuthorizationModelInput(
				fmt.Errorf("type 'user' defines a relation with an empty string for a name"),
			),
		},
	}

		t.Run(test.name, func(t *testing.T) {
			cmd := commands.NewWriteAuthorizationModelCommand(datastore, logger)
			resp, err := cmd.Execute(ctx, test.request)
			require.ErrorIs(t, err, test.err)

			if err == nil {
				_, err = ulid.Parse(resp.AuthorizationModelId)
		}
	}

	return nil
}

	return fmt.Errorf("the relation type '%s' on '%s' in object type '%s' is not valid", relationType, relation, objectType)
}

// getAllTupleToUsersetsDefinitions returns a map where the key is the object type and the value
// is another map where key=relationName, value=list of tuple to usersets declared in that relation
func (t *TypeSystem) getAllTupleToUsersetsDefinitions() map[string]map[string][]*openfgav1.TupleToUserset {
	"github.com/stretchr/testify/require"
)

func TestNewAndValidate(t *testing.T) {

	tests := []struct {
