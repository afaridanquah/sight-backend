package valueobject

import (
	"errors"
	"net/url"
)

type Domain struct {
	a string
}

var (
	ErrUrlIsRequired = errors.New("url is required")
)

func NewDomain(a string) (Domain, error) {
	if a == "" {
		return Domain{}, ErrUrlIsRequired
	}

	url, err := url.Parse(a)
	if err != nil {
		return Domain{}, err
	}

	return Domain{
		a: url.String(),
	}, nil
}
