package customerbus

import (
	"context"
	"errors"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/valueobject"
)

var (
	ErrCustomerNotFound       = errors.New("the customer was not found in repository")
	ErrFailedToAddCustomer    = errors.New("failed to add customer to respository")
	ErrFailedToUpdateCustomer = errors.New("failed to update customer in the respository")
	ErrFailedToDeleteCustomer = errors.New("failed to delete customer in the respository")
)

type Repository interface {
	QueryByIDAndOrgID(ctx context.Context, id valueobject.ID, orgID valueobject.ID) (Customer, error)
	Add(ctx context.Context, c Customer) error
	Update(ctx context.Context, cust Customer) error
	QueryByOrgID(ctx context.Context, orgID valueobject.ID) ([]Customer, error)
}

type SearchRepository interface {
	Search(ctx context.Context, sc SearchCustomer) ([]Customer, error)
}

type MessageBrokerRepository interface {
	Created(ctx context.Context, c Customer) error
	// Updated(ctx context.Context, c Customer) error
}
