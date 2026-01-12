package govbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/govbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	repo Repository
	log  *logger.Logger
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

func (srv *Service) Create(ctx context.Context, napp NewIdentification) (Identification, error) {
	ctx, span := otel.AddSpan(ctx, "identificationbus.service.create")
	defer span.End()

	now := time.Now()
	bus := Identification{
		ID:           valueobject.NewID(),
		CustomerID:   napp.CustomerID,
		FirstName:    napp.FirstName,
		LastName:     napp.LastName,
		MiddleName:   napp.MiddleName,
		Sex:          napp.Sex,
		Pin:          napp.Pin,
		PlaceOfBirth: napp.PlaceOfBirth,
		DateOfBirth:  napp.DateOfBirth,
		Nationality:  napp.Nationality,
		IssuedDate:   napp.IssuedDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// if err := srv.repo.Add(ctx, bus); err != nil {
	// 	return Identification{}, err
	// }

	return bus, nil
}
