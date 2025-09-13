package postgres

import (
	"context"
	"fmt"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/otpbus"
	db "bitbucket.org/msafaridanquah/sight-backend/business/sdk/postgres/out"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/vaulti"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const vaultKey = "pii_key"

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

func (repo *Repository) Add(ctx context.Context, bus otpbus.OTP) error {
	ctx, span := otel.AddSpan(ctx, "otpbus.postgres.Add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)

	defer span.End()

	pin, err := repo.vaulti.TransitEncrypt(bus.Code, vaultKey)
	if err != nil {
		return err
	}

	if err := repo.queries.InsertOTP(ctx, db.InsertOTPParams{
		ID: uuid.NullUUID{
			UUID:  bus.ID,
			Valid: true,
		},
		CustomerID: uuid.NullUUID{
			UUID:  bus.CustomerID,
			Valid: true,
		},
		HashedCode: pgtype.Text{
			String: bus.Hash.String(),
			Valid:  true,
		},
		Code: pgtype.Text{
			String: pin.Ciphertext,
			Valid:  true,
		},
		Channel: db.NullChannel{
			Channel: db.Channel(bus.Channel.String()),
			Valid:   true,
		},
		ExpiresAt: pgtype.Timestamp{
			Time:  bus.ExpiresAt,
			Valid: true,
		},
		Destination: bus.Destination,
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) FindByCustomerIDAndHash(ctx context.Context, id uuid.UUID, hashed string) (otpbus.OTP, error) {
	ctx, span := otel.AddSpan(ctx, "otpbus.postgres.findbycustomeridandhash")
	span.SetAttributes(semconv.DBSystemPostgreSQL)

	defer span.End()

	res, err := repo.queries.GetOTPByCustomerIDAndCode(ctx, db.GetOTPByCustomerIDAndCodeParams{
		CustomerID: uuid.NullUUID{
			UUID:  id,
			Valid: true,
		},
		HashedCode: pgtype.Text{
			String: hashed,
			Valid:  true,
		},
	})
	if err != nil {
		return otpbus.OTP{}, err
	}

	fmt.Printf("code from db %s", res.Code.String)

	rawcode, err := repo.vaulti.TransitDecrypt(res.Code.String, vaultKey)
	if err != nil {
		return otpbus.OTP{}, err
	}

	bus := otpbus.OTP{
		ID:         res.ID.UUID,
		Code:       rawcode.Plaintext,
		CustomerID: res.CustomerID.UUID,
		VerifiedAt: res.VerifiedAt.Time,
		ExpiresAt:  res.VerifiedAt.Time,
	}

	return bus, nil
}
