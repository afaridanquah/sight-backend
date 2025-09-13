package valueobject

import (
	"fmt"
	"strings"
)

type VerificationType struct {
	a string
}

var (
	DocumentInsight     = VerificationType{"DOCUMENT_INSIGHT"}
	AddressVerification = VerificationType{"ADDRESS_VERIFICATION"}
	GovVendor           = VerificationType{"GOV_VERIFICATION"}
	AMLScreening        = VerificationType{"AML_SCREENING"}
	AdverseMedia        = VerificationType{"ADVERSE_MEDIA"}
	PhoneNumber         = VerificationType{"PHONENUMBER"}
	Selfie              = VerificationType{"SELFIE"}
	Liveness            = VerificationType{"LIVENESS"}
)

var VerificationTypes = []VerificationType{DocumentInsight, AddressVerification, GovVendor, AMLScreening, AdverseMedia, PhoneNumber, Selfie, Liveness}

func ParseVerificationType(v string) (VerificationType, error) {
	if v == "" {
		return VerificationType{}, fmt.Errorf("parseverification :%s", v)
	}
	upper := strings.ToUpper(v)

	switch upper {
	case "DOCUMENT_INSIGHT":
		return DocumentInsight, nil
	case "ADDRESS_VERIFICATION":
		return AddressVerification, nil
	case "GOV_VERIFICATION":
		return GovVendor, nil
	case "AML_SCREENING":
		return AMLScreening, nil
	case "ADVERSE_MEDIA":
		return AdverseMedia, nil
	case "PHONENUMBER":
		return PhoneNumber, nil
	case "SELFIE":
		return Selfie, nil
	case "LIVENESS":
		return Liveness, nil
	default:
		return VerificationType{}, fmt.Errorf("parseverification invalid")
	}
}

func (v *VerificationType) String() string {
	return v.a
}
