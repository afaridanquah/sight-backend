package addressbus

import "context"

type Repository interface {
	Add(ctx context.Context, bus Address) error
}
