package valueobject

import "fmt"

type Business struct {
	Name   string
	Status string
}

func NewBusiness(name, status string) (Business, error) {
	if name == "" {
		return Business{}, fmt.Errorf("business name is required")
	}
	if status == "" {
		return Business{}, fmt.Errorf("business status is required")
	}

	return Business{
		Name:   name,
		Status: status,
	}, nil
}
