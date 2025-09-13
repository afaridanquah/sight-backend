package valueobject

import (
	"errors"
	"strings"
)

type Side struct {
	a string
}

var (
	ErrSideCannotBeEmpty = errors.New("status name cannot be empty")
)

var (
	FRONT = Side{"FRONT"}
	BACK  = Side{"BACK"}
)

var Sidees = []Side{FRONT, BACK}

func ParseSide(name string) (Side, error) {
	if name == "" {
		return Side{}, ErrSideCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "FRONT":
		return FRONT, nil
	case "BACK":
		return BACK, nil
	default:
		return Side{}, errors.New(" invalid name")
	}
}

func MustParseSide(name string) Side {
	status, err := ParseSide(name)
	if err != nil {
		panic(err)
	}
	return status
}

func (s Side) String() string {
	return s.a
}
