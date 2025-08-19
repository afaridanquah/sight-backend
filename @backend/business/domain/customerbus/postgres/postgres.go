package postgres

import (
	"context"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus/valueobject"
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

func (r *Repository) Add(ctx context.Context, bus customerbus.Customer) error {
	ctx, span := otel.AddSpan(ctx, "customerbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	dob, _ := time.Parse(time.DateOnly, bus.DateOfBirth.String())

	createParams := db.CreateCustomerParams{
		ID:        bus.ID,
		FirstName: bus.Person.FirstName,
		LastName:  bus.Person.LastName,
		MiddleName: pgtype.Text{
			String: bus.Person.MiddleName,
			Valid:  true,
		},
		DateOfBirth: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
		Email: pgtype.Text{
			String: bus.Email.String(),
			Valid:  true,
		},
		BirthCountry: pgtype.Text{
			String: bus.BirthCountry.Alpha2(),
			Valid:  true,
		},
		BusinessID: uuid.NullUUID{
			UUID:  bus.BusinessID,
			Valid: true,
		},
		CreatorID: uuid.NullUUID{
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
	}
	identificationsJson, err := toDBIdentifications(bus.Identifications, r.vaulti)
	if err != nil {
		return err
	}
	createParams.Identifications = identificationsJson

	if bus.PhoneNumber != (valueobject.Phone{}) {
		createParams.PhoneNumber = pgtype.Text{
			String: bus.PhoneNumber.E164Format,
			Valid:  true,
		}
	}

	if err := r.queries.CreateCustomer(ctx, createParams); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, bus customerbus.Customer) error {
	ctx, span := otel.AddSpan(ctx, "business.customerbus.postgres.update")
	defer span.End()

	dob, _ := time.Parse(time.DateOnly, bus.DateOfBirth.String())

	if err := r.queries.UpdateCustomer(ctx, db.UpdateCustomerParams{
		ID:       bus.ID,
		LastName: bus.Person.LastName,
		MiddleName: pgtype.Text{
			String: bus.Person.MiddleName,
			Valid:  true,
		},
		DateOfBirth: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
		Email: pgtype.Text{
			String: bus.Email.String(),
			Valid:  true,
		},
		BirthCountry: pgtype.Text{
			String: bus.BirthCountry.Alpha2(),
			Valid:  true,
		},

		UpdatedAt: pgtype.Timestamp{
			Time:  bus.UpdatedAt,
			Valid: true,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (r *Repository) QueryByCustomerAndBusinessID(ctx context.Context, id uuid.UUID, businessID uuid.UUID) (customerbus.Customer, error) {
	ctx, span := otel.AddSpan(ctx, "business.customerbus.postgres.querybycustomerandbusinessid")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	resp, err := r.queries.QueryCustomerByAndBusinessID(ctx, db.QueryCustomerByAndBusinessIDParams{
		ID: id,
		BusinessID: uuid.NullUUID{
			UUID:  businessID,
			Valid: true,
		},
	})

	if err != nil {
		return customerbus.Customer{}, err
	}

	customer, err := toBusCustomer(resp, r.vaulti)
	if err != nil {
		return customerbus.Customer{}, err
	}

	return customer, nil
}
