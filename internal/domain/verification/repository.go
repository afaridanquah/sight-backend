package verification

import (
	"context"
	"errors"

	"bitbucket.org/msafaridanquah/verifylab-service/internal/domain/verification/valueobject"
)

var (
	ErrVerificationCannotBeAdded = errors.New("cannot be added to repo")
	ErrVerificationCannotFound   = errors.New("cannot be found in repo")
)

type Repository interface {
	Add(ctx context.Context, v Verification) error
	Find(ctx context.Context, id valueobject.ID) (Verification, error)
}
