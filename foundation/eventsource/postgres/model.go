package postgres

import (
	db "bitbucket.org/msafaridanquah/verifylab-service/business/sdk/postgres/out"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/eventsource"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toDBStoredEvent(a eventsource.Event) db.InsertEventParams {
	return db.InsertEventParams{
		ID: uuid.NullUUID{
			UUID:  a.ID(),
			Valid: true,
		},
		Type: pgtype.Text{
			String: a.Type().String(),
			Valid:  true,
		},
		AggregateID: uuid.NullUUID{
			UUID:  a.AggregateID(),
			Valid: true,
		},
	}
}
