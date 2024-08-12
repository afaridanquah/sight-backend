package params

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateCustomerRequest struct {
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Email       string `json:"email"`
	Country     string `json:"country"`
}

func (cr CreateCustomerRequest) Validate() error {
	err := validation.ValidateStruct(&cr,
		validation.Field(&cr.FirstName, validation.Required),
		validation.Field(&cr.LastName, validation.Required),
		validation.Field(&cr.MiddleName),
		validation.Field(&cr.Country, validation.Required, is.CountryCode2),
		validation.Field(&cr.DateOfBirth, validation.Length(2, 2), validation.Date(time.DateOnly)),
		validation.Field(&cr.Email, is.Email),
	)

	if err != nil {
		return err
		// return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "validation.Validate")
	}

	return nil
}
