package valueobject

import (
	"errors"
	"strings"
)

type Status struct {
	a string
}

var (
	ErrStatusCannotBeEmpty = errors.New("status name : cannot be empty")
)

var (
	ACTIVE    = Status{"ACTIVE"}
	INACTIVE  = Status{"INACTIVE"}
	BLOCKED   = Status{"BLOCKED"}
	SUSPENDED = Status{"SUSPENDED"}
)

var Statuses = []Status{ACTIVE, INACTIVE, BLOCKED, SUSPENDED}

func ParseStatus(name string) (Status, error) {
	if name == "" {
		return Status{}, ErrStatusCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "ACTIVE":
		return ACTIVE, nil
	case "INACTIVE":
		return INACTIVE, nil
	case "BLOCKED":
		return BLOCKED, nil
	case "SUSPENDED":
		return SUSPENDED, nil
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
