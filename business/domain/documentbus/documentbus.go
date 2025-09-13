package documentbus

import (
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus/valueobject"
	"github.com/google/uuid"
)

type Document struct {
	ID             uuid.UUID
	Parent         uuid.UUID
	DocumentType   valueobject.DocumentType
	Side           valueobject.Side
	Classification valueobject.Classification
	CustomerID     uuid.UUID
	OriginalName   string
	FileName       string
	Customer       valueobject.User
	BusinessID     uuid.UUID
	OrgID          uuid.UUID
	Status         valueobject.Status
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type NewDocument struct {
	DocumentType   valueobject.DocumentType
	Classification valueobject.Classification
	CustomerID     uuid.UUID
	BusinessID     uuid.UUID
	Side           valueobject.Side
	File           valueobject.File
}
