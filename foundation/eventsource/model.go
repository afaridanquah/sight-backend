package eventsource

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Version uint64

type AggregateType string
type AggregateID string

func (t AggregateType) String() string {
	return string(t)
}

func (t AggregateType) IsZero() bool {
	return len(strings.TrimSpace(string(t))) == 0
}

type Aggregate interface {
	ID() uuid.UUID
	Type() AggregateType
	Version() uint64
	Changes() []Event
	StackChange(change Event)
	SetVersion(version Version)
	IncrementVersion()
	PrepareForLoading()
	ParseEvents(context.Context, ...EventReadState) []Event
}

// Event holding meta data and the application specific event in the Data property
type AggregateState struct {
	id            uuid.UUID
	aggregateType AggregateType
	version       Version
	changes       []Event
}

func NewAggregate(id uuid.UUID, t AggregateType) *AggregateState {
	return &AggregateState{
		id:            id,
		aggregateType: t,
		version:       0,
		changes:       make([]Event, 0),
	}
}

func (a *AggregateState) PrepareForLoading() {
	a.version = 0
	a.changes = make([]Event, 0)
}

func (a *AggregateState) ID() uuid.UUID {
	return a.id
}

func (a *AggregateState) Version() Version {
	return a.version
}

func (a *AggregateState) Changes() []Event {
	return a.changes
}

// =====================================================================================================

type EventType string

type Metadata map[string]interface{}

type Event interface {
	fmt.Stringer
	ApplyTo(ctx context.Context, aggregate Aggregate) // ApplyTo applies the event to the aggregate
	ID() uuid.UUID                                    // ID returns the id of the event.
	Type() EventType                                  // Type returns the type of the event.
	OccurredAt() time.Time                            // OccurredAt of when the event was created.
	AggregateID() uuid.UUID                           // AggregateID is the id of the aggregate that the event belongs to.
	AggregateType() AggregateType                     // AggregateType is the type of the aggregate that the event can be applied to.
	AggregateVersion() Version                        // AggregateVersion is the version of the aggregate after the event has been applied.
	SetVersion(Version)                               // SetVersion sets the aggregate version of the event
	Metadata() Metadata                               // Metadata is app-specific metadata such as request AggregateID, originating user etc.
}

func (t EventType) String() string {
	return string(t)
}

func (t EventType) IsZero() bool {
	return len(strings.TrimSpace(t.String())) == 0
}

type EventState struct {
	id               uuid.UUID
	occurredAt       time.Time
	aggregateID      uuid.UUID
	aggregateType    AggregateType
	aggregateVersion Version
	metadata         map[string]interface{}
}

func (e *EventState) ID() uuid.UUID {
	return e.id
}

func (e *EventState) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *EventState) AggregateID() uuid.UUID {
	return e.aggregateID
}

func (e *EventState) AggregateType() AggregateType {
	return e.aggregateType
}

func (e *EventState) AggregateVersion() Version {
	return e.aggregateVersion
}

func (e *EventState) Metadata() Metadata {
	return e.metadata
}

func (e *EventState) String() string {
	return fmt.Sprintf("event '%s' occurred on aggregate '%s' (v%d => v%d) with id '%s'", e.ID(), e.AggregateType(), e.aggregateVersion-1, e.aggregateVersion, e.AggregateID())
}

func (e *EventState) SetVersion(version Version) {
	e.aggregateVersion = version
}

func newBaseEvent(id uuid.UUID, occurredAt time.Time, aggregateID uuid.UUID, aggregateType AggregateType, version Version, metadata Metadata) *BaseEvent {
	if metadata == nil {
		metadata = make(Metadata, 0)
	}

	return &EventState{
		id:               id,
		occurredAt:       occurredAt,
		aggregateID:      aggregateID,
		aggregateType:    aggregateType,
		aggregateVersion: version,
		metadata:         metadata,
	}
}

// ==============================================================

type EventReadState struct {
	ID               uuid.UUID              `json:"id"`
	Type             EventType              `json:"type"`
	OccurredAt       time.Time              `json:"occurred_at"`
	AggregateID      uuid.UUID              `json:"aggregate_id"`
	AggregateType    AggregateType          `json:"aggregate_type"`
	AggregateVersion Version                `json:"aggregate_version"`
	Metadata         map[string]interface{} `json:"metadata"`
	Data             json.RawMessage        `json:"data"`
}

func (r *EventReadState) NewBaseEvent() *EventState {
	return newBaseEvent(
		r.ID,
		r.OccurredAt,
		r.AggregateID,
		r.AggregateType,
		r.AggregateVersion,
		r.Metadata,
	)
}
