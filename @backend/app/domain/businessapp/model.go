package businessapp

import (
	"fmt"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
)

type Business struct {
	ID              string   `json:"id"`
	LegalName       string   `json:"legal_name"`
	DoingBusinessAs string   `json:"dba"`
	Entity          string   `json:"entity"`
	EmailAddresses  []string `json:"email_addresses"`
	PhoneNumbers    []string `json:"phone_numbers"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

type NewBusiness struct {
	LegalName       string   `json:"legal_name" validate:"required"`
	DoingBusinessAs string   `json:"dba" validate:"required"`
	Entity          string   `json:"entity" validate:"required"`
	EmailAddresses  []string `json:"email_addresses" validate:"required"`
	PhoneNumbers    []string `json:"phone_numbers" validate:"required"`
}

func (o NewBusiness) Validate() error {
	if err := ierr.Check(o); err != nil {
		return fmt.Errorf("validate new customer failed: %w", err)
	}

	return nil
}

func toBusBusiness(napp NewBusiness) (businessbus.NewBusiness, error) {
	emails, err := valueobject.NewEmails(napp.EmailAddresses)
	if err != nil {
		return businessbus.NewBusiness{}, err
	}

	entity, err := valueobject.ParseBusinessEntity(napp.Entity)
	if err != nil {
		return businessbus.NewBusiness{}, err
	}

	return businessbus.NewBusiness{
		LegalName:       napp.LegalName,
		DoingBusinessAs: napp.DoingBusinessAs,
		Entity:          entity,
		EmailAddresses:  emails,
	}, nil
}

func toAppBusiness(bus businessbus.Business) Business {
	emails := make([]string, len(bus.EmailAddresses))
	if len(bus.EmailAddresses) > 0 {
		for i, e := range bus.EmailAddresses {
			emails[i] = e.String()
		}
	}
	return Business{
		ID:              bus.ID.String(),
		LegalName:       bus.LegalName,
		DoingBusinessAs: bus.DoingBusinessAs,
		EmailAddresses:  emails,
		Entity:          bus.Entity.String(),
		CreatedAt:       bus.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       bus.UpdatedAt.Format(time.RFC3339),
	}
}
