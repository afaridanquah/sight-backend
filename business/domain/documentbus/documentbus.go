package documentbus

import (
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus/valueobject"
)

type Document struct {
	ID             valueobject.ID
	Parent         valueobject.ID
	DocumentType   valueobject.DocumentType
	Side           valueobject.Side
	Classification valueobject.Classification
	CustomerID     string
	OriginalName   string
	FileName       string
	Customer       valueobject.User
	BusinessID     string
	OrgID          string
	Status         valueobject.Status
	MetaData       any
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type NewDocument struct {
	DocumentType   valueobject.DocumentType
	Classification valueobject.Classification
	CustomerID     string
	BusinessID     string
	Side           valueobject.Side
	File           valueobject.File
}

type UpdateDocumentStatus struct {
	ID     valueobject.ID
	Status valueobject.Status
}
