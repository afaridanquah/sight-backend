package identificationbus

import "context"

type Repository interface {
	Add(ctx context.Context, idv Identification) (Identification, error)
}
