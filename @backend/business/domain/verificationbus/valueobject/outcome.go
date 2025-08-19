package valueobject

import (
	"errors"
	"strings"
)

type Outcome struct {
	a string
}

var (
	ErrOutcomeCannotBeEmpty = errors.New("status name cannot be empty")
)

var (
	CLEARED           = Outcome{"CLEARED"}
	ATTENTION_NEEDED  = Outcome{"ATTENTION_NEEDED"}
	FAILED_OR_UNKNOWN = Outcome{"UNKNOWN"}
)

var Outcomees = []Outcome{CLEARED, ATTENTION_NEEDED}

func ParseOutcome(name string) (Outcome, error) {
	if name == "" {
		return Outcome{}, ErrOutcomeCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "CLEARED":
		return CLEARED, nil
	case "ATTENTION_NEEDED":
		return ATTENTION_NEEDED, nil
	case "UNKNOWN":
		return FAILED_OR_UNKNOWN, nil
	default:
		return Outcome{}, errors.New("status name : invalid name")
	}
}

func MustParseOutcome(name string) Outcome {
	status, err := ParseOutcome(name)
	if err != nil {
		panic(err)
	}
	return status
}

func (s Outcome) String() string {
	return s.a
}
