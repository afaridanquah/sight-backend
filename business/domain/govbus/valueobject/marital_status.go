package valueobject

import (
	"errors"
	"strings"
)

type MaritalStatus struct {
	a string
}

var (
	single    = MaritalStatus{"SINGLE"}
	married   = MaritalStatus{"MARRIED"}
	divorced  = MaritalStatus{"DIVORCED"}
	widowed   = MaritalStatus{"WIDOWED"}
	separated = MaritalStatus{"SEPARATED"}
)

func ParseMaritalStatus(a string) (MaritalStatus, error) {
	if a == "" {
		return MaritalStatus{}, errors.New("status is required")
	}

	uppercase := strings.ToUpper(a)

	switch uppercase {
	case "SINGLE":
		return single, nil
	case "MARRIED":
		return married, nil
	case "DIVORCED":
		return divorced, nil
	case "WIDOWED":
		return widowed, nil
	case "SEPARATED":
		return separated, nil
	default:
		return MaritalStatus{}, errors.New("invalid")
	}
}

func (m MaritalStatus) String() string {
	return m.a
}
