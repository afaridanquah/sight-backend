package valueobject

import (
	"errors"
	"strings"
)

type DocumentType struct {
	a string
}

var (
	ErrDocumentTypeCannotBeEmpty = errors.New("status name cannot be empty")
)

var (
	PASSPORT               = DocumentType{"PASSPORT"}
	NATIONAL_IDENTITY_CARD = DocumentType{"NATIONAL_IDENTITY_CARD"}
	DRIVING_LICENSE        = DocumentType{"DRIVING_LICENSE"}
	RESIDENCE_PERMIT       = DocumentType{"RESIDENCE_PERMIT"}
	VISA                   = DocumentType{"VISA"}
	OTHER                  = DocumentType{"OTHER"}
)

var DocumentTypees = []DocumentType{PASSPORT, NATIONAL_IDENTITY_CARD}

func ParseDocumentType(name string) (DocumentType, error) {
	if name == "" {
		return DocumentType{}, ErrDocumentTypeCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "PASSPORT":
		return PASSPORT, nil
	case "NATIONAL_IDENTITY_CARD":
		return NATIONAL_IDENTITY_CARD, nil
	case "DRIVING_LICENSE":
		return DRIVING_LICENSE, nil
	case "RESIDENCE_PERMIT":
		return RESIDENCE_PERMIT, nil
	case "VISA":
		return VISA, nil
	case "OTHER":
		return OTHER, nil
	default:
		return DocumentType{}, errors.New("status name : invalid name")
	}
}

func MustParseDocumentType(name string) DocumentType {
	status, err := ParseDocumentType(name)
	if err != nil {
		panic(err)
	}
	return status
}

func (s DocumentType) String() string {
	return s.a
}
