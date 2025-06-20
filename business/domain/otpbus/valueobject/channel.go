package valueobject

import (
	"fmt"
	"strings"
)

type Channel struct {
	a string
}

var (
	Sms          = Channel{"SMS"}
	EmailAddress = Channel{"EMAIL"}
)

func ParseChannel(s string) (Channel, error) {
	upper := strings.ToUpper(s)

	switch upper {
	case "SMS":
		return Sms, nil
	case "EMAIL":
		return EmailAddress, nil
	default:
		return Channel{}, fmt.Errorf("parse channel :%s is invalid", s)
	}
}

func MustParseChannel(s string) Channel {
	channel, err := ParseChannel(s)
	if err != nil {
		panic(err)
	}

	return channel
}

func (c Channel) String() string {
	return c.a
}
