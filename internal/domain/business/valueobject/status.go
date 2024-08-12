package valueobject

import (
	"errors"
	"strings"
)

type Status struct {
	a string
}

var (
	ErrStatusNameCannotBeEmpty = errors.New("status name cannot be empty")
	ErrStatusNotValid          = errors.New("status is invalid")
)

var (
	Active   = Status{"Active"}
	Inactive = Status{"Inactive"}
)

var StatusList = []Status{
	Active,
	Inactive,
}

func NewStatus(s string) (Status, error) {
	if s == "" {
		return Status{}, ErrStatusNameCannotBeEmpty
	}

	for _, ss := range StatusList {
		if strings.EqualFold(ss.String(), s) {
			return ss, nil
		}
	}

	return Status{}, ErrStatusNotValid
}

func (s Status) String() string {
	return s.a
}
