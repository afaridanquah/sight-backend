package userbus

import (
	"context"
)

type Repository interface {
	Add(ctx context.Context, u User) error
	// Query(ctx context.Context) ([]User, error)
	// GetByEmail(ctx context.Context, e valueobject.Email) (User, error)
	// GetByID(ctx context.Context, id valueobject.ID) (User, error)
	// Delete(ctx context.Context, id valueobject.ID) error
	Update(ctx context.Context, u User) error
}
