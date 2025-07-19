package verificationbus

import "context"

type Repository interface {
	Add(ctx context.Context, ver Verification) error
}
