package govbus

import (
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/govbus/valueobject"
)

type Identification struct {
	ID                   valueobject.ID
	CustomerID           string
	FirstName            string
	LastName             string
	MiddleName           string
	OtherNames           string
	Sex                  string
	Occupation           string
	MaritalStatus        valueobject.MaritalStatus
	MotherMaidenName     string
	Pin                  string
	CardNumber           string
	PlaceOfBirth         string
	BirthCountry         valueobject.Country
	DateOfBirth          time.Time
	Nationality          valueobject.Country
	IssuedDate           time.Time
	ExpDate              time.Time
	StateOrRegion        string
	PrimaryPhoneNumber   string
	SecondaryPhoneNumber string
	Email                string
	GoogleAddressAlias   string
	City                 string
	Country              valueobject.Country
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type NewIdentification struct {
	CustomerID           string
	FirstName            string
	LastName             string
	MiddleName           string
	OtherNames           string
	Sex                  string
	MaritalStatus        valueobject.MaritalStatus
	MotherMaidenName     string
	Pin                  string
	Occupation           string
	CardNumber           string
	PlaceOfBirth         string
	BirthCountry         valueobject.Country
	DateOfBirth          time.Time
	Nationality          valueobject.Country
	IssuedDate           time.Time
	ExpDate              time.Time
	PrimaryPhoneNumber   string
	SecondaryPhoneNumber string
	Email                string
	GoogleAddressAlias   string
	StateOrRegion        string
	City                 string
	Country              valueobject.Country
}
