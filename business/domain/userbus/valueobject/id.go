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

func NewUserID() ID {
	return ID{
		prefix: "usr_",
		value:  ksuid.New(),
	}
}

func (id ID) String() string {
	return strings.Join([]string{id.prefix, id.value.String()}, "_")
}

func ParseUserID(s string) (ID, error) {
	if s == "" {
		return ID{}, fmt.Errorf("arg cannot be empty")
	}

	if !strings.HasPrefix(s, "usr_") {
		return ID{}, fmt.Errorf("%s is not a valid prefix", s)
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
