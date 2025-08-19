package valueobject

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

type Phone struct {
	Country        Country
	E164Format     string
	NationalFormat string
	Carrier        string
	PhoneType      string
}

func ParsePhone(countryCode string, digits string) (Phone, error) {
	if countryCode == "" {
		return Phone{}, fmt.Errorf("country iso 2 is required")
	}

	if digits == "" {
		return Phone{}, fmt.Errorf("phone digits is required")
	}

	country, err := NewCountry(countryCode)
	if err != nil {
		return Phone{}, err
	}

	num, err := phonenumbers.Parse(digits, country.String())
	if err != nil {
		return Phone{}, err
	}

	e164 := phonenumbers.Format(num, phonenumbers.E164)
	national := phonenumbers.Format(num, phonenumbers.NATIONAL)
	carrier, err := phonenumbers.GetSafeCarrierDisplayNameForNumber(num, "en")
	if err != nil {
		fmt.Printf("GetSafeCarrierDisplayNameForNumber: %v", err)
		return Phone{}, err
	}

	return Phone{
		Country:        country,
		NationalFormat: national,
		E164Format:     e164,
		Carrier:        carrier,
	}, nil
}

func ParseIntlPhone(digits string) (Phone, error) {
	if digits == "" {
		return Phone{}, fmt.Errorf("phone digits is required")
	}

	num, err := phonenumbers.Parse(digits, "")
	if err != nil {
		return Phone{}, err
	}

	country, err := NewCountry(phonenumbers.GetRegionCodeForCountryCode(int(*num.CountryCode)))
	if err != nil {
		return Phone{}, err
	}

	e164 := phonenumbers.Format(num, phonenumbers.E164)
	national := phonenumbers.Format(num, phonenumbers.NATIONAL)
	carrier, err := phonenumbers.GetSafeCarrierDisplayNameForNumber(num, "en")
	if err != nil {
		fmt.Printf("GetSafeCarrierDisplayNameForNumber: %v", err)
		return Phone{}, err
	}

	return Phone{
		Country:        country,
		NationalFormat: national,
		E164Format:     e164,
		Carrier:        carrier,
	}, nil
}

func (p Phone) IsZero() bool {
	if p == (Phone{}) {
		return true
	}
	return false
}
