package valueobject

import "fmt"

type Person struct {
	FirstName  string
	LastName   string
	MiddleName string
}

var (
	ErrPersonFirstNameCannotBeEmpty = fmt.Errorf("person: first name cannot be empty")
	ErrPersonLastNameCannotBeEmpty  = fmt.Errorf("person: last name cannot be empty")
)

func NewPerson(fn string, ln string, mn *string) (Person, error) {
	if fn == "" {
		return Person{}, ErrPersonFirstNameCannotBeEmpty
	}
	if ln == "" {
		return Person{}, ErrPersonLastNameCannotBeEmpty
	}

	var person = Person{
		FirstName: fn,
		LastName:  ln,
	}

	if mn != nil {
		person.MiddleName = *mn
	}

	return person, nil
}
