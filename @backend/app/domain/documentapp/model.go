package documentapp

import (
	"fmt"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
)

type Document struct {
	ID       string `json:"id"`
	FileName string `json:"filename"`
}

type NewDocument struct {
	Type         string `json:"type"`
	FileName     string `json:"file_name" validate:"required"`
	DocumentType string `json:"document_type" validate:"required"`
}

func (o NewDocument) Validate() error {
	if err := ierr.Check(o); err != nil {
		return fmt.Errorf("validate new document failed: %w", err)
	}

	return nil
}

func toAppDocument(bus documentbus.Document) Document {
	return Document{
		ID:       bus.ID.String(),
		FileName: bus.OriginalName,
	}
}

func toBusNewDocument(napp NewDocument) (documentbus.NewDocument, error) {
	dt, err := valueobject.ParseDocumentType(napp.DocumentType)
	if err != nil {
		return documentbus.NewDocument{}, err
	}

	return documentbus.NewDocument{
		DocumentType: dt,
	}, nil

}
