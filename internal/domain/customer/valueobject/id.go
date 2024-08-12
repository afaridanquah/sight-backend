package valueobject

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
)

var (
	ErrCustomerIDRequired = errors.New("id cannot be empty")
)

type ID struct {
	prefix string
	value  ulid.ULID
}

func NewID() ID {
	return ID{
		prefix: "cus",
		value:  ulid.Make(),
	}
}

func (id ID) String() string {
	return fmt.Sprintf("%s_%s", id.prefix, id.value.String())
}

func ParseID(s string) (ID, error) {
	if s == "" {
		return ID{}, ErrCustomerIDRequired
	}
	d := strings.Split(s, "_")

	return ID{
		prefix: d[0],
		value:  ulid.MustParse(d[1]),
	}, nil
}
