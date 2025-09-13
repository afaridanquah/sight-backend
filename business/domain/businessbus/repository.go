package businessbus

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Add(ctx context.Context, bus Business) error
	Update(ctx context.Context, bus Business) error
	QueryByIDAndOrgID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (Business, error)
	Delete(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error
}
