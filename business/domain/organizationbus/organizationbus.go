package organizationbus

import (
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus/valueobject"
	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID
	Name      string
	Status    valueobject.Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewOrganization struct {
	Name   string
	Status valueobject.Status
}
