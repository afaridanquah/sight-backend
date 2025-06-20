package eventsource

import "context"

type Store interface {
	Save(ctx context.Context, events ...Event) error
	// Load(ctx context.Context, a Aggregate) (Aggregate, error)
	// EventsHistory(ctx context.Context, aggregateID, aggregateType string, fromVersion int, limit int) ([]EventReadModel, error)
}
