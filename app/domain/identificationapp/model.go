package identificationapp

import (
	"fmt"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/identificationbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/identificationbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
)

type Identification struct {
	Number             string  `json:"number"`
	PlaceOfBirth       string  `json:"place_of_birth"`
	DateOfBirth        string  `json:"date_of_birth"`
	IdentificationType string  `json:"identification_type"`
	IssuedCountry      Country `json:"country"`
	Person             Person  `json:"person"`
	Address            Address `json:"address"`
	IssueDate          string  `json:"issued_date"`
	ExpDate            string  `json:"exp_date"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type Address struct {
	Street  string
	City    string
	State   string
	Country Country
}

type Country struct {
	AlphaCode2 string `json:"code"`
	Name       string `json:"name"`
}

type Person struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type NewIdentification struct {
	Number             string `json:"pin" validate:"required"`
	PlaceOfBirth       string `json:"place_of_birth"`
	DateOfBirth        string `json:"date_of_birth" validate:"required"`
	IdentificationType string `json:"identification_type" validate:"required"`
	IssuedCountry      string `json:"country_code" validate:"required"`
	Person             Person `json:"person"`
	IssueDate          string `json:"issued_date"`
	ExpDate            string `json:"exp_date"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

func toBusNewIdentification(napp NewIdentification) (identificationbus.NewIdentification, error) {
	idType, err := valueobject.ParseIdentificationType(napp.IdentificationType)
	if err != nil {
		return identificationbus.NewIdentification{}, err
	}

	dateOfBirth, err := time.Parse(time.DateOnly, napp.DateOfBirth)
	issuedDate, err := time.Parse(time.DateOnly, napp.IssueDate)
	person, err := valueobject.NewPerson(napp.Person.FirstName, napp.Person.FirstName, &napp.Person.MiddleName)

	return identificationbus.NewIdentification{
		Number:             napp.Number,
		PlaceOfBirth:       napp.PlaceOfBirth,
		DateOfBirth:        dateOfBirth,
		IssuedDate:         issuedDate,
		IdentificationType: idType,
		Person:             person,
	}, nil
}

func toAppIdentification(bus identificationbus.Identification) Identification {
	return Identification{
		Number: bus.Number,
	}
}

func (o NewIdentification) Validate() error {
	if err := ierr.Check(o); err != nil {
		return fmt.Errorf("validate new rate failed: %w", err)
	}

	return nil
}
