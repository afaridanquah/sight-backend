package valueobject

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	FirstName     string
	LastName      string
	OtherNames    *string
	PreferredName *string
}

var (
	ErrUserIDNameRequired    = errors.New("id required")
	ErrUserFirstNameRequired = errors.New("last name required")
	ErrUserLastNameRequired  = errors.New("last name required")
)

func ParseUser(id uuid.UUID, firstName string, lastName string, otherNames *string, preferredNames *string) (User, error) {
	if id == uuid.Nil {
		return User{}, ErrUserIDNameRequired
	}

	if firstName == "" {
		return User{}, ErrUserFirstNameRequired
	}

	if lastName == "" {
		return User{}, ErrUserLastNameRequired
	}

	return User{
		ID:            id,
		FirstName:     firstName,
		LastName:      lastName,
		OtherNames:    otherNames,
		PreferredName: preferredNames,
	}, nil
}
