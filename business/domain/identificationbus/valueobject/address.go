package valueobject

import "fmt"

type Address struct {
	Address1 string
	Address2 string
	ZipCode  string
	City     string
	State    string
	Country  Country
}

func ParseAddress(add1 string, add2 *string, zip string, city string, state *string, countryCode string) (Address, error) {
	if add1 == "" {
		return Address{}, fmt.Errorf("address 1 is required")
	}
	if zip == "" {
		return Address{}, fmt.Errorf("zip code is required")
	}

	if city == "" {
		return Address{}, fmt.Errorf("city is required")
	}

	if countryCode == "" {
		return Address{}, fmt.Errorf("country is required")
	}

	country, err := NewCountry(countryCode)
	if err != nil {
		return Address{}, err
	}

	return Address{
		Address1: add1,
		Address2: *add2,
		City:     city,
		State:    *state,
		ZipCode:  zip,
		Country:  country,
	}, nil
}
