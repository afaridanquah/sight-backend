package identificationbus

import (
	"context"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus/valueobject"
)

type Repository interface {
	Add(ctx context.Context, identification Identification) error
	AddMany(ctx context.Context, identifications []Identification) error
	QueryByCustomerID(ctx context.Context, id valueobject.ID) ([]Identification, error)
}
