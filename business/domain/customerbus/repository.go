package customerbus

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCustomerNotFound       = errors.New("the customer was not found in repository")
	ErrFailedToAddCustomer    = errors.New("failed to add customer to respository")
	ErrFailedToUpdateCustomer = errors.New("failed to update customer in the respository")
	ErrFailedToDeleteCustomer = errors.New("failed to delete customer in the respository")
)

type Repository interface {
	QueryByCustomerAndBusinessID(ctx context.Context, id uuid.UUID, businessID uuid.UUID) (Customer, error)
	Add(ctx context.Context, c Customer) error
}
