package businessbus

import "context"

type Repository interface {
	Add(ctx context.Context, bus Business) (Business, error)
}
