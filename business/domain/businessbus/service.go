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
	businesses Repository
	cb         *circuitbreaker.CircuitBreaker
}

type ServiceConfig func(*Service) error

func New(businesses Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{
		businesses: businesses,
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

	bus := Business{
		ID:        uuid.New(),
		Name:      nbus.Name,
		Status:    nbus.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := ser.businesses.Add(ctx, bus)
	if err != nil {
		return Business{}, err
	}

	return bus, nil
}
