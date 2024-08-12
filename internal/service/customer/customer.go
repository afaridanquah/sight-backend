package customer

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/afaridanquah/verifylab-backend/internal"
	"github.com/afaridanquah/verifylab-backend/internal/domain/customer"
	"github.com/afaridanquah/verifylab-backend/internal/domain/customer/memory"
	"github.com/afaridanquah/verifylab-backend/internal/domain/customer/valueobject"
	"github.com/afaridanquah/verifylab-backend/internal/params"
	ivo "github.com/afaridanquah/verifylab-backend/internal/valueobject"
	"github.com/mercari/go-circuitbreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type CustomerService struct {
	customers customer.Repository
	cb        *circuitbreaker.CircuitBreaker
}

const otelName = "github.com/afaridanquah/verifylab-backend/internal/domain/customer/service"

type CustomerServiceConfig func(*CustomerService) error

func New(logger *slog.Logger, cfgs ...CustomerServiceConfig) (*CustomerService, error) {
	var ser = &CustomerService{}

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

func WithRepository(cr customer.Repository) CustomerServiceConfig {
	return func(cs *CustomerService) error {
		cs.customers = cr
		return nil
	}
}

func WithMemoryRepository() CustomerServiceConfig {
	mr := memory.New()

	return WithRepository(mr)
}

func (cs *CustomerService) CreateCustomer(ctx context.Context, req params.CreateCustomerRequest) (customer.Customer, error) {
	defer newOTELSpan(ctx, "Customer.Create").End()

	err := req.Validate()

	log.Printf("error from validation: %v", err)

	if err != nil {
		return customer.Customer{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Validate")
	}
	country, err := ivo.NewCountry(req.Country)
	if err != nil {
		return customer.Customer{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "params.Validate")
	}

	id := valueobject.NewID()
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

func (cs *CustomerService) FindCustomer(ctx context.Context, id string) (customer.Customer, error) {
	defer newOTELSpan(ctx, "Customer.Find").End()

	custID, err := valueobject.ParseID(id)
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
