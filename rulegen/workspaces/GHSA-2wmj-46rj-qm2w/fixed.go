package main

// Eventstore abstracts all functions needed to store valid events
// and filters the stored events
type Eventstore struct {
	// TODO: get rid of this mutex,
	// or if we scale to >4vCPU use a sync.Map
	interceptorMutex  sync.RWMutex
	eventInterceptors map[EventType]eventTypeInterceptors
	eventTypes        []string
	aggregateTypes    []string
func NewEventstore(config *Config) *Eventstore {
	return &Eventstore{
		eventInterceptors: map[EventType]eventTypeInterceptors{},
		PushTimeout:       config.PushTimeout,

		pusher:  config.Pusher,

// Filter filters the stored events based on the searchQuery
// and maps the events to the defined event structs
//
// Deprecated: Use [FilterToQueryReducer] instead to avoid allocations.
func (es *Eventstore) Filter(ctx context.Context, searchQuery *SearchQueryBuilder) ([]Event, error) {
	events := make([]Event, 0, searchQuery.GetLimit())
	searchQuery.ensureInstanceID(ctx)
	err := es.querier.FilterToReducer(ctx, searchQuery, func(event Event) error {
		event, err := es.mapEvent(event)
		if err != nil {
			return err
		}
		events = append(events, event)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (es *Eventstore) mapEvents(events []Event) (mappedEvents []Event, err error) {
	mappedEvents = make([]Event, len(events))

	es.interceptorMutex.RLock()
	defer es.interceptorMutex.RUnlock()

	for i, event := range events {
		mappedEvents[i], err = es.mapEventLocked(event)
		if err != nil {
			return nil, err
		}
}

func (es *Eventstore) mapEvent(event Event) (Event, error) {
	es.interceptorMutex.RLock()
	defer es.interceptorMutex.RUnlock()
	return es.mapEventLocked(event)
}

func (es *Eventstore) mapEventLocked(event Event) (Event, error) {
	interceptors, ok := es.eventInterceptors[event.Type()]
	if !ok || interceptors.eventMapper == nil {
		return BaseEventFromRepo(event), nil
	return interceptors.eventMapper(event)
}

// TODO: refactor so we can change to the following interface:
/*
type reducer interface {
	// Reduce applies an event on the object.
	Reduce(Event) error
}
*/

type reducer interface {
	//Reduce handles the events of the internal events list
	// it only appends the newly added events

// FilterToReducer filters the events based on the search query, appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToReducer(ctx context.Context, searchQuery *SearchQueryBuilder, r reducer) error {
	searchQuery.ensureInstanceID(ctx)
	return es.querier.FilterToReducer(ctx, searchQuery, func(event Event) error {
		event, err := es.mapEvent(event)
		if err != nil {
			return err
		}
		r.AppendEvents(event)
		return r.Reduce()
	})
}

// LatestSequence filters the latest sequence for the given search query
// FilterToQueryReducer filters the events based on the search query of the query function,
// appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToQueryReducer(ctx context.Context, r QueryReducer) error {
	return es.FilterToReducer(ctx, r.Query(), r)
}

// RegisterFilterEventMapper registers a function for mapping an eventstore event to an event
	return es
}

type Reducer func(event Event) error

type Querier interface {
	// Health checks if the connection to the storage is available
	Health(ctx context.Context) error
	// FilterToReducer calls r for every event returned from the storage
	FilterToReducer(ctx context.Context, searchQuery *SearchQueryBuilder, r Reducer) error
	// LatestSequence returns the latest sequence found by the search query
	LatestSequence(ctx context.Context, queryFactory *SearchQueryBuilder) (float64, error)
	// InstanceIDs returns the instance ids found by the search query
package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 16.sql
	uniqueConstraintLower string
)

type UniqueConstraintToLower struct {
	dbClient *database.DB
}

func (mig *UniqueConstraintToLower) Execute(ctx context.Context) error {
	res, err := mig.dbClient.ExecContext(ctx, uniqueConstraintLower)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	logging.WithFields("count", count).Info("unique constraints updated")
	return err
}

func (mig *UniqueConstraintToLower) String() string {
	return "16_unique_constraint_lower"
}
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"

	host, err := HostFromRequest(r, headerName)
	if err != nil {
		return nil, caos_errors.ThrowNotFound(err, "INST-zWq7X", "Errors.Instance.NotFound")
	}

	instance, err := verifier.InstanceByHost(authCtx, host)
	return authz.WithInstance(ctx, instance), nil
}

func HostFromRequest(r *http.Request, headerName string) (host string, err error) {
	if headerName != "host" {
		return hostFromSpecialHeader(r, headerName)
	}
	return hostFromOrigin(r.Context())
}

func hostFromSpecialHeader(r *http.Request, headerName string) (host string, err error) {
	host = r.Header.Get(headerName)
	if host == "" {
		return "", fmt.Errorf("host header `%s` not found", headerName)
	}
	return host, nil
}

func hostFromOrigin(ctx context.Context) (host string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("invalid origin: %w", err)
		}
	}()
	origin := zitadel_http.ComposedOrigin(ctx)
	u, err := url.Parse(origin)
	if err != nil {
		return "", err
	}
	host = u.Hostname()
	if host == "" {
		err = errors.New("empty host")
	}
	return host, err
}

func newZitadelTranslator() *i18n.Translator {
	dir, err := fs.NewWithNamespace("zitadel")
	logging.WithFields("namespace", "zitadel").OnError(err).Panic("unable to get namespace")
