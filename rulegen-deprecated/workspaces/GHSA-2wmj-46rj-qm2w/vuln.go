package main

// Eventstore abstracts all functions needed to store valid events
// and filters the stored events
type Eventstore struct {
	interceptorMutex  sync.Mutex
	eventInterceptors map[EventType]eventTypeInterceptors
	eventTypes        []string
	aggregateTypes    []string
func NewEventstore(config *Config) *Eventstore {
	return &Eventstore{
		eventInterceptors: map[EventType]eventTypeInterceptors{},
		interceptorMutex:  sync.Mutex{},
		PushTimeout:       config.PushTimeout,

		pusher:  config.Pusher,

// Filter filters the stored events based on the searchQuery
// and maps the events to the defined event structs
func (es *Eventstore) Filter(ctx context.Context, queryFactory *SearchQueryBuilder) ([]Event, error) {
	// make sure that the instance id is always set
	if queryFactory.instanceID == nil && authz.GetInstance(ctx).InstanceID() != "" {
		queryFactory.InstanceID(authz.GetInstance(ctx).InstanceID())
	}

	events, err := es.querier.Filter(ctx, queryFactory)
	if err != nil {
		return nil, err
	}

	return es.mapEvents(events)
}

func (es *Eventstore) mapEvents(events []Event) (mappedEvents []Event, err error) {
	mappedEvents = make([]Event, len(events))

	es.interceptorMutex.Lock()
	defer es.interceptorMutex.Unlock()

	for i, event := range events {
		mappedEvents[i], err = es.mapEvent(event)
		if err != nil {
			return nil, err
		}
}

func (es *Eventstore) mapEvent(event Event) (Event, error) {
	interceptors, ok := es.eventInterceptors[event.Type()]
	if !ok || interceptors.eventMapper == nil {
		return BaseEventFromRepo(event), nil
	return interceptors.eventMapper(event)
}

type reducer interface {
	//Reduce handles the events of the internal events list
	// it only appends the newly added events

// FilterToReducer filters the events based on the search query, appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToReducer(ctx context.Context, searchQuery *SearchQueryBuilder, r reducer) error {
	events, err := es.Filter(ctx, searchQuery)
	if err != nil {
		return err
	}

	r.AppendEvents(events...)

	return r.Reduce()
}

// LatestSequence filters the latest sequence for the given search query
// FilterToQueryReducer filters the events based on the search query of the query function,
// appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToQueryReducer(ctx context.Context, r QueryReducer) error {
	events, err := es.Filter(ctx, r.Query())
	if err != nil {
		return err
	}
	r.AppendEvents(events...)

	return r.Reduce()
}

// RegisterFilterEventMapper registers a function for mapping an eventstore event to an event
	return es
}

type Querier interface {
	// Health checks if the connection to the storage is available
	Health(ctx context.Context) error
	// Filter returns all events matching the given search query
	Filter(ctx context.Context, searchQuery *SearchQueryBuilder) (events []Event, err error)
	// LatestSequence returns the latest sequence found by the search query
	LatestSequence(ctx context.Context, queryFactory *SearchQueryBuilder) (float64, error)
	// InstanceIDs returns the instance ids found by the search query
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"

	host, err := HostFromRequest(r, headerName)
	if err != nil {
		return nil, err
	}

	instance, err := verifier.InstanceByHost(authCtx, host)
	return authz.WithInstance(ctx, instance), nil
}

func HostFromRequest(r *http.Request, headerName string) (string, error) {
	host := r.Host
	if headerName != "host" {
		host = r.Header.Get(headerName)
	}
	if host == "" {
		return "", fmt.Errorf("host header `%s` not found", headerName)
	}
	return host, nil
}

func newZitadelTranslator() *i18n.Translator {
	dir, err := fs.NewWithNamespace("zitadel")
	logging.WithFields("namespace", "zitadel").OnError(err).Panic("unable to get namespace")
