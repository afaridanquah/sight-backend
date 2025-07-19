package valueobject

import (
	"fmt"
	"strings"
	"time"
)

type Identification struct {
	IdentificationType IdentificationType
	Pin                string
	Nationality        *Country
	ExpDate            *time.Time
	IssuedDate         time.Time
	IssedCountry       Country
}

func NewIdentification(idt string, pin string, issuedCountry string, issuedDate *string, nationality *string, expDate *string) (Identification, error) {
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
		IssedCountry:       parsedIssuedCountry,
		Pin:                pin,
	}

	if issuedDate != nil {
		toIssuedDate, err := time.Parse(time.DateOnly, *issuedDate)
		if err != nil {
			return Identification{}, fmt.Errorf("parse issued date:%s, %w", *issuedDate, err)
		}

		identification.IssuedDate = toIssuedDate
	}

	if nationality != nil {
		parsedNationality, err := NewCountry(*nationality)
		if err != nil {
			return Identification{}, fmt.Errorf("parse issued date:%s, %w", *issuedDate, err)
		}

		identification.Nationality = &parsedNationality
	}

	if expDate != nil {
		toExpDate, err := time.Parse(time.DateOnly, *expDate)
		if err != nil {
			return Identification{}, fmt.Errorf("parse exp date %w", err)
		}

		identification.ExpDate = &toExpDate
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
