package businessapp

import (
	"fmt"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus/valueobject"
	"github.com/google/uuid"
)

type Business struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type NewBusiness struct {
	ID      string `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	OwnerID string `json:"owner_id" validate:"required"`
	Status  string `json:"status" validate:"required"`
}

func toBusBusiness(napp NewBusiness) (businessbus.NewBusiness, error) {
	id, err := uuid.Parse(napp.ID)
	if err != nil {
		return businessbus.NewBusiness{}, fmt.Errorf("parse id %w", err)
	}

	ownerID, err := uuid.Parse(napp.OwnerID)
	if err != nil {
		return businessbus.NewBusiness{}, fmt.Errorf("parse owner id %w", err)
	}

	status, err := valueobject.ParseStatus(napp.Status)
	if err != nil {
		return businessbus.NewBusiness{}, fmt.Errorf("parse status %w", err)
	}

	return businessbus.NewBusiness{
		ID:      id,
		Name:    napp.Name,
		OwnerID: ownerID,
		Status:  status,
	}, nil
}

func toAppBusiness(bus businessbus.Business) Business {
	return Business{
		ID:        bus.ID.String(),
		Name:      bus.Name,
		OwnerID:   bus.OwnerID.String(),
		CreatedAt: bus.CreatedAt.Format(time.RFC3339),
		UpdatedAt: bus.UpdatedAt.Format(time.RFC3339),
	}
}
