package verification

import (
	"context"
	"log"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/internal"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/domain/customer"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/domain/verification"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/domain/verification/memory"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/domain/verification/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/params"
	ivo "bitbucket.org/msafaridanquah/verifylab-service/internal/valueobject"
	"github.com/mercari/go-circuitbreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	verifications verification.Repository
	customers     customer.Repository
	cb            *circuitbreaker.CircuitBreaker
}

const otelName = "bitbucket.org/msafaridanquah/verifylab-service/internal/domain/verification/service"

type ServiceConfig func(*Service) error

func New(logger *slog.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{}
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

func WithVerficationRepo(vr verification.Repository) ServiceConfig {
	return func(s *Service) error {
		s.verifications = vr
		return nil
	}
}

func WithMemoryVerificationRepository() ServiceConfig {
	mr, _ := memory.New()

	return WithVerficationRepo(mr)
}

func (s *Service) CreateVerification(ctx context.Context, req params.CreateVerificationRequest) (verification.Verification, error) {
	defer newOTELSpan(ctx, "Verification.Create").End()

	err := req.Validate()
	if err != nil {
		return verification.Verification{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Validate")
	}

	vt, err := valueobject.NewVerificationType(req.VerificationType)
	if err != nil {
		return verification.Verification{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Validate Verification Type")
	}

	verID := ivo.NewID("ver")

	if req.CustomerId != "" {
		parsedCustomerID, err := ivo.ParseID(req.CustomerId)

		if err != nil {
			return verification.Verification{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Customer ID Parsing")
		}

		cus, err := s.customers.Find(ctx, parsedCustomerID)

		log.Printf("cus %v", cus)

		if err != nil {
			return verification.Verification{}, internal.WrapErrorf(err, internal.ErrorCodeNotFound, "customer.Find")
		}

		person, err := valueobject.NewPerson(cus.ID().String(), cus.FirstName(), cus.MiddleName(), cus.Lastname(), cus.DateOfBirth(), cus.Email(), cus.Country())
		if err != nil {
			return verification.Verification{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "person.Create")
		}

		ver, err := verification.New(verID, vt, person)
		if err != nil {
			return verification.Verification{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "verification.Create")
		}

		err = s.verifications.Add(ctx, *ver)
		if err != nil {
			return verification.Verification{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Create")
		}

		return *ver, nil
	}

	return verification.Verification{}, nil
}

func newOTELSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)

	return span
}
