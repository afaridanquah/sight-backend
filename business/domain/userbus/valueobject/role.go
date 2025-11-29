package valueobject

import "fmt"

type Role struct {
	Name  string
	Alias string
}

func NewRole(name, alias string) (Role, error) {
	if name == "" {
		return Role{}, fmt.Errorf("name is required")
	}

	if alias == "" {
		return Role{}, fmt.Errorf("alias is required")
	}

	return Role{
		Name:  name,
		Alias: alias,
	}, nil
}
