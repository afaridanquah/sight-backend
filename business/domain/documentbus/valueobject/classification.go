package valueobject

import (
	"fmt"
	"strings"
)

type Classification struct {
	a string
}

var (
	PROOF_OF_IDENTITY = Classification{"PROOF_OF_IDENTITY"}
	PROOF_OF_ADDRESS  = Classification{"PROOF_OF_ADDRESS"}
)

func ParseClassification(a string) (Classification, error) {
	if a == "" {
		return Classification{}, fmt.Errorf("cannot be empty")
	}
	uppercase := strings.ToUpper(a)

	switch uppercase {
	case "PROOF_OF_IDENTITY":
		return PROOF_OF_IDENTITY, nil
	case "PROOF_OF_ADDRESS":
		return PROOF_OF_ADDRESS, nil
	default:
		return Classification{}, fmt.Errorf("not valid")
	}
}
