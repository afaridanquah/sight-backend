package valueobject

import (
	"errors"
	"strings"
)

type DocumentSide struct {
	name string
}

var (
	Front = DocumentSide{"Front"}
	Back  = DocumentSide{"Back"}
)

var (
	ErrDocumentSideNameRequired = errors.New("side name required")
	ErrDocumentSideIsInvalid    = errors.New("side does not exist")
)

var ListOfDocumentSides = []DocumentSide{
	Front,
	Back,
}

func NewDocumentSide(s string) (DocumentSide, error) {
	if s == "" {
		return DocumentSide{}, ErrDocumentSideNameRequired
	}

	for _, dt := range ListOfDocumentSides {
		if strings.EqualFold(dt.Value(), s) {
			return dt, nil
		}
	}

	return DocumentSide{}, ErrDocumentSideIsInvalid
}

func (ds DocumentSide) Value() string {
	return ds.name
}
