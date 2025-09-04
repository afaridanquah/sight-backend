package organizationbus

import "context"

type Repository interface {
	Add(ctx context.Context, org Organization) error
}
