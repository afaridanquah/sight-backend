package valueobject

import (
	"errors"
	"time"
)

type DateOfBirth string

var (
	ErrDateCannotBeEmpty   = errors.New("date cannot be empty")
	ErrDateNotValidBeEmpty = errors.New("date is invalid be empty")
)

func NewDateOfBirth(d string) (DateOfBirth, error) {
	if d == "" {
		return "", ErrDateCannotBeEmpty
	}

	data, err := time.Parse(time.DateOnly, d)

	if err != nil {
		return "", ErrDateNotValidBeEmpty
	}

	return DateOfBirth(data.Format(time.DateOnly)), nil
}
