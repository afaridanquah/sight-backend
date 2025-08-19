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
	STATUSCLEARED        = Status{"CLEARED"}
	STATUSACTIONREQUIRED = Status{"ACTION_REQUIRED"}
	FAILED               = Status{"FAILED"}
	STARTED              = Status{"STARTED"}
	COMPLETED            = Status{"COMPLETED"}
)

var Statuses = []Status{STATUSCLEARED, STATUSACTIONREQUIRED}

func ParseStatus(name string) (Status, error) {
	if name == "" {
		return Status{}, ErrStatusCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "CLEARED":
		return STATUSCLEARED, nil
	case "ACTION_REQUIRED":
		return STATUSACTIONREQUIRED, nil
	case "FAILED":
		return FAILED, nil
	case "STARTED":
		return STARTED, nil
	case "COMPLETED":
		return COMPLETED, nil
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
