package businessapp

import (
	"fmt"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
)

type Business struct {
	ID                 string   `json:"id"`
	LegalName          string   `json:"legal_name"`
	Jurisdiction       string   `json:"jurisdiction"`
	RegistrationNumber string   `json:"registration_number"`
	DoingBusinessAs    string   `json:"dba"`
	Entity             string   `json:"entity"`
	Website            string   `json:"website"`
	EmailAddresses     []string `json:"email_addresses"`
	PhoneNumbers       []string `json:"phone_numbers"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
	Address            Address  `json:"address,omitzero"`
}

type NewBusiness struct {
	LegalName          string   `json:"legal_name" validate:"required"`
	RegistrationNumber string   `json:"registration_number"`
	Website            string   `json:"website" validate:"required,url"`
	DoingBusinessAs    string   `json:"dba" validate:"required"`
	Jurisdiction       string   `json:"jurisdiction" validate:"required"`
	Entity             string   `json:"entity" validate:"required"`
	EmailAddresses     []string `json:"email_addresses" validate:"required"`
	PhoneNumbers       []string `json:"phone_numbers" validate:"required"`
	Address            Address  `json:"address,omitzero"`
}

type NewDocument struct {
	Classification string `json:"classification" validate:"required"`
	DocumentType   string `json:"document_type" validate:"required"`
}

type Address struct {
	Line1   string `json:"line_1" validate:"required"`
	Line2   string `json:"line_2"`
	City    string `json:"city" validate:"required"`
	State   string `json:"state" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	Country string `json:"country" validate:"required"`
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

	phones, err := valueobject.ParseIntlPhoneNumbers(napp.PhoneNumbers)
	if err != nil {
		return businessbus.NewBusiness{}, err
	}

	entity, err := valueobject.ParseBusinessEntity(napp.Entity)
	if err != nil {
		return businessbus.NewBusiness{}, err
	}

	jurisdiction, err := valueobject.NewCountry(napp.Jurisdiction)
	if err != nil {
		return businessbus.NewBusiness{}, err
	}

	address, err := valueobject.ParseAddress(napp.Address.Line1, &napp.Address.Line2, napp.Address.City, napp.Address.State, napp.Address.Country)
	if err != nil {
		return businessbus.NewBusiness{}, err
	}

	return businessbus.NewBusiness{
		LegalName:              napp.LegalName,
		DoingBusinessAs:        napp.DoingBusinessAs,
		Website:                napp.Website,
		Entity:                 entity,
		EmailAddresses:         emails,
		PhoneNumbers:           phones,
		Address:                address,
		CountryOfIncorporation: jurisdiction,
		RegistrationNumber:     napp.RegistrationNumber,
	}, nil
}

func toAppBusiness(bus businessbus.Business) Business {
	emails := make([]string, len(bus.EmailAddresses))
	if len(bus.EmailAddresses) > 0 {
		for i, e := range bus.EmailAddresses {
			emails[i] = e.String()
		}
	}

	phones := make([]string, len(bus.PhoneNumbers))
	if len(bus.PhoneNumbers) > 0 {
		for i, e := range bus.PhoneNumbers {
			phones[i] = e.E164Format
		}
	}

	return Business{
		ID:              bus.ID.String(),
		LegalName:       bus.LegalName,
		DoingBusinessAs: bus.DoingBusinessAs,
		EmailAddresses:  emails,
		PhoneNumbers:    phones,
		Address: Address{
			Line1:   bus.Address.Line1,
			Line2:   bus.Address.Line2,
			City:    bus.Address.City,
			State:   bus.Address.StateOrRegion,
			Country: bus.Address.Country.Alpha2(),
		},
		RegistrationNumber: bus.RegistrationNumber,
		Jurisdiction:       bus.CountryOfIncorporation.Alpha2(),
		Website:            bus.Website,
		Entity:             bus.Entity.String(),
		CreatedAt:          bus.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          bus.UpdatedAt.Format(time.RFC3339),
	}
}

// ============================================================================================

type UpdateBusiness struct {
	LegalName       *string   `json:"legal_name"`
	Website         *string   `json:"website"`
	DoingBusinessAs *string   `json:"dba"`
	Entity          *string   `json:"entity"`
	EmailAddresses  *[]string `json:"email_addresses"`
	PhoneNumbers    *[]string `json:"phone_numbers"`
	Address         *Address  `json:"address"`
}

func (o UpdateBusiness) Validate() error {
	if err := ierr.Check(o); err != nil {
		return fmt.Errorf("validate update customer failed: %w", err)
	}

	return nil
}

func toBusUpdateBusiness(up UpdateBusiness) (businessbus.UpdateBusiness, error) {
	var bus businessbus.UpdateBusiness

	if up.EmailAddresses != nil {
		emails, err := valueobject.NewEmails(*up.EmailAddresses)
		if err != nil {
			return businessbus.UpdateBusiness{}, err
		}
		bus.EmailAddresses = &emails
	}

	if up.PhoneNumbers != nil {
		phones, err := valueobject.ParseIntlPhoneNumbers(*up.PhoneNumbers)
		if err != nil {
			return businessbus.UpdateBusiness{}, err
		}
		bus.PhoneNumbers = &phones
	}

	if up.Entity != nil {
		entity, err := valueobject.ParseBusinessEntity(*up.Entity)
		if err != nil {
			return businessbus.UpdateBusiness{}, err
		}
		bus.Entity = &entity
	}

	if up.Address != nil {
		address, err := valueobject.ParseAddress(up.Address.Line1, &up.Address.Line2, up.Address.City, up.Address.State, up.Address.Country)
		if err != nil {
			return businessbus.UpdateBusiness{}, err
		}
		bus.Address = &address
	}

	return bus, nil
}
