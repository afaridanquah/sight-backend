package postgres

import (
	db "bitbucket.org/msafaridanquah/verifylab-service/business/sdk/postgres/out"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	queries *db.Queries
	conn    *pgxpool.Pool
}

func New(d db.DBTX, conn *pgxpool.Pool) *Repository {
	return &Repository{
		queries: db.New(d),
		conn:    conn,
	}
}

// func (es *EventStore) Save(ctx context.Context, ) error {
// 	ctx, span := otel.AddSpan(ctx, "eventsource.postgres.Save")
// 	span.SetAttributes(semconv.DBSystemPostgreSQL)

// 	defer span.End()

// 	tx, err := es.conn.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	defer tx.Rollback(ctx)

// 	qtx := es.queries.WithTx(tx)

// 	qtx.InsertEvent(ctx, db.InsertEventParams{
// 		ID: uuid.NullUUID{
// 			UUID: string() a.ID(),
// 		},

// 	})

// }
