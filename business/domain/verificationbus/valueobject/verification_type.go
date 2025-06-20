package valueobject

import (
	"fmt"
	"strings"
)

type VerificationType struct {
	a string
}

var (
	DocumentVerification = VerificationType{"DOCUMENT_VERIFCATION"}
	AddressVerification  = VerificationType{"ADDRESS_VERIFICATION"}
	IDVerifciation       = VerificationType{"ID_VERIFCIATION"}
	AMLScreening         = VerificationType{"AML_SCREENING"}
	AdverseMedia         = VerificationType{"ADVERSE_MEDIA"}
	PhoneNumber          = VerificationType{"PHONENUMBER"}
)

func ParseVerificationType(v string) (VerificationType, error) {
	if v == "" {
		return VerificationType{}, fmt.Errorf("parseverification :%s", v)
	}
	upper := strings.ToUpper(v)

	switch upper {
	case "DOCUMENT_VERIFCATION":
		return DocumentVerification, nil
	case "ADDRESS_VERIFICATION":
		return AddressVerification, nil
	case "ID_VERIFCIATION":
		return IDVerifciation, nil
	case "AML_SCREENING":
		return AMLScreening, nil
	case "ADVERSE_MEDIA":
		return AdverseMedia, nil
	case "PHONENUMBER":
		return PhoneNumber, nil
	default:
		return VerificationType{}, fmt.Errorf("parseverification invalid")
	}
}
