package otpbus

import (
	"context"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus/valueobject"
	"github.com/google/uuid"
)

type Repository interface {
	Add(ctx context.Context, o OTP) (OTP, error)
	Update(ctx context.Context, id uuid.UUID, o OTP) error
	Find(ctx context.Context, id uuid.UUID) (OTP, error)
	FindByCustomerIDAndHash(ctx context.Context, customerID uuid.UUID, hashCode valueobject.HashCode) (OTP, error)
}
