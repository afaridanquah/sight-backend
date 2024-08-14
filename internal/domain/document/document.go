package document

import (
	"errors"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/internal/domain/document/valueobject"
	ivo "bitbucket.org/msafaridanquah/verifylab-service/internal/valueobject"
)

var (
	ErrDocumentIdIsRequired           = errors.New("id is required")
	ErrDocumentDocumentTypeIsRequired = errors.New("document type is required")
	ErrDocumentTimestampRequired      = errors.New("document timestamp is required")
	ErrDocumentCountryIsRequired      = errors.New("document issuing country is required")
)

type DocumentOptions func(*Document) error

type Document struct {
	id             valueobject.ID
	documentType   valueobject.DocumentType
	issuingCountry ivo.Country
	images         []valueobject.Image
	createdAt      time.Time
	updatedAt      time.Time
}

func New(id valueobject.ID, dt valueobject.DocumentType, opts ...DocumentOptions) (*Document, error) {
	if id == (valueobject.ID{}) {
		return &Document{}, ErrDocumentIdIsRequired
	}

	if dt == (valueobject.DocumentType{}) {
		return &Document{}, ErrDocumentIdIsRequired
	}

	doc := &Document{
		id:           id,
		documentType: dt,
	}

	for _, opt := range opts {
		err := opt(doc)
		if err != nil {
			return &Document{}, err
		}
	}

	return doc, nil
}

func (doc *Document) WithTimestamps(c, u time.Time) DocumentOptions {
	return func(d *Document) error {
		if c != (time.Time{}) {
			return ErrDocumentTimestampRequired
		}
		if u != (time.Time{}) {
			return ErrDocumentTimestampRequired
		}
		d.createdAt = c
		d.updatedAt = u
		return nil
	}
}

func (doc *Document) WithIssuingCountry(c ivo.Country) DocumentOptions {
	return func(d *Document) error {
		if c != (ivo.Country{}) {
			return ErrDocumentCountryIsRequired
		}
		d.issuingCountry = c
		return nil
	}
}

func (doc *Document) WithImages(imgs []valueobject.Image) DocumentOptions {
	return func(d *Document) error {
		d.images = imgs
		return nil
	}
}
