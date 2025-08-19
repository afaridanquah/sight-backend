package valueobject

import (
	"errors"
	"strings"
)

type BusinessEntity struct {
	a string
}

var (
	ErrBusinessEntityCannotBeEmpty = errors.New("status name : cannot be empty")
)

var (
	ESTE = BusinessEntity{"ESTATE"}
	SPRO = BusinessEntity{"SOLE_PROPRIETOR"}
	CORP = BusinessEntity{"CORPORATION"}
	EORG = BusinessEntity{"EXEMPT_ORGANIZATION"}
)

var BusinessEntities = []BusinessEntity{ESTE, SPRO, CORP, EORG}

func ParseBusinessEntity(name string) (BusinessEntity, error) {
	if name == "" {
		return BusinessEntity{}, ErrBusinessEntityCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "ESTATE":
		return ESTE, nil
	case "SOLE_PROPRIETOR":
		return SPRO, nil
	case "CORPORATION":
		return CORP, nil
	case "EXEMPT_ORGANIZATION":
		return EORG, nil
	default:
		return BusinessEntity{}, errors.New("status name : invalid name")
	}
}

func (e BusinessEntity) String() string {
	return e.a
}

func MustParseBusinessEntity(name string) BusinessEntity {
	busType, err := ParseBusinessEntity(name)
	if err != nil {
		panic(err)
	}
	return busType
}
