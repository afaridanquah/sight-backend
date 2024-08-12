package entity

import (
	"github.com/google/uuid"
)

type Person struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
}

// func (p Person) Validate() error {
// 	return validation.ValidateStruct(&p,
// 		validation.Field(&p.FirstName, validation.Required),
// 		validation.Field(&p.LastName, validation.Required),
// 	)
// }
