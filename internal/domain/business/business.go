package business

import (
	"errors"

	"bitbucket.org/msafaridanquah/verifylab-service/internal/domain/business/valueobject"
)

type Business struct {
	id     string
	name   string
	status valueobject.Status
}

var (
	ErrBusinessIDIsRequired   = errors.New("business id is required")
	ErrBusinessNameIsRequired = errors.New("business name is required")
)

func New(id string, name string, status valueobject.Status) (*Business, error) {
	if id == "" {
		return &Business{}, ErrBusinessIDIsRequired
	}
	if name == "" {
		return &Business{}, ErrBusinessNameIsRequired
	}

	if status == (valueobject.Status{}) {
		return &Business{}, ErrBusinessNameIsRequired
	}

	return &Business{
		id:     id,
		name:   name,
		status: status,
	}, nil
}

func (b Business) IsActive() bool {
	return b.status == valueobject.Active
}
