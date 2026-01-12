package identificationbus

import "context"

type Repository interface {
	Add(ctx context.Context, identification Identification) error
	AddMany(ctx context.Context, identifications []Identification) error
	QueryByCustomerID(ctx context.Context, id string) ([]Identification, error)
}
