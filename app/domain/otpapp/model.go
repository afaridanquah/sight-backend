package otpapp

import (
	"fmt"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"github.com/google/uuid"
)

type OTP struct {
	ID         string `json:"id"`
	Channel    string `json:"channel"`
	CustomerID string `json:"customer_id"`
}

type NewOTP struct {
	Channel     string `json:"channel" validate:"required,oneof=sms email"`
	Destination string `json:"destination"`
}

func toBusNewOTP(newapp NewOTP) (otpbus.NewOTP, error) {
	channel, err := valueobject.ParseChannel(newapp.Channel)
	if err != nil {
		return otpbus.NewOTP{}, fmt.Errorf("parseChannel :%w", err)
	}

	return otpbus.NewOTP{
		Channel: channel,
	}, nil
}

func (o NewOTP) Validate() error {
	if err := ierr.Check(o); err != nil {
		return fmt.Errorf("validate new otp request: %w", err)
	}

	return nil
}

func toAppOTP(bus otpbus.OTP) OTP {
	return OTP{
		ID:         bus.ID.String(),
		Channel:    bus.Channel.String(),
		CustomerID: bus.CustomerID.String(),
	}
}

// ===================================================================

type VerifyOTP struct {
	Code string `json:"code" validate:"required"`
}

func (o VerifyOTP) Validate() error {
	if err := ierr.Check(o); err != nil {
		return fmt.Errorf("validate otp: %w", err)
	}

	return nil
}

func toBusVerifyOTP(newapp VerifyOTP, customerID uuid.UUID) (otpbus.VerifyOTP, error) {
	return otpbus.VerifyOTP{
		Code:       newapp.Code,
		CustomerID: customerID,
	}, nil
}
