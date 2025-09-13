package postgres

import (
	"encoding/json"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/verificationbus"
	db "bitbucket.org/msafaridanquah/sight-backend/business/sdk/postgres/out"
	"github.com/jackc/pgx/v5/pgtype"
)

type Customer struct {
	ID              string           `json:"id"`
	FirstName       string           `json:"first_name"`
	MiddleName      string           `json:"middle_name"`
	LastName        string           `json:"last_name"`
	OtherNames      string           `json:"other_names"`
	DateOfBirth     string           `json:"date_of_birth"`
	Email           string           `json:"email"`
	PhoneNumber     string           `json:"phone_number"`
	BirthCountry    string           `json:"birth_country"`
	Identifications []Identification `json:"identifications,omitempty"`
	Address         Address          `json:"address"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type Identification struct {
	IdentificationType string `json:"identification_type"`
	Pin                string `json:"pin"`
	PlaceOfBirth       string `json:"place_of_birth"`
	DateOfBirth        string `json:"date_of_birth"`
	Nationality        string `json:"nationality"`
	IssuedDate         string `json:"issued_date"`
	ExpDate            string `json:"exp_date"`
	CountryIssued      string `json:"country_issued"`
}

func toDBInsertVerification(bus verificationbus.Verification) (db.InsertVerificationParams, error) {
	customer := Customer{
		ID:           bus.Customer.ID.String(),
		FirstName:    bus.Customer.Person.FirstName,
		LastName:     bus.Customer.Person.LastName,
		MiddleName:   bus.Customer.Person.MiddleName,
		OtherNames:   bus.Customer.Person.OtherNames,
		DateOfBirth:  bus.Customer.DateOfBirth.String(),
		Email:        bus.Customer.Email.String(),
		PhoneNumber:  bus.Customer.Phone.E164Format,
		BirthCountry: bus.Customer.BirthCountry.Alpha2(),
	}
	customerJson, err := json.Marshal(customer)
	if err != nil {
		return db.InsertVerificationParams{}, err
	}

	return db.InsertVerificationParams{
		ID:         bus.ID,
		CustomerID: bus.CustomerID,
		Customer:   customerJson,
		VerificationType: pgtype.Text{
			String: bus.VerificationType.String(),
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  bus.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  bus.UpdatedAt,
			Valid: true,
		},
	}, nil
}
