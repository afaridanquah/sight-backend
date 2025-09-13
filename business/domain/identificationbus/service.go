package identificationbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
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

}
