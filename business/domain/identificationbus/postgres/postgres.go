package postgres

import (
	"context"
	"errors"
	"fmt"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus/valueobject"
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

const vaultKey = "pii_key"

func New(d db.DBTX, conn *pgxpool.Pool, vault *vaulti.Vaulty) *Repository {
	return &Repository{
		conn:    conn,
		queries: db.New(d),
		vaulti:  vault,
	}
}

func (r *Repository) Add(ctx context.Context, bus identificationbus.Identification) error {
	ctx, span := otel.AddSpan(ctx, "identificationbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	p, err := r.vaulti.TransitEncrypt(bus.Pin, vaultKey)
	if err != nil {
		return err
	}
	if err := r.queries.InsertCustomerIdentification(ctx, db.InsertCustomerIdentificationParams{
		ID: bus.ID.String(),
		FirstName: pgtype.Text{
			String: bus.FirstName,
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: bus.LastName,
			Valid:  true,
		},
		MiddleName: pgtype.Text{
			String: bus.MiddleName,
			Valid:  true,
		},
		Pin:                p.Ciphertext,
		IdentificationType: db.IdentificationType(bus.IdentificationType.String()),
		IssuedCountry: pgtype.Text{
			String: bus.Country.Alpha2(),
			Valid:  true,
		},
		PlaceOfBirth: pgtype.Text{
			String: bus.PlaceOfBirth,
			Valid:  true,
		},
		DateOfBirth: pgtype.Date{
			Time:  bus.DateOfBirth,
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
		return err
	}

	return nil
}

func (r *Repository) AddMany(ctx context.Context, bus []identificationbus.Identification) error {
	ctx, span := otel.AddSpan(ctx, "identificationbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	params := make([]db.BulkInsertCustomerIdentificationsParams, len(bus))

	if len(bus) > 0 {
		for k, v := range bus {
			p, err := r.vaulti.TransitEncrypt(v.Pin, vaultKey)
			if err != nil {
				fmt.Printf("vault transit %v", err)
				return err
			}
			params[k] = db.BulkInsertCustomerIdentificationsParams{
				ID: v.ID.String(),
				CustomerID: pgtype.Text{
					String: v.CustomerID,
					Valid:  true,
				},
				FirstName: pgtype.Text{
					String: v.FirstName,
					Valid:  true,
				},
				LastName: pgtype.Text{
					String: v.LastName,
					Valid:  true,
				},
				MiddleName: pgtype.Text{
					String: v.MiddleName,
					Valid:  true,
				},
				Pin:                p.Ciphertext,
				IdentificationType: db.IdentificationType(v.IdentificationType.String()),
				IssuedCountry: pgtype.Text{
					String: v.Country.Alpha2(),
					Valid:  true,
				},
				PlaceOfBirth: pgtype.Text{
					String: v.PlaceOfBirth,
					Valid:  true,
				},
				DateOfBirth: pgtype.Date{
					Time:  v.DateOfBirth,
					Valid: true,
				},
				CreatedAt: pgtype.Timestamp{
					Time:  v.CreatedAt,
					Valid: true,
				},
				UpdatedAt: pgtype.Timestamp{
					Time:  v.UpdatedAt,
					Valid: true,
				},
			}
		}
	}

	_, err := r.queries.BulkInsertCustomerIdentifications(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) QueryByCustomerID(ctx context.Context, cusID valueobject.ID) ([]identificationbus.Identification, error) {
	ctx, span := otel.AddSpan(ctx, "identificationbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	resp, err := r.queries.QueryIdentificationsByCustomerID(ctx, pgtype.Text{
		String: cusID.String(),
		Valid:  true,
	})

	if err != nil {
		return []identificationbus.Identification{}, err
	}

	bus := make([]identificationbus.Identification, len(resp))
	var errs []error
	for k, v := range resp {
		idx, err := valueobject.ParseID(v.ID)
		if err != nil {
			errs = append(errs, err)
		}

		ic, err := valueobject.NewCountry(v.IssuedCountry.String)
		if err != nil {
			errs = append(errs, err)
		}

		idt, err := valueobject.ParseIdentificationType(string(v.IdentificationType))
		if err != nil {
			errs = append(errs, err)
		}
		// dob, err := valueobject.
		i := identificationbus.Identification{
			ID:                 idx,
			CustomerID:         v.CustomerID.String,
			Pin:                v.Pin,
			FirstName:          v.FirstName.String,
			LastName:           v.LastName.String,
			MiddleName:         v.MiddleName.String,
			Country:            ic,
			DateOfBirth:        v.DateOfBirth.Time,
			IdentificationType: idt,
			CreatedAt:          v.CreatedAt.Time,
			UpdatedAt:          v.UpdatedAt.Time,
		}
		bus[k] = i
	}

	if len(errs) > 0 {
		return []identificationbus.Identification{}, errors.Join(errs...)
	}

	return bus, nil
}
