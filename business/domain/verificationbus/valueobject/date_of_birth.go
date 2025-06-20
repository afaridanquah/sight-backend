package valueobject

import (
	"errors"
	"time"
)

type DateOfBirth struct {
	date string
}

var (
	ErrDateCannotBeEmpty   = errors.New("date cannot be empty")
	ErrDateNotValidBeEmpty = errors.New("date is invalid be empty")
)

func NewDateOfBirth(d string) (DateOfBirth, error) {
	if d == "" {
		return DateOfBirth{}, ErrDateCannotBeEmpty
	}

	data, err := time.Parse(time.DateOnly, d)

	if err != nil {
		return DateOfBirth{}, ErrDateNotValidBeEmpty
	}

	return DateOfBirth{data.Format(time.DateOnly)}, nil
}

func (d DateOfBirth) String() string {
	return d.date
}
