package customerapp

import (
	"fmt"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/app/sdk/errs"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus"
	dvo "bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus"
	ivo "bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus/valueobject"
)

type Customer struct {
	ID              string           `json:"id"`
	FirstName       string           `json:"first_name"`
	MiddleName      string           `json:"middle_name"`
	LastName        string           `json:"last_name"`
	DateOfBirth     string           `json:"date_of_birth"`
	Email           string           `json:"email"`
	PhoneNumber     string           `json:"phone_number"`
	BirthCountry    Country          `json:"birth_country"`
	Identifications []Identification `json:"identifications" validate:"omitempty,dive"`
}

type Country struct {
	AlphaCode2 string `json:"code"`
	Name       string `json:"name"`
}

type PhoneNumber struct {
	Digits  string `json:"digits" validate:"required"`
	Country string `json:"country" validate:"required"`
}

type NewCustomer struct {
	FirstName       string              `json:"first_name" validate:"required"`
	MiddleName      string              `json:"middle_name"`
	LastName        string              `json:"last_name" validate:"required"`
	OtherNames      string              `json:"other_names"`
	DateOfBirth     string              `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Email           string              `json:"email" validate:"required_if=PhoneNumber ''"`
	PhoneNumber     *PhoneNumber        `json:"phone_number" validate:"required_if=Email ''"`
	BirthCountry    string              `json:"birth_country" validate:"required"`
	Identifications []NewIdentification `json:"identifications" validate:"omitempty,dive"`
}

func (o NewCustomer) Validate() error {
	if err := errs.Check(o); err != nil {
		return fmt.Errorf("validate new customer failed: %w", err)
	}

	return nil
}

func toAppCustomer(cus customerbus.Customer) Customer {
	identifications := make([]Identification, len(cus.Identifications))
	if len(cus.Identifications) > 0 {
		for i, idx := range cus.Identifications {
			expDate := idx.ExpDate.Format(time.DateOnly)
			identifications[i] = Identification{
				Pin: idx.Pin,
				IssuedCountry: Country{
					AlphaCode2: idx.IssedCountry.Alpha2(),
					Name:       idx.IssedCountry.Name(),
				},
				IdentificationType: idx.IdentificationType.String(),
				Nationality: Country{
					AlphaCode2: idx.Nationality.Alpha2(),
					Name:       idx.Nationality.Name(),
				},
				ExpDate:    expDate,
				IssuedDate: idx.IssuedDate.Format(time.DateOnly),
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
		PhoneNumber:     cus.PhoneNumber.E164Format,
		BirthCountry: Country{
			AlphaCode2: cus.BirthCountry.Alpha2(),
			Name:       cus.BirthCountry.Name(),
		},
	}
}

// func toAppCustomers(cuss []customerbus.Customer) []Customer {
// 	app := make([]Customer, len(cuss))
// 	for i, cus := range cuss {
// 		app[i] = toAppCustomer(cus)
// 	}

// 	return app
// }

func toBusNewCustomer(c NewCustomer) (customerbus.NewCustomer, error) {
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

	person, err := valueobject.NewPerson(c.FirstName, c.LastName, &c.MiddleName, &c.OtherNames)
	if err != nil {
		return customerbus.NewCustomer{}, fmt.Errorf("newPerson: %w", err)
	}

	customer := customerbus.NewCustomer{
		Person:       person,
		BirthCountry: country,
		Email:        email,
		DateOfBirth:  dob,
		// Identifications: identifications,
	}

	if c.PhoneNumber != nil {
		phone, err := valueobject.ParsePhone(c.PhoneNumber.Country, c.PhoneNumber.Digits)
		if err != nil {
			return customerbus.NewCustomer{}, fmt.Errorf("parse phone: %w", err)
		}
		customer.PhoneNumber = phone
	}

	return customer, nil
}

// -----------------------------------------------------------------------------------
// Parse identification request to business logic

type NewIdentification struct {
	Pin                string  `json:"pin" validate:"required"`
	IssuedCountry      string  `json:"issued_country" validate:"required"`
	IdentificationType string  `json:"identification_type" validate:"required"`
	Nationality        *string `json:"nationality,omitempty"`
	IssuedDate         *string `json:"issued_date,omitempty" validate:"datetime=2006-01-02"`
	ExpDate            *string `json:"exp_date,omitempty" validate:"datetime=2006-01-02"`
}

type NewIdentifications struct {
	Identifications []NewIdentification `json:"identifications" validate:"required,dive"`
}

type Identification struct {
	ID                 string  `json:"id"`
	Pin                string  `json:"pin"`
	IssuedCountry      Country `json:"issued_country"`
	IdentificationType string  `json:"identification_type"`
	Nationality        Country `json:"nationality"`
	IssuedDate         string  `json:"issued_date"`
	ExpDate            string  `json:"exp_date"`
}

func toAppIdentification(bus identificationbus.Identification) Identification {
	return Identification{
		ID:  bus.ID.String(),
		Pin: bus.Pin,
		IssuedCountry: Country{
			AlphaCode2: bus.Country.Alpha2(),
			Name:       bus.Country.Name(),
		},
		IdentificationType: bus.IdentificationType.String(),
		Nationality: Country{
			AlphaCode2: bus.Nationality.Alpha2(),
			Name:       bus.Nationality.Name(),
		},
		IssuedDate: bus.IssuedDate.Format(time.RFC3339),
		ExpDate:    bus.ExpDate.Format(time.RFC3339),
	}
}

func (o NewIdentifications) Validate() error {
	if err := errs.Check(o); err != nil {
		return fmt.Errorf("validate new identication failed: %w", err)
	}

	return nil
}

func toBusNewIdentifications(idx NewIdentifications, bcus customerbus.Customer) ([]identificationbus.NewIdentification, error) {
	identifications := make([]identificationbus.NewIdentification, len(idx.Identifications))
	if len(idx.Identifications) > 0 {
		for k, v := range idx.Identifications {
			idt, err := ivo.ParseIdentificationType(v.IdentificationType)
			if err != nil {
				return []identificationbus.NewIdentification{}, err
			}

			issued, err := time.Parse(time.DateOnly, *v.IssuedDate)
			if err != nil {
				return []identificationbus.NewIdentification{}, err
			}

			exp, err := time.Parse(time.DateOnly, *v.ExpDate)
			if err != nil {
				return []identificationbus.NewIdentification{}, err
			}

			issc, err := ivo.NewCountry(v.IssuedCountry)
			if err != nil {
				return []identificationbus.NewIdentification{}, err
			}

			identifications[k] = identificationbus.NewIdentification{
				CustomerID:         bcus.ID.String(),
				Pin:                v.Pin,
				IdentificationType: idt,
				IssuedDate:         issued,
				ExpDate:            exp,
				Country:            issc,
			}

			if v.Nationality != nil {
				national, err := ivo.NewCountry(*v.Nationality)
				if err != nil {
					return []identificationbus.NewIdentification{}, err
				}
				identifications[k].Nationality = national
			}
		}
	}

	return identifications, nil
}

// func toBusNewIdentifications(c NewCustomer, bcus customerbus.Customer) ([]identificationbus.NewIdentification, error) {
// 	identifications := make([]identificationbus.NewIdentification, len(c.Identifications))
// 	if c.Identifications != nil {
// 		if len(c.Identifications) > 0 {
// 			for k, v := range c.Identifications {
// 				idt, err := ivo.ParseIdentificationType(v.IdentificationType)
// 				if err != nil {
// 					return []identificationbus.NewIdentification{}, err
// 				}

// 				issued, err := time.Parse(time.DateOnly, *v.IssuedDate)
// 				if err != nil {
// 					return []identificationbus.NewIdentification{}, err
// 				}

// 				exp, err := time.Parse(time.DateOnly, *v.ExpDate)
// 				if err != nil {
// 					return []identificationbus.NewIdentification{}, err
// 				}

// 				issc, err := ivo.NewCountry(v.IssuedCountry)
// 				if err != nil {
// 					return []identificationbus.NewIdentification{}, err
// 				}

// 				identifications[k] = identificationbus.NewIdentification{
// 					CustomerID:         bcus.ID.String(),
// 					Pin:                v.Pin,
// 					IdentificationType: idt,
// 					IssuedDate:         issued,
// 					ExpDate:            exp,
// 					Country:            issc,
// 				}
// 			}
// 		}
// 	}
// 	return identifications, nil
// }

// =======================================================================================

type Address struct {
	Address1 string `json:"line_1"`
	Address2 string `json:"line_2"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
	Country  string `json:"country"`
}

type UpdateCustomer struct {
	FirstName   *string  `json:"first_name"`
	MiddleName  *string  `json:"middle_name"`
	LastName    *string  `json:"last_name"`
	OtherNames  *string  `json:"other_names"`
	DateOfBirth *string  `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
	Email       *string  `json:"email"`
	Country     *string  `json:"birth_country" validate:"omitempty"`
	Address     *Address `json:"address" validate:"omitempty"`
}

func (o UpdateCustomer) Validate() error {
	if err := errs.Check(o); err != nil {
		return fmt.Errorf("validate update customer failed: %w", err)
	}

	return nil
}

func toBusUpdateCustomer(uapp UpdateCustomer) (customerbus.UpdateCustomer, error) {
	var ubus customerbus.UpdateCustomer
	if uapp.FirstName != nil || uapp.LastName != nil || uapp.MiddleName != nil || uapp.OtherNames != nil {
		p, err := valueobject.NewPerson(*uapp.FirstName, *uapp.LastName, uapp.MiddleName, uapp.OtherNames)
		if err != nil {
			return customerbus.UpdateCustomer{}, err
		}
		ubus.Person = &p
	}

	bc, err := valueobject.NewCountry(*uapp.Country)
	if err != nil {
		return customerbus.UpdateCustomer{}, err
	}
	ubus.BirthCountry = &bc

	email, err := valueobject.NewEmail(*uapp.Email)
	if err != nil {
		return customerbus.UpdateCustomer{}, err
	}
	ubus.Email = &email

	return ubus, nil
}

// ===========================================================================================
// Attach a document to a customer.
type NewDocument struct {
	Classification string `json:"classification" form:"classification" validate:"required"`
	DocumentType   string `json:"document_type" form:"document_type" validate:"required"`
}

type Business struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Dba  string `json:"dba"`
}

type Document struct {
	ID        string   `json:"id"`
	FileName  string   `json:"filename"`
	MimeType  string   `json:"mime_type"`
	Business  Business `json:"business,omitzero"`
	Side      string   `json:"side"`
	Customer  Customer `json:"customer,omitzero"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func (o NewDocument) Validate() error {
	if err := errs.Check(o); err != nil {
		return fmt.Errorf("validate new document failed: %w", err)
	}

	return nil
}

func toAppDocument(bus documentbus.Document) Document {
	return Document{
		ID:        bus.ID.String(),
		FileName:  bus.FileName,
		Side:      bus.Side.String(),
		CreatedAt: bus.CreatedAt.Format(time.RFC3339),
		UpdatedAt: bus.UpdatedAt.Format(time.RFC3339),
	}
}

func toBusNewDocument(napp NewDocument) (documentbus.NewDocument, error) {
	dt, err := dvo.ParseDocumentType(napp.DocumentType)
	if err != nil {
		return documentbus.NewDocument{}, err
	}

	classification, err := dvo.ParseClassification(napp.Classification)
	if err != nil {
		return documentbus.NewDocument{}, err
	}

	return documentbus.NewDocument{
		DocumentType:   dt,
		Classification: classification,
	}, nil

}
