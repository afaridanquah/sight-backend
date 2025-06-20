package valueobject

import (
	"errors"

	"github.com/biter777/countries"
)

type Country struct {
	alphaCode2 string
	name       string
}

var (
	ErrCountryCodeCannotBeEmpty = errors.New("country code cannot be empty")
	ErrCountryCodeNotValid      = errors.New("country code not valid, must be exactly two letters")
)

func NewCountry(c string) (Country, error) {
	if c == "" {
		return Country{}, ErrCountryCodeCannotBeEmpty
	}

	if len(c) != 2 {
		return Country{}, ErrCountryCodeNotValid
	}

	country := countries.ByName(c)
	if country.IsValid() {
		return Country{
			alphaCode2: country.Alpha2(),
			name:       country.Info().Name,
		}, nil
	}

	return Country{}, ErrCountryCodeNotValid
}

func (c Country) String() string {
	return c.alphaCode2
}

func (c Country) Name() string {
	return c.name
}

func (c Country) Alpha2() string {
	return c.alphaCode2
}
