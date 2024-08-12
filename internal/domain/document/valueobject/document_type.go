package valueobject

import (
	"errors"
	"strings"
)

type DocumentType struct {
	name string
}

var (
	ErrDocumentTypeNameRequired = errors.New("name is required")
	ErrDocumentTypeIsInvalid    = errors.New("document type does not exist")
)

var (
	Passport             = DocumentType{"Passport"}
	DriverLicense        = DocumentType{"DriverLicense"}
	UtilityBill          = DocumentType{"UtilityBill"}
	TaxDocument          = DocumentType{"TaxDocument"}
	NationalIdentityCard = DocumentType{"NationalIdentityCard"}
	Visa                 = DocumentType{"Visa"}
	PollingCard          = DocumentType{"PollingCard"}
	BirthCerificate      = DocumentType{"BirthCerificate"}
	Other                = DocumentType{"Other"}
)

var ListOfDocumentTypes = []DocumentType{
	Passport,
	DriverLicense,
	UtilityBill,
	TaxDocument,
	NationalIdentityCard,
	Visa,
	PollingCard,
	BirthCerificate,
	Other,
}

func NewDocumentType(s string) (DocumentType, error) {
	if s == "" {
		return DocumentType{}, ErrDocumentTypeNameRequired
	}

	for _, dt := range ListOfDocumentTypes {
		if strings.EqualFold(dt.Value(), s) {
			return dt, nil
		}
	}

	return DocumentType{}, ErrDocumentTypeIsInvalid
}

func (dt DocumentType) Value() string {
	return dt.name
}
