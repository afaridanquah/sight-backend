package verificationbus

import (
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus/valueobject"
	"github.com/google/uuid"
)

type Verification struct {
	ID               uuid.UUID
	CustomerID       uuid.UUID
	BusinessID       uuid.UUID
	Business         valueobject.Business
	VerificationType valueobject.VerificationType
	Status           valueobject.Status
	Customer         valueobject.Customer
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type NewVerification struct {
	VerificationType valueobject.VerificationType
	BusinessID       uuid.UUID
	CustomerID       uuid.UUID
}
