package verification

import (
	"errors"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/internal/domain/verification/valueobject"
	ivo "bitbucket.org/msafaridanquah/verifylab-service/internal/valueobject"
)

var (
	ErrVerificationIdIsRequired   = errors.New("verification id cannot be empty")
	ErrVerificationTypeIsRequired = errors.New("verification type cannot be empty")
	ErrPersonIsRequired           = errors.New("person cannot be empty")
)

type Verification struct {
	id               ivo.ID
	referenceId      string
	status           valueobject.Status
	person           valueobject.Person
	verificationType valueobject.VerificationType
	createdAt        time.Time
	updatedAt        time.Time
}

type VerificationOption func(*Verification)

func New(id ivo.ID, vt valueobject.VerificationType, p valueobject.Person, opts ...VerificationOption) (*Verification, error) {
	if id == (ivo.ID{}) {
		return &Verification{}, ErrVerificationIdIsRequired
	}

	if vt == (valueobject.VerificationType{}) {
		return &Verification{}, ErrVerificationTypeIsRequired
	}

	if p == (valueobject.Person{}) {
		return &Verification{}, ErrPersonIsRequired
	}

	var ver = &Verification{
		id:               id,
		verificationType: vt,
		person:           p,
	}

	for _, opt := range opts {
		opt(ver)
	}

	return ver, nil
}

func (ver *Verification) WithTimestamp() VerificationOption {
	return func(v *Verification) {
		v.createdAt = time.Now()
		v.updatedAt = time.Now()
	}
}

func (ver *Verification) WithReferenceID(id string) VerificationOption {
	return func(v *Verification) {
		v.referenceId = id
	}
}

func (ver *Verification) WithStatus(s valueobject.Status) VerificationOption {
	return func(v *Verification) {
		ver.status = s
	}
}

func (ver Verification) ID() ivo.ID {
	return ver.id
}

func (ver Verification) StringID() string {
	return string(ver.id.String())
}
