package postgres

import (
	"context"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/identificationbus"
	db "bitbucket.org/msafaridanquah/verifylab-service/business/sdk/postgres/out"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Repository struct {
	queries *db.Queries
}

func New(d db.DBTX) *Repository {
	return &Repository{
		queries: db.New(d),
	}
}

func (r *Repository) Add(ctx context.Context, bus identificationbus.Identification) (identificationbus.Identification, error) {
	ctx, span := otel.AddSpan(ctx, "identificationbus.postgres.Add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	_, err := r.queries.CreateIdentification(ctx, db.CreateIdentificationParams{
		ID: bus.ID,
	})

	if err != nil {
		return identificationbus.Identification{}, err
	}
	return bus, nil
}
