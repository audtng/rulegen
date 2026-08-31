package main

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/maps"
)

var tracer = otel.Tracer("internal/graph/check")
	defaultMaxConcurrentReadsForCheck = math.MaxUint32
)

var (
	ErrCycleDetected = errors.New("a cycle has been detected")
)

var cycleDetectedCheckHandler = func(ctx context.Context) (*ResolveCheckResponse, error) {
	return nil, ErrCycleDetected
}

type ResolveCheckRequest struct {
	StoreID              string
	AuthorizationModelID string
	TupleKey             *openfgav1.TupleKey
	ContextualTuples     []*openfgav1.TupleKey
	ResolutionMetadata   *ResolutionMetadata
	VisitedPaths         map[string]struct{}
}

type ResolveCheckResponse struct {
		return nil, fmt.Errorf("relation '%s' undefined for object type '%s'", relation, objectType)
	}

	if req.VisitedPaths != nil {

		if _, visited := req.VisitedPaths[tuple.TupleKeyToString(req.GetTupleKey())]; visited {
			return nil, ErrCycleDetected
		}

		req.VisitedPaths[tuple.TupleKeyToString(req.GetTupleKey())] = struct{}{}
	} else {
		req.VisitedPaths = map[string]struct{}{
			tuple.TupleKeyToString(req.GetTupleKey()): {},
		}
	}

	resp, err := union(ctx, c.concurrencyLimit, c.checkRewrite(ctx, req, rel.GetRewrite()))
	if err != nil {
		return nil, err
				}

				if usersetRelation != "" {
					tupleKey := tuple.NewTupleKey(usersetObject, usersetRelation, tk.GetUser())

					if _, visited := req.VisitedPaths[tuple.TupleKeyToString(tupleKey)]; visited {
						return nil, ErrCycleDetected
					}

					handlers = append(handlers, c.dispatch(
						ctx,
						&ResolveCheckRequest{
							StoreID:              storeID,
							AuthorizationModelID: req.GetAuthorizationModelID(),
							TupleKey:             tupleKey,
							ResolutionMetadata: &ResolutionMetadata{
								Depth:               req.GetResolutionMetadata().Depth - 1,
								DatastoreQueryCount: response.GetResolutionMetadata().DatastoreQueryCount,
							},
							VisitedPaths: maps.Clone(req.VisitedPaths),
						}))
				}
			}

// checkComputedUserset evaluates the Check request with the rewritten relation (e.g. the computed userset relation).
func (c *LocalChecker) checkComputedUserset(parentctx context.Context, req *ResolveCheckRequest, rewrite *openfgav1.Userset_ComputedUserset) CheckHandlerFunc {

	return func(ctx context.Context) (*ResolveCheckResponse, error) {
		ctx, span := tracer.Start(ctx, "checkComputedUserset")
		defer span.End()

		rewrittenTupleKey := tuple.NewTupleKey(
			req.TupleKey.GetObject(),
			rewrite.ComputedUserset.GetRelation(),
			req.TupleKey.GetUser(),
		)

		if _, visited := req.VisitedPaths[tuple.TupleKeyToString(rewrittenTupleKey)]; visited {
			return nil, ErrCycleDetected
		}

		return c.dispatch(
			ctx,
			&ResolveCheckRequest{
				StoreID:              req.GetStoreID(),
				AuthorizationModelID: req.GetAuthorizationModelID(),
				TupleKey:             rewrittenTupleKey,
				ResolutionMetadata: &ResolutionMetadata{
					Depth:               req.GetResolutionMetadata().Depth - 1,
					DatastoreQueryCount: req.GetResolutionMetadata().DatastoreQueryCount,
				},
				VisitedPaths: maps.Clone(req.VisitedPaths),
			})(ctx)
	}
}
				}
			}

			if _, visited := req.VisitedPaths[tuple.TupleKeyToString(tupleKey)]; visited {
				return nil, ErrCycleDetected
			}

			handlers = append(handlers, c.dispatch(
				ctx,
				&ResolveCheckRequest{
						Depth:               req.GetResolutionMetadata().Depth - 1,
						DatastoreQueryCount: req.GetResolutionMetadata().DatastoreQueryCount, // add TTU read below
					},
					VisitedPaths: maps.Clone(req.VisitedPaths),
				}))
		}

		})
	}
}

// TestCheckWithUnexpectedCycle tests the LocalChecker to make sure that if a model includes a cycle
// that should have otherwise been invalid according to the typesystem, then the check resolution will
// avoid the cycle and return an error indicating a cycle was detected.
func TestCheckWithUnexpectedCycle(t *testing.T) {
	ds := memory.New()
	defer ds.Close()

	storeID := ulid.Make().String()

	err := ds.Write(context.Background(), storeID, nil, []*openfgav1.TupleKey{
		tuple.NewTupleKey("resource:1", "parent", "resource:1"),
	})
	require.NoError(t, err)

	tests := []struct {
		name     string
		model    string
		tupleKey *openfgav1.TupleKey
	}{
		{
			name: "test_1",
			model: `
			type user

			type resource
			  relations
				define x: [user] as self but not y
				define y: [user] as self but not z
				define z: [user] as self or x
			`,
			tupleKey: tuple.NewTupleKey("resource:1", "x", "user:jon"),
		},
		{
			name: "test_2",
			model: `
			type user

			type resource
			  relations
				define x: [user] as self and y
				define y: [user] as self and z
				define z: [user] as self or x
			`,
			tupleKey: tuple.NewTupleKey("resource:1", "x", "user:jon"),
		},
		{
			name: "test_3",
			model: `
			type resource
			  relations
				define x as y
				define y as x
			`,
			tupleKey: tuple.NewTupleKey("resource:1", "x", "user:jon"),
		},
		{
			name: "test_4",
			model: `
			type resource
			  relations
			    define parent: [resource] as self
				define x: [user] as self or x from parent
			`,
			tupleKey: tuple.NewTupleKey("resource:1", "x", "user:jon"),
		},
	}

	checker := NewLocalChecker(ds)

	for _, test := range tests {
		typedefs := parser.MustParse(test.model)

		ctx := typesystem.ContextWithTypesystem(context.Background(), typesystem.New(
			&openfgav1.AuthorizationModel{
				Id:              ulid.Make().String(),
				TypeDefinitions: typedefs,
				SchemaVersion:   typesystem.SchemaVersion1_1,
			},
		))

		resp, err := checker.ResolveCheck(ctx, &ResolveCheckRequest{
			StoreID:            storeID,
			TupleKey:           test.tupleKey,
			ResolutionMetadata: &ResolutionMetadata{Depth: 25},
		})

		// if the branch producing the cycle is reached first, then an error is returned, otherwise
		// a result is returned if some other terminal path of evaluation was reached before the cycle
		if err != nil {
			require.ErrorIs(t, err, ErrCycleDetected)
		} else {
			require.False(t, resp.GetAllowed())
			require.GreaterOrEqual(t, resp.ResolutionMetadata.DatastoreQueryCount, uint32(1)) // min of 1 (x) if x isn't found and it returns quickly
			require.LessOrEqual(t, resp.ResolutionMetadata.DatastoreQueryCount, uint32(3))    // max of 3 (x, y, z) before the cycle
		}
	}
}
				ContextualTuples: req.GetContextualTuples().GetTupleKeys(),
			}, reverseExpandResultsChan, resolutionMetadata)
			if err != nil {
				if errors.Is(err, graph.ErrResolutionDepthExceeded) || errors.Is(err, graph.ErrCycleDetected) {
					resultsChan <- ListObjectsResult{Err: serverErrors.AuthorizationModelResolutionTooComplex}
					return
				}

				resultsChan <- ListObjectsResult{Err: err}
			}

					},
				})
				if err != nil {
					if errors.Is(err, graph.ErrResolutionDepthExceeded) || errors.Is(err, graph.ErrCycleDetected) {
						resultsChan <- ListObjectsResult{Err: serverErrors.AuthorizationModelResolutionTooComplex}
						return
					}

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
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/storage/storagewrappers"
	"github.com/openfga/openfga/pkg/tuple"
		ctx = graph.ContextWithResolutionDepth(ctx, 0)
	} else {
		if depth >= c.resolveNodeLimit {
			return graph.ErrResolutionDepthExceeded
		}

		ctx = graph.ContextWithResolutionDepth(ctx, depth+1)
		},
	})
	if err != nil {
		if errors.Is(err, graph.ErrResolutionDepthExceeded) || errors.Is(err, graph.ErrCycleDetected) {
			return nil, serverErrors.AuthorizationModelResolutionTooComplex
		}

	}
}

func TestListObjects_ErrorCases(t *testing.T) {
	ctx := context.Background()
	store := ulid.Make().String()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	t.Run("database_errors", func(t *testing.T) {
		mockDatastore := mockstorage.NewMockOpenFGADatastore(mockController)

		s := MustNewServerWithOpts(
			WithDatastore(mockDatastore),
		)

		modelID := ulid.Make().String()

		mockDatastore.EXPECT().ReadAuthorizationModel(gomock.Any(), store, modelID).AnyTimes().Return(&openfgav1.AuthorizationModel{
			SchemaVersion: typesystem.SchemaVersion1_1,
			TypeDefinitions: parser.MustParse(`
			type user
	
			type document
			  relations
				define viewer: [user, user:*] as self
			`),
		}, nil)

		mockDatastore.EXPECT().ReadStartingWithUser(gomock.Any(), store, storage.ReadStartingWithUserFilter{
			ObjectType: "document",
			Relation:   "viewer",
			UserFilter: []*openfgav1.ObjectRelation{
				{Object: "user:*"},
				{Object: "user:bob"},
			}}).AnyTimes().Return(nil, errors.New("error reading from storage"))

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
	})

	t.Run("graph_resolution_errors", func(t *testing.T) {

		s := MustNewServerWithOpts(
			WithDatastore(memory.New()),
			WithResolveNodeLimit(2),
		)

		writeModelResp, err := s.WriteAuthorizationModel(ctx, &openfgav1.WriteAuthorizationModelRequest{
			StoreId:       store,
			SchemaVersion: typesystem.SchemaVersion1_1,
			TypeDefinitions: parser.MustParse(`
			type user

			type group
			  relations
			    define member: [user, group#member] as self

			type document
			  relations
				define viewer: [group#member] as self
			`),
		})
		require.NoError(t, err)

		_, err = s.Write(ctx, &openfgav1.WriteRequest{
			StoreId: store,
			Writes: &openfgav1.TupleKeys{
				TupleKeys: []*openfgav1.TupleKey{
					tuple.NewTupleKey("document:1", "viewer", "group:1#member"),
					tuple.NewTupleKey("group:1", "member", "group:2#member"),
					tuple.NewTupleKey("group:2", "member", "group:3#member"),
					tuple.NewTupleKey("group:3", "member", "user:jon"),
				},
			},
		})
		require.NoError(t, err)

		t.Run("resolution_depth_exceeded_error_unary", func(t *testing.T) {
			res, err := s.ListObjects(ctx, &openfgav1.ListObjectsRequest{
				StoreId:              store,
				AuthorizationModelId: writeModelResp.GetAuthorizationModelId(),
				Type:                 "document",
				Relation:             "viewer",
				User:                 "user:jon",
			})

			require.Nil(t, res)
			require.ErrorIs(t, err, serverErrors.AuthorizationModelResolutionTooComplex)
		})

		t.Run("resolution_depth_exceeded_error_streaming", func(t *testing.T) {
			err := s.StreamedListObjects(&openfgav1.StreamedListObjectsRequest{
				StoreId:              store,
				AuthorizationModelId: writeModelResp.GetAuthorizationModelId(),
				Type:                 "document",
				Relation:             "viewer",
				User:                 "user:jon",
			}, NewMockStreamServer())

			require.ErrorIs(t, err, serverErrors.AuthorizationModelResolutionTooComplex)
		})
	})
}

	parser "github.com/craigpastro/openfga-dsl-parser/v2"
	"github.com/oklog/ulid/v2"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/internal/graph"
	"github.com/openfga/openfga/pkg/server/commands/reverseexpand"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/tuple"
	"github.com/openfga/openfga/pkg/typesystem"
				tuple.NewTupleKey("folder:folder2", "parent", "folder:folder1"),
				tuple.NewTupleKey("folder:folder3", "parent", "folder:folder2"),
			},
			expectedError:        graph.ErrResolutionDepthExceeded,
			expectedDSQueryCount: 0,
		},
		{

import (
	"context"
	"fmt"
	"testing"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/pkg/logger"
	"github.com/openfga/openfga/pkg/server/commands"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/typesystem"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WriteAuthorizationModelTest(t *testing.T, datastore storage.OpenFGADatastore) {
		name          string
		request       *openfgav1.WriteAuthorizationModelRequest
		allowSchema10 bool
		errCode       codes.Code
	}{
		{
			name: "fails_if_too_many_types",
				SchemaVersion:   typesystem.SchemaVersion1_1,
			},
			allowSchema10: false,
			errCode:       codes.Code(openfgav1.ErrorCode_exceeded_entity_limit),
		},
		{
			name: "fails_if_a_relation_is_not_defined",
				SchemaVersion: typesystem.SchemaVersion1_1,
			},
			allowSchema10: false,
			errCode:       codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "Fails_if_type_info_metadata_is_omitted_in_1.1_model",
				},
			},
			allowSchema10: false,
			errCode:       codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "Fails_if_writing_1_0_model_because_it_will_be_interpreted_as_1_1",
				},
			},
			allowSchema10: true,
			errCode:       codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "Works_if_no_schema_version",
				`),
				SchemaVersion: typesystem.SchemaVersion1_1,
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "self_referencing_type_restriction_without_entrypoint_2",
				`),
				SchemaVersion: typesystem.SchemaVersion1_1,
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "self_referencing_type_restriction_without_entrypoint_3",
				`),
				SchemaVersion: typesystem.SchemaVersion1_1,
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "rewritten_relation_in_intersection_unresolvable",
				`),
				SchemaVersion: typesystem.SchemaVersion1_1,
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
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
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "no_entrypoint_3a",
				    define editor: [user] as self
				`),
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "no_entrypoint_3b",
			request: &openfgav1.WriteAuthorizationModelRequest{
				    define editor: [user] as self
				`),
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "no_entrypoint_4",
				    define viewer as editor from parent
				`),
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "self_referencing_type_restriction_with_entrypoint_1",
					},
				},
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "relation_name_is_empty_string",
					},
				},
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "many_circular_computed_relations",
			request: &openfgav1.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: parser.MustParse(`
				type user

				type canvas
				  relations
					define can_edit as editor or owner
					define editor: [user, account#member] as self
					define owner: [user] as self
					define viewer: [user, account#member] as self
	  
				type account
				  relations
					define admin: [user] as self or member or super_admin or owner
					define member: [user] as self or owner or admin or super_admin
					define owner: [user] as self
					define super_admin: [user] as self or admin or member
				`),
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "circular_relations_involving_intersection",
			request: &openfgav1.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: parser.MustParse(`
				type user

				type other
				  relations
					define x: [user] as self and y
					define y: [user] as self and z
					define z: [user] as self or x
				`),
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
		{
			name: "circular_relations_involving_exclusion",
			request: &openfgav1.WriteAuthorizationModelRequest{
				StoreId: storeID,
				TypeDefinitions: parser.MustParse(`
				type user

				type other
				  relations
					define x: [user] as self but not y
					define y: [user] as self but not z
					define z: [user] as self or x
				`),
			},
			errCode: codes.Code(openfgav1.ErrorCode_invalid_authorization_model),
		},
	}

		t.Run(test.name, func(t *testing.T) {
			cmd := commands.NewWriteAuthorizationModelCommand(datastore, logger)
			resp, err := cmd.Execute(ctx, test.request)
			status, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, test.errCode, status.Code())

			if err == nil {
				_, err = ulid.Parse(resp.AuthorizationModelId)
		}
	}

	hasCycle, err := t.HasCycle(typeName, relationName)
	if err != nil {
		return err
	}

	if hasCycle {
		return &InvalidRelationError{
			ObjectType: typeName,
			Relation:   relationName,
			Cause:      ErrCycle,
		}
	}

	return nil
}

	return fmt.Errorf("the relation type '%s' on '%s' in object type '%s' is not valid", relationType, relation, objectType)
}

func (t *TypeSystem) hasCycle(
	objectType, relationName string,
	rewrite *openfgav1.Userset,
	visited map[string]struct{},
) (bool, error) {

	visited[fmt.Sprintf("%s#%s", objectType, relationName)] = struct{}{}

	visitedCopy := maps.Clone(visited)

	var children []*openfgav1.Userset

	switch rw := rewrite.Userset.(type) {
	case *openfgav1.Userset_This, *openfgav1.Userset_TupleToUserset:
		return false, nil
	case *openfgav1.Userset_ComputedUserset:
		rewrittenRelation := rw.ComputedUserset.Relation

		if _, ok := visited[fmt.Sprintf("%s#%s", objectType, rewrittenRelation)]; ok {
			return true, nil
		}

		rewrittenRewrite, err := t.GetRelation(objectType, rewrittenRelation)
		if err != nil {
			return false, err
		}

		return t.hasCycle(objectType, rewrittenRelation, rewrittenRewrite.GetRewrite(), visitedCopy)
	case *openfgav1.Userset_Union:
		children = append(children, rw.Union.GetChild()...)
	case *openfgav1.Userset_Intersection:
		children = append(children, rw.Intersection.GetChild()...)
	case *openfgav1.Userset_Difference:
		children = append(children, rw.Difference.GetBase(), rw.Difference.GetSubtract())
	}

	for _, child := range children {

		hasCycle, err := t.hasCycle(objectType, relationName, child, visitedCopy)
		if err != nil {
			return false, err
		}

		if hasCycle {
			return true, nil
		}
	}

	return false, nil
}

// HasCycle runs a cycle detection test on the provided `objectType#relation` to see if the relation
// defines a rewrite rule that is self-referencing in any way (through computed relationships).
func (t *TypeSystem) HasCycle(objectType, relationName string) (bool, error) {
	visited := map[string]struct{}{}

	relation, err := t.GetRelation(objectType, relationName)
	if err != nil {
		return false, err
	}

	return t.hasCycle(objectType, relationName, relation.GetRewrite(), visited)
}

// getAllTupleToUsersetsDefinitions returns a map where the key is the object type and the value
// is another map where key=relationName, value=list of tuple to usersets declared in that relation
func (t *TypeSystem) getAllTupleToUsersetsDefinitions() map[string]map[string][]*openfgav1.TupleToUserset {
	"github.com/stretchr/testify/require"
)

func TestHasCycle(t *testing.T) {

	tests := []struct {
		name       string
		model      string
		objectType string
		relation   string
		expected   bool
	}{
		{
			name: "test_1",
			model: `
			type resource
			  relations
			    define x as y
			    define y as x
			`,
			objectType: "resource",
			relation:   "x",
			expected:   true,
		},
		{
			name: "test_2",
			model: `
			type resource
			  relations
			    define x as y
			    define y as z
				define z as x
			`,
			objectType: "resource",
			relation:   "y",
			expected:   true,
		},
		{
			name: "test_3",
			model: `
			type user

			type resource
			  relations
			    define x: [user] as self or y
			    define y: [user] as self or z
				define z: [user] as self or x
			`,
			objectType: "resource",
			relation:   "z",
			expected:   true,
		},
		{
			name: "test_4",
			model: `
			type user

			type resource
			  relations
			    define x: [user] as self or y
			    define y: [user] as self or z
				define z: [user] as self or x
			`,
			objectType: "resource",
			relation:   "z",
			expected:   true,
		},
		{
			name: "test_5",
			model: `
			type user

			type resource
			  relations
				define x: [user] as self but not y
				define y: [user] as self but not z
				define z: [user] as self or x
			`,
			objectType: "resource",
			relation:   "x",
			expected:   true,
		},
		{
			name: "test_6",
			model: `
			type user

			type group
			  relations
				define member: [user] as self or memberA or memberB or memberC
				define memberA: [user] as self or member or memberB or memberC
				define memberB: [user] as self or member or memberA or memberC
				define memberC: [user] as self or member or memberA or memberB
			`,
			objectType: "group",
			relation:   "member",
			expected:   true,
		},
		{
			name: "test_7",
			model: `
			type user

			type account
			relations
				define admin: [user] as self or member or super_admin or owner
				define member: [user] as self or owner or admin or super_admin
				define super_admin: [user] as self or admin or member or owner
				define owner: [user] as self
			`,
			objectType: "account",
			relation:   "member",
			expected:   true,
		},
		{
			name: "test_8",
			model: `
			type user

			type account
			relations
				define admin: [user] as self or member or super_admin or owner
				define member: [user] as self or owner or admin or super_admin
				define super_admin: [user] as self or admin or member or owner
				define owner: [user] as self
			`,
			objectType: "account",
			relation:   "owner",
			expected:   false,
		},
		{
			name: "test_9",
			model: `
			type user

			type document
			  relations
				define editor: [user] as self
				define viewer: [document#viewer] as self or editor
			`,
			objectType: "document",
			relation:   "viewer",
			expected:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			typesys := New(&openfgav1.AuthorizationModel{
				SchemaVersion:   SchemaVersion1_1,
				TypeDefinitions: parser.MustParse(test.model),
			})

			hasCycle, err := typesys.HasCycle(test.objectType, test.relation)
			require.Equal(t, test.expected, hasCycle)
			require.NoError(t, err)
		})
	}
}

func TestNewAndValidate(t *testing.T) {

	tests := []struct {
