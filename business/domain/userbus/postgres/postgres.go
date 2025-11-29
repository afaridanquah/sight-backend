package postgres

import (
	"context"
	"fmt"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/userbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/userbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
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

func (r *Repository) Add(ctx context.Context, bus userbus.User) error {
	ctx, span := otel.AddSpan(ctx, "userbus.postgres.add")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	if err := r.queries.InsertUser(ctx, db.InsertUserParams{
		ID:        bus.ID.String(),
		FirstName: bus.FirstName,
		LastName:  bus.LastName,
		OtherNames: pgtype.Text{
			String: bus.OtherNames,
			Valid:  true,
		},
		Email:    bus.Email.String(),
		Password: bus.Password.String(),
		TenantID: pgtype.Text{
			String: bus.TenantID,
			Valid:  true,
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

func (r *Repository) GetByEmail(ctx context.Context, email valueobject.Email) (userbus.User, error) {
	ctx, span := otel.AddSpan(ctx, "userbus.postgres.getuserbyemail")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	res, err := r.queries.QueryUserByEmail(ctx, email.String())
	if err != nil {
		return userbus.User{}, err
	}

	parsedEmail, _ := valueobject.NewEmail(res.Email)

	id, err := valueobject.ParseID(res.ID)
	if err != nil {
		return userbus.User{}, err
	}

	parsedPassword, err := valueobject.ParsePassword(res.Password)
	if err != nil {
		return userbus.User{}, err
	}

	return userbus.User{
		ID:         id,
		FirstName:  res.FirstName,
		LastName:   res.LastName,
		OtherNames: res.OtherNames.String,
		TenantID:   res.TenantID.String,
		Password:   parsedPassword,
		Email:      parsedEmail,
	}, nil
}

func (r *Repository) GetByID(ctx context.Context, id valueobject.ID) (userbus.User, error) {
	ctx, span := otel.AddSpan(ctx, "userbus.postgres.getuserbyid")
	span.SetAttributes(semconv.DBSystemPostgreSQL)

	defer span.End()

	res, err := r.queries.QueryUserByID(ctx, id.String())
	if err != nil {
		return userbus.User{}, err
	}

	parsedEmail, _ := valueobject.NewEmail(res.Email)

	return userbus.User{
		ID:         id,
		FirstName:  res.FirstName,
		LastName:   res.LastName,
		OtherNames: res.OtherNames.String,
		TenantID:   res.TenantID.String,
		Email:      parsedEmail,
	}, nil
}

func (r *Repository) Update(ctx context.Context, bus userbus.User) error {
	ctx, span := otel.AddSpan(ctx, "userbus.postgres.update")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	if err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		FirstName: bus.FirstName,
		LastName:  bus.LastName,
		OtherNames: pgtype.Text{
			String: bus.OtherNames,
			Valid:  true,
		},
		Email: bus.Email.String(),
	}); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id valueobject.ID) error {
	ctx, span := otel.AddSpan(ctx, "userbus.postgres.delete")
	span.SetAttributes(semconv.DBSystemPostgreSQL)

	defer span.End()

	if err := r.queries.DeleteUser(ctx, id.String()); err != nil {
		return fmt.Errorf("%s: %w", id.String(), err)
	}
	return nil
}

func (r *Repository) Query(ctx context.Context) ([]userbus.User, error) {
	ctx, span := otel.AddSpan(ctx, "userbus.postgres.query")
	span.SetAttributes(semconv.DBSystemPostgreSQL)
	defer span.End()

	res, err := r.queries.QueryUsers(ctx)
	if err != nil {
		return []userbus.User{}, err
	}

	users := make([]userbus.User, len(res))

	for k, v := range res {
		id, _ := valueobject.ParseID(v.ID)
		email, _ := valueobject.NewEmail(v.Email)
		users[k] = userbus.User{
			ID:         id,
			FirstName:  v.FirstName,
			LastName:   v.LastName,
			OtherNames: v.OtherNames.String,
			Email:      email,
			TenantID:   v.TenantID.String,
			CreatedAt:  v.CreatedAt.Time,
			UpdatedAt:  v.UpdatedAt.Time,
		}
	}

	return users, nil
}
