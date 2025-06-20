package identificationbus

import (
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/identificationbus/valueobject"
	"github.com/google/uuid"
)

type Identification struct {
	ID                 uuid.UUID
	Number             string
	PlaceOfBirth       string
	DateOfBirth        time.Time
	IdentificationType valueobject.IdentificationType
	IssuedCountry      valueobject.Country
	Person             valueobject.Person
	Address            valueobject.Address
	IssuedDate         time.Time
	ExpDate            time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type NewIdentification struct {
	Number             string
	PlaceOfBirth       string
	DateOfBirth        time.Time
	IdentificationType valueobject.IdentificationType
	Person             valueobject.Person
	Address            valueobject.Address
	IssuedCountry      valueobject.Country
	IssuedDate         time.Time
	ExpDate            time.Time
}
