package valueobject

import (
	"fmt"

	"github.com/google/uuid"
)

type Business struct {
	ID                 uuid.UUID
	Name               string
	Country            Country
	RegistrationNumber string
}

func NewBusiness(id uuid.UUID, name string, c *string, regNumber *string) (Business, error) {
	if name == "" {
		return Business{}, fmt.Errorf("business name is required")
	}

	bus := Business{
		ID:   id,
		Name: name,
	}

	if c != nil {
		country, err := NewCountry(*c)
		if err != nil {
			return Business{}, err
		}
		bus.Country = country
	}

	if regNumber != nil {
		bus.RegistrationNumber = *regNumber
	}

	return bus, nil

}
