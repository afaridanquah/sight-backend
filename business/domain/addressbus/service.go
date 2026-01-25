package addressbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/addressbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	repo Repository
	cb   *circuitbreaker.CircuitBreaker
	pool *pgxpool.Pool
}

type ServiceConfig func(*Service) error

func New(repo Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var srv = &Service{
		repo: repo,
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

func (srv *Service) Create(ctx context.Context, na NewAddress) (Address, error) {
	ctx, span := otel.AddSpan(ctx, "addressbus.service.create")
	defer span.End()
	now := time.Now()

	addr := Address{
		ID:        valueobject.NewID(),
		EntityID:  na.EntityID,
		Line1:     na.Line1,
		Line2:     na.Line2,
		State:     na.State,
		City:      na.City,
		Country:   na.Country,
		IsPrimary: na.IsPrimary,
		Zip:       na.Zip,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return addr, nil
}
