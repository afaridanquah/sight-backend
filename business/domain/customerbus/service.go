package customerbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"github.com/google/uuid"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	customers Repository
	cb        *circuitbreaker.CircuitBreaker
}

type ServiceConfig func(*Service) error

func New(customers Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{
		customers: customers,
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

func WithRepository(cr Repository) ServiceConfig {
	return func(s *Service) error {
		s.customers = cr
		return nil
	}
}

func (cs *Service) Create(ctx context.Context, nc NewCustomer) (Customer, error) {
	ctx, span := otel.AddSpan(ctx, "customerbus.service.create")
	defer span.End()
	now := time.Now()

	userID := uuid.New()
	businessID := uuid.New()

	customer := Customer{
		ID:              uuid.New(),
		Person:          nc.Person,
		UserID:          userID,
		BusinessID:      businessID,
		DateOfBirth:     nc.DateOfBirth,
		Email:           nc.Email,
		BirthCountry:    nc.BirthCountry,
		Address:         nc.Address,
		Identifications: nc.Identifications,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	_, err := cs.customers.Add(ctx, customer)

	if err != nil {
		return Customer{}, ierr.WrapErrorf(err, ierr.ErrorCodeUnknown, "repo.Create")
	}

	return customer, nil
}

// func (cs *Service) FindByID(ctx context.Context, id string) (customerbus.Customer, error) {
// 	defer itel.NewOTELSpan(ctx, otelName, "Customer.Find").End()

// 	custID, err := uuid.Parse(id)
// 	if err != nil {
// 		return customerbus.Customer{}, ierr.WrapErrorf(err, ierr.ErrorCodeUnknown, "Find")
// 	}

// 	cus, err := cs.customers.Find(ctx, custID)

// 	if err != nil {
// 		return customerbus.Customer{}, ierr.WrapErrorf(err, ierr.ErrorCodeUnknown, "repo.Find")
// 	}

// 	return cus, nil
// }

// func (cs *Service) All(ctx context.Context) ([]customerbus.Customer, error) {
// 	defer itel.NewOTELSpan(ctx, otelName, "Customer.All").End()

// 	results, err := cs.customers.All(ctx)

// 	if err != nil {
// 		return []customerbus.Customer{}, ierr.WrapErrorf(err, ierr.ErrorCodeUnknown, "repo.All")
// 	}

// 	return results, nil
// }
