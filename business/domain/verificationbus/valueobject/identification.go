package valueobject

import (
	"fmt"
	"strings"
)

type Identification struct {
	IdentificationType IdentificationType
	Pin                string
	CountryIssued      Country
}

func NewIdentification(idt, pin, c string) (Identification, error) {
	if idt == "" {
		return Identification{}, fmt.Errorf("new identification %q", idt)
	}
	if pin == "" {
		return Identification{}, fmt.Errorf("new identification %q", pin)
	}

	if c == "" {
		return Identification{}, fmt.Errorf("new identification %q", c)
	}

	idtype, err := ParseIdentificationType(idt)
	if err != nil {
		return Identification{}, fmt.Errorf("new identification %q:%w", idt, err)
	}

	country, err := NewCountry(c)
	if err != nil {
		return Identification{}, fmt.Errorf("new identification %q:%w", c, err)
	}

	return Identification{
		IdentificationType: idtype,
		CountryIssued:      country,
		Pin:                pin,
	}, nil
}

// =============================================
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

func (idt *IdentificationType) String() string {
	return idt.a
}
