package identificationbus

import "context"

type Repository interface {
	Add(ctx context.Context, identification Identification) error
}
