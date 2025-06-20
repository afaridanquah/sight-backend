package businessbus

import (
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus/valueobject"
	"github.com/google/uuid"
)

type Business struct {
	ID        uuid.UUID
	Name      string
	OwnerID   uuid.UUID
	Status    valueobject.Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewBusiness struct {
	ID      uuid.UUID
	Name    string
	OwnerID uuid.UUID
	Status  valueobject.Status
}
