package valueobject

import (
	"errors"
	"strings"
)

type Summary struct {
	a string
}

var (
	ErrSummaryCannotBeEmpty = errors.New("status name cannot be empty")
)

var (
	APPROVED = Summary{"APPROVED"}
	UNKNOWN  = Summary{"UNKNOWN"}
	REVIEW   = Summary{"DECLINED"}
	DECLINED = Summary{"DECLINED"}
)

var Summaryes = []Summary{APPROVED, REVIEW, UNKNOWN, DECLINED}

func ParseSummary(name string) (Summary, error) {
	if name == "" {
		return Summary{}, ErrSummaryCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "APPROVED":
		return APPROVED, nil
	case "REVIEW":
		return REVIEW, nil
	case "DECLINED":
		return DECLINED, nil
	case "UNKNOWN":
		return UNKNOWN, nil
	default:
		return Summary{}, errors.New("status name : invalid name")
	}
}

func MustParseSummary(name string) Summary {
	status, err := ParseSummary(name)
	if err != nil {
		panic(err)
	}
	return status
}

func (s Summary) String() string {
	return s.a
}
