package otpbus

import (
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus/valueobject"
	"github.com/google/uuid"
)

type OTP struct {
	ID          uuid.UUID
	Issuer      string
	Channel     valueobject.Channel
	CustomerID  uuid.UUID
	Destination string
	SentAt      time.Time
	ExpiresAt   time.Time
	VerifiedAt  time.Time
	HashedCode  valueobject.HashCode
}

type NewOTP struct {
	Channel     valueobject.Channel
	Destination string
}

type VerifyOTP struct {
	HashedCode valueobject.HashCode
}
