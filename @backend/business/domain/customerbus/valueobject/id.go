package valueobject

import (
	"errors"
	"fmt"
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
		prefix: "cus",
		value:  ksuid.New(),
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

	k, err := ksuid.Parse(d[1])
	if err != nil {
		return ID{}, err
	}

	return ID{
		prefix: d[0],
		value:  k,
	}, nil
}
