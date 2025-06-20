package otpapp

import (
	"fmt"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
)

type OTP struct {
	Channel     string `json:"channel"`
	CustomerID  string `json:"customer_id"`
	Destination string `json:"destination"`
	ExpiresAt   string `json:"expires_at"`
	VerifiedAt  string `json:"verified_at"`
}

type NewOTP struct {
	Channel     string `json:"channel"`
	Destination string `json:"destination"`
}

func toBusNewOTP(newapp NewOTP) (otpbus.NewOTP, error) {
	channel, err := valueobject.ParseChannel(newapp.Channel)
	if err != nil {
		return otpbus.NewOTP{}, fmt.Errorf("parseChannel :%w", err)
	}

	return otpbus.NewOTP{
		Channel:     channel,
		Destination: newapp.Destination,
	}, nil
}

func (o NewOTP) Validate() error {
	if err := ierr.Check(o); err != nil {
		return fmt.Errorf("validate new customer: %w", err)
	}

	return nil
}

func toAppOTP(bus otpbus.OTP) OTP {
	return OTP{
		Channel:     bus.Channel.String(),
		CustomerID:  bus.CustomerID.String(),
		Destination: bus.Destination,
		ExpiresAt:   bus.ExpiresAt.Format(time.RFC3339),
	}
}

type UpdateOTP struct {
	Code string `json:"code" validate:"required"`
}

func toBusUpdateOTP(newapp UpdateOTP) (otpbus.VerifyOTP, error) {
	hashed, err := valueobject.ParseToHashCode(newapp.Code)
	if err != nil {
		return otpbus.VerifyOTP{}, fmt.Errorf("parse hash code: %w", err)
	}

	return otpbus.VerifyOTP{
		HashedCode: hashed,
	}, nil
}
