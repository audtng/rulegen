package main

	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/openfga/openfga/internal/concurrency"
	"github.com/openfga/openfga/internal/condition"
	"github.com/openfga/openfga/internal/condition/eval"
	openfgaErrors "github.com/openfga/openfga/internal/errors"
	serverconfig "github.com/openfga/openfga/internal/server/config"
	"github.com/openfga/openfga/internal/validation"
	"github.com/openfga/openfga/pkg/logger"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/telemetry"
	"github.com/openfga/openfga/pkg/tuple"
	delegate           CheckResolver
	concurrencyLimit   uint32
	maxConcurrentReads uint32
	usersetBatchSize   uint32
	logger             logger.Logger
}

type LocalCheckerOption func(d *LocalChecker)
	}
}

// WithUsersetBatchSize see server.WithUsersetBatchSize.
func WithUsersetBatchSize(usersetBatchSize uint32) LocalCheckerOption {
	return func(d *LocalChecker) {
		d.usersetBatchSize = usersetBatchSize
	}
}

// WithMaxConcurrentReads see server.WithMaxConcurrentReadsForCheck.
func WithMaxConcurrentReads(limit uint32) LocalCheckerOption {
	return func(d *LocalChecker) {
	}
}

func WithLocalCheckerLogger(logger logger.Logger) LocalCheckerOption {
	return func(d *LocalChecker) {
		d.logger = logger
	}
}

// NewLocalChecker constructs a LocalChecker that can be used to evaluate a Check
// request locally.
//
// Developers wanting a LocalChecker with other optional layers (e.g caching and others)
// are encouraged to use [[NewOrderedCheckResolvers]] instead.
func NewLocalChecker(opts ...LocalCheckerOption) *LocalChecker {
	checker := &LocalChecker{
		concurrencyLimit:   serverconfig.DefaultResolveNodeBreadthLimit,
		maxConcurrentReads: serverconfig.DefaultMaxConcurrentReadsForCheck,
		usersetBatchSize:   serverconfig.DefaultUsersetBatchSize,
		logger:             logger.NewNoopLogger(),
	}
	// by default, a LocalChecker delegates/dispatchs subproblems to itself (e.g. local dispatch) unless otherwise configured.
	checker.delegate = checker
// handled concurrently relative to one another.
func exclusion(ctx context.Context, concurrencyLimit uint32, handlers ...CheckHandlerFunc) (*ResolveCheckResponse, error) {
	if len(handlers) != 2 {
		return nil, fmt.Errorf("%w, expected two rewrite operands for exclusion operator, but got '%d'", openfgaErrors.ErrUnknown, len(handlers))
	}

	span := trace.SpanFromContext(ctx)
var _ CheckResolver = (*LocalChecker)(nil)

// ResolveCheck implements [[CheckResolver.ResolveCheck]].
func (c *LocalChecker) ResolveCheck(
	ctx context.Context,
	req *ResolveCheckRequest,
	ctx, span := tracer.Start(ctx, "ResolveCheck", trace.WithAttributes(
		attribute.String("store_id", req.GetStoreID()),
		attribute.String("resolver_type", "LocalChecker"),
		attribute.String("tuple_key", tuple.TupleKeyWithConditionToString(req.GetTupleKey())),
	))
	defer span.End()

		return nil, ErrResolutionDepthExceeded
	}

	cycle := c.hasCycle(req)
	if cycle {
		span.SetAttributes(attribute.Bool("cycle_detected", true))
		return &ResolveCheckResponse{
			Allowed: false,
			ResolutionMetadata: &ResolveCheckResponseMetadata{
				CycleDetected: true,
			},
		}, nil
	}

	tupleKey := req.GetTupleKey()
		}, nil
	}

	typesys, ok := typesystem.TypesystemFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%w: typesystem missing in context", openfgaErrors.ErrUnknown)
	}
	_, ok = storage.RelationshipTupleReaderFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%w: relationship tuple reader datastore missing in context", openfgaErrors.ErrUnknown)
	}

	objectType, _ := tuple.SplitObject(object)
	rel, err := typesys.GetRelation(objectType, relation)
	if err != nil {
	return resp, nil
}

// hasCycle returns true if a cycle has been found. It modifies the request object.
func (c *LocalChecker) hasCycle(req *ResolveCheckRequest) bool {
	key := tuple.TupleKeyToString(req.GetTupleKey())
	if req.VisitedPaths == nil {
		req.VisitedPaths = map[string]struct{}{}
	}

	_, cycleDetected := req.VisitedPaths[key]
	if cycleDetected {
		return true
	}

	req.VisitedPaths[key] = struct{}{}
	return false
}

// usersetsMapType is a map where the key is object#relation and the value is a sorted set (no duplicates allowed).
// For example, given [group:1#member, group:2#member, group:1#owner, group:3#owner] it will be stored as:
// [group#member][1, 2]
// [group#owner][1, 3].
// nolint:unused
type usersetsMapType map[string]storage.SortedSet

// nolint:unused
func (c *LocalChecker) buildCheckAssociatedObjects(req *ResolveCheckRequest, objectRel string, objectIDs storage.SortedSet) CheckHandlerFunc {
	return func(ctx context.Context) (*ResolveCheckResponse, error) {
		ctx, span := tracer.Start(ctx, "checkAssociatedObjects")
		defer span.End()

		typesys, _ := typesystem.TypesystemFromContext(ctx)

		ds, _ := storage.RelationshipTupleReaderFromContext(ctx)

		storeID := req.GetStoreID()
		reqTupleKey := req.GetTupleKey()

// checkUsersetSlowPath will check userset or public wildcard path.
// This is the slow path as it requires dispatch on all its children.
func (c *LocalChecker) checkUsersetSlowPath(ctx context.Context, req *ResolveCheckRequest, iter *storage.ConditionsFilteredTupleKeyIterator) (*ResolveCheckResponse, error) {
	ctx, span := tracer.Start(ctx, "checkUsersetSlowPath")
	defer span.End()
	var handlers []CheckHandlerFunc

	typesys, _ := typesystem.TypesystemFromContext(ctx)

	response := &ResolveCheckResponse{
		Allowed: false,
	}
}

// nolint:unused
type usersetDetailsFunc func(*openfgav1.TupleKey) (string, string, error)

// buildUsersetDetails given tuple doc:1#viewer@group:2#member will return group#member, 2, nil.
// This util takes into account pre-computed relationships, otherwise it will resolve it from the target UserType.
// nolint:unused
func buildUsersetDetails(typesys *typesystem.TypeSystem, computedRelation, userType string) usersetDetailsFunc {
	return func(t *openfgav1.TupleKey) (string, string, error) {
		cr := computedRelation
		object, relation := tuple.SplitObjectRelation(t.GetUser())
		objectType, objectID := tuple.SplitObject(object)
		if cr == "" {
			terminalRelations := typesys.GetTerminalRelations(objectType, relation, userType)
			cr = terminalRelations[0]
		}

		return tuple.ToObjectRelationString(objectType, cr), objectID, nil
	}
}

// checkUsersetFastPath is the fast path to evaluate userset.
// The general idea of the algorithm is that it tries to find intersection on the objects as identified in the userset
// with the objects the user has the specified relation with.
// Finally, find the intersection between the two.
// To use the fast path, we will need to ensure that the userset and all the children associated with the userset are
// exclusively directly assignable. In our case, group member must be directly exclusively assignable.
// nolint:unused
func (c *LocalChecker) checkUsersetFastPath(ctx context.Context, req *ResolveCheckRequest, iter *storage.ConditionsFilteredTupleKeyIterator) (*ResolveCheckResponse, error) {
	ctx, span := tracer.Start(ctx, "checkUsersetFastPath")
	defer span.End()
	// Caller already verified typesys
	typesys, _ := typesystem.TypesystemFromContext(ctx)
	usersetDetails := buildUsersetDetails(typesys, "", tuple.GetType(req.GetTupleKey().GetUser()))
	return c.checkMembership(ctx, req, iter, usersetDetails)
}

// nolint:unused
type usersetsChannelType struct {
	err            error
	objectRelation string            // e.g. group#member
	objectIDs      storage.SortedSet // eg. [1,2,3] (no duplicates allowed, sorted)
}

// checkMembership for this model
//
// type user
// type org
//
//	relations
//		define viewer: [user]
//
// type folder
//
//	relations
//		define viewer: [user]
//
// type doc
//
//	relations
//		define viewer: viewer from parent
//		define parent: [folder, org]
//
// works as follows.
// If the request is Check(user:maria, viewer, doc:1).
// 1. We build a map with folder#viewer:[1...N], org#viewer:[1...M] that are parents of doc:1. We send those through a channel.
// 2. The consumer of the channel finds all the folders (and orgs) by looking at tuples of the form folder:X#viewer@user:maria (and org:Y#viewer@user:maria).
// 3. If there is one folder or org found in step (2) that appears in the map found in step (1), it returns allowed=true immediately.
// nolint:unused
func (c *LocalChecker) checkMembership(ctx context.Context, req *ResolveCheckRequest, iter *storage.ConditionsFilteredTupleKeyIterator, usersetDetails usersetDetailsFunc) (*ResolveCheckResponse, error) {
	ctx, span := tracer.Start(ctx, "checkMembership")
	defer span.End()

	// since this is an unbuffered channel, producer will be blocked until consumer catches up
	// TODO: when implementing set math operators, change to buffered. consider using the number of sets as the concurrency limit
	usersetsChan := make(chan usersetsChannelType)

	cancellableCtx, cancelFunc := context.WithCancel(ctx)
	// sending to channel in batches up to a pre-configured value to subsequently checkMembership for.
	pool := concurrency.NewPool(cancellableCtx, 1)
	defer func() {
		cancelFunc()
		// We need to wait always to avoid a goroutine leak.
		_ = pool.Wait()
	}()
	pool.Go(func(ctx context.Context) error {
		c.produceUsersets(ctx, usersetsChan, iter, usersetDetails)
		return nil
	})

	resp, err := c.consumeUsersets(ctx, req, usersetsChan)
	if err != nil {
		telemetry.TraceError(span, err)
	}

	return resp, err
}

// nolint:unused
func (c *LocalChecker) consumeUsersets(ctx context.Context, req *ResolveCheckRequest, usersetsChan chan usersetsChannelType) (*ResolveCheckResponse, error) {
	var finalErr error
	dbReads := req.GetRequestMetadata().DatastoreQueryCount

ConsumerLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case newBatch, channelOpen := <-usersetsChan:
			if !channelOpen {
				break ConsumerLoop
			}
			if newBatch.err != nil {
				// Irrecoverable error when fetching usersets, so we abort.
				finalErr = newBatch.err
				break ConsumerLoop
			}
			objectRel := newBatch.objectRelation
			objectIDs := newBatch.objectIDs

			resp, err := c.buildCheckAssociatedObjects(req, objectRel, objectIDs)(ctx)
			dbReads++
			if err != nil {
				// We don't exit because we do a best effort to find the objectId that will give `allowed=true`.
				// If that doesn't happen, we will return this error down below.
				finalErr = err
			} else if resp.Allowed {
				resp.ResolutionMetadata.DatastoreQueryCount = dbReads
				return resp, nil
			}
		}
	}

	// context cancellation from upstream (e.g. client)
	if ctx.Err() != nil {
		finalErr = ctx.Err()
	}

	if finalErr != nil {
		return nil, finalErr
	}

	return &ResolveCheckResponse{
		Allowed: false,
		ResolutionMetadata: &ResolveCheckResponseMetadata{
			DatastoreQueryCount: dbReads,
		},
	}, nil
}

// nolint:unused
func (c *LocalChecker) produceUsersets(ctx context.Context, usersetsChan chan usersetsChannelType, iter *storage.ConditionsFilteredTupleKeyIterator, usersetDetails usersetDetailsFunc) {
	usersetsMap := make(usersetsMapType)
	defer close(usersetsChan)
	for {
		t, err := iter.Next(ctx)
		if err != nil {
			// cancelled doesn't need to flush nor send errors back to main routine
			if !errors.Is(err, storage.ErrIteratorDone) && !errors.Is(err, context.Canceled) {
				trySendUsersetsError(ctx, err, usersetsChan)
			}
			break
		}

		objectRel, objectID, err := usersetDetails(t)
		if err != nil {
			trySendUsersetsError(ctx, err, usersetsChan)
			break
		}

		if _, ok := usersetsMap[objectRel]; !ok {
			if len(usersetsMap) > 0 {
				// Flush results from a previous objectRel it begin processing immediately.
				// The assumption (which may not be true) is that the datastore yields objectRel in order.
				trySendUsersetsAndDeleteFromMap(ctx, usersetsMap, usersetsChan)
			}
			usersetsMap[objectRel] = storage.NewSortedSet()
		}

		usersetsMap[objectRel].Add(objectID)

		if usersetsMap[objectRel].Size() > int(c.usersetBatchSize) {
			trySendUsersetsAndDeleteFromMap(ctx, usersetsMap, usersetsChan)
		}
	}

	trySendUsersetsAndDeleteFromMap(ctx, usersetsMap, usersetsChan)
}

// nolint:unused
func trySendUsersetsError(ctx context.Context, err error, errorChan chan usersetsChannelType) {
	select {
	case <-ctx.Done():
	case errorChan <- usersetsChannelType{err: err}:
	}
}

// nolint:unused
func trySendUsersetsAndDeleteFromMap(ctx context.Context, usersetsMap usersetsMapType, usersetsChan chan usersetsChannelType) {
	for k, v := range usersetsMap {
		select {
		case <-ctx.Done():
			return
		case usersetsChan <- usersetsChannelType{
			objectRelation: k,
			objectIDs:      v,
		}:
			delete(usersetsMap, k)
		}
	}
}

// checkDirect composes two CheckHandlerFunc which evaluate direct relationships with the provided
			return nil, ctx.Err()
		}

		typesys, _ := typesystem.TypesystemFromContext(parentctx) // note: use of 'parentctx' not 'ctx' - this is important

		ds, _ := storage.RelationshipTupleReaderFromContext(parentctx)

		storeID := req.GetStoreID()
		reqTupleKey := req.GetTupleKey()

		// TODO(jpadilla): can we lift this function up?
		checkDirectUserTuple := func(ctx context.Context) (*ResolveCheckResponse, error) {
			ctx, span := tracer.Start(ctx, "checkDirectUserTuple",
				trace.WithAttributes(attribute.String("tuple_key", tuple.TupleKeyWithConditionToString(reqTupleKey))))
			defer span.End()

			response := &ResolveCheckResponse{
			)
			defer filteredIter.Stop()

			return c.checkUsersetSlowPath(ctx, req, filteredIter)
		}

		var checkFuncs []CheckHandlerFunc
}

// checkComputedUserset evaluates the Check request with the rewritten relation (e.g. the computed userset relation).
func (c *LocalChecker) checkComputedUserset(_ context.Context, req *ResolveCheckRequest, rewrite *openfgav1.Userset) CheckHandlerFunc {
	rewrittenTupleKey := tuple.NewTupleKey(
		req.GetTupleKey().GetObject(),
		rewrite.GetComputedUserset().GetRelation(),
		req.GetTupleKey().GetUser(),
	)

	childRequest := clone(req)
	childRequest.TupleKey = rewrittenTupleKey

	return func(ctx context.Context) (*ResolveCheckResponse, error) {
		ctx, span := tracer.Start(ctx, "checkComputedUserset")
		defer span.End()
		// No dispatch here, as we don't want to increase resolution depth.
		return c.ResolveCheck(ctx, childRequest)
	}
}

// checkTTUSlowPath is the slow path for checkTTU where we cannot short-circuit TTU evaluation and

	var handlers []CheckHandlerFunc

	typesys, _ := typesystem.TypesystemFromContext(ctx)

	computedRelation := rewrite.GetTupleToUserset().GetComputedUserset().GetRelation()
	tk := req.GetTupleKey()
//
// check(user, viewer, doc) will find the intersection of all group assigned to the doc's parent AND
// all group where the user is a member of.
// nolint:unused
func (c *LocalChecker) checkTTUFastPath(ctx context.Context, req *ResolveCheckRequest, _ *openfgav1.Userset, iter *storage.ConditionsFilteredTupleKeyIterator) (*ResolveCheckResponse, error) {
	ctx, span := tracer.Start(ctx, "checkTTUFastPath")
	defer span.End()
	// Caller already verified typesys
	typesys, _ := typesystem.TypesystemFromContext(ctx)

	terminalRelations := typesys.GetTerminalRelations(tuple.GetType(req.GetTupleKey().GetObject()), req.GetTupleKey().GetRelation(), tuple.GetType(req.GetTupleKey().GetUser()))

	computedRelation := terminalRelations[0]
	usersetDetails := buildUsersetDetails(typesys, computedRelation, tuple.GetType(req.GetTupleKey().GetUser()))

	return c.checkMembership(ctx, req, iter, usersetDetails)
}

// checkTTU looks up all tuples of the target tupleset relation on the provided object and for each one
			return nil, ctx.Err()
		}

		typesys, _ := typesystem.TypesystemFromContext(parentctx) // note: use of 'parentctx' not 'ctx' - this is important

		ds, _ := storage.RelationshipTupleReaderFromContext(parentctx)

		ctx = typesystem.ContextWithTypesystem(ctx, typesys)
		ctx = storage.ContextWithRelationshipTupleReader(ctx, ds)
		)
		defer filteredIter.Stop()

		return c.checkTTUSlowPath(ctx, req, rewrite, filteredIter)
	}
}

			handlers = append(handlers, c.checkRewrite(ctx, req, child))
		}
	default:
		return func(ctx context.Context) (*ResolveCheckResponse, error) {
			return nil, fmt.Errorf("%w: unexpected set operator type encountered", openfgaErrors.ErrUnknown)
		}
	}

	return func(ctx context.Context) (*ResolveCheckResponse, error) {
	case *openfgav1.Userset_Difference:
		return c.checkSetOperation(ctx, req, exclusionSetOperator, exclusion, rw.Difference.GetBase(), rw.Difference.GetSubtract())
	default:
		return func(ctx context.Context) (*ResolveCheckResponse, error) {
			return nil, fmt.Errorf("%w: unexpected set operator type encountered", openfgaErrors.ErrUnknown)
		}
	}
}

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/oklog/ulid/v2"
	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	parser "github.com/openfga/language/pkg/go/transformer"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openfga/openfga/pkg/testutils"

	"github.com/openfga/openfga/pkg/server/commands"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/tuple"
)

func TestWriteAndReadAssertions(t *testing.T, datastore storage.OpenFGADatastore) {
	store := ulid.Make().String()

	model := parser.MustTransformDSLToProto(`
	model
		schema 1.1
	type user

	type repo
		relations
			define reader: [user, user with condX]
			define can_read: reader

	condition condX(x :int) {
		x > 0
	}
	`)

	writeAuthzModelCmd := commands.NewWriteAuthorizationModelCommand(datastore)

	writeModelResponse, err := writeAuthzModelCmd.Execute(context.Background(), &openfgav1.WriteAuthorizationModelRequest{
		StoreId:         store,
		TypeDefinitions: model.GetTypeDefinitions(),
		SchemaVersion:   model.GetSchemaVersion(),
		Conditions:      model.GetConditions(),
	})
	require.NoError(t, err)
	modelID := writeModelResponse.GetAuthorizationModelId()

	type writeAssertionsTestSettings struct {
		_name                string
		inputModelID         string
		assertions           []*openfgav1.Assertion
		expectErrWhenWriting string
	}

	var tests = []writeAssertionsTestSettings{
		{
			_name:        "writing_assertion_succeeds",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{{
				TupleKey:    tuple.NewAssertionTupleKey("repo:test", "reader", "user:elbuo"),
				Expectation: false,
			}},
		},
		{
			_name:        "writing_assertion_succeeds_when_it_is_not_directly_assignable",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{{
				TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
				Expectation: false,
			}},
		},
		{
			_name:        "writing_multiple_assertions_succeeds",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "reader", "user:elbuo"),
			},
		},
		{
			_name:        "writing_multiple_assertions_succeeds_when_it_is_not_directly_assignable",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
			},
		},
		{
			_name:        "writing_empty_assertions_succeeds",
			inputModelID: modelID,
			assertions:   []*openfgav1.Assertion{},
		},
		{
			_name:        "writing_assertion_with_contextual_tuple_succeeds",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
					Expectation: false,
					ContextualTuples: []*openfgav1.TupleKey{
						tuple.NewTupleKey("repo:test", "reader", "user:elbuo"),
					},
				},
			},
		},
		{
			_name:        "writing_assertion_with_contextual_tuple_with_condition_succeeds",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
					Expectation: false,
					ContextualTuples: []*openfgav1.TupleKey{
						tuple.NewTupleKeyWithCondition("repo:test", "reader", "user:elbuo", "condX",
							testutils.MustNewStruct(t, map[string]interface{}{"x": 0})),
					},
				},
			},
		},
		{
			_name:        "writing_assertion_with_contextual_tuple_with_condition_fails_because_invalid_context_parameter",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
					Expectation: false,
					ContextualTuples: []*openfgav1.TupleKey{
						tuple.NewTupleKeyWithCondition("repo:test", "reader", "user:elbuo", "condX",
							testutils.MustNewStruct(t, map[string]interface{}{"unknownparam": 0})),
					},
				},
			},
			expectErrWhenWriting: "Invalid tuple 'repo:test#reader@user:elbuo (condition condX)'. Reason: found invalid context parameter: unknownparam",
		},
		{
			_name:        "writing_assertion_with_contextual_tuple_with_condition_fails_because_undefined_condition",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
					Expectation: false,
					ContextualTuples: []*openfgav1.TupleKey{
						tuple.NewTupleKeyWithCondition("repo:test", "reader", "user:elbuo", "condundefined", nil),
					},
				},
			},
			expectErrWhenWriting: "Invalid tuple 'repo:test#reader@user:elbuo (condition condundefined)'. Reason: undefined condition",
		},
		{
			_name:        "writing_assertion_with_contextual_tuple_fails_because_contextual_tuple_is_not_directly_assignable",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
					Expectation: false,
					ContextualTuples: []*openfgav1.TupleKey{
						tuple.NewTupleKey("repo:test", "can_read", "user:elbuo"),
					},
				},
			},
			expectErrWhenWriting: "Invalid tuple 'repo:test#can_read@user:elbuo'. Reason: type 'user' is not an allowed type restriction for 'repo#can_read'",
		},
		{
			_name:        "writing_assertion_with_contextual_tuple_fails_because_invalid_relation_in_contextual_tuple",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
					Expectation: false,
					ContextualTuples: []*openfgav1.TupleKey{
						tuple.NewTupleKey("repo:test", "invalidrelation", "user:elbuo"),
					},
				},
			},
			expectErrWhenWriting: "Invalid tuple 'repo:test#invalidrelation@user:elbuo'. Reason: relation 'repo#invalidrelation' not found",
		},
		{
			_name:        "writing_assertion_with_contextual_tuple_fails_because_invalid_type_in_contextual_tuple",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
					Expectation: false,
					ContextualTuples: []*openfgav1.TupleKey{
						tuple.NewTupleKey("unknown:test", "reader", "user:elbuo"),
					},
				},
			},
			expectErrWhenWriting: "Invalid tuple 'unknown:test#reader@user:elbuo'. Reason: type 'unknown' not found",
		},
		{
			_name:        "writing_assertion_with_invalid_relation_fails",
			inputModelID: modelID,
			assertions: []*openfgav1.Assertion{
				{
					TupleKey: tuple.NewAssertionTupleKey(
					Expectation: false,
				},
			},
			expectErrWhenWriting: "relation 'repo#invalidrelation' not found",
		},
		{
			_name:        "writing_assertion_with_invalid_model_id",
			inputModelID: "not_valid_id",
			assertions: []*openfgav1.Assertion{
				{
					TupleKey:    tuple.NewAssertionTupleKey("repo:test", "can_read", "user:elbuo"),
					Expectation: false,
				},
			},
			expectErrWhenWriting: "Authorization Model 'not_valid_id' not found",
		},
	}

	for _, test := range tests {
		t.Run(test._name, func(t *testing.T) {
			_, err := commands.NewWriteAssertionsCommand(datastore).Execute(context.Background(), &openfgav1.WriteAssertionsRequest{
				StoreId:              store,
				Assertions:           test.assertions,
				AuthorizationModelId: test.inputModelID,
			})
			if test.expectErrWhenWriting != "" {
				require.ErrorContains(t, err, test.expectErrWhenWriting)
			} else {
				require.NoError(t, err)

				actualResponse, err := commands.NewReadAssertionsQuery(datastore).Execute(context.Background(), store, test.inputModelID)
				require.NoError(t, err)

				expectedResponse := &openfgav1.ReadAssertionsResponse{
					AuthorizationModelId: test.inputModelID,
					Assertions:           test.assertions,
				}
				if diff := cmp.Diff(expectedResponse, actualResponse, protocmp.Transform()); diff != "" {
					t.Errorf("store mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
package graph

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"golang.org/x/time/rate"

	"github.com/openfga/openfga/pkg/logger"
	"github.com/openfga/openfga/pkg/tuple"
)

const (
	trackerLogLines    = 15
	trackerMaxInterval = 60
	trackerMinInterval = 30
	trackerLogInterval = time.Duration(500) * time.Millisecond
)

type TrackerCheckResolverOpt func(checkResolver *TrackerCheckResolver)

type trackerKey struct {
	store string
	model string
}

type resolutionNode struct {
	tm   time.Time
	hits *atomic.Uint64
}

func (r *resolutionNode) expired(trackerInterval time.Duration) bool {
	return time.Since(r.tm) > trackerInterval
}

type TrackerCheckResolver struct {
	delegate       CheckResolver
	ticker         *time.Ticker
	tickerInterval time.Duration
	logger         logger.Logger
	limiter        *rate.Limiter
	ctx            context.Context
	cancel         context.CancelFunc

	nodes sync.Map // key is a string store#model and value is another map where key has the shape doc:1#viewer@user and value is the number of hits.
}

var _ CheckResolver = (*TrackerCheckResolver)(nil)

func WithTrackerLogger(logger logger.Logger) TrackerCheckResolverOpt {
	return func(t *TrackerCheckResolver) {
		t.logger = logger
	}
}

func WithTrackerContext(ctx context.Context) TrackerCheckResolverOpt {
	return func(t *TrackerCheckResolver) {
		t.ctx = ctx
	}
}

func WithTrackerInterval(d time.Duration) TrackerCheckResolverOpt {
	return func(t *TrackerCheckResolver) {
		t.ticker.Stop() // clear the default
		t.tickerInterval = d
		t.ticker = time.NewTicker(d)
	}
}

func NewTrackCheckResolver(opts ...TrackerCheckResolverOpt) *TrackerCheckResolver {
	randomInterval := rand.Intn(trackerMaxInterval-trackerMinInterval+1) + trackerMinInterval
	defaultTickerInterval := time.Duration(randomInterval) * time.Second

	t := &TrackerCheckResolver{
		logger:         logger.NewNoopLogger(),
		limiter:        rate.NewLimiter(rate.Every(trackerLogInterval), trackerLogLines),
		ticker:         time.NewTicker(defaultTickerInterval),
		tickerInterval: defaultTickerInterval,
	}

	for _, opt := range opts {
		opt(t)
	}

	t.delegate = t

	if t.ctx == nil {
		t.ctx, t.cancel = context.WithCancel(context.Background())
	} else {
		t.cancel = func() { // no op
		}
	}

	t.launchFlush()

	return t
}

// LaunchFlush starts the execution path logging and removal of entries.
func (t *TrackerCheckResolver) launchFlush() {
	go func() {
		for {
			select {
			case <-t.ctx.Done():
				t.ticker.Stop()
				return
			case <-t.ticker.C:
				t.logExecutionPaths(false)
				t.ticker.Reset(t.tickerInterval)
			}
		}
	}()
}

func (t *TrackerCheckResolver) SetDelegate(delegate CheckResolver) {
	t.delegate = delegate
}

func (t *TrackerCheckResolver) GetDelegate() CheckResolver {
	return t.delegate
}

func (t *TrackerCheckResolver) Close() {
	t.cancel()
	t.logExecutionPaths(true)
}

func (t *TrackerCheckResolver) logExecutionPaths(flush bool) {
	t.nodes.Range(func(k, v any) bool {
		paths := v.(*sync.Map)
		storeModel := k.(trackerKey)
		paths.Range(func(k, v any) bool {
			path := k.(string)
			node := v.(*resolutionNode)
			if node.expired(t.tickerInterval) || flush {
				if !t.limiter.Allow() && !flush {
					return false
				}
				t.logger.Info("hits",
					zap.String("store_id", storeModel.store),
					zap.String("model_id", storeModel.model),
					zap.String("path", path),
					zap.Uint64("hits", node.hits.Load()))

				paths.Delete(path)
			}
			return true
		})
		return true
	})
}

// getTupleKeyAsPath for a tuple like (user:anne, viewer, doc:1) returns doc:1#viewer@user
// for a tuple like (group:fga#member, viewer, doc:1), returns doc:1#viewer@userset
// for a tuple like (user:*, viewer, doc:1), returns doc:1#viewer@userset.
func getTupleKeyAsPath(tk *openfgav1.TupleKey) string {
	return fmt.Sprintf("%s#%s@%s", tk.GetObject(), tk.GetRelation(), string(tuple.GetUserTypeFromUser(tk.GetUser())))
}

func (t *TrackerCheckResolver) loadModel(r *ResolveCheckRequest) (value any, ok bool) {
	key := trackerKey{store: r.GetStoreID(), model: r.GetAuthorizationModelID()}
	value, ok = t.nodes.Load(key)
	if !ok {
		value = &sync.Map{}
		value.(*sync.Map).Store(getTupleKeyAsPath(r.GetTupleKey()), &resolutionNode{tm: time.Now().UTC(), hits: &atomic.Uint64{}})
		t.nodes.Store(key, value)
	}
	return value, ok
}

func (t *TrackerCheckResolver) addPathHits(r *ResolveCheckRequest) {
	path := getTupleKeyAsPath(r.GetTupleKey())

	value, ok := t.loadModel(r)
	if ok {
		paths, _ := value.(*sync.Map)
		if _, ok := paths.Load(path); !ok {
			paths.Store(path, &resolutionNode{tm: time.Now(), hits: &atomic.Uint64{}})
		}
	}

	paths, ok := value.(*sync.Map)
	if ok {
		value, ok := paths.Load(path)
		if ok {
			value.(*resolutionNode).hits.Add(1)
		}
	}
}

func (t *TrackerCheckResolver) ResolveCheck(
	ctx context.Context,
	req *ResolveCheckRequest,
) (*ResolveCheckResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Bool("track_execution", true))

	resp, err := t.delegate.ResolveCheck(ctx, &ResolveCheckRequest{
		StoreID:              req.GetStoreID(),
		AuthorizationModelID: req.GetAuthorizationModelID(),
		TupleKey:             req.GetTupleKey(),
		ContextualTuples:     req.GetContextualTuples(),
		RequestMetadata:      req.GetRequestMetadata(),
		VisitedPaths:         req.VisitedPaths,
		Context:              req.GetContext(),
	})

	if err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		t.addPathHits(req)
	}

	return resp, err
}
