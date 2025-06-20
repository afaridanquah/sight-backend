package postgres

import (
	"context"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	db "bitbucket.org/msafaridanquah/verifylab-service/business/sdk/postgres/out"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Repository struct {
	queries *db.Queries
	conn    *pgxpool.Pool
}

func New(d db.DBTX, conn *pgxpool.Pool) *Repository {
	return &Repository{
		conn:    conn,
		queries: db.New(d),
	}
}

func (r *Repository) Add(ctx context.Context, bus customerbus.Customer) (customerbus.Customer, error) {
	ctx, span := otel.AddSpan(ctx, "customerbus.postgres.Add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return customerbus.Customer{}, err
	}

	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	if _, err := qtx.CreateCustomer(ctx, db.CreateCustomerParams{
		ID:        bus.ID,
		FirstName: bus.Person.FirstName,
		LastName:  bus.Person.LastName,
		MiddleName: pgtype.Text{
			String: bus.Person.MiddleName,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: bus.Email.String(),
			Valid:  true,
		},
		Country: bus.BirthCountry.Alpha2(),
		BusinessID: uuid.NullUUID{
			UUID:  bus.BusinessID,
			Valid: true,
		},
		UserID: uuid.NullUUID{
			UUID:  bus.UserID,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  bus.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  bus.UpdatedAt,
			Valid: true,
		},
	}); err != nil {
		return customerbus.Customer{}, err
	}

	identifications := make([]db.CreateCustomerIdentificationsParams, len(bus.Identifications))

	if len(bus.Identifications) > 0 {
		for i, idt := range bus.Identifications {
			identifications[i] = db.CreateCustomerIdentificationsParams{
				ID: uuid.NullUUID{
					UUID:  uuid.New(),
					Valid: true,
				},
				CustomerID:         bus.ID,
				IdentificationType: db.IdentificationType(idt.IdentificationType.String()),
				IssuedCountry: pgtype.Text{
					String: idt.CountryIssued.Alpha2(),
					Valid:  true,
				},
				Pin: pgtype.Text{
					String: idt.Pin,
					Valid:  true,
				},
			}
		}
	}

	if _, err := qtx.CreateCustomerIdentifications(ctx, identifications); err != nil {
		return customerbus.Customer{}, err
	}

	tx.Commit(ctx)

	return bus, nil
}
