package postgres

import (
	"context"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus/valueobject"
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

func (repo *Repository) Add(ctx context.Context, bus otpbus.OTP) (otpbus.OTP, error) {
	ctx, span := otel.AddSpan(ctx, "otpbus.postgres.Add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)

	defer span.End()

	_, err := repo.queries.InsertOTP(ctx, db.InsertOTPParams{
		ID: uuid.NullUUID{
			UUID:  bus.ID,
			Valid: true,
		},
		CustomerID: uuid.NullUUID{
			UUID:  bus.CustomerID,
			Valid: true,
		},
		HashedCode: pgtype.Text{
			String: bus.HashedCode.String(),
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
	})

	if err != nil {
		return otpbus.OTP{}, err
	}

	return bus, nil
}

func (repo *Repository) Find(ctx context.Context, id uuid.UUID) (otpbus.OTP, error) {
	ctx, span := otel.AddSpan(ctx, "otpbus.postgres.Find")
	span.SetAttributes(semconv.DBSystemPostgreSQL)

	defer span.End()

	res, err := repo.queries.GetOTP(ctx, uuid.NullUUID{
		UUID:  id,
		Valid: true,
	})

	if err != nil {
		return otpbus.OTP{}, err
	}
	channel, err := valueobject.ParseChannel(string(res.Channel.Channel))
	if err != nil {
		return otpbus.OTP{}, err
	}

	bus := otpbus.OTP{
		ID:          res.ID.UUID,
		Channel:     channel,
		CustomerID:  res.CustomerID.UUID,
		Destination: res.Destination,
		VerifiedAt:  res.VerifiedAt.Time,
		ExpiresAt:   res.VerifiedAt.Time,
	}

	return bus, nil
}
