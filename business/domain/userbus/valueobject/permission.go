package valueobject

import "fmt"

type Permission struct {
	Name  string
	Alias string
}

func NewPermission(name, alias string) (Permission, error) {
	if name == "" {
		return Permission{}, fmt.Errorf("name is required")
	}

	if alias == "" {
		return Permission{}, fmt.Errorf("alias is required")
	}

	return Permission{
		Name:  name,
		Alias: alias,
	}, nil
}
