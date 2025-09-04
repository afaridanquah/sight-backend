package valueobject

import "errors"

type Address struct {
	Line1         string
	Line2         string
	City          string
	StateOrRegion string
	Country       Country
}

var (
	ErrAddressLine1Required   = errors.New("line 1 is required")
	ErrAddressCityRequired    = errors.New("city is required")
	ErrAddressStateRequired   = errors.New("state or region is required")
	ErrAddressCountryRequired = errors.New("country is required")
)

func ParseAddress(line1 string, line2 *string, city string, state string, cc string) (Address, error) {
	if line1 == "" {
		return Address{}, ErrAddressLine1Required
	}
	if city == "" {
		return Address{}, ErrAddressCityRequired
	}
	if state == "" {
		return Address{}, ErrAddressStateRequired
	}

	if cc == "" {
		return Address{}, ErrAddressCountryRequired
	}

	country, err := NewCountry(cc)
	if err != nil {
		return Address{}, err
	}

	address := Address{
		Line1:         line1,
		City:          city,
		StateOrRegion: state,
		Country:       country,
	}

	if line2 != nil {
		address.Line2 = *line2
	}

	return address, nil
}

func (a *Address) IsEmpty() bool {
	return a == &Address{}
}
