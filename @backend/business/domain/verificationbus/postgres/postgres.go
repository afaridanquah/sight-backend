package postgres

import (
	"context"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus"
	db "bitbucket.org/msafaridanquah/verifylab-service/business/sdk/postgres/out"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/vaulti"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *Repository) Add(ctx context.Context, bus verificationbus.Verification) error {
	ctx, span := otel.AddSpan(ctx, "verificationbus.postgres.add")

	defer span.End()

	arg, err := toDBInsertVerification(bus, r.vaulti)
	if err != nil {
		return err
	}

	if err := r.queries.InsertVerification(ctx, arg); err != nil {
		return err
	}

	return nil
}
