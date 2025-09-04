package postgres

import (
	"context"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/organizationbus"
	db "bitbucket.org/msafaridanquah/verifylab-service/business/sdk/postgres/out"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/vaulti"
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

func (r *Repository) Add(ctx context.Context, org organizationbus.Organization) error {
	ctx, span := otel.AddSpan(ctx, "customerbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	if err := r.queries.InsertOrg(ctx, db.InsertOrgParams{
		ID:   org.ID,
		Name: org.Name,
	}); err != nil {
		return err
	}

	return nil
}
