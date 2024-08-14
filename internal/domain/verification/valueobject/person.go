package valueobject

import (
	"errors"

	vo "bitbucket.org/msafaridanquah/verifylab-service/internal/valueobject"
)

type Person struct {
	id          string
	firstName   string
	middleName  string
	lastName    string
	dateOfBirth vo.DateOfBirth
	email       vo.Email
	country     vo.Country
}

var (
	ErrPersonIdCannotBeEmpty        = errors.New("person id cannot be empty")
	ErrPersonFirstNameCannotBeEmpty = errors.New("person first name cannot be empty")
	ErrPersonLastNameCannotBeEmpty  = errors.New("person last name cannot be empty")
)

func NewPerson(id string, fn string, middle string, ln string, dob vo.DateOfBirth, email vo.Email, c vo.Country) (Person, error) {
	if id == "" {
		return Person{}, ErrPersonIdCannotBeEmpty
	}
	if fn == "" {
		return Person{}, ErrPersonFirstNameCannotBeEmpty
	}
	if ln == "" {
		return Person{}, ErrPersonLastNameCannotBeEmpty
	}
	var per = Person{
		id:          id,
		firstName:   fn,
		middleName:  middle,
		lastName:    ln,
		dateOfBirth: dob,
		email:       email,
		country:     c,
	}

	return per, nil
}
