package main

	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/openfga/openfga/internal/condition"
	"github.com/openfga/openfga/internal/condition/eval"
	serverconfig "github.com/openfga/openfga/internal/server/config"
	"github.com/openfga/openfga/internal/validation"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/telemetry"
	"github.com/openfga/openfga/pkg/tuple"
	delegate           CheckResolver
	concurrencyLimit   uint32
	maxConcurrentReads uint32
}

type LocalCheckerOption func(d *LocalChecker)
	}
}

// WithMaxConcurrentReads see server.WithMaxConcurrentReadsForCheck.
func WithMaxConcurrentReads(limit uint32) LocalCheckerOption {
	return func(d *LocalChecker) {
	}
}

// NewLocalChecker constructs a LocalChecker that can be used to evaluate a Check
// request locally.
//
// The constructed LocalChecker is not wrapped with cycle detection. Developers
// wanting a LocalChecker without other wrapped layers (e.g caching and others)
// are encouraged to use [[NewOrderedCheckResolvers]] instead.
func NewLocalChecker(opts ...LocalCheckerOption) *LocalChecker {
	checker := &LocalChecker{
		concurrencyLimit:   serverconfig.DefaultResolveNodeBreadthLimit,
		maxConcurrentReads: serverconfig.DefaultMaxConcurrentReadsForCheck,
	}
	// by default, a LocalChecker delegates/dispatchs subproblems to itself (e.g. local dispatch) unless otherwise configured.
	checker.delegate = checker
// handled concurrently relative to one another.
func exclusion(ctx context.Context, concurrencyLimit uint32, handlers ...CheckHandlerFunc) (*ResolveCheckResponse, error) {
	if len(handlers) != 2 {
		panic(fmt.Sprintf("expected two rewrite operands for exclusion operator, but got '%d'", len(handlers)))
	}

	span := trace.SpanFromContext(ctx)
var _ CheckResolver = (*LocalChecker)(nil)

// ResolveCheck implements [[CheckResolver.ResolveCheck]].
// If the typesystem isn't set in the context, it will panic.
func (c *LocalChecker) ResolveCheck(
	ctx context.Context,
	req *ResolveCheckRequest,
	ctx, span := tracer.Start(ctx, "ResolveCheck", trace.WithAttributes(
		attribute.String("store_id", req.GetStoreID()),
		attribute.String("resolver_type", "LocalChecker"),
		attribute.String("tuple_key", req.GetTupleKey().String()),
	))
	defer span.End()

		return nil, ErrResolutionDepthExceeded
	}

	typesys, ok := typesystem.TypesystemFromContext(ctx)
	if !ok {
		panic("typesystem missing in context")
	}

	tupleKey := req.GetTupleKey()
		}, nil
	}

	objectType, _ := tuple.SplitObject(object)
	rel, err := typesys.GetRelation(objectType, relation)
	if err != nil {
	return resp, nil
}

// usersetsMapType is a map where the key is object#relation and the value is a sorted set (no duplicates allowed).
// For example, given [group:1#member, group:2#member, group:1#owner, group:3#owner] it will be stored as:
// [group#member][1, 2]
// [group#owner][1, 3].
type usersetsMapType map[string]storage.SortedSet

func (c *LocalChecker) buildCheckAssociatedObjects(req *ResolveCheckRequest, objectRel string, objectIDs storage.SortedSet) CheckHandlerFunc {
	return func(ctx context.Context) (*ResolveCheckResponse, error) {
		ctx, span := tracer.Start(ctx, "checkAssociatedObjects")
		defer span.End()

		typesys, ok := typesystem.TypesystemFromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("typesystem missing in context")
		}

		ds, ok := storage.RelationshipTupleReaderFromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("relationship tuple reader datastore missing in context")
		}

		storeID := req.GetStoreID()
		reqTupleKey := req.GetTupleKey()

// checkUsersetSlowPath will check userset or public wildcard path.
// This is the slow path as it requires dispatch on all its children.
func (c *LocalChecker) checkUsersetSlowPath(ctx context.Context, iter *storage.ConditionsFilteredTupleKeyIterator, req *ResolveCheckRequest) (*ResolveCheckResponse, error) {
	ctx, span := tracer.Start(ctx, "checkUsersetSlowPath")
	defer span.End()
	var handlers []CheckHandlerFunc

	typesys, ok := typesystem.TypesystemFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("typesystem missing in context")
	}

	response := &ResolveCheckResponse{
		Allowed: false,
	}
}

// checkUsersetFastPath is the fast path to evaluate userset.
// The general idea of the algorithm is that it tries to find intersection on the objects as identified in the userset
// with the objects the user has the specified relation with.
// Finally, find the intersection between the two.
// To use the fast path, we will need to ensure that the userset and all the children associated with the userset are
// exclusively directly assignable. In our case, group member must be directly exclusively assignable.
func (c *LocalChecker) checkUsersetFastPath(ctx context.Context, iter *storage.ConditionsFilteredTupleKeyIterator, req *ResolveCheckRequest) (*ResolveCheckResponse, error) {
	ctx, span := tracer.Start(ctx, "checkUsersetFastPath")
	defer span.End()

	typesys, _ := typesystem.TypesystemFromContext(ctx)
	// We had just checked for the existence of the typesys in the caller.
	// So, we are guaranteed for the presence of the context.

	reqUserType := tuple.GetType(req.GetTupleKey().GetUser())

	usersetsMap := make(usersetsMapType)

	for {
		// NOTE: For the future, once we observe a new ObjectRelation in the usersetsMap that means that all objectIDs
		// needed for the intersection lookup for that ObjectRelation are in memory. There is no need to build the full
		// usersetMap by draining the full iterator, it can start processing in batches of ObjectRelation.
		// Consider doing it when we have better concurrency tooling for "streaming" operations.
		t, err := iter.Next(ctx)
		if err != nil {
			if errors.Is(err, storage.ErrIteratorDone) {
				break
			}
			telemetry.TraceError(span, err)
			return nil, err
		}

		object, relation := tuple.SplitObjectRelation(t.GetUser())
		objectType, objectID := tuple.SplitObject(object)
		terminalRelations := typesys.GetTerminalRelations(objectType, relation, reqUserType)
		// the terminalRelations is expected to be 1 (as we checked earlier in typesys.UsersetCanFastPath)
		if len(terminalRelations) != 1 {
			return nil, fmt.Errorf("expected exactly one terminal relation for fast path, received %d", len(terminalRelations))
		}
		computedRelation := terminalRelations[0]
		objectRel := tuple.ToObjectRelationString(objectType, computedRelation)
		if _, ok := usersetsMap[objectRel]; !ok {
			usersetsMap[objectRel] = storage.NewSortedSet()
		}
		usersetsMap[objectRel].Add(objectID)
	}

	// Next, for all the ObjectRelation, compare the associated objectIDs
	// to the users associated objects
	// all of this can likely bee its own function
	handlers := make([]CheckHandlerFunc, 0, len(usersetsMap))
	for objectRel, objectIDs := range usersetsMap {
		handler := c.buildCheckAssociatedObjects(req, objectRel, objectIDs)
		handlers = append(handlers, handler)
	}
	resp, err := union(ctx, c.concurrencyLimit, handlers...)
	if err != nil {
		telemetry.TraceError(span, err)
		return nil, err
	}

	resp.GetResolutionMetadata().DatastoreQueryCount++

	return resp, nil
}

// checkDirect composes two CheckHandlerFunc which evaluate direct relationships with the provided
			return nil, ctx.Err()
		}

		typesys, ok := typesystem.TypesystemFromContext(parentctx) // note: use of 'parentctx' not 'ctx' - this is important
		if !ok {
			return nil, fmt.Errorf("typesystem missing in context")
		}

		ds, ok := storage.RelationshipTupleReaderFromContext(parentctx)
		if !ok {
			return nil, fmt.Errorf("relationship tuple reader datastore missing in context")
		}

		storeID := req.GetStoreID()
		reqTupleKey := req.GetTupleKey()

		// TODO(jpadilla): can we lift this function up?
		checkDirectUserTuple := func(ctx context.Context) (*ResolveCheckResponse, error) {
			ctx, span := tracer.Start(ctx, "checkDirectUserTuple", trace.WithAttributes(attribute.String("tuple_key", reqTupleKey.String())))
			defer span.End()

			response := &ResolveCheckResponse{
			)
			defer filteredIter.Stop()

			resolver := c.checkUsersetSlowPath

			if typesys.UsersetCanFastPath(directlyRelatedUsersetTypes, tuple.GetType(reqTupleKey.GetUser())) {
				resolver = c.checkUsersetFastPath
			}

			return resolver(ctx, filteredIter, req)
		}

		var checkFuncs []CheckHandlerFunc
}

// checkComputedUserset evaluates the Check request with the rewritten relation (e.g. the computed userset relation).
func (c *LocalChecker) checkComputedUserset(ctx context.Context, req *ResolveCheckRequest, rewrite *openfgav1.Userset) CheckHandlerFunc {
	_, span := tracer.Start(ctx, "checkComputedUserset")
	defer span.End()

	seen := map[string]struct{}{}
	typesys, _ := typesystem.TypesystemFromContext(ctx)

	tk := tuple.NewTupleKey(
		req.GetTupleKey().GetObject(),
		req.GetTupleKey().GetRelation(),
		req.GetTupleKey().GetUser(),
	)

	for {
		tk = tuple.NewTupleKey(
			tk.GetObject(),
			rewrite.GetComputedUserset().GetRelation(),
			tk.GetUser(),
		)
		key := tuple.TupleKeyToString(tk)
		if _, cycleDetected := seen[key]; cycleDetected {
			return func(ctx context.Context) (*ResolveCheckResponse, error) {
				return &ResolveCheckResponse{
					Allowed: false,
					ResolutionMetadata: &ResolveCheckResponseMetadata{
						CycleDetected: true,
					},
				}, nil
			}
		}
		seen[key] = struct{}{}

		userObject, userRelation := tuple.SplitObjectRelation(tk.GetUser())

		// Check(document:1#viewer@document:1#viewer) will always return true
		if tk.GetRelation() == userRelation && tk.GetObject() == userObject {
			return func(ctx context.Context) (*ResolveCheckResponse, error) {
				return &ResolveCheckResponse{
					Allowed: true,
					ResolutionMetadata: &ResolveCheckResponseMetadata{
						DatastoreQueryCount: req.GetRequestMetadata().DatastoreQueryCount,
					},
				}, nil
			}
		}

		rw, err := typesys.GetRelation(tuple.GetType(tk.GetObject()), tk.GetRelation())
		if err != nil {
			return func(ctx context.Context) (*ResolveCheckResponse, error) {
				return nil, fmt.Errorf("relation '%s' undefined for object type '%s'", tk.GetRelation(), tuple.GetType(tk.GetObject()))
			}
		}
		rewrite = rw.GetRewrite()
		if _, isComputed := rewrite.GetUserset().(*openfgav1.Userset_ComputedUserset); !isComputed {
			break
		}
	}

	childRequest := clone(req)
	childRequest.TupleKey = tk

	return c.checkRewrite(ctx, childRequest, rewrite)
}

// checkTTUSlowPath is the slow path for checkTTU where we cannot short-circuit TTU evaluation and

	var handlers []CheckHandlerFunc

	typesys, ok := typesystem.TypesystemFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("typesystem missing in context")
	}

	computedRelation := rewrite.GetTupleToUserset().GetComputedUserset().GetRelation()
	tk := req.GetTupleKey()
//
// check(user, viewer, doc) will find the intersection of all group assigned to the doc's parent AND
// all group where the user is a member of.

func (c *LocalChecker) checkTTUFastPath(ctx context.Context, req *ResolveCheckRequest, rewrite *openfgav1.Userset, iter *storage.ConditionsFilteredTupleKeyIterator) (*ResolveCheckResponse, error) {
	ctx, span := tracer.Start(ctx, "checkTTUFastPath")
	defer span.End()

	typesys, ok := typesystem.TypesystemFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("typesystem missing in context")
	}

	terminalRelations := typesys.GetTerminalRelations(
		tuple.GetType(req.GetTupleKey().GetObject()), req.GetTupleKey().GetRelation(), tuple.GetType(req.GetTupleKey().GetUser()),
	)
	if len(terminalRelations) != 1 {
		return nil, fmt.Errorf("expected exactly one terminal relation for fast path, received %d", len(terminalRelations))
	}

	computedRelation := terminalRelations[0]
	// usersetsMap is a map of all ObjectRelations and its Ids. For example,
	// [group:1#member, group:2#member, group:1#owner, group:3#owner] will be stored as
	// [group#member][1, 2]
	// [group#owner][1, 3]
	usersetsMap := make(usersetsMapType)

	for {
		t, err := iter.Next(ctx)
		if err != nil {
			if errors.Is(err, storage.ErrIteratorDone) {
				break
			}
			telemetry.TraceError(span, err)
			return nil, err
		}

		object, _ := tuple.SplitObjectRelation(t.GetUser())
		objectType, objectID := tuple.SplitObject(object)
		objectRel := tuple.ToObjectRelationString(objectType, computedRelation)
		if _, ok := usersetsMap[objectRel]; !ok {
			usersetsMap[objectRel] = storage.NewSortedSet()
		}
		usersetsMap[objectRel].Add(objectID)
	}

	// Next, for each of the type in tuplesetRelationUserMap, look up what object is in computedRelation for the specified user.
	// We will then try to see if there are any intersection.
	// Return true if user is in any of the computedRelation.  For example,
	// type group
	//   define member: [user]
	// type doc
	//   define parent: [group]
	//   define viewer: member from parent
	// we want to find out which group user:bob is a member of.
	// After that, we will find the intersection.
	handlers := make([]CheckHandlerFunc, 0, len(usersetsMap))
	for objectRel, objectIDs := range usersetsMap {
		handler := c.buildCheckAssociatedObjects(req, objectRel, objectIDs)
		handlers = append(handlers, handler)
	}

	resp, err := union(ctx, c.concurrencyLimit, handlers...)
	if err != nil {
		telemetry.TraceError(span, err)
		return nil, err
	}

	resp.GetResolutionMetadata().DatastoreQueryCount++

	return resp, nil
}

// checkTTU looks up all tuples of the target tupleset relation on the provided object and for each one
			return nil, ctx.Err()
		}

		typesys, ok := typesystem.TypesystemFromContext(parentctx) // note: use of 'parentctx' not 'ctx' - this is important
		if !ok {
			return nil, fmt.Errorf("typesystem missing in context")
		}

		ds, ok := storage.RelationshipTupleReaderFromContext(parentctx)
		if !ok {
			return nil, fmt.Errorf("relationship tuple reader datastore missing in context")
		}

		ctx = typesystem.ContextWithTypesystem(ctx, typesys)
		ctx = storage.ContextWithRelationshipTupleReader(ctx, ds)
		)
		defer filteredIter.Stop()

		resolver := c.checkTTUSlowPath

		// TODO: optimize the case where user is an userset.
		// If the user is a userset, we will not be able to use the shortcut because the algo
		// will look up the objects associated with user.
		if !tuple.IsObjectRelation(tk.GetUser()) {
			if canFastPath := typesys.TTUCanFastPath(
				tuple.GetType(object), req.GetTupleKey().GetRelation(), tuple.GetType(req.GetTupleKey().GetUser())); canFastPath {
				resolver = c.checkTTUFastPath
			}
		}

		return resolver(ctx, req, rewrite, filteredIter)
	}
}

			handlers = append(handlers, c.checkRewrite(ctx, req, child))
		}
	default:
		panic("unexpected set operator type encountered")
	}

	return func(ctx context.Context) (*ResolveCheckResponse, error) {
	case *openfgav1.Userset_Difference:
		return c.checkSetOperation(ctx, req, exclusionSetOperator, exclusion, rw.Difference.GetBase(), rw.Difference.GetSubtract())
	default:
		panic("unexpected userset rewrite encountered")
	}
}

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	parser "github.com/openfga/language/pkg/go/transformer"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openfga/openfga/pkg/server/commands"
	serverErrors "github.com/openfga/openfga/pkg/server/errors"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/testutils"
	"github.com/openfga/openfga/pkg/tuple"
	"github.com/openfga/openfga/pkg/typesystem"
)

func TestWriteAndReadAssertions(t *testing.T, datastore storage.OpenFGADatastore) {
	type writeAssertionsTestSettings struct {
		_name      string
		assertions []*openfgav1.Assertion
	}

	store := testutils.CreateRandomString(10)

	githubModelReq := &openfgav1.WriteAuthorizationModelRequest{
		StoreId: store,
		TypeDefinitions: parser.MustTransformDSLToProto(`
			model
				schema 1.1
			type user

			type repo
				relations
					define reader: [user]
					define can_read: reader`).GetTypeDefinitions(),
		SchemaVersion: typesystem.SchemaVersion1_1,
	}

	var tests = []writeAssertionsTestSettings{
		{
			_name: "writing_assertions_succeeds",
			assertions: []*openfgav1.Assertion{{
				TupleKey:    tuple.NewAssertionTupleKey("repo:test", "reader", "user:elbuo"),
				Expectation: false,
			}},
		},
		{
			_name: "writing_assertions_succeeds_when_it_is_not_directly_assignable",
			assertions: []*openfgav1.Assertion{{
				TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
				Expectation: false,
			}},
		},
		{
			_name: "writing_multiple_assertions_succeeds",
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "reader", "user:elbuo"),
			},
		},
		{
			_name: "writing_multiple_assertions_succeeds_when_it_is_not_directly_assignable",
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
			},
		},
		{
			_name:      "writing_empty_assertions_succeeds",
			assertions: []*openfgav1.Assertion{},
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test._name, func(t *testing.T) {
			model := githubModelReq

			writeAuthzModelCmd := commands.NewWriteAuthorizationModelCommand(datastore)

			modelID, err := writeAuthzModelCmd.Execute(ctx, model)
			require.NoError(t, err)
			request := &openfgav1.WriteAssertionsRequest{
				StoreId:              store,
				Assertions:           test.assertions,
				AuthorizationModelId: modelID.GetAuthorizationModelId(),
			}

			writeAssertionCmd := commands.NewWriteAssertionsCommand(datastore)
			_, err = writeAssertionCmd.Execute(ctx, request)
			require.NoError(t, err)
			query := commands.NewReadAssertionsQuery(datastore)
			actualResponse, actualError := query.Execute(ctx, store, modelID.GetAuthorizationModelId())
			require.NoError(t, actualError)

			expectedResponse := &openfgav1.ReadAssertionsResponse{
				AuthorizationModelId: modelID.GetAuthorizationModelId(),
				Assertions:           test.assertions,
			}
			if diff := cmp.Diff(expectedResponse, actualResponse, protocmp.Transform()); diff != "" {
				t.Errorf("store mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWriteAssertionsFailure(t *testing.T, datastore storage.OpenFGADatastore) {
	type writeAssertionsTestSettings struct {
		_name      string
		assertions []*openfgav1.Assertion
		modelID    string
		err        error
	}

	store := testutils.CreateRandomString(10)

	githubModelReq := &openfgav1.WriteAuthorizationModelRequest{
		StoreId: store,
		TypeDefinitions: parser.MustTransformDSLToProto(`
			model
				schema 1.1
			type user

			type repo
				relations
					define reader: [user]
					define can_read: reader`).GetTypeDefinitions(),
		SchemaVersion: typesystem.SchemaVersion1_1,
	}
	ctx := context.Background()

	writeAuthzModelCmd := commands.NewWriteAuthorizationModelCommand(datastore)
	modelID, err := writeAuthzModelCmd.Execute(ctx, githubModelReq)
	require.NoError(t, err)

	var tests = []writeAssertionsTestSettings{
		{
			_name: "writing_assertion_with_invalid_relation_fails",
			assertions: []*openfgav1.Assertion{
				{
					TupleKey: tuple.NewAssertionTupleKey(
					Expectation: false,
				},
			},
			modelID: modelID.GetAuthorizationModelId(),
			err: serverErrors.ValidationError(
				fmt.Errorf("relation 'repo#invalidrelation' not found"),
			),
		},
		{
			_name: "writing_assertion_with_not_found_id",
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
					Expectation: false,
				},
			},
			modelID: "not_valid_id",
			err: serverErrors.AuthorizationModelNotFound(
				"not_valid_id",
			),
		},
	}

	for _, test := range tests {
		t.Run(test._name, func(t *testing.T) {
			request := &openfgav1.WriteAssertionsRequest{
				StoreId:              store,
				Assertions:           test.assertions,
				AuthorizationModelId: test.modelID,
			}

			writeAssertionCmd := commands.NewWriteAssertionsCommand(datastore)
			_, err = writeAssertionCmd.Execute(ctx, request)
			require.ErrorIs(t, test.err, err)
		})
	}
}
