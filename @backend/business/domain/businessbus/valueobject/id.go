package valueobject

import (
	"errors"
	"strings"

	"github.com/segmentio/ksuid"
)

var (
	ErrCustomerIDRequired = errors.New("id cannot be empty")
)

type ID struct {
	prefix string
	value  ksuid.KSUID
}

func NewID() ID {
	return ID{
		prefix: "bus",
		value:  ksuid.New(),
	}
}

func (id ID) String() string {
	return strings.Join([]string{id.prefix, id.value.String()}, "_")
}

func ParseID(s string) (ID, error) {
	if s == "" {
		return ID{}, ErrCustomerIDRequired
	}
	d := strings.Split(s, "_")

	k, err := ksuid.Parse(d[1])
	if err != nil {
		return ID{}, err
	}

	return ID{
		prefix: d[0],
		value:  k,
	}, nil
}
