package customer

import (
	"errors"

	"github.com/afaridanquah/verifylab-service/internal/domain/customer/valueobject"
	vo "github.com/afaridanquverifylab-serviceend/internal/valueobject"
)

var (
	ErrInvalidPerson      = errors.New("a customer requires a valid name")
	ErrPersonIdIsRequired = errors.New("a customer id is required")
)

type Customer struct {
	id          valueobject.ID
	firstName   string
	middleName  string
	lastName    string
	dateOfBirth vo.DateOfBirth
	email       vo.Email
	country     vo.Country
}

type CustomerOption func(*Customer)

func New(id valueobject.ID, fn string, ln string, c vo.Country, opts ...CustomerOption) (*Customer, error) {
	if id == (valueobject.ID{}) {
		return &Customer{}, ErrInvalidPerson
	}

	if fn == "" || ln == "" {
		return &Customer{}, ErrInvalidPerson
	}

	return &Customer{
		id:        id,
		firstName: fn,
		lastName:  ln,
		country:   c,
	}, nil
}

func (c *Customer) WithMiddleName(m string) CustomerOption {
	return func(c *Customer) {
		c.middleName = m
	}
}

func (c *Customer) WithDateOfBirth(dob vo.DateOfBirth) CustomerOption {
	return func(c *Customer) {
		c.dateOfBirth = dob
	}
}

func (c *Customer) WithEmail(e vo.Email) CustomerOption {
	return func(c *Customer) {
		c.email = e
	}
}

func (c Customer) ID() valueobject.ID {
	return c.id
}

func (c Customer) FirstName() string {
	return c.firstName
}

func (c Customer) Lastname() string {
	return c.lastName
}

func (c Customer) MiddleName() string {
	return c.middleName
}

func (c Customer) DateOfBirth() vo.DateOfBirth {
	return c.dateOfBirth
}

func (c Customer) Email() vo.Email {
	return c.email
}

func (c Customer) Country() vo.Country {
	return c.country
}
