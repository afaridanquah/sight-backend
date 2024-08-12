package memory

import (
	"context"
	"sync"

	"github.com/afaridanquah/verifylab-backend/internal/domain/verification"
	"github.com/afaridanquah/verifylab-backend/internal/domain/verification/valueobject"
)

type Respository struct {
	verifications map[string]verification.Verification
	sync.Mutex
}

func New() (*Respository, error) {
	return &Respository{
		verifications: make(map[string]verification.Verification),
	}, nil
}

func (mr *Respository) Add(ctx context.Context, v verification.Verification) error {
	mr.Lock()
	defer mr.Unlock()

	if _, ok := mr.verifications[v.ID().String()]; ok {
		return verification.ErrVerificationCannotBeAdded
	}

	mr.verifications[string(v.ID())] = v

	return nil
}

func (mr *Respository) Find(ctx context.Context, id valueobject.ID) (verification.Verification, error) {
	mr.Lock()
	defer mr.Unlock()

	if ver, ok := mr.verifications[id.String()]; ok {
		return ver, nil
	}

	return verification.Verification{}, verification.ErrVerificationCannotFound

}
