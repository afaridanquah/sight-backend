package govbus

import (
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/govbus/valueobject"
)

type Identification struct {
	ID            valueobject.ID
	CustomerID    string
	FirstName     string
	LastName      string
	MiddleName    string
	Sex           string
	Pin           string
	PlaceOfBirth  string
	DateOfBirth   time.Time
	Nationality   valueobject.Country
	IssuedDate    time.Time
	ExpDate       time.Time
	StateOrRegion string
	City          string
	Country       valueobject.Country
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type NewIdentification struct {
	CustomerID    string
	FirstName     string
	LastName      string
	MiddleName    string
	Sex           string
	Pin           string
	PlaceOfBirth  string
	DateOfBirth   time.Time
	Nationality   valueobject.Country
	IssuedDate    time.Time
	ExpDate       time.Time
	StateOrRegion string
	City          string
	Country       valueobject.Country
}
