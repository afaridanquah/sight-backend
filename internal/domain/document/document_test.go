package document_test

import (
	"testing"

	"github.com/afaridanquah/verifylab-backend/internal/domain/document"
	"github.com/afaridanquah/verifylab-backend/internal/domain/document/valueobject"
)

func TestNewDocument(t *testing.T) {
	t.Parallel()
	dt, _ := valueobject.NewDocumentType("NationalIdentityCard")

	testCases := []struct {
		name         string
		id           valueobject.ID
		documentType valueobject.DocumentType
		expectedErr  error
	}{
		{
			name:         "Can Create New Document",
			id:           valueobject.NewID(),
			documentType: dt,
			expectedErr:  nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := document.New(tt.id, tt.documentType)
			if err != tt.expectedErr {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
			t.Logf("response: %v", doc)
		})
	}
}
