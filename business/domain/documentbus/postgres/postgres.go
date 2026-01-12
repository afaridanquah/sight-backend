package postgres

import (
	"context"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus"
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

func (r *Repository) Add(ctx context.Context, bus documentbus.Document) error {
	ctx, span := otel.AddSpan(ctx, "documentbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	dbParams := db.InsertDocumentParams{
		ID: bus.ID.String(),
		OriginalName: pgtype.Text{
			String: bus.OriginalName,
			Valid:  true,
		},
		Filename: pgtype.Text{
			String: bus.FileName,
			Valid:  true,
		},
		Status: db.NullDocumentStatus{
			DocumentStatus: db.DocumentStatus(bus.Status.String()),
			Valid:          true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  bus.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  bus.UpdatedAt,
			Valid: true,
		},
	}

	if bus.BusinessID != "" {
		dbParams.BusinessID = pgtype.Text{
			String: bus.BusinessID,
			Valid:  true,
		}
	}

	if bus.CustomerID != "" {
		dbParams.CustomerID = pgtype.Text{
			String: bus.CustomerID,
			Valid:  true,
		}
	}

	if err := r.queries.InsertDocument(ctx, dbParams); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, bus documentbus.Document) error {
	ctx, span := otel.AddSpan(ctx, "documentbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	dbParams := db.UpdateDocumentParams{
		ID: bus.ID.String(),
		Status: db.NullDocumentStatus{
			DocumentStatus: db.DocumentStatus(bus.Status.String()),
		},
	}

	if err := r.queries.UpdateDocument(ctx, dbParams); err != nil {
		return err
	}

	return nil
}
