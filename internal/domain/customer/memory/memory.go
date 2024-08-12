package memory

import (
	"context"
	"sync"

	"github.com/afaridanquah/verifylab-backend/internal/domain/customer"
	"github.com/afaridanquah/verifylab-backend/internal/domain/customer/valueobject"
)

type MemoryRepository struct {
	customers map[string]customer.Customer
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		customers: make(map[string]customer.Customer),
	}
}

func (mr *MemoryRepository) Find(ctx context.Context, id valueobject.ID) (customer.Customer, error) {
	if client, ok := mr.customers[id.String()]; ok {
		return client, nil
	}
	return customer.Customer{}, customer.ErrCustomerNotFound
}

func (mr *MemoryRepository) Add(ctx context.Context, newCustomer customer.Customer) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.customers[newCustomer.ID().String()]; ok {
		return customer.ErrFailedToAddCustomer
	}
	mr.customers[newCustomer.ID().String()] = newCustomer

	return nil
}

func (mr *MemoryRepository) Update(ctx context.Context, updateCustomer customer.Customer) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.customers[updateCustomer.ID().String()]; !ok {
		return customer.ErrFailedToUpdateCustomer
	}
	mr.customers[updateCustomer.ID().String()] = updateCustomer

	return nil
}

func (mr *MemoryRepository) Delete(ctx context.Context, id valueobject.ID) error {
	mr.Lock()
	defer mr.Unlock()

	if _, ok := mr.customers[id.String()]; !ok {
		return customer.ErrFailedToDeleteCustomer
	}

	delete(mr.customers, id.String())

	return nil
}
