package organizationapp

import (
	"fmt"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/app/sdk/errs"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus/valueobject"
)

type Organization struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type NewOrganization struct {
	Name string `json:"name" validate:"required"`
}

func (o NewOrganization) Validate() error {
	if err := errs.Check(o); err != nil {
		return fmt.Errorf("validate new organization failed: %w", err)
	}

	return nil
}

func toBusNewOrganization(napp NewOrganization) (organizationbus.NewOrganization, error) {
	status, err := valueobject.ParseStatus("active")
	if err != nil {
		return organizationbus.NewOrganization{}, err
	}

	return organizationbus.NewOrganization{
		Name:   napp.Name,
		Status: status,
	}, nil
}

func toAppOrganization(bus organizationbus.Organization) Organization {
	return Organization{
		ID:        bus.ID.String(),
		Name:      bus.Name,
		CreatedAt: bus.CreatedAt.Format(time.RFC3339),
		UpdatedAt: bus.UpdatedAt.Format(time.RFC3339),
	}
}
