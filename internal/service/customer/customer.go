package customer

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/internal"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/domain/customer"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/domain/customer/memory"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/params"
	ivo "bitbucket.org/msafaridanquah/verifylab-service/internal/valueobject"
	"github.com/mercari/go-circuitbreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	customers customer.Repository
	cb        *circuitbreaker.CircuitBreaker
}

const otelName = "bitbucket.org/msafaridanquah/verifylab-service/internal/domain/customer/service"

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

func WithRepository(cr customer.Repository) ServiceConfig {
	return func(cs *Service) error {
		cs.customers = cr
		return nil
	}
}

func WithMemoryRepository() ServiceConfig {
	mr := memory.New()

	return WithRepository(mr)
}

func (cs *Service) CreateCustomer(ctx context.Context, req params.CreateCustomerRequest) (customer.Customer, error) {
	defer newOTELSpan(ctx, "Customer.Create").End()

	err := req.Validate()
	if err != nil {
		return customer.Customer{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Validate")
	}

	country, err := ivo.NewCountry(req.Country)
	if err != nil {
		return customer.Customer{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Validate")
	}

	id := ivo.NewID("cus")
	cus, err := customer.New(id, req.FirstName, req.LastName, country)
	if err != nil {
		return customer.Customer{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Validate")
	}

	if req.DateOfBirth != "" {
		dob, err := ivo.NewDateOfBirth(req.DateOfBirth)
		if err != nil {
			return customer.Customer{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Validate")
		}
		cus.WithDateOfBirth(dob)
	}

	if req.Email != "" {
		email, err := ivo.NewEmail(req.Email)
		if err != nil {
			return customer.Customer{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Validate")
		}
		cus.WithEmail(email)
	}

	err = cs.customers.Add(ctx, *cus)
	if err != nil {
		return customer.Customer{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Create")
	}

	return *cus, err
}

func (cs *Service) FindCustomer(ctx context.Context, id string) (customer.Customer, error) {
	defer newOTELSpan(ctx, "Customer.Find").End()

	custID, err := ivo.ParseID(id)
	if err != nil {
		return customer.Customer{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "Find")
	}

	cus, err := cs.customers.Find(ctx, custID)

	if err != nil {
		return customer.Customer{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo.Find")
	}

	return cus, nil

}

func newOTELSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)

	return span
}
