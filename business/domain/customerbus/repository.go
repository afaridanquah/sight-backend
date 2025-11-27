package customerbus

import (
	"context"
	"errors"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/valueobject"
	"github.com/google/uuid"
)

var (
	ErrCustomerNotFound       = errors.New("the customer was not found in repository")
	ErrFailedToAddCustomer    = errors.New("failed to add customer to respository")
	ErrFailedToUpdateCustomer = errors.New("failed to update customer in the respository")
	ErrFailedToDeleteCustomer = errors.New("failed to delete customer in the respository")
)

type Repository interface {
	QueryByIDAndOrgID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (Customer, error)
	Add(ctx context.Context, c Customer) error
	Update(ctx context.Context, cust Customer) error
}

type SearchRepository interface {
	Search(ctx context.Context, sc SearchCustomer) ([]Customer, error)
}

type CustomerMessageBrokerPublisher interface {
	Created(ctx context.Context, c Customer) error
	Updated(ctx context.Context, id valueobject.ID) error
}
