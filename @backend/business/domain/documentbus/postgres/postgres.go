package postgres

import (
	db "bitbucket.org/msafaridanquah/verifylab-service/business/sdk/postgres/out"
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
