package businessbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"github.com/google/uuid"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	repo Repository
	cb   *circuitbreaker.CircuitBreaker
}

type ServiceConfig func(*Service) error

func New(repo Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{
		repo: repo,
	}

	for _, cfg := range cfgs {
		err := cfg(ser)
		if err != nil {
			return nil, err
		}
	}

	ser.cb = circuitbreaker.New(
		circuitbreaker.WithOpenTimeout(time.Minute*2),
		circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncConsecutiveFailures(3)),
		circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
			logger.Info(context.Background(), "state changed",
				slog.String("old", string(oldState)),
				slog.String("new", string(newState)),
			)
		}),
	)

	return ser, nil
}

func (ser *Service) Create(ctx context.Context, nbus NewBusiness) (Business, error) {
	ctx, span := otel.AddSpan(ctx, "businessbus.service.create")
	defer span.End()
	now := time.Now()

	orgID := uuid.MustParse("6fe9cace-7c71-4e4b-b943-dd2f5bb21693")

	bus := Business{
		ID:              uuid.New(),
		LegalName:       nbus.LegalName,
		Entity:          nbus.Entity,
		DoingBusinessAs: nbus.DoingBusinessAs,
		EmailAddresses:  nbus.EmailAddresses,
		PhoneNumbers:    nbus.PhoneNumbers,
		Website:         nbus.Website,
		OrgID:           orgID,
		Address:         nbus.Address,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := ser.repo.Add(ctx, bus); err != nil {
		return Business{}, err
	}

	return bus, nil
}

// Update modifies information about a business.
func (ser *Service) Update(ctx context.Context, bus Business, up UpdateBusiness) (Business, error) {
	ctx, span := otel.AddSpan(ctx, "businessbus.service.update")
	defer span.End()

	if up.DoingBusinessAs != nil {
		bus.DoingBusinessAs = *up.DoingBusinessAs
	}

	if up.LegalName != nil {
		bus.LegalName = *up.LegalName
	}

	if up.EmailAddresses != nil {
		bus.EmailAddresses = *up.EmailAddresses
	}

	if up.Website != nil {
		bus.Website = *up.Website
	}

	if up.LegalName != nil {
		bus.LegalName = *up.LegalName
	}

	if up.PhoneNumbers != nil {
		bus.PhoneNumbers = *up.PhoneNumbers
	}

	if up.Entity != nil {
		bus.Entity = *up.Entity
	}

	if err := ser.repo.Update(ctx, bus); err != nil {
		return Business{}, err
	}

	return bus, nil
}

func (ser *Service) FindByID(ctx context.Context, id uuid.UUID) (Business, error) {
	ctx, span := otel.AddSpan(ctx, "businessbus.service.findbyid")
	defer span.End()

	//Get user id and org id from mid
	orgID := uuid.MustParse("6fe9cace-7c71-4e4b-b943-dd2f5bb21693")

	bus, err := ser.repo.QueryByIDAndOrgID(ctx, id, orgID)
	if err != nil {
		return Business{}, err
	}

	return bus, nil
}
