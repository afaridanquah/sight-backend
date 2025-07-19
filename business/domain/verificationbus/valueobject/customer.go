package valueobject

import (
	"fmt"

	"github.com/google/uuid"
)

type Customer struct {
	ID              uuid.UUID
	Person          Person
	DateOfBirth     DateOfBirth
	Identifications []Identification
	Email           Email
	BirthCountry    Country
	Address         Address
	Phone           Phone
}

func ParseCustomer(id uuid.UUID, person Person, dob *DateOfBirth, birthcountry *Country, identifications *[]Identification, email *Email, phone *Phone) (Customer, error) {
	if id == (uuid.Nil) {
		return Customer{}, fmt.Errorf("id is required")
	}

	if person == (Person{}) {
		return Customer{}, fmt.Errorf("person is required")
	}

	customer := Customer{
		ID:     id,
		Person: person,
	}

	if dob != nil {
		customer.DateOfBirth = *dob
	}

	if email != nil {
		customer.Email = *email
	}

	if birthcountry != nil {
		customer.BirthCountry = *birthcountry
	}

	if phone != nil {
		customer.Phone = *phone
	}

	if identifications != nil {
		customer.Identifications = *identifications
	}

	return customer, nil
}
