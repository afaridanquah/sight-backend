package identificationbus

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
	identifications Repository
	cb              *circuitbreaker.CircuitBreaker
}

type ServiceConfig func(*Service) error

func New(identifications Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{
		identifications: identifications,
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

func (ser *Service) Create(ctx context.Context, nbus NewIdentification) (Identification, error) {
	ctx, span := otel.AddSpan(ctx, "identificationbus.service.Create")
	defer span.End()

	now := time.Now()
	bus := Identification{
		ID:                 uuid.New(),
		Number:             nbus.Number,
		PlaceOfBirth:       nbus.PlaceOfBirth,
		DateOfBirth:        nbus.DateOfBirth,
		IdentificationType: nbus.IdentificationType,
		IssuedCountry:      nbus.IssuedCountry,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return bus, nil
}
