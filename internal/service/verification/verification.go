package verification

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/afaridanquah/verifylab-backend/internal"
	"github.com/afaridanquah/verifylab-backend/internal/domain/verification"
	"github.com/afaridanquah/verifylab-backend/internal/domain/verification/memory"
	"github.com/afaridanquah/verifylab-backend/internal/params"
	"github.com/mercari/go-circuitbreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type VerificationService struct {
	verifications verification.Respository
	cb            *circuitbreaker.CircuitBreaker
}

const otelName = "github.com/afaridanquah/verifylab-backend/internal/domain/verification/service"

type VerificationServiceConfig func(*VerificationService) error

func New(logger *slog.Logger, cfgs ...VerificationServiceConfig) (*VerificationService, error) {
	var ser = &VerificationService{}
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
			logger.Info("state changed",
				slog.String("old", string(oldState)),
				slog.String("new", string(newState)),
			)
		}),
	)

	return ser, nil
}

func WithVerficationRepo(vr verification.Respository) VerificationServiceConfig {
	return func(s *VerificationService) error {
		s.verifications = vr
		return nil
	}
}

func WithMemoryVerificationRepository() VerificationServiceConfig {
	mr, _ := memory.New()

	return WithVerficationRepo(mr)
}

func (vs *VerificationService) CreateVerification(ctx context.Context, req params.CreateVerificationRequest) (verification.Verification, error) {
	defer newOTELSpan(ctx, "Verification.Create").End()

	err := req.Validate()

	log.Printf("error from validation: %v", err)

	if err != nil {
		return verification.Verification{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Validate")
	}

	return verification.Verification{}, nil
}

func newOTELSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)

	return span
}
