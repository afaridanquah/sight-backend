package valueobject

import "fmt"

type Person struct {
	FirstName  string
	LastName   string
	MiddleName string
	OtherNames string
}

var (
	ErrPersonFirstNameCannotBeEmpty = fmt.Errorf("person: first name cannot be empty")
	ErrPersonLastNameCannotBeEmpty  = fmt.Errorf("person: last name cannot be empty")
)

func NewPerson(fn string, ln string, mn *string, on *string) (Person, error) {
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

	if on != nil {
		person.MiddleName = *on
	}

	return person, nil
}

func (p Person) FullName() string {
	return fmt.Sprintf("%s %s %s", p.FirstName, p.MiddleName, p.LastName)
}
