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
	Identifications []valueobject.Identification
	Email           valueobject.Email
	BirthCountry    valueobject.Country
	Address         valueobject.Address
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         uint64
	Changes         []Event
}

type NewCustomer struct {
	Person          valueobject.Person
	BirthCountry    valueobject.Country
	DateOfBirth     valueobject.DateOfBirth
	Email           valueobject.Email
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

func NewFromEvents(events []Event) *Customer {
	var customer = &Customer{}

	for _, event := range events {
		customer.On(event, false)
	}
}

func (customer *Customer) On(event Event, new bool) {
	// switch e := event.(type) {
	// case *PatientAdmitted:
	// 	p.id = e.ID
	// 	p.age = e.Age
	// 	p.ward = e.Ward

	// case *PatientDischarged:
	// 	p.discharged = true

	// case *PatientTransferred:
	// 	p.ward = e.NewWardNumber

	// }

	// if !new {
	// 	p.version++
	// }
}

func (customer *Customer) raise(event Event) {
	customer.Changes = append(customer.Changes, event)
	customer.On(event, true)
}
