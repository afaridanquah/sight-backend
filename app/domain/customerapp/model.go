package customerapp

import (
	"fmt"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
)

type Customer struct {
	ID              string           `json:"id"`
	FirstName       string           `json:"first_name"`
	MiddleName      string           `json:"middle_name"`
	LastName        string           `json:"last_name"`
	DateOfBirth     string           `json:"date_of_birth"`
	Email           string           `json:"email"`
	Country         Country          `json:"birth_country"`
	Identifications []Identification `json:"identifications" validate:"omitempty"`
}

type Country struct {
	AlphaCode2 string `json:"code"`
	Name       string `json:"name"`
}

type Identification struct {
	Pin                string `json:"pin" validate:"required"`
	IssuedCountry      string `json:"issued_country" validate:"required"`
	IdentificationType string `json:"identification_type" validate:"required"`
}

type NewCustomer struct {
	FirstName       string           `json:"first_name" validate:"required"`
	MiddleName      string           `json:"middle_name"`
	LastName        string           `json:"last_name" validate:"required"`
	DateOfBirth     string           `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Email           string           `json:"email" validate:"required"`
	BirthCountry    string           `json:"birth_country" validate:"required"`
	Identifications []Identification `json:"identifications" validate:"omitempty"`
}

func (o NewCustomer) Validate() error {
	if err := ierr.Check(o); err != nil {
		return fmt.Errorf("validate new customer failed: %w", err)
	}

	return nil
}

func toAppCustomer(cus customerbus.Customer) Customer {
	identifications := make([]Identification, len(cus.Identifications))
	if len(cus.Identifications) > 0 {
		for i, idx := range cus.Identifications {
			identifications[i] = Identification{
				Pin:                idx.Pin,
				IssuedCountry:      idx.CountryIssued.Alpha2(),
				IdentificationType: idx.IdentificationType.String(),
			}
		}
	}
	return Customer{
		ID:              cus.ID.String(),
		FirstName:       cus.Person.FirstName,
		MiddleName:      cus.Person.MiddleName,
		LastName:        cus.Person.LastName,
		DateOfBirth:     cus.DateOfBirth.String(),
		Email:           cus.Email.String(),
		Identifications: identifications,
		Country: Country{
			AlphaCode2: cus.BirthCountry.Alpha2(),
			Name:       cus.BirthCountry.Name(),
		},
	}
}

func toAppCustomers(cuss []customerbus.Customer) []Customer {
	app := make([]Customer, len(cuss))
	for i, cus := range cuss {
		app[i] = toAppCustomer(cus)
	}

	return app
}

func toBusinessNewCustomer(c NewCustomer) (customerbus.NewCustomer, error) {
	country, err := valueobject.NewCountry(c.BirthCountry)
	if err != nil {
		return customerbus.NewCustomer{}, fmt.Errorf("newCountry: %w", err)
	}

	dob, err := valueobject.NewDateOfBirth(c.DateOfBirth)
	if err != nil {
		return customerbus.NewCustomer{}, fmt.Errorf("newDateOfBirth: %w", err)
	}

	email, err := valueobject.NewEmail(c.Email)
	if err != nil {
		return customerbus.NewCustomer{}, fmt.Errorf("newEmail: %w", err)
	}

	person, err := valueobject.NewPerson(c.FirstName, c.LastName, &c.MiddleName)
	if err != nil {
		return customerbus.NewCustomer{}, fmt.Errorf("newPerson: %w", err)
	}

	identifications := make([]valueobject.Identification, len(c.Identifications))

	if len(c.Identifications) != 0 {
		for i, idx := range c.Identifications {
			id, err := valueobject.NewIdentification(idx.IdentificationType, idx.Pin, idx.IssuedCountry)
			if err != nil {
				return customerbus.NewCustomer{}, fmt.Errorf("newIdentification: %w", err)
			}
			identifications[i] = id
		}
	}

	return customerbus.NewCustomer{
		Person:          person,
		BirthCountry:    country,
		Email:           email,
		DateOfBirth:     dob,
		Identifications: identifications,
	}, nil

}

// =======================================================================================

type Address struct {
	Address1 string `json:"address_1"`
	Address2 string `json:"address_2"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
	Country  string `json:"country"`
}

type UpdateCustomer struct {
	FirstName   string  `json:"first_name" validate:"required"`
	MiddleName  string  `json:"middle_name"`
	LastName    string  `json:"last_name" validate:"required"`
	DateOfBirth string  `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Email       string  `json:"email" validate:"required"`
	Country     string  `json:"country" validate:"required"`
	Address     Address `json:"address"`
}
