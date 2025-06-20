package valueobject

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

type Phone struct {
	Country Country
	Digits  string
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

	return Phone{
		Country: country,
		Digits:  num.String(),
	}, nil
}
