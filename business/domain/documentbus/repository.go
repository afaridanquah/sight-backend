package documentbus

import (
	"context"
)

type Repository interface {
	Add(ctx context.Context, bus Document) error
	Update(ctx context.Context, bus Document) error
}
