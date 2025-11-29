package valueobject

import (
	"fmt"

	"github.com/google/uuid"
)

type Tenant struct {
	ID     uuid.UUID
	Name   string
	Alias  string
	Domain Domain
}

func NewTenant(id string, name string, alias string, domain string) (Tenant, error) {
	if id == "" {
		return Tenant{}, fmt.Errorf("id cannot be empty")
	}

	if name == "" {
		return Tenant{}, fmt.Errorf("name cannot be empty")
	}

	if alias == "" {
		return Tenant{}, fmt.Errorf("alias cannot be empty")
	}

	if domain == "" {
		return Tenant{}, fmt.Errorf("alias cannot be empty")
	}

	dn, err := NewDomain(domain)
	if err != nil {
		return Tenant{}, err
	}

	tid, err := uuid.Parse(id)
	if err != nil {
		return Tenant{}, err
	}

	return Tenant{
		ID:     tid,
		Name:   name,
		Alias:  alias,
		Domain: dn,
	}, nil
}
