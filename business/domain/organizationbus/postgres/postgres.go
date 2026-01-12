package postgres

import (
	"context"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus"
	db "bitbucket.org/msafaridanquah/sight-backend/business/sdk/postgres/out"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/vaulti"
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

func (r *Repository) Add(ctx context.Context, org organizationbus.Organization) error {
	ctx, span := otel.AddSpan(ctx, "customerbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	if err := r.queries.InsertOrg(ctx, db.InsertOrgParams{
		ID:   org.ID.String(),
		Name: org.Name,
		Status: db.NullStatus{
			Status: db.Status(org.Status.String()),
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  org.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  org.UpdatedAt,
			Valid: true,
		},
	}); err != nil {
		return err
	}

	return nil
}
