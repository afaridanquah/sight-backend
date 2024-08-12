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
	Approved   = Status{"Approved"}
	Completed  = Status{"Completed"}
	Created    = Status{"Created"}
	Declined   = Status{"Declined"}
	Expired    = Status{"Expired"}
	Failed     = Status{"Failed"}
	NeedReview = Status{"Needs Review"}
	Pending    = Status{"Pending"}
)

var StatusList = []Status{
	Approved,
	Completed,
	Created,
	Declined,
	Expired,
	Failed,
	NeedReview,
	Pending,
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
