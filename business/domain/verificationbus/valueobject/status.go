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
	CLEARED        = Status{"CLEARED"}
	ACTIONREQUIRED = Status{"ACTION_REQUIRED"}
	FAILED         = Status{"FAILED"}
)

var Statuses = []Status{CLEARED, ACTIONREQUIRED}

func ParseStatus(name string) (Status, error) {
	if name == "" {
		return Status{}, ErrStatusCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "CLEARED":
		return CLEARED, nil
	case "ACTION_REQUIRED":
		return ACTIONREQUIRED, nil
	case "FAILED":
		return FAILED, nil
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
