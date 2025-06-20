package verificationbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	cb  *circuitbreaker.CircuitBreaker
	log *logger.Logger
}

type ServiceConfig func(*Service) error

func New(logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{
		log: logger,
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
	ctx, span := otel.AddSpan(ctx, "verficationbus.Create")
	defer span.End()

}
