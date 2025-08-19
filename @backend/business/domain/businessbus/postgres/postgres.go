package postgres

import (
	"context"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus"
	db "bitbucket.org/msafaridanquah/verifylab-service/business/sdk/postgres/out"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/vaulti"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Repository struct {
	queries *db.Queries
	conn    *pgxpool.Pool
	vaulti  *vaulti.Vaulty
}

func New(d db.DBTX, conn *pgxpool.Pool, vault *vaulti.Vaulty) *Repository {
	return &Repository{
		conn:    conn,
		queries: db.New(d),
		vaulti:  vault,
	}
}

func (r *Repository) Add(ctx context.Context, bus businessbus.Business) error {
	ctx, span := otel.AddSpan(ctx, "customerbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	now := time.Now()
	if err := r.queries.InsertBusiness(ctx, db.InsertBusinessParams{
		ID:        uuid.New(),
		LegalName: bus.LegalName,
		TaxID: pgtype.Text{
			String: bus.TaxID,
			Valid:  true,
		},
		Entity:       db.Entity(bus.Entity.String()),
		Jurisdiction: bus.CountryOfIncorporation.Alpha2(),
		Dba:          bus.DoingBusinessAs,
		AdminID:      uuid.New(), //Get from middleware
		CreatedAt: pgtype.Timestamp{
			Time:  now,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  now,
			Valid: true,
		},
	}); err != nil {
		return err
	}

	return nil
}
