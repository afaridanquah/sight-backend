package otpbus

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Add(ctx context.Context, o OTP) error
	// Update(ctx context.Context, id uuid.UUID, o OTP) error
	// Find(ctx context.Context, id uuid.UUID) (OTP, error)
	FindByCustomerIDAndHash(ctx context.Context, customerID uuid.UUID, hash string) (OTP, error)
}
