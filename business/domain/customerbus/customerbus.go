package customerbus

import (
	"errors"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus/valueobject"
	"github.com/google/uuid"
)

var (
	ErrInvalidPerson      = errors.New("a customer requires a valid name")
	ErrPersonIdIsRequired = errors.New("a customer id is required")
)

type Customer struct {
	ID              uuid.UUID
	Person          valueobject.Person
	UserID          uuid.UUID
	BusinessID      uuid.UUID
	DateOfBirth     valueobject.DateOfBirth
	CityOfBirth     string
	Identifications []valueobject.Identification
	Email           valueobject.Email
	PhoneNumber     valueobject.Phone
	BirthCountry    valueobject.Country
	Address         valueobject.Address
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type NewCustomer struct {
	Person          valueobject.Person
	BirthCountry    valueobject.Country
	CityOfBirth     string
	DateOfBirth     valueobject.DateOfBirth
	Email           valueobject.Email
	PhoneNumber     valueobject.Phone
	Address         valueobject.Address
	Identifications []valueobject.Identification
}

type UpdateCustomer struct {
	Person          valueobject.Person
	BirthCountry    valueobject.Country
	DateOfBirth     valueobject.DateOfBirth
	Email           valueobject.Email
	Addresses       valueobject.Address
	Identifications []valueobject.Identification
}
