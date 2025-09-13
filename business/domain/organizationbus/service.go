package organizationbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"github.com/google/uuid"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	repo   Repository
	logger *logger.Logger
	cb     *circuitbreaker.CircuitBreaker
}

type ServiceConfig func(*Service) error

func New(repo Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var srv = &Service{
		logger: logger,
		repo:   repo,
	}

	for _, cfg := range cfgs {
		err := cfg(srv)
		if err != nil {
			return nil, err
		}
	}

	srv.cb = circuitbreaker.New(
		circuitbreaker.WithOpenTimeout(time.Minute*2),
		circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncConsecutiveFailures(3)),
		circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
			logger.Info(context.Background(), "state changed",
				slog.String("old", string(oldState)),
				slog.String("new", string(newState)),
			)
		}),
	)

	return srv, nil
}

func (srv *Service) Create(ctx context.Context, nbus NewOrganization) (Organization, error) {
	ctx, span := otel.AddSpan(ctx, "documentbus.service.create")
	defer span.End()
	now := time.Now()

	org := Organization{
		ID:        uuid.New(),
		Name:      nbus.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := srv.repo.Add(ctx, org); err != nil {
		return Organization{}, nil
	}

	return org, nil
}
