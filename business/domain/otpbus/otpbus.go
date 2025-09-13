package otpbus

import (
	"crypto/rand"
	"math/big"
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
	Code        string
	Hash        valueobject.HashCode
}

type NewOTP struct {
	Channel     valueobject.Channel
	Destination string
}

type VerifyOTP struct {
	CustomerID uuid.UUID
	Code       string
}

//====================================================================================================

func generateOTPCode() (string, error) {
	const charset = "0123456789"
	charsetLen := big.NewInt(int64(len(charset)))
	otp := make([]byte, 6)
	for i := range otp {
		randInt, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		otp[i] = charset[randInt.Int64()]
	}

	return string(otp), nil
}
