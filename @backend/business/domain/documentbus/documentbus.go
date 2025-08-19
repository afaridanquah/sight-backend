package documentbus

import (
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus/valueobject"
	"github.com/google/uuid"
)

type Document struct {
	ID           uuid.UUID
	OriginalName string
	FilePath     string
	DocumentType valueobject.DocumentType
	UserID       uuid.UUID
	User         valueobject.User
	BusinessID   uuid.UUID
	Status       valueobject.Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type NewDocument struct {
	OriginalName string
	File         valueobject.File
	DocumentType valueobject.DocumentType
	UserID       uuid.UUID
	BusinessID   uuid.UUID
}
