package documentbus

import "context"

type Repository interface {
	Add(ctx context.Context, bus Document) error
}
