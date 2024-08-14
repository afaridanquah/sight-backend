package customer

import (
	"context"
	"errors"

	"bitbucket.org/msafaridanquah/verifylab-service/internal/valueobject"
)

var (
	ErrCustomerNotFound       = errors.New("the customer was not found in repository")
	ErrFailedToAddCustomer    = errors.New("failed to add customer to respository")
	ErrFailedToUpdateCustomer = errors.New("failed to update customer in the respository")
	ErrFailedToDeleteCustomer = errors.New("failed to delete customer in the respository")
)

type Repository interface {
	Find(ctx context.Context, id valueobject.ID) (Customer, error)
	Add(ctx context.Context, c Customer) error
	Update(ctx context.Context, c Customer) error
	Delete(ctx context.Context, id valueobject.ID) error
}
