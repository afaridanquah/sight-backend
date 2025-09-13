package postgres

import (
	"context"
	"encoding/json"
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
	ctx, span := otel.AddSpan(ctx, "businessbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	now := time.Now()

	emails := make([]string, len(bus.EmailAddresses))
	if len(bus.EmailAddresses) > 0 {
		for i, e := range bus.EmailAddresses {
			emails[i] = e.String()
		}
	}

	emailAddresses, err := json.Marshal(emails)
	if err != nil {
		return err
	}

	phones := make([]string, len(bus.PhoneNumbers))
	if len(bus.PhoneNumbers) > 0 {
		for i, e := range bus.PhoneNumbers {
			phones[i] = e.E164Format
		}
	}
	phoneNumbers, err := json.Marshal(phones)
	if err != nil {
		return err
	}

	address, err := toDBAddress(bus.Address)
	if err != nil {
		return err
	}

	if err := r.queries.InsertBusiness(ctx, db.InsertBusinessParams{
		ID:        bus.ID,
		LegalName: bus.LegalName,
		TaxID: pgtype.Text{
			String: bus.TaxID,
			Valid:  true,
		},
		Entity:       db.Entity(bus.Entity.String()),
		Jurisdiction: bus.CountryOfIncorporation.Alpha2(),
		Dba:          bus.DoingBusinessAs,
		AdminID:      uuid.New(), // Get from middleware
		OrgID: uuid.NullUUID{
			UUID:  bus.OrgID,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  now,
			Valid: true,
		},
		EmailAddresses: emailAddresses,
		PhoneNumbers:   phoneNumbers,
		Website: pgtype.Text{
			String: bus.Website,
			Valid:  true,
		},
		Address: address,
		RegistrationNumber: pgtype.Text{
			String: bus.RegistrationNumber,
			Valid:  true,
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

func (r *Repository) Update(ctx context.Context, bus businessbus.Business) error {
	ctx, span := otel.AddSpan(ctx, "businessbus.postgres.update")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	now := time.Now()

	emails := make([]string, len(bus.EmailAddresses))
	if len(bus.EmailAddresses) > 0 {
		for i, e := range bus.EmailAddresses {
			emails[i] = e.String()
		}
	}

	emailAddresses, err := json.Marshal(emails)
	if err != nil {
		return err
	}

	phones := make([]string, len(bus.PhoneNumbers))
	if len(bus.PhoneNumbers) > 0 {
		for i, e := range bus.PhoneNumbers {
			phones[i] = e.E164Format
		}
	}
	phoneNumbers, err := json.Marshal(phones)
	if err != nil {
		return err
	}

	addr, err := toDBAddress(bus.Address)
	if err != nil {
		return err
	}

	if err := r.queries.UpdateBusinessByID(ctx, db.UpdateBusinessByIDParams{
		ID:        bus.ID,
		LegalName: bus.LegalName,
		TaxID: pgtype.Text{
			String: bus.TaxID,
			Valid:  true,
		},
		Entity:         db.Entity(bus.Entity.String()),
		Jurisdiction:   bus.CountryOfIncorporation.Alpha2(),
		Dba:            bus.DoingBusinessAs,
		EmailAddresses: emailAddresses,
		Address:        addr,
		PhoneNumbers:   phoneNumbers,
		Website: pgtype.Text{
			String: bus.Website,
			Valid:  true,
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

func (r *Repository) QueryByIDAndOrgID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (businessbus.Business, error) {
	ctx, span := otel.AddSpan(ctx, "businessbus.postgres.querybyidandorgid")
	span.SetAttributes(semconv.DBSystemPostgreSQL)

	defer span.End()

	res, err := r.queries.GetBusinessByID(ctx, db.GetBusinessByIDParams{
		ID: id,
		OrgID: uuid.NullUUID{
			UUID:  orgID,
			Valid: true,
		},
	})
	if err != nil {
		return businessbus.Business{}, err
	}

	bus, err := toBusBusiness(res)
	if err != nil {
		return businessbus.Business{}, err
	}

	return bus, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error {
	ctx, span := otel.AddSpan(ctx, "businessbus.postgres.delete")
	span.SetAttributes(semconv.DBSystemPostgreSQL)

	defer span.End()

	if err := r.queries.DeleteByID(ctx, db.DeleteByIDParams{
		ID: id,
		OrgID: uuid.NullUUID{
			UUID:  orgID,
			Valid: true,
		},
	}); err != nil {
		return err
	}

	return nil
}
