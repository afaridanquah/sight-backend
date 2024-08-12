package params

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateVerificationRequest struct {
	VerificationType string                            `json:"verification_type"`
	Customer         CreateVerificationRequestCustomer `json:"customer"`
	CustomerId       string                            `json:"customer_id"`
}

type CreateVerificationRequestCustomer struct {
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Country     string `json:"country"`
}

func (ra CreateVerificationRequest) Validate() error {
	return validation.ValidateStruct(&ra,
		validation.Field(&ra.VerificationType, validation.Required),
		validation.Field(&ra.Customer, validation.When(ra.CustomerId == "", validation.Required, validation.By(func(customer interface{}) error {
			cus := ra.Customer
			return validation.ValidateStruct(&cus,
				validation.Field(&cus.FirstName, validation.Required),
				validation.Field(&cus.LastName, validation.Required),
				validation.Field(&cus.Country, validation.Required, is.CountryCode2),
			)
		}))),
		validation.Field(&ra.CustomerId, validation.When(ra.Customer == CreateVerificationRequestCustomer{}, validation.Required)),
	)

}
