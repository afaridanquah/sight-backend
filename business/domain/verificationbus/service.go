package verificationbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/sdk/yenti"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"github.com/google/uuid"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	cb     *circuitbreaker.CircuitBreaker
	log    *logger.Logger
	envvar *envvar.Configuration
	repo   Repository
}

type ServiceConfig func(*Service) error

func New(repo Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{
		log:  logger,
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

func (srv *Service) Create(ctx context.Context, newbus NewVerification) (Verification, error) {
	ctx, span := otel.AddSpan(ctx, "verficationbus.service.create")
	defer span.End()

	businessID := uuid.MustParse("6fe9cace-7c71-4e4b-b943-dd2f5bb21693")

	now := time.Now()
	verification := Verification{
		ID:               uuid.New(),
		VerificationType: newbus.VerificationType,
		CustomerID:       newbus.CustomerID,
		Customer:         newbus.Customer,
		BusinessID:       businessID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	srv.log.Info(ctx, "bus create verification", verification.Customer)

	yentiClient, err := yenti.New(yenti.Config{
		Env: srv.envvar,
		Log: srv.log,
	})

	if err != nil {
		return Verification{}, err
	}

	if err := verification.OpenSanctionMatch(yentiClient); err != nil {
		return Verification{}, err
	}

	if err := srv.repo.Add(ctx, verification); err != nil {
		return Verification{}, err
	}

	return verification, nil
}
