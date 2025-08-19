package valueobject

import (
	"fmt"
)

type Owner struct {
	// Person holds the legal name of the owner
	Person Person
	// A percentage that indicate the shares this person holds
	OwnershipPercentage float32
	HomeAddress         Address
	CountryOfResident   Country
}

func NewOwner(p Person, percentage float32, addr Address, countryCode string) (Owner, error) {
	if p == (Person{}) {
		return Owner{}, fmt.Errorf("person is required")
	}

	if percentage == 0.0 {
		return Owner{}, fmt.Errorf("percentage is required")
	}

	if addr == (Address{}) {
		return Owner{}, fmt.Errorf("address is required")
	}

	if countryCode == "" {
		return Owner{}, fmt.Errorf("address is required")
	}

	country, err := NewCountry(countryCode)
	if err != nil {
		return Owner{}, fmt.Errorf("country is required")
	}

	return Owner{
		Person:              p,
		OwnershipPercentage: percentage,
		HomeAddress:         addr,
		CountryOfResident:   country,
	}, nil
}
