package valueobject

import (
	"errors"
	"strings"
)

type VerificationType struct {
	name string
}

var (
	ErrVerificationTypeNameRequired = errors.New("verification name is required")
	ErrVerificationTypeInvalid      = errors.New("verification invalid")
)

var (
	AgeVerification   = VerificationType{"Age-Verification"}
	CarrierLookup     = VerificationType{"Carrier-Lookup"}
	GhanaGovermentIDV = VerificationType{"Ghana-Goverment-ID"}
	GlobalCheck       = VerificationType{"Global-Check"}
	PhoneNumber       = VerificationType{"Phone-Number-Verification"}
)

var VerificationTypeOptions = []VerificationType{
	AgeVerification,
	CarrierLookup,
	GhanaGovermentIDV,
	GlobalCheck,
	PhoneNumber,
}

func NewVerificationType(payload string) (VerificationType, error) {
	if payload == "" {
		return VerificationType{}, ErrVerificationTypeNameRequired
	}

	for _, s := range VerificationTypeOptions {
		if strings.EqualFold(s.String(), payload) {
			return s, nil
		}
	}

	return VerificationType{}, ErrVerificationTypeInvalid
}

func (vt VerificationType) String() string {
	return vt.name
}
