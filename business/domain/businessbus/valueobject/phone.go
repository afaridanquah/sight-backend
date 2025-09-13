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
	PhoneType      any
	Location       string
	IsValidNumber  bool
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
	location, err := phonenumbers.GetGeocodingForNumber(num, "en")
	isValid := phonenumbers.IsValidNumber(num)

	if err != nil {
		return Phone{}, err
	}

	return Phone{
		Country:        country,
		NationalFormat: national,
		E164Format:     e164,
		Location:       location,
		IsValidNumber:  isValid,
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
	location, err := phonenumbers.GetGeocodingForNumber(num, "en")
	isValid := phonenumbers.IsValidNumber(num)

	if err != nil {
		return Phone{}, err
	}

	return Phone{
		Country:        country,
		NationalFormat: national,
		E164Format:     e164,
		Location:       location,
		IsValidNumber:  isValid,
	}, nil
}

func ParseIntlPhoneNumbers(phones []string) ([]Phone, error) {
	if len(phones) > 0 {
		parsed := make([]Phone, len(phones))
		for i, v := range phones {
			e, err := ParseIntlPhone(v)
			if err != nil {
				return []Phone{}, err
			}
			parsed[i] = e
		}

		return parsed, nil
	}

	return []Phone{}, fmt.Errorf("no phone number provided")
}

func (p Phone) IsZero() bool {
	return p == (Phone{})
}

func (p *Phone) IsEmpty() bool {
	return p == nil
}
