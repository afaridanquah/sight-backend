package valueobject

import (
	"errors"
	"strings"
)

type Status struct {
	a string
}

var (
	ErrStatusCannotBeEmpty = errors.New("status name cannot be empty")
)

var (
	VERIFIED = Status{"VERIFIED"}
	REJECTED = Status{"REJECTED"}
	PENDING  = Status{"PENDING"}
	DRAFT    = Status{"DRAFT"}
)

var Statuses = []Status{VERIFIED, REJECTED}

func ParseStatus(name string) (Status, error) {
	if name == "" {
		return Status{}, ErrStatusCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "VERIFIED":
		return VERIFIED, nil
	case "REJECTED":
		return REJECTED, nil
	case "PENDING":
		return PENDING, nil
	case "DRAFT":
		return DRAFT, nil
	default:
		return Status{}, errors.New("status name : invalid name")
	}
}

func MustParseStatus(name string) Status {
	status, err := ParseStatus(name)
	if err != nil {
		panic(err)
	}
	return status
}

func (s Status) String() string {
	return s.a
}
