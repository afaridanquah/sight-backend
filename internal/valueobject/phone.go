package valueobject

import "fmt"

type Phone struct {
	CountryPrefix string
	AreadCode     string
	Number        string
}

func (p Phone) FullNumber() string {
	return fmt.Sprintf("%s %s %s", p.CountryPrefix, p.AreadCode, p.Number)
}
