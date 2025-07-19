package identificationbus

import (
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/identificationbus/valueobject"
	"github.com/google/uuid"
)

type Identification struct {
	ID            uuid.UUID
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
