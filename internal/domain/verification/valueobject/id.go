package valueobject

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

type ID string

func NewID() ID {
	data := ulid.Make()
	d, _ := fmt.Printf("ver_%v", data)

	return ID(d)
}

func (id ID) String() string {
	return string(id)
}
