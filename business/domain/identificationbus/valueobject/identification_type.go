package valueobject

import (
	"fmt"
	"strings"
)

type IdentificationType struct {
	a string
}

var (
	PASSPORT       = IdentificationType{"PASSPORT"}
	DRIVERSLICENSE = IdentificationType{"DRIVERS_LICENSE"}
	NATIONALID     = IdentificationType{"NATIONAL_ID"}
	RESIDENTPERMIT = IdentificationType{"RESIDENT_PERMIT"}
	SSN            = IdentificationType{"SSN"}
)

func ParseIdentificationType(s string) (IdentificationType, error) {
	lower := strings.ToUpper(s)
	switch lower {
	case "PASSPORT":
		return PASSPORT, nil
	case "DRIVERS_LICENSE":
		return DRIVERSLICENSE, nil
	case "NATIONAL_ID":
		return NATIONALID, nil
	case "RESIDENT_PERMIT":
		return RESIDENTPERMIT, nil
	case "SSN":
		return SSN, nil
	default:
		return IdentificationType{}, fmt.Errorf("%s : is not valid", s)
	}
}
