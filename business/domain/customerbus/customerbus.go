package customerbus

import (
	"errors"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/valueobject"
)

var (
	ErrInvalidPerson      = errors.New("a customer requires a valid name")
	ErrPersonIdIsRequired = errors.New("a customer id is required")
)

type Customer struct {
	ID              valueobject.ID
	Person          valueobject.Person
	UserID          valueobject.ID
	OrgID           valueobject.ID
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
	Person          *valueobject.Person
	BirthCountry    *valueobject.Country
	DateOfBirth     *valueobject.DateOfBirth
	Email           *valueobject.Email
	Address         *valueobject.Address
	Identifications *[]valueobject.Identification
}

type SearchCustomer struct {
	FirstName  *string
	LastName   *string
	MiddleName *string
	FromDate   *time.Time
	ToDate     *time.Time
}
