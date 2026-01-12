package govbus

import "context"

type Repository interface {
	Add(ctx context.Context, bus Identification) error
}
