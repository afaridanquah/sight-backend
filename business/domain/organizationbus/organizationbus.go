package organizationbus

import (
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus/valueobject"
)

type Organization struct {
	ID        valueobject.ID
	Name      string
	Status    valueobject.Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewOrganization struct {
	Name   string
	Status valueobject.Status
}
