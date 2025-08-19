package valueobject

import "fmt"

type Classification struct {
	a string
}


func ParseClassification(a string) (Classification, error) {
	if a == "" {
		return Classification{}, fmt.Errorf("cannot be empty")
	}
	uppercase :=
}
