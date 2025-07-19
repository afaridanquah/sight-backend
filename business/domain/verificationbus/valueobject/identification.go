package valueobject

import (
	"fmt"
	"strings"
	"time"
)

type Identification struct {
	IdentificationType IdentificationType
	Pin                string
	PlaceOfBirth       string
	DateOfBirth        time.Time
	Nationality        Country
	IssuedDate         time.Time
	ExpDate            time.Time
	CountryIssued      Country
}

func NewIdentification(idt string, pin string, issuedCountry string, issueDate time.Time, nationality *string, expDate *string) (Identification, error) {
	if idt == "" {
		return Identification{}, fmt.Errorf("identification type %q", idt)
	}
	if pin == "" {
		return Identification{}, fmt.Errorf("pin %q", pin)
	}

	if issuedCountry == "" {
		return Identification{}, fmt.Errorf("issued country %q", issuedCountry)
	}

	idtype, err := ParseIdentificationType(idt)
	if err != nil {
		return Identification{}, fmt.Errorf("parse identification type %q:%w", idt, err)
	}

	parsedIssuedCountry, err := NewCountry(issuedCountry)
	if err != nil {
		return Identification{}, fmt.Errorf("parse country code %q:%w", issuedCountry, err)
	}

	identification := Identification{
		IdentificationType: idtype,
		CountryIssued:      parsedIssuedCountry,
		Pin:                pin,
	}

	return identification, nil
}

func (v *Identification) HasExpired() bool {
	if v.ExpDate.Before(time.Now()) {
		return true
	}
	return false
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
