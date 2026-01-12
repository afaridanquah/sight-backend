package organizationbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
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
	ctx, span := otel.AddSpan(ctx, "oganizationbus.service.create")
	defer span.End()
	now := time.Now()

	status, err := valueobject.ParseStatus("pending")
	if err != nil {
		return Organization{}, nil
	}

	org := Organization{
		ID:        valueobject.NewOrgID(),
		Name:      nbus.Name,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := srv.repo.Add(ctx, org); err != nil {
		return Organization{}, nil
	}

	return org, nil
}
